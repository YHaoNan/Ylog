---
title: VSCode 扩展开发(七) 实战-翻译插件
date: 2019-06-08 11:35:54
tags: [vscode,editor]
categories: VSCode
---

## 学习目标
本章和接下来一章都是实战相关的插件。

本章主要制作一个翻译插件，提供悬停翻译和划词翻译，整篇翻译和注释翻译，并且会学到一些TextEditor的API。

## 注册API服务
首先我们使用百度翻译API进行翻译工作，所以请前往[百度翻译开放平台](http://api.fanyi.baidu.com/api/trans/product/index)去注册开发者账号。

注册后到管理控制台->开发者信息中去查看APP Id和秘钥。


通用翻译API HTTP地址：

[http://api.fanyi.baidu.com/api/trans/vip/translate](http://api.fanyi.baidu.com/api/trans/vip/translate)

通用翻译API HTTPS地址：

[https://fanyi-api.baidu.com/api/trans/vip/translate](https://fanyi-api.baidu.com/api/trans/vip/translate)

API参数：

<table class="info-table">
<tbody><tr>
<th>字段名</th>
<th>类型</th>
<th>必填参数</th>
<th>描述</th>
<th class="last-co">备注</th>
</tr>
<tr>
<td>q</td>
<td>TEXT</td>
<td>Y</td>
<td>请求翻译query</td>
<td>UTF-8编码</td>
</tr>
<tr>
<td>from</td>
<td>TEXT</td>
<td>Y</td>
<td>翻译源语言</td>
<td><a href="#languageList" class="doc-nav-item" data-index="1">语言列表</a>(可设置为auto)</td>
</tr>
<tr>
<td>to</td>
<td>TEXT</td>
<td>Y</td>
<td>译文语言</td>
<td><a href="#languageList" class="doc-nav-item" data-index="1">语言列表</a>(不可设置为auto)</td>
</tr>
<tr>
<td>appid</td>
<td>INT</td>
<td>Y</td>
<td>APP ID</td>
<td>可在<a class="goto-console" href="/api/trans/product/desktop?req=developer">管理控制台</a>查看</td>
</tr>
<tr>
<td>salt</td>
<td>INT</td>
<td>Y</td>
<td>随机数</td>
<td></td>
</tr>
<tr>
<td>sign</td>
<td>TEXT</td>
<td>Y</td>
<td>签名</td>
<td>appid+q+salt+密钥 的MD5值</td>
</tr>
</tbody></table>

签名生成方法如下：

1. 将请求参数中的 APPID(appid), 翻译query(q, 注意为UTF-8编码), 随机数(salt), 以及平台分配的密钥(可在管理控制台查看)按照 appid+q+salt+密钥 的顺序拼接得到字符串1。

2. 对字符串1做md5，得到32位小写的sign。


知道了这些，我们就可以开始编写代码了，新建一个项目，地址在[这里](https://github.com/YHaoNan/vscode-tutorial/tree/master/vsc-extensions-tutorial-7)

返回的格式请自己查看API吧。
## 编写代码
我们先把请求的库写好，我们统一使用Https的POST接口，然后使用axios库进行请求。

请在工程目录下执行`npm install axios --save`安装axios库。

因为appid和秘钥，还有目标语言都需要用户进行设置，所以我们编写`package.json`提供设置项：
```json
"configuration": [{
    "title": "Translator",
    "properties": {
        "ts.to": {
        "type": "string",
        "default": "zh",
        "description": "The language that you wanna translate to"
        },
        "ts.appid": {
        "type": [
            "integer",
            "null"
        ],
        "default": null,
        "description": "Your appid"
        },
        "ts.key": {
        "type": [
            "string",
            "null"
        ],
        "default": null,
        "description": "Your key"
        }
    }
}]
```
然后提供一个util工具类，用于获取设置、生成随机数字和加密：

```typescript
import * as vscode from 'vscode'
import {Md5} from "md5-typescript";
export let getConfig = <T>(key: string): T=><T>vscode.workspace.getConfiguration().get<T>(key);

//生成指定长度的随机数字
export let genRandomNumber = (len: number): number=>{
    return getRandomInt(Math.pow(10,len-1),Math.pow(10,len)-1);
}

//生成Sign，也就是加密信息
export let genSign = (appid: string,q: string,salt: string,key:string): string => Md5.init(appid+q+salt+key);

function getRandomInt(min: number, max: number): number {  
    var Range = max - min;  
    var Rand = Math.random();  
    return(min + Math.round(Rand * Range));  
}
```

编写请求库
```typescript
import axios from 'axios'
import * as utils from './utils'

//返回码映射到错误信息
const ERROR_MSG: {[index:string]:string}= {
    '52001': '请求超时，请重试',
    '52002': '系统错误',
    '52003': '请检查您的appid是否正确',
    '54000': '参数不完整',
    '54001': '签名错误',
    '58000': '客户端IP非法',
    '54005': '长query请求过于频繁，请稍后再试',
    '58001': '译文语言不支持',
    '58002': '您的服务已关闭'
}

const API_URL = 'https://fanyi-api.baidu.com/api/trans/vip/translate'

export async function translate(source:string): Promise<string>{
    //准备数据
    let data:{[index:string]:string} = {
        q: source,
        from: 'auto',
        to: utils.getConfig('ts.to'),
        appid: utils.getConfig<string>('ts.appid'),
        salt: utils.genRandomNumber(10).toString(),
        sign: ''
    }
    

    data.sign = utils.genSign(data.appid,data.q,data.salt,utils.getConfig<string>('ts.key'));

    //设置header为 form-urlencoded
    axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
    //设置transformRequest用于格式化数据 转成querystring 因为API规定不能直接用JSON格式的数据所以要如此苦逼的写这些
    axios.defaults.transformRequest = [function (data) {
        let ret = '';
        for (let it in data) {
        ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&';
        }
        return ret;
    }];

    try{
        //发送请求
        const result = await axios.post(API_URL,data);
        
        //如果有结果 则返回
        if(result.data.trans_result){
            let results: any[] = result.data.trans_result;
            let targetStr:string = results.map(obj=>obj.dst).join('\n');
            return Promise.resolve(targetStr);
        //如果没有结果 检查有没有错误码 有的话调用reject
        }else if(result.data.error_code){
            let msg = ERROR_MSG[result.data.error_code as string];
            msg=msg?msg:'未知错误';
            return Promise.reject(msg);
        }
    }catch(e){
        return Promise.reject(e.message);
    }
}
```
编写`extensions.ts`，暂时提供两种翻译手段，一个是通过`HoverProvider`提供悬停事件，在悬停的单词上翻译，一个是通过命令翻译当前编辑器的选中部分。
```typescript
import * as vscode from 'vscode';
import {TranslateHoverProvider} from './hover'
import {translate} from './request'
export function activate(context: vscode.ExtensionContext) {
    //注册Hover
	context.subscriptions.push(
		vscode.languages.registerHoverProvider('*',new TranslateHoverProvider())
	);
    //注册翻译选中范围的命令
	context.subscriptions.push(
		vscode.commands.registerTextEditorCommand('ns.translateSelection', async(editor)=>{
            let outputChannel = vscode.window.createOutputChannel('Translate Result');
            //兼容多光标
			editor.selections.forEach(async(selection)=>{
				if(!selection.isEmpty){
                    let source = editor.document.getText(selection);
					translateAndShow(outputChannel,source);
				}
			});
		})
	);
}

async function translateAndShow(outputChannel:vscode.OutputChannel,source:string){
	try{
		let result = await translate(source);
		outputChannel.appendLine(result);
		outputChannel.show();	
	}catch(e){
		vscode.window.showErrorMessage('翻译器：'+e);
	}
}
export function deactivate() {}
```

editor的`selections`是一个`Selection`对象的数组，记录的是所有用户选中的地方。`Selection`是`Range`的子类，`Range`表示一个范围，它有`start`和`end`属性，分别代表起始位置和结束位置，这两个属性是`Position`类型，它代表一个位置，有`line`和`character`属性，分别表示这个位置所在的行和字符在本行的位置。

所以我们这里判断每一个selection是否是空，如果不是就翻译那一段文字。

请大家自行在`package.json`中注册，我们的翻译插件要求vscode打开时就激活，所以在`package.json`中写：
```json
"activationEvents": [
    "*"
]
```
写星之后打开VSCode直接激活插件。

## 扩展
添加翻译全文功能和翻译注释并替换原注释功能。

翻译全文功能比较简单了，我们当用户没有选中段落的时候认为它想翻译全文，你可能会这样判断：
```typescript
if(editor.selections.length == 0)
```
或者
```typescript
//selection是当前的选中区域对象
if(editor.selection == null)
```
可是这是不行的，因为上面说了，就算没选中也会有一个为空的selection，而且我们不能用`editor.selection`判断，因为用户可能有多个选中区域呢。

所以我们打算在`editor.selections.forEach`里记录是否有被翻译的段落，如果没有那么就是全文了。

```typescript
let isHasTranslatedSelection = false;
editor.selections.forEach(async(selection)=>{
    if(!selection.isEmpty){
        let source = editor.document.getText(selection);
        translateAndShow(outputChannel,source);
        isHasTranslatedSelection = true;
    }
});
if(!isHasTranslatedSelection){
    translateAndShow(outputChannel,editor.document.getText());
}
```

翻译全文很简单也很实用，下面一个就比较难了，下面是翻译注释并替换，这个也是很实用的功能，尤其在使用别人的库的时候，但这在不同的语言中需要不同的行为，因为不同的语言的注释块是不一样的。

比如Python你需要兼容如下的注释：
```python
# Single Line Comment

'''
Multi Line Comment
'''

"""
Multi Line Comment
"""
```

而在Java中你需要兼容如下的注释：
```java
// Single Line Comment

/**
* Multi Line Comment
*/
```
其实大部分语言都是用`/**/`实现多行注释的，Java也一样，上面的Java多行注释只是Java程序员约定俗成的习惯而已。

我们的示例仅仅支持单行注释，其余的留给各位自行扩展，哈哈其实我觉得常见到的应该是多行。

对于单行注释，我们有如下的正则：
```text/plain
(\/\/|#)(.*)
```
```typescript
context.subscriptions.push(
    vscode.commands.registerTextEditorCommand('ns.translateComment', async(editor)=>{
        let outputChannel = vscode.window.createOutputChannel('Translate Result');
        for(let i=0;i<editor.document.lineCount;i++){
            let lineText = editor.document.lineAt(i).text;
            let commentMatch = lineText.match('(\/\/|#)(.*)');
            if(commentMatch){
                let result = await translate(commentMatch[2]);
                let start = lineText.indexOf(commentMatch[2]);
                let end = start + commentMatch[2].length;
                editor.edit(editorEdit=>{
                    editorEdit.replace(new vscode.Range(new vscode.Position(i,start),new vscode.Position(i,end)),result);
                })
        }
    }})
);
```

上面的代码我们使用了`editor.edit`来替换编辑器中的内容，`editor.edit`传入一个`callback`，这个回调方法里有个参数是`TextEditorEdit`类型的，它有`delete`、`insert`、`replace`等方法可以让我们操作编辑器。

好了，差不多了，今天就到这里了。

其实这离一个成熟的翻译插件还差得远呢，比如不能根据语言识别不同的注释，注释需要发送多次请求等，代码复用性不高等，这些都是问题。不过我们的目的不是开发插件，是学开发插件，哈哈，真正的成熟的插件还得靠您各位以后慢慢的开发呐～

**[点此返回目录](/post/vscode-ext1#目录)**
