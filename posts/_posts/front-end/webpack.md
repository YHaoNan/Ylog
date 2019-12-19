---
title: Webpack入门（一） 基本概念
date: 2019-08-18 14:55:27
tags: [前端,webpack]
categories: Webpack
---

在前端逐渐模块化的今天，不会nodejs和webpack简直就像外星人一样...

## webpack项目的创建
先创建一个文件夹，之后运行：
```bash
npm init
```

Webpack4版本中需要安装`webpack`和`webpack-cli`两个工具：
```bash
npm i -D webpack webpack-cli
B
```

这样一来项目的依赖就全都下载好了，我们可以开始创建项目结构了。

```
.
├── dist
│   └── index.html
├── package.json
├── package-lock.json
├── src
│   └── index.js
└── webpack.config.js
```
A
这是一个webpack项目最基本的目录结构，`dist`中是生成的目标文件和SPA应用的主html文件，`src`是我们编写的代码，`webpack.config.js`是webpack的配置文件。

我们首先编写`src`中的`index.js`：
```js
function createElement(text){
    let div = document.createElement('div');
    div.innerText = text;
    return div;
}

document.body.append(createElement('Hello,World!'))
```

代码很简单，就是往body中添加一个div。

接下来是`dist`中的`index.html`
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    <script src="./main.js"></script>
</body>
</html>
```
这个html中引用一个`main.js`，可是并没有这个文件，这个就是webpack给我们生成的文件。我们需要在`webpack.config.js`中编写配置：

```js
var path = require('path');

module.exports = {
    entry: './src/index.js',
    mode: 'development',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname,'dist')
    }
}
```
这个配置文件导出了一个对象，这个对象就是webpack项目的配置，`entry`指定了入口文件，我们的入口文件就是刚刚编写的`index.js`，mode是模式，我们指定开发模式，output则是输出配置，这里就看到了`main.js`的身影，我们把他导出到`dist/main.js`。

在浏览器中打开这个`index.html`，我们编写的div显示了出来。

## loader
在webpack中，万物都是模块，甚至图片和字体都是模块。一个web工程肯定会包括很多资源文件，比如CSS、JS、图片、模板、字体等等，webpack本身只能处理js，所以需要各式各样的loader来处理这些文件。

其实loader的本质就是一个函数，你传入一个东西，它给你返回处理后的东西。比如使用`babel-loader`将es6的js语法转换为es5，可以理解为其工作原理是这样的：
```
ES5 = babel-loader(ES6)
```
loader也可以是链式的，例如我们在使用scss时需要将它编译成css，然后再加载到html的style标签中，就像如下：
```
Style Element = style-loader(css-loader(scss-loader(SCSS)))
```

下面开始在webpack中使用loader工作，我们在src中创建一个css文件，并且引入它：

`style.css`
```css
body{
    color: #f00;
    font-size: 1.5em;
}
```
`index.js`
```js
import './style.css'
//...
```

运行示例，发现报错了：
```
ERROR in ./src/style.css 1:4
Module parse failed: Unexpected token (1:4)
You may need an appropriate loader to handle this file type, currently no loaders are configured to process this file. See https://webpack.js.org/concepts#loaders
> body{
|     color: #f00;
|     font-size: 1.5em;
 @ ./src/index.js 1:0-20
```
从异常信息看，是因为没有一个loader来处理css类型的文件，所以我们需要安装并配置loader。

安装`css-loader`和`style-loader`，他们分别用来解析css文件和把css插入到html的style标签中。

```
npm i -D css-loader style-loader
```

安装完后，修改`webpack.config.js`去针对css文件设置使用这两个loader。

```js
var path = require('path');

module.exports = {
    entry: './src/index.js',
    mode: 'development',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname,'dist')
    },
    module: {
        rules: [{
            test: /\.css$/,
            use: ['style-loader','css-loader']
        }]
    }
}
```
在`module`中设置一个rules，里面是针对引入的资源文件的扩展类型而使用的loader，上面我们指定了一个css类型的文件，并且使用了这两个loader，使用数组方式声明loader会从后向前链式调用。

## loader-options
有的loader提供了预置项，可以通过`options`设置：
```js
module: {
    rules: [{
        test: /\.css$/,
        use: ['style-loader',{
            loader:'css-loader',
            options: {
                //配置项
            }
        }]
    }]
}
```
## 更多配置
### exclude和include
很多时候loader需要排除一些目录或者仅包含一些目录，这时可以用这两个选项。

```js
module: {
    rules: [{
        test: /\.css$/,
        use: ['style-loader',{
            loader:'css-loader',
            options: {
                //配置项
            }
        }],
        exclude: /node_modules/
    }]
}
```
### enforce
用来指定一个loader的触发时机，`pre`代表比所有loader先触发，post表示在所有模块之后触发。

## 参考资料
* Webpack实战：入门、进阶与调优
