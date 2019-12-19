---
title: Nodejs中的异步编程 -- Promise
date: 2019-07-25 14:33:46
tags: [nodejs,async]
categories: NodeJS
---
Promise是一个用于链式编写异步代码的规范。嘶。我感觉没有`async`好使，不过可能各有各的好处吧！VSCode API整个都使用了Promise规范。

Promise从语义上看就是一个承诺，承诺有两种状态，一个是兑现了(resolve 解决)、一个是没有兑现(reject 拒绝)，Promise通过这两种状态控制异步操作。

一个链式操作大概是这样的：  
`task1().then(task2).then(task3).catch(err=>{处理错误})`

其中每一个task都是一个方法，这个方法返回一个Proimse，方法中如果处理完成调用`resolve`，然后`then`也就是下一个任务函数被调用，如果处理出错调用`reject`，整个链条中断，`catch`被调用。


比如：
```javascript
const fs = require('fs');

function getFile(path){
    return new Promise(function(resolve,reject){
        fs.readFile(path,function(err,data){
            if(err)
                reject(err);
            else
                resolve(data);
        })
    })
}

getFile('./file1').then(data=>{
    console.log('file1 => \n'+data);
    return getFile('./file2.js');
}).then(data=>{
    console.log('file2 => \n'+data);
}).catch(err=>{
    console.log('get an err => '+err)
})
```

上面代码提供了一个公共方法`getFile`返回promise对象，下面就可以通过`then-catch`进行链式调用了。它们是同步执行的。

除了`catch`还可以接`finally`，不管结果怎样都会被执行。


异步执行的话，可以使用`Promise.all`方法，这是一个静态方法。

```javascript
Promise.all([promise1,promise2,promise3]).then(data=>{
    console.log(data);
}).catch(err){
    console.log(err);
}
```

有任何一个promise返回了reject，都会回调`catch`，如果全部执行成功，`then`会被执行，并且所有的任务的结果会被封装成数组传递进去。
