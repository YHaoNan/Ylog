---
title: Vuex状态管理
date: 2019-08-21 8:23:48
tags: [前端,vue,js]
categories: vue
---

Vue中可以通过子组件可以通过`prop`获取父组件传入的数据，子组件可以通过`$emit`来发射事件与父组件交互，但是很多情况下不单单是父子组件之间的交互，还需要跨级交互和兄弟组件交互。Vuex可以很好的实现这一功能。

## 安装
`npm i -P vuex`

## 引用
它也是一个插件，和路由插件一样，需要先`Vue.use`，然后再在根节点的配置中指定一下。

目录结构如下，完全复制自路由那篇笔记。
```
.
├── index.html
├── package.json
├── package-lock.json
├── src
│   ├── about.vue
│   ├── App.vue
│   ├── index.js
│   ├── index.vue
│   └── user.vue
└── webpack.config.js

1 directory, 9 files
```

在`index.js`中引用Vuex：
```js
import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'
import Vuex from 'vuex'

Vue.use(Vuex);
Vue.use(VueRouter);

//...

const store = new Vuex.Store({
    
})

new Vue({
    el: '#app',
    router: router,
    store: store,
    render: h => h(App)
});

```

我们通过这个`store`来存储全局状态，比如：
```js
const store = new Vuex.Store({
    state: {
        count: 0
    }              
})
```

`state`是全局可访问的状态，这里面的变量都是响应式的，可以在其他组件这样引用它：

`index.vue`：
```html
<template>
    <div>
        {{ $store.state.count }}
        This is the index page.
    </div>
</template>
```
如果在所有需要引用`count`的地方都这样写有些麻烦，可以用`computed`：
```html
<template>
    <div>
        {{ count }}
        This is the index page.
    </div>
</template>

<script>
export default {
    computed: {
        count () {
            return this.$store.state.count;
        }
    }
}
</script>
```

如果想写入计数器中的数据，不能直接写，而是要提交`mutations`。

```js
const store = new Vuex.Store({
    state: {
        count: 0
    },
    mutations: {
        increment (state) {
            state.count++
        },
        decrement (state) {
            state.count--
        }
    }
})
```
我们给store添加了两个`mutations`，在组件中我们可以通过`commit`调用它：
```html
<template>
    <div>
        {{ count }}
        This is the index page.
        <button @click="handleIncrement">+</button>
        <button @click="handleDecrement">-</button>
    </div>
</template>

<script>
export default {
    computed: {
        count () {
            return this.$store.state.count;
        }
    },
    methods: {
        handleIncrement () {
            this.$store.commit('increment');
        },
        handleDecrement () {
            this.$store.commit('decrement');
        }
    }
}
</script>
```

`commit`也可以传参
```js
mutations: {
    changeCount (state , n=1) {
        state.count += n
    },
}
```

```html
<template>
    <div>
        {{ count }}
        This is the index page.
        <button @click="handleIncrement">+</button>
        <button @click="handleDecrement">-</button>
        <br>
        Step: <input type="number" v-model.number="step">
    </div>
</template>

<script>
export default {
    data () {
        return {
            step: 1
        }
    },
    computed: {
        count () {
            return this.$store.state.count;
        }
    },
    methods: {
        handleIncrement () {
            this.$store.commit('changeCount',this.step);
        },
        handleDecrement () {
            this.$store.commit('changeCount',-this.step);
        }
    }
}
</script>
```

## getters
Vuex还可以提供一个`getters`，它类似于组件的`computed`，假如我们的状态中有一个列表，业务中经常要用到取列表中的正数，但是在每一个组件中编写这个逻辑代码就太冗余了，所以可以在`vuex`的getters中写:

```js
getters: {
    filteredList (state) {
        return state.list.filter(item=>item>0);
    }
}
```

在组件中可以这样用：
`this.$store.getters.filteredList`

## actions
actions是用来做异步操作的。

```js
actions: {
    asyncFunction (context) {
        return new Promise(resolve => {
            setTimeout(()=>{
                resolve('ok');
            },1000)
        })
    }
}
```

```js
callIt () {
    this.$store.dispatch('asyncFunction').then(res=>{
        console.log(res);
    })
}
```

## modules
如果状态过多，可以分模块。

```js
const moduleA = {
    state: {}
}
//..
const store = new Vuex.Store({
    modules: {
        a: moduleA
    }
})
```
子节点的getters和mutation可以接受第三个参数rootState来访问root的状态。


