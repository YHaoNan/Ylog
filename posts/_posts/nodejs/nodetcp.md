---
title: NodeJS中的网络操作
date: 2019-07-17 08:49:31
tags: [nodejs]
categories: NodeJS
---

唉，睁眼睛看一会书就要去练车的感觉真不怎么样，没啥时间干自己的事，本来计划好的假期学习计划几乎被全盘推翻。哈哈，还好十九号就考科二了，考完之后就不会去得那么勤了。

这几天在驾校无聊等待的过程中一直在刷《NodeJS高级编程》，然后回家还要把所有活都干完后，有空了才能把代码实现一遍，所以也一直没写笔记，今天补点，嗯补一点。

## 构建TCP服务器
nodejs的事件发射器模式让本来在其它语言中生涩难懂的socket API变得简单灵活起来。

在nodejs中进行网络操作，你需要引入`net`包，下面创建一个最简单的回声服务器。

`echo_server.js`

```javascript
require('net').createServer(socket=>{
    socket.on('data',data=> socket.write(data))
}).listen(4000);
```

使用`net`模块中的`createServer`方法创建一个服务器对象，然后传入一个回调，这个回调会在有连接接入时被回调，它的参数就是连接的socket对象，下面我们监听这个对象的`data`事件，当接受到客户端传来的数据时，就会发射`data`事件，它的回调方法的参数是发送过来的数据，然后我们通过socket对象的`write`方法把接收到的数据发送回去。

最后别忘了调用服务器对象的`listen`方法才实际的把服务绑定在了4000端口上进行监听，这样客户端才能从4000端口连接进来。

我们使用`nc`进行测试：
```bash
~ » nc 127.0.0.1 4000
Hello
Hello
What is your name
What is your name
```

## 创建一个mini的聊天服务器
```javascript
var socketList = [];
require('net').createServer((socket)=>{
    
    socketList.push(socket);
    socket.write('Welcome to the chat room , online : '+socketList.length+'\n');
    sendBroadCastMessage(socket,socket.address()['address']+' join the room!'+'\n');

    socket.on('data',(data)=>{
        sendBroadCastMessage(socket,socket.address()['address']+' say :'+data);
    })
}).listen(10404);

function sendBroadCastMessage(besideSocket,message){
    socketList.forEach((s)=>{
        if(besideSocket!=s)
            s.write(message);
    })
}
```
上面代码创建了一个聊天服务器，不管它代表了什么意思，我们先用nc测试一下：
`用户1`

```bash
/mnt/d/blog(dev*) » nc 127.0.0.1 10404
Welcome to the chat room , online : 1
::ffff:127.0.0.1 join the room!
::ffff:127.0.0.1 say :Hi
Helo~
```

`用户2`

```bash
~ » nc 127.0.0.1 10404
Welcome to the chat room , online : 2
Hi
::ffff:127.0.0.1 say :Helo~
```

这个聊天室确实可以运行，只不过因为我都是在本地连接，所以所有人的ip地址都是一样的。。。啊哈哈。

看上面的代码，我们维护了一个`socketList`来保存所有的套接字，每当进来一个连接就push到这个列表中。我还编写了一个`sendBroadCastMessage`方法来广播数据，只不过很多时候广播数据需要排除某个人，比如在第6行的代码，我们向所有房间里的人广播有一个新用户进入房间时就要排除这个新用户，所以`sendBroadCastMessage`方法的第一个参数理所当然是排除的那个套接字对象。

每当从某个套接字中接收到数据时，也就是说接收到`data`事件时，向所有人转发这条消息，当然排除发送消息的用户。

## 构建tcp客户端实现更强大的聊天室
```javascript
const printUsage = require('./usage').printUsage;

const host = '127.0.0.1';
const port = 4000;  
var conn = require('net').createConnection(port,host,()=>{
    process.stdin.resume();
    process.stdin.on('data',data=>{
        data = convertSingleLine(data.toString());
        if(data=='/quit'||data=='/exit'){
            process.stdin.pause();
            conn.end();
            console.log('Bye!Have a nice day!')
        }else if(data == '/help'){
            printUsage();
        }else
            conn.write(data);
    });
});

conn.on('data',data=>{
    process.stdout.write(data);
});
conn.on('error',err=>{
    console.log('Got an error : '+err.message);
});
conn.on('close',()=>{
    console.log('Bye!');
    process.stdin.pause();
});

function convertSingleLine(string){
    return string.replace(/[\r\n]/g,"");
}
```

上面的代码创建了一个TCP客户端，使用`createConnection`方法，当连接成功后回调参数被调用。

我们通过process.stdin来接受用户的输入，并且发送给服务器。

`createConnection`返回了一个连接对象，我们监听这个连接对象发出的事件，它的事件其实和tcp服务器差不多。

