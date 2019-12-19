---
title: Vue前端路由
date: 2019-08-20 15:19:48
tags: [前端,vue,js]
categories: vue
---

Vue擅长SPA应用，就是单页面应用，而且近几年SPA应用越来越流行，很多站点的html打开后都只有一个空壳，并且这个空壳会被动态填充，贯穿全局。

只用一个页面贯穿整个网站产生了一个问题，那就是路由问题，这部分问题以前都是在服务端处理的，根据不同的URL返回不同的页面，实际上，这还没做到前后端分离，因为后端还是要处理路由。SPA应用把路由功能转移到了前端JS上，通过监听url变化填充不同的页面内容，这样开发前端应用就很像开发固定平台的APP，页面相关的功能由JS实现，后端只提供一套API即可。

## 目录结构
先把本篇笔记的目录放出来：
```
.
├── index.html
├── package.json
├── package-lock.json
├── src
│   ├── about.vue
│   ├── App.vue
│   ├── index.js
│   └── index.vue
└── webpack.config.js
```

如果需要跟本篇的思路一起练习编写代码，请按上面的目录结构创建文件。

## vue-router
Vue提供了`vue-router`插件来实现前端路由功能。

```bash
npm i -P vue-router
```
然后我们可以通过`Vue.use`方法来使用它:

`index.js`：
```
import Vue from 'vue';
import VueRouter from 'vue-router'

Vue.use(VueRouter);
```

我们先编写`index.html`：
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title></title>
</head>
<body>
    <div id="app"></div>
    <script src="/dist/main.js"></script>
</body>
</html>
```

这肯定没什么好说的，然后就是App.vue：
```html
<template>
    <div>
        APP
        <router-view></router-view>
    </div>
</template>

<script>
    export default {
    }
</script>

<style scoped></style>
```
这里面唯一一个新面孔就是`router-view`，这是vue提供的一个组件，它最终会被渲染成其他的路由页面。

我们为这个示例编写两个路由页面，分别是`index.vue`和`about.vue`，并且在访问其他不存在的页面时统统跳转到`index`，这些逻辑我们需要在`index.js`中编写：
```js
import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

var router = new VueRouter({
    mode: 'history',
    routes: [
        {
            path: '/index',
            components: require('./index.vue')
        },
        {
            path: '/about',
            components: require('./about.vue')
        },
        {
            path: '*',
            redirect: '/index'
        }
    ]
})

new Vue({
    el: '#app',
    router: router,
    render: h => h(App)
});
```
这里面用`VueRouter`来初始化了一个路由，并且传入了创建的Vue根实例中。注册路由时，用`path`制定一个表达式，如果满足的就会渲染定义的组件。我们来看看`index.vue`和`about.vue`：

```html
<template>
    <div>
        This is the index page.
    </div>
</template>

<script>
export default {

}
</script>
```

```html
<template>
    <div>
        This is the about page.
    </div>
</template>

<script>
export default {

}
</script>
```
这样在浏览器中运行时输入`index`和`about`就会看到相应的页面了。

## 带参路由
路由的path可以携带一个参数，比如`/user/123456`，其中的`123456`为用户id，这在业务中非常常见。

```js
{
    path: '/user/:id',
    components: require('./user.vue')
}
```

这样就注册了一个带参数的路由，我们在组件中可以通过`$route.params`来获取它：
```html
<template>
    <div>
        {{ $route.params.id }}
    </div>
</template>
<script>
export default {

}
</script>
```

## 跳转
vue-router提供了`router-link`用于跳转页面。

```html
<div>
    APP
    <router-view></router-view>
    <router-link to="/index">首页</router-link>
    <router-link to="/about">关于</router-link>
    <router-link :to="'/user/'+id">我</router-link>
    <br>
    你的id
    <br>
    <input type="number" v-model="id">
</div>
```
修改`App.vue`，让它支持跳转到其他的页面。

默认情况下`router-link`会被渲染成a标签，你也可以通过指定`tag`prop来决定渲染成什么标签：
```html
<router-link to='/about' tag='li'>About</router-link>
```

如果使用`replace`，则不会留下历史记录，不能使用回退：
```html
<router-link to='/about' replace>About</router-link>
```

默认情况下如果`router-link`对应路由匹配成功，会自动添加一个`router-link-active`的类，可以通过给它设置`active-class`prop修改这个默认值。

除了使用Vue的`router-link`组件，还可以使用`this.$router.push('/user/123')`跳转页面。

- `this.$router.replace('/index')`  不会留下历史记录
- `this.$router.go(-1)` 前进n页 示例为后退一页

## 修改网页标题
Vue提供了两个钩子，`beforeEach`和`afterEach`分别是路由改变之前和路由改变之后触发。修改标题我们需要用到`beforeEach`。

```js
var router = new VueRouter({
    mode: 'history',
    routes: [
        {
            path: '/index',
            components: require('./index.vue'),
            meta: {
                title: 'Index'
            }
        },
        /...
    ]
});

router.beforeEach((to,from,next) => {
    console.log(to.meta.title)
    window.document.title = to.meta.title;
    next();
})
```

必须要调用`next()`进入下一个钩子。

可以用这个`beforeEach`判断是否登录：
```js
router.beforeEach((to,from,next)=>{
    if(window.localStorage.getItem('token'))
        next()
    else
        next('/login')
})
```
`next`的参数为false可以取消导航


