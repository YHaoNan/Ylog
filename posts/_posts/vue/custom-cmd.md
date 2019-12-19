---
title: Vue自定义命令
date: 2019-08-11 18:07:48
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

Vue除了`v-bind`、`v-if`这些优秀的内置命令外，还允许我们根据业务需求自定义命令。自定义命令和自定义组件的过程非常相似

Vue中这样注册一个命令：
```javascript
Vue.directive('clickoutside',{
    //...
})
```

因为命令需要写在一个html元素上，所以可以注册很多html元素生命周期的钩子：

- bind: 只调用一次，指令第一此绑定到元素时调用
- inserted: 被绑定元素插入父节点时调用
- update: 被绑定元素所在模板更新时调用
- componentUpdate: 被绑定元素模板完成一次更新周期时调用
- unbind: 指令与元素解绑时调用

我们可以通过以上的内容来编写一个自动获取焦点的输入框：

```html
<html lang="en">
<!-- ... -->
    <div id="app">
        <input type="text" v-focus>
    </div>

    <script>
        Vue.directive('focus',{
            inserted: function (el) {
                el.focus();
            }
        });
        var vm = new Vue({
            el: '#app'
        });
    </script>
<!-- ... -->
</html>
```
我们监听了inserted钩子，当input被插入时立即获取焦点

这些钩子还可以有其它的参数：

- el: 命令所绑定的元素
- binding: 一个对象，这个对象包含很多属性
    - name: 指令名
    - value: 指令绑定值 如果是表达式则是计算后的值
    - oldValue: 指令绑定的前一个值
    - expression: 指令值的字符串形式
    - arg: 传给指令的参数，如`v-bind:msg`，msg为arg
    - modifiers: 指令的修饰符，如`@click.stop`，此时modifiers的值为{stop: true}
- vnode: Vue编译生成的虚拟节点
- oldVnode: 上一个虚拟节点

## 下拉菜单
我们开发一个下拉菜单，这个菜单当点击菜单以外的位置时会关闭。

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/pop_menu.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/pop_menu.html"><button>在浏览器中运行</button></a>
</div>

分析需求，我们写一个`clickoutside`的命令，并且传入一个函数，如果点击`document`就调用这个函数，在这个函数中进行关闭。

这是html代码：

```html
<div id="app" v-clickoutside="handleClose">
    <button @click="show = !show">菜 单</button>
    <br>
    <div id="menu" v-show="show">
        <ul>
            <li>ITEM1</li>
            <li>ITEM2</li>
            <li>ITEM3</li>
        </ul>
    </div>
</div>
```

下面是js代码，原理就是给`document`添加和移除事件。
```javascript
Vue.directive('clickoutside',{
    inserted: function (el,binding,vnode) {
        function documentHandler(e){
            if (el.contains(e.target)){
                return false;
            }
            if (binding.expression){
                binding.value(e);
            }
        }
        el.__vueClickOutSide__ = documentHandler;
        document.addEventListener('click',documentHandler);
    },
    unbind: function (el,binding) {
        document.removeEventListener('click',el.__vueClickOutSide__);
        delete el.__vueClickOutSide__;
    }
});
var vm = new Vue({
    el: '#app',
    data: {
        show: false
    },
    methods: {
        handleClose: function () {
            this.show = false;
        }
    }
});
```

