---
title: 事件发射器-NodeJS
date: 2019-07-02 05:44:12
tags: [nodejs]
categories: NodeJS
---

## 标准回调模式
回调肯定都用过，就是函数不使用返回值，而是让调用者提供一个回调函数，可以有一些参数，如果函数执行完了再调用这个回调函数通知调用者处理结果。

比如：
```javascript
var fs = require('fs');
fs.readFile('/etc/passwd/',function(err,fileContent){
    if(err){
        console.log('发生错误');
        return;
    }
    console.log(fileContent.toString());
});
```

回调有一个好处，它适合异步编程，因为异步编程时你所在的主线程并不会等函数中的子线程执行完毕，函数没有执行完毕就返回了，这会返回不正确的结果，而回调则不会。

## 事件发射器
回调也有局限性，比如我们想接受多种事件的回调，而且有的事件可能会被触发多次，比如读取socket中的所有数据，这用标准回调模式时我们需要自己写很多代码，至少要是这样的：
```javascript
xxx(function(event){
    if(event == A)
        doA();
    else if(event == B)
        doB();
})
```
这里我们提供了一个event来确定回调的类型，这样很不方便，nodejs提供了一个事件发射器的模型，可以方便的实现这些功能，比如http模块就使用了它：
```javascript
var req = http.request(options,function(response){
    response.on('data',function(data){
        console.log(data);
    });
    response.on('end',function(){
        console.log('Response ended!');
    });
});
req.end();
```
在这里我们可以用`on`来接受一个事件，当接收到`data`事件时，后面的回调就会被触发。

```javascript
util = require('util');
var EventEmitter = require('events').EventEmitter;

var Ticker = function(){
    var self = this;
    setInterval(function(){
        self.emit('tick');
    }, 1000);
}
//使Ticker继承自EventEmitter
util.inherits(Ticker,EventEmitter);

var ticker = new Ticker();
ticker.on('tick',function(){
    console.log('tiktok!');
});
```
使用EventEmitter的emit方法可以发射一个事件，我们在上面的代码中每秒发射一个tick事件，然后在下面监听tick事件，一旦被触发就向控制台输出一行字符串。

运行它，你会发现每隔一秒控制台都会输出一行`tiktok！`

EventEmitter有如下方法：
- on(event,callback) 注册一个事件监听器
- addListener(event,callback) 同on
- removeListener(event,callback) 移除一个监听器
- removeAllListener(event) 移除一个事件下的所有监听器
- once(event,callback) 注册一个事件监听器，这个监听器只会被触发一次，之后就会被移除
