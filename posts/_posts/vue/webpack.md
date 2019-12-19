---
title: Vue和Webpack
date: 2019-08-20 13:59:48
tags: [前端,vue,js]
categories: vue
---

在前端模块化的今天，只要出现一个前端技术貌似就会和webpack挂上钩，Vue当然也是了。

下面将从零开始使用wabpack搭建一个Vue小示例。

## 安装需要的环境
创建项目文件夹后，先`npm init`初始化。

然后按照如下目录创建文件：
```
├── dist
│   └── index.html
├── package.json
├── package-lock.json
├── src
│   ├── App.vue
│   └── index.js
└── webpack.config.js
```

```bash
npm i -D webpack webpack webpack-dev-server
npm i -P vue
npm i -D vue-loader
npm i -D vue-style-loader
npm i -D vue-template-compiler
npm i -D vue-hot-reload-api
npm i -D babel
npm i -D babel-loader
npm i -D babel-core
npm i -D babel-plugin-transform-runtime
npm i -D babel-preset-es2015
npm i -D babel-runtime
```

以上环境全部安装完毕后进行babel的配置，在项目根目录下创建`.babelrc`，写入如下内容：
```json
{
    "presets": ["es2015"],
    "plugins": ["transform-runtime"],
    "comments": false
}
```
这个配置用来解析es6，然后我们就可以在项目中用es6去编写代码了。

在`webpack.config.js`中编写：
```
var path = require('path');
var VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
    mode: 'development',
    entry: {
        main: './src/index.js'
    },
    output: {
        path: path.join(__dirname,'./dist'),
        publicPath: '/dist',
        filename: 'main.js'
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader',
            },
            {
                test: /\.js$/,
                loader: 'babel-loader',
                exclude: /node_modules/
            }
        ]
    },
    plugins:[
        new VueLoaderPlugin()
    ]
}
```
这个代码也没啥好解释的，就是在webpack中配置入口出口和一些loader规则。

然后编写`App.vue`，这是我们的根组件。

vue文件分为三块，分别用`template`、`style`和`script`包裹。分别代表组件的html，样式和基于js的组件配置。

```html
<template>
    <div>
        Helo,{{ name }}.
    </div>
</template>

<script>
    export default {
        data () {
            return {
                name: 'Vue.js'
            }
        }
    }
</script>

<style scoped></style>
```
这代码也很简单。

我们再看`index.js`：
```js
import Vue from 'vue'
import App from './App.vue'

new Vue({
    el: '#app',
    render: h => h(App)
});
```
在`index.js`中引入Vue和我们编写的根组件`App.vue`，然后初始化vue实例并挂载到主页面的`#app`节点上。

我们再看`dist/index.html`：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title></title>
</head>
<body>
    <div id="app"></div>   
    <script src="./main.js"></script>
</body>
</html>
```
提供了一个`app`节点供vue挂载，并引入了经过webpack打包后的`main.js`文件。

输入`npx webpack-dev-server --open --config webpack.config.js`运行，也可以在`package.json`中配置script。

运行后浏览器中显示`Helo,Vue.js.`

## 开发之前的yesorno
yesorno就是一个小工具，你可以提一个问题，他会告诉你是或者否。

我们这次打算用组件的形式开发，把问题的输入框抽象成一个组件，然后它自己监听自己的输入事件，当输入`?`时向父组件发射`submit`事件。父组件接收到这个事件后向服务器发送请求，获取数据并渲染到页面。

`App.vue`：
```html
<template>
    <div>
        I'll tell you yes or no.
        <yesorno-input @submit="answerIt"></yesorno-input>
        {{ message }}
    </div>
</template>

<script>
import axios from 'axios'
import yesornoInput from './Dialog.vue'
    export default {
        data () {
            return {
                message: 'Waitting for you.'
            }
        },
        components: {
            yesornoInput  
        },
        methods: {
            answerIt (question) {
                this.message = 'I\'m thinking...'
                axios.get('https://yesno.wtf/api').then(resp=>{
                    if(resp.data.answer){
                       this.message = resp.data.answer; 
                    }else{
                        this.message = 'Api Error!'
                    }
                }).catch(err=>{
                    this.message = 'Network Error!'
                })
            }
        }
    }
</script>

<style scoped></style>
```
我们使用`axios`向服务器发送请求，使用一个`message`作为用户提示和答案，并且我们可以发现组件`yesorno-input`，也就是我们要编写的子组件，它来自`Dialog.vue`，接下来我们编写这个文件。

`Dialog.vue`：
```html
<template>
    <input type="text" :placeholder="pholder" @input="handleInput">
</template>

<script>
export default {
    data () {
        return {
            pholder: 'Ask me a question!'
        }
    },
    methods: {
        handleInput: function(event){
            var question = event.target.value;
            if(question.endsWith('?') || question.endsWith('？')){
                this.$emit('submit',question);
            }
        }
    }
}
</script>
```
这个组件的内容也很简单，就是监测自己的`input`并发送事件给父组件。



