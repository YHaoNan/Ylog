---
title: VSCode 扩展开发(八) 实战-Markdown图片上传
date: 2019-06-18 8:35:54
tags: [vscode,editor]
categories: VSCode
---
## 学习目标
本章主要制作一个Markdown图片上传插件，使用UCloud的对象存储作为服务端。

大概是这样的效果
![效果](https://s2.ax1x.com/2019/06/18/VbHBuD.gif)

## 开篇寒暄
第一篇中说的，我们在编写Markdown的时候，图片上传是一个问题，如果不整合到编辑器中的话会很麻烦，总要在编辑器和图片上传窗口中跳来跳去，还需要手动复制上传后的URL，并且粘到编辑器中。  
先不说是否浪费时间，首先这用户体验就不是很好。所以我们编写一个自动上传图片的工具。

## 注册
首先去UCloud中注册并通过实名认证。

然后去开通对象存储服务，在`全部产品->存储->对象存储 UFile`中，然后选择`创建存储空间`，创建一个公开的空间，你填写的存储空间域名就是你的桶名。

然后点击创建就好了，UCloud每月有免费的20G流量和存储空间，并且给每个账户50元的金额，个人记录笔记的博客足够用了。

创建好后，点击你的头像，然后里面有一个`API秘钥`，进去之后获取你的公钥和私钥。

至此，注册工作全部完成。

## 获取Node SDK
UFile的NodeSDK在这里`https://github.com/ufilesdk-dev/ufile-nodejssdk`

下载下来，有三个目录：
```
./ufile     -- UFile的nodejs SDK，其依赖的模块可以查看./ufile/package.json
./sdk_test  -- 几个借助nodejs SDK的测试脚本
./raw_test  -- 几个不借助nodejs SDK的测试脚本，适合阅读人群：需要探究细节，采用原生HTTP请求
```
我们只需要`ufile`这个文件夹，把他复制到我们的项目中：
```bash
cp -r ./ufile ~/vsc-extensions-tutorial-8 
```
项目[git地址](https://github.com/YHaoNan/vscode-tutorial/tree/master/vsc-extensions-tutorial-8)

按理来说，我们需要运行`npm`来将这个安装到项目中，但是还不行，我们先分析下这个sdk。
```text/plain
├── ufile
│   ├── auth_client.js
│   ├── config.json
│   ├── helper.js
│   ├── http_request.js
│   ├── index.js
│   ├── multi_upload.js
│   ├── package.json
│   └── user_agent.json
```
这是sdk的整个目录结构，也挺简单，不过作为sdk的用户，我们肯定一眼就看到了和配置相关的`config.json`，看看里面有什么内容。

```json
{
    "ucloud_public_key" : "你的公钥",
    "ucloud_private_key" : "你的私钥",
    "proxy_suffix" : ".cn-bj.ufileos.com"
}
```
每个人下载的sdk中都把各自的公钥和私钥直接写进去了，还有域名后缀，我们不希望这三个参数是死的，因为如果是死的，别人用我们编写的插件时就会把图片上传到你的UFile服务中，这是不合理的。所以我们要把这三个参数留出来，让用户配置。

既然官方提供了config文件，那就说明SDK的代码里肯定去读取了这些配置，我们在插件的代码里不方便修改这个配置文件，对于插件来说，如果SDK的上传方法能提供几个参数让我们直接把公钥和私钥什么的填写进去，对我们来说是个好事。

但是显然官方没有提供，就需要我们自己阅读并修改官方的SDK代码了，我们在sdk的`auth_client.js`中发现了这几行：
```javascript
var config = require('./config');

var AuthClient = function(req) {
    this.req = req;
    this.public_key = config['ucloud_public_key'];
    this.private_key = config['ucloud_private_key'];
    this.proxy_suffix = config['proxy_suffix'];
    this.authorization = "";
}
```
因为json就是js的原生对象，可以直接用`require`引入，然后构造AuthClient的时候去引用了这里面的值，通过字面意义和官方文档，知道了这个AuthClient就是用作和服务端做认证的一个类。

所以我们直接修改这个类的定义，添加几个参数，而不是读取`config`即可。

```javascript
var AuthClient = function(req,public_key,private_key,proxy_suffix) {
    this.req = req;
    this.public_key = public_key;
    this.private_key = private_key;
    this.proxy_suffix = proxy_suffix;
    this.authorization = "";
}
```
现在我们就可以通过构造AuthClient传入各种认证需要的参数了，并且整个auth_client.js中对`config`没有别的引用了，所以可以把`var config = require('./config');`这行代码删除了。

一切完成，我们开始向项目中安装这个库。
```bash
npm install ./ufile --save
```
安装完成，现在可以在我们的项目中引用这个库了。

## 编写代码
我们打算提供图片上传并自动向编辑器中写入URL，图片压缩、备份功能。

首先我们需要提供个命令让用户选择图片，然后才能进行下一步的上传操作：


```typescript
context.subscriptions.push(
    vscode.commands.registerTextEditorCommand('lilpig.pickandupload',(editor=>{
        vscode.window.showOpenDialog(options).then(uris=>{
            if(uris&&uris.length>0){
                vscode.window.showInformationMessage(uris[0].fsPath);
            }
        })
    }))
);
```
`package.json`中的注册代码请自行编写。

现在当你调用这个命令时，就会打开文件选择窗口，并且当你选中一个文件就会弹出这个文件的路径。

下一步就是压缩和上传图片，在这之前，先给用户提供一些设置项：
```json
"img.ucloud_private_key":{
    "type":"string",
    "description": "The private key of your ucloud account."
},
"img.ucloud_public_key":{
    "type":"string",
    "description": "The public key of your ucloud account."
},
"img.bucket_name":{
    "type":"string",
    "description": "Your bucket name."
},
"img.domain":{
    "type":"string",
    "description": "The domain of your ufile space."
},
"img.backupdir":{
    "type":"string",
    "description": "The path that the backup of the image your uploaded."
},
"img.imghandlerpath":{
    "type":"string",
    "description": "The path of the script that can compress your pick image." 
}
```
然后我们就可以在插件中读取这些配置了。

关于压缩，由于没在Node上看到一个好用的，跨平台的图片压缩库，我们决定使用Python实现，在本项目下有一个叫`img-handler.py`的python脚本，用它可以实现图片的压缩，对应参数为：
```
python3 img-handler.py 原始路径 压缩等级1-100 输出路径
```
因为jpg的压缩效果比较好，所以所有的图片都会被压缩成jpeg，该脚本依赖于`pillow`。

脚本返回格式是`返回码:返回信息`，压缩成功的返回是`200:图片存储路径`，失败的则是`错误码:错误信息`。

该脚本在[这里](https://github.com/YHaoNan/vscode-tutorial/blob/master/vsc-extensions-tutorial-8/src/img-handler.py)

再来看看Ufile的SDK使用方法：
```
//先使用HttpRequest构造一个请求信息，参数是协议，上传后存储的路径，桶名，上传图片的名字，上传的图片在本地的路径
const req = new HttpRequest('POST','/',bucketName,upload_name,local_path);
//然后创建一个AuthClient，参数分别是请求对象，公钥，私钥和域名
const client = new AuthClient(req,ucloudPublicKey,ucloudPrivateKey,ucloudDomain);

client.SendRequest((res:any)=>{
    console.log(res);
    if(res instanceof Error){
        //上传错误 代码执行错误
        return;
    }
    if(res.statusCode != 200){
        //上传错误 返回了错误信息
        return;
    }
    //否则上传成功
});
```

然后我们就可以编写代码了，先编写一个简单的工具类读取配置信息：

```typescript
import * as vscode from 'vscode'
export let getConfig = <T>(key: string): T=><T>vscode.workspace.getConfiguration().get<T>(key);
```

然后编写`img-upload.ts`用于处理图像并上传：
```typescript
const HttpRequest = require('ufile').HttpRequest;
const AuthClient = require('ufile').AuthClient;
import {exec} from 'child_process'
import {getConfig} from './utils'
import * as vscode from 'vscode'
import { pathToFileURL } from 'url';
import {join} from 'path';

export function handle(source:string,compressLevel:number,callback:(err:string|undefined,data:string|undefined)=>void){
    //读取相关配置
    let handler_path = <string>getConfig('img.imghandlerpath');
    let ucloudPublicKey = <string>getConfig('img.ucloud_public_key');
    let ucloudPrivateKey = <string>getConfig('img.ucloud_private_key');
    let bucketName = <string>getConfig('img.bucketname');
    let ucloudDomain = <string>getConfig('img.domain');
    let backupDir = <string>getConfig('img.backupdir');

    //生成图片上传后的名字 img-unix时间戳.jpg
    let upload_name = 'img-'+new Date().getTime()+'.jpg';

    //执行压缩脚本的pyhton命令 完整命令是 python3 压缩脚本路径 待上传的文件路径 压缩等级 备份文件夹/$upload_name
    //用join是因为不同平台的文件夹分隔符不一样
    let command = 'python3 "'+handler_path+'" "'+ source+'" '+compressLevel+' "'+join(backupDir,upload_name)+'"';

    //执行命令
    exec(command,(err,stdout,stderr)=>{
        if(err){
            callback("调用压缩脚本出错："+stderr,undefined);
            return;
        }
        if(stdout){
            let upload_img_path = '';
            let result = stdout.split(':');
            if(result.length<2){
                callback("脚本返回结果出错，返回："+stdout,undefined);
                return;
            }
            //如果返回的是正确结果
            if(result[0]=='200'){
                //获取压缩后的图片路径
                upload_img_path = result[1];
                console.log(upload_name,upload_img_path);
                //构造请求
                const req = new HttpRequest('POST','/',bucketName,upload_name,backupDir+'/'+upload_name);
                const client = new AuthClient(req,ucloudPublicKey,ucloudPrivateKey,ucloudDomain);
                client.SendRequest((res:any)=>{
                    console.log(res);
                    if(res instanceof Error){
                        callback('上传错误：'+res.message,undefined);
                        return;
                    }
                    if(res.statusCode != 200){
                        callback('上传错误：'+res.statusCode,undefined);
                        return;
                    }
                    callback(undefined,'http://'+bucketName+ucloudDomain+'/'+upload_name);
                    
                });
            }else{
                callback("压缩错误："+result[1],undefined);
            }
        }
    });
    
}
```
`img-upload.ts`中向外提供了一个handle函数用于压缩并上传图像，并要求传入一个callback，这个callback和nodejs大部分的callback差不多，就是`(err,data)=>{}`这种的，如果有错误，err就有值，data没有值，如果没有错误，err无值，data有值。

这样就编写好了，然后再在`extension.ts`中调用该方法：
```typescript
context.subscriptions.push(
    vscode.commands.registerTextEditorCommand('lilpig.pickandupload',(editor=>{
        vscode.window.showOpenDialog(options).then(uris=>{
            if(uris&&uris.length>0){
                vscode.window.showInputBox({prompt:'压缩等级1-99',value:'50'}).then(str=>{
                    handle(uris[0].fsPath,parseInt(str as string),(err,data)=>{
                        if(err){
                            vscode.window.showErrorMessage(err);
                        }else{
                            editor.insertSnippet(new vscode.SnippetString('![${1:图片}]('+data+')'));
                        }
                    })
                })
            }
        })
    }))
);
```
这括号真吓人，回调地狱啊哈哈哈哈。

这样就写好了，你还可以根据前面学到的知识扩展它，比如加个加载条等。

**[点此返回目录](/post/vscode-ext1#目录)**