接下来放上服务器端的代码：
```javascript
const getCommand = require('./command').getCommand;
const commandCallback = {
    rename:(socket,newnick)=>{
        newnick = newnick.replace(/\ +/g,"");
        newnick = convertSingleLine(newnick);
        var i = connectionList.indexOf(socket);
        if(nickNameList.indexOf(newnick)==-1 && i!=-1 && newnick!=''){
            var oldName = nickNameList[i];
            nickNameList[i] = newnick;
            sendMessage(socket,'Your name has changed!');
            postMessage(socket,oldName +' changed his/her name to '+newnick);
        }else{
            sendMessage(socket,'Change name faild,may be this name is already used or your input in not correct!')
        }
    },
    gname:(socket,arg)=>{
        socket.write('Your name is '+getNameBySocket(socket)+'\n');
    },
    kickout:(foo,name)=>{
        var index = nickNameList.indexOf(name);
        nickNameList.splice(index,1);
        var socket = connectionList.splice(index,1)[0];
        sendMessage(socket,'You were kick out from the room!');
        socket.end();
        postMessage(socket,name+' was kick out from the room!');
    },
    list:(foo,bar)=>{
        nickNameList.forEach(str=>{
            console.log(str);
        })
    }
}
var connectionList = [];
var nickNameList = [];
const server = require('net').createServer(socket=>{
    var tempName = Date.now();
    connectionList.push(socket);
    nickNameList.push(tempName);
    sendMessage(socket,'Welcome to the chatroom!\nYour default nick name is '+tempName+'.\nType "help" to show the usage!\n'+'Online:'+connectionList.length);
    postMessage(socket,tempName+' join the room!');

    socket.on('data',data=>{
        var command = getCommand(data.toString());
        if(command!=null){
            commandCallback[command[0]](socket,command[1]);
        }else{
            postMessage(socket,getNameBySocket(socket)+' says : '+data);
        }
    })
    socket.on('close',()=>{
        var index = connectionList.indexOf(socket);
        if(index==-1)return;
        nickNameList.splice(index,1);
        var socket = connectionList.splice(index,1)[0];
        postMessage(socket,name+' leave from the chatroom!');
    });
    socket.on('error',err=>{
        console.log('Got an error : '+err.message);
    })
}).listen(4000);

function getSocketByName(name){
    var i = nickNameList.indexOf(name);
    return connectionList[i];
}
function getNameBySocket(socket){
    var i = connectionList.indexOf(socket);
    return nickNameList[i];
}
function postMessage(besides,message){
    connectionList.forEach(socket=>{
        if(socket != besides){
            socket.write(message+'\n');
        }
    })
}

function sendMessage(socket,message){
    socket.write(message+'\n');
}

function convertSingleLine(string){
    return string.replace(/[\r\n]/g,"");
}
server.on('listening',()=>{
    process.stdin.resume();
    process.stdin.on('data',data=>{
        var command = getCommand(convertSingleLine(data.toString()));
        if(command){
            commandCallback[command[0]](null,command[1]);            
        }else{
            console.log('Unknow command '+data);
        }
    })
})
```

其实只看`createServer`方法的回调中的内容，和最初的聊天服务器差不多，只是更加完整，加了错误处理什么的。

然后其他的代码基本都是为用户名系统和命令系统服务，这个高级聊天室支持使用命令系统，用户可以使用命令修改自己的用户名，服务器能使用命令踢人等等，不过用户名和命令并非一个聊天服务器必须存在的元素，所以不做解释，因为他们并没有那么重要。

然后就是其他的工具类的代码：
`command.js`
```javascript
"use strict";
exports.__esModule = true;
function getCommand(commandStr) {
    if (commandStr.charAt(0) != '/')
        return null;
    else {
        var firstWhiteSpaceIndex = commandStr.indexOf(' ');
        if (firstWhiteSpaceIndex != -1 && firstWhiteSpaceIndex < commandStr.length) {
            var command = commandStr.substr(1, firstWhiteSpaceIndex - 1);
            var arg = commandStr.substr(firstWhiteSpaceIndex + 1);
            return [command, arg];
        }
        return null;
    }
}
exports.getCommand = getCommand;
```

`usage.js`
```javascript
function printUsage(){
    console.log('Command:\n\t1. /rename <Name> -- Modify your nickname in chatroom.\n\t 2. /exit or /quit -- Exit from chatroom.\n\t 3. /help -- Show usage.\nAdmin Command:\n\t4. /gname <any string> -- Get your nickname.Note that the option <any string> is must exits.\n\t'+
    '1. /kickout <Name> -- Kickout a user by nickname.\n\t2. /close -- Close chatroom.\n\t3. /post <Message> -- Post some message.');
}
exports.printUsage = printUsage;
```

以下是该聊天服务器的测试：
`启动server`
```cmd
PS D:\learn\node\tcp_cs_test> node .\server.js
```

`user1`
```cmd
PS D:\learn\node\tcp_cs_test> node .\client.js
Welcome to the chatroom!
Your default nick name is 1563328661593.
Type "help" to show the usage!
Online:1
1563328672755 join the room!
1563328672755 says : Anyone here>
Yep , I am here
1563328672755 says : How to modify my nickname?
Tryna type '/help'.
1563328672755 changed his/her name to LILPIG
LILPIG says : Oh,got it ! Thanks!
Nothing... xD
```

`user2`
```cmd
PS D:\learn\node\tcp_cs_test> node .\client.js
Welcome to the chatroom!
Your default nick name is 1563328672755.
Type "help" to show the usage!
Online:2
Anyone here>
1563328661593 says : Yep , I am here
How to modify my nickname?
1563328661593 says : Tryna type '/help'.
/help
Command:
        1. /rename <Name> -- Modify your nickname in chatroom.
         2. /exit or /quit -- Exit from chatroom.
         3. /help -- Show usage.
Admin Command:
        4. /gname <any string> -- Get your nickname.Note that the option <any string> is must exits.
        1. /kickout <Name> -- Kickout a user by nickname.
        2. /close -- Close chatroom.
        3. /post <Message> -- Post some message.
/rename LILPIG
Your name has changed!
Oh,got it ! Thanks!
1563328661593 says : Nothing... xD
```
