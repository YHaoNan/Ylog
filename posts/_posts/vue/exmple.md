---
title: 初识Vue —— 制作 I'll tell you
date: 2019-08-02 09:34:28
tags: [前端,vue,js]
categories: Vue
---

<style>
    #exm{
        text-align: left;
    }
    #exm iframe{
        border: 2px #333 solid;
        outline: none;
        width:100%;
    }
    #exm button{
        outline: none;
        border: 2px #333 solid;
        margin-top:10px;
        background-color: #ddd;
        padding: 10px;
        
    }
</style>

这几天一直在看《Java并发编程实战》，那书看得很累，今天想换换口味，之前一直想学vue也没行动，今天学学。就跟着官方教程学了起来，并做了下官方的Demo。

先看下示例效果：[I'll tell you](http://host.lilpig.site/ity)

因为没对手机适配，但是也能用，电脑效果更佳。



## 准备
先说一下vue的一些基本语法。

首先看看一个基本的Vue示例：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>E01</title>
    <!--  使用cdn引入vue  -->
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
</head>
<body>
    <div id="app">
        {{message}}
    </div>

    <script>
        var app = new Vue({
            el: '#app',
            data: {
                'message': 'Hello,World!'
            }
        });
    </script>
</body>
</html>
```
效果：
<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e01.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e01.html"><button>在浏览器中运行</button></a>
</div>


首先要初始化vue的构造函数，并传入一些数据，然后vue会操作DOM，把这些数据渲染进去，而且这些数据是实时的，点击在浏览器中运行，打开开发者工具的`console`，输入`app.message="ABC"`，那么你会发现对应的文字也会变成ABC，因为Vue是响应式的，所以你做的所有修改都能实时传递到DOM中。

Vue不仅能通过双大括号把数据渲染进标签文本，对于标签属性，vue也提供了支持。我们看一个更复杂的示例：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>


</head>
<body>
    <div id="helo">
        {{ message }}
        <input type="text" v-model="message">
    </div>
    <script>
            var app = new Vue({
                el: "#helo",
                data:{
                    message: "",
                },
            });
    </script>
</body>
</html>
```

这个示例和刚刚的只差了多了一个input标签，`v-model`可以把表单元素和vue中的一个值绑定，当表单元素改变这个值也会跟着改变，当这个值改变，表单元素也会改。

于是上面的示例我们会发现，在输入框中输入啥就会出现啥文字，试一下。

效果：
<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e02.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e02.html"><button>在浏览器中运行</button></a>
</div>

`v-on`则会把元素的一个监听器和一个声明在`methods`中的方法绑定，比如：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>


</head>
<body>
<div id="helo">
    <select v-model="selected">
        <option>百度</option>
        <option>必应</option>
    </select>
    <input type="text" v-model="searchKW">
    <button v-on:click="search">搜索一下</button>

</div>
<script>
    var app = new Vue({
        el: "#helo",
        data:{
            urls: {
                '百度':'https://www.baidu.com/s?wd=',
                '必应':'https://cn.bing.com/search?q='
            }
        },
        methods:{
            search:function () {
                window.location = this.urls[this.selected]+this.searchKW;
            }
        }
    });
</script>
</body>
</html>
```

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e03.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e03.html"><button>在浏览器中运行</button></a>
</div>

`v-if`可以决定一个元素是否显示：

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>


</head>
<body>
<div id="helo">
    <input type="checkbox" v-model="visiable">{{visiable ? '展示' : '隐藏'}}
    <div v-if="visiable">我是div</div>
</div>
<script>
    var app = new Vue({
        el: "#helo",
        data:{
            visiable: true
        },
    });
</script>
</body>
</html>
```

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e04.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e04.html"><button>在浏览器中运行</button></a>
</div>

`watch`可以监听一个变量的改变：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>


</head>
<body>
<div id="helo">
    {{ message }}
    <input type="text" v-model="message">
    <span>{{ log }}</span>
</div>
<script>
    var app = new Vue({
        el: "#helo",
        data:{
            message: "Hello,World!",
            log: '',
            to:null
        },
        watch: {
            message: function(oldV,newV){
                var vm = this;
                if (vm.to==null){
                    vm.to = setTimeout(function(){
                        vm.log = '';
                    },500)
                } else {
                    clearTimeout(vm.to);
                    vm.to = null;
                }
                vm.log = 'Change message to '+newV
            }
        }
    });
</script>
</body>
</html>
```

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e05.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e05.html"><button>在浏览器中运行</button></a>
</div>

## I'll tell you
准备工作做完了，我们来开始搞这个I'll tell you

当然，我们只学了Vue的冰山一角，不过搞这个够了。

依赖：
    - vue
    - axios@0.12.0  
        异步请求库，用于请求yesorno api
    - lodash@4.13.1  
        用于实现等待输入功能

html结构如下：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Yes Or No</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=0.5, maximum-scale=2.0, user-scalable=yes" />
    <link rel="stylesheet" href="./res/style.css">
    <!-- 引入VUE -->
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
    <!-- 引入axios -->
    <script src="https://cdn.jsdelivr.net/npm/axios@0.12.0/dist/axios.min.js"></script>
    <!-- 引入lodash -->
    <script src="https://cdn.jsdelivr.net/npm/lodash@4.13.1/lodash.min.js"></script>
</head>
<body>
    <div id="app">
        <img v-bind:src="imgUrl"/>
        <transition name="fade">
            <div id="loading-dialog" v-if="loading">I'm thinking...</div>
        </transition>
        <div id="center">
            <p><span class="title">I'll tell U yes or no  <span class="bigger">:)</span></span></p>
            <input v-model="question"/>
            <p><span id="answer">{{answer}}</span></p>
        </div>
    </div>
    <div id="link">LINK &nbsp;: &nbsp;<a href="http://lilpig.site" target="_blank">LILPIG的博客</a> &nbsp;|&nbsp; <a href="http://host.lilpig.site/gameoflife" target="_blank">生命游戏试玩</a></div>
    <!-- 构造VUE对象 -->
    <script src="./res/vue-binding-elements.js"></script>
</body>
</html>
```

css我就不放上来了，占地方。接下来就是`vue-binding-elements.js`

```js
var vm = new Vue({
    el: '#app',
    data:{
        question: '',
        answer: 'Ask me a question!',
        imgUrl:'',
        loading: false
    },
    watch:{
        question:function (oldV,newV) {
            this.answer = "Waiting for your typing...";
            this.doDebounce();
        }
    },
    created: function () {
        this.doDebounce = _.debounce(this.getAnswer, 500);
        this.getAnswer(true);
    },
    methods:{
        getAnswer: function(isInit){
            if (this.question.endsWith("?") || this.question.endsWith("？") || isInit){
                var vm = this;
                vm.loading = true;
                axios.get('https://yesno.wtf/api')
                    .then(function (response) {
                        if (!isInit)
                            vm.answer = _.capitalize(response.data.answer);
                        var respImgUrl = response.data.image;
                        if (respImgUrl != vm.imgUrl)
                            vm.imgUrl = response.data.image;
                    })
                    .catch(function (err) {
                        vm.answer = 'I\' so sorry.I don\'t know.T_T'
                    }).finally(function () {
                        vm.loading = false;
                    })
            }else {
                this.answer = 'A question must end with a ?'
            }
        }
    }
});
```

嘶，，代码挺简单的，结合官方文档看吧。