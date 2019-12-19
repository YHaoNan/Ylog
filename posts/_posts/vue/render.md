---
title: Vue组件的render函数
date: 2019-08-17 8:27:00
tags: [前端,vue,js]
categories: Vue
---

之前我们都是用`template`来给组件提供一个模板，但是这并不一定适合所有情况，而且它最后也会被编译为`render`的形式，效率比直接用`render`稍低。

`render`函数提供了一个虚拟dom，也就是上一篇中提到的那个东西，你可以用这个虚拟dom来创建节点，绑定数据...

直接用`render`写的代码会比较冗余，因为是用js而不是html，所以说`render`函数要看情况使用。


下面的示例就需要用到render，我们需要一个标题组件，这个组件实现的功能是点击标题后自动给url添加锚点，但因为`template`中无法访问组件的作用域中的变量，所以这个组件无法工作。

```html
<div id="app" v-cloak>
    <anchor :level="5" title="特性">特性</anchor>
</div>
<script>
    Vue.component('anchor',{
        template:'<a :href="\'#\'+title" @mouseenter="setHover(true)" @mouseleave="setHover(false)"><h'+this.level+'><slot></slot>' +
            '<span v-show="isHover">#</span>' +
            '</h'+this.level+'></a>',
        props: {
            level: {
                type: Number,
                default: 1
            },
            title: {
                type: String,
                required: true
            }
        },
        data: function(){
            return {
                isHover: false,
            }
        },
        methods: {
            setHover: function (bol) {
                this.isHover = bol;
            }
        },
    })
    var vm = new Vue({
        el: '#app'
    });
</script>
```
当然也可以使用`v-if`来判断`level`，生成不同的标题，但是这样写代码太冗余，所以这时就要用到render函数。

```html
<div id="app" v-cloak>
    <anchor :level="1" title="特性1">特性</anchor>
    <anchor :level="1" title="特性2">特性</anchor>
    <anchor :level="1" title="特性3">特性</anchor>
</div>

<script>
    Vue.component('anchor',{
        render: function(h){
            var _this = this;
            if(this.isHover){
                var span = h('span','#')
            }
            return h('a',{
                attrs:{
                    href: '#'+this.title
                },
                on:{
                    mouseenter: function () {
                        _this.isHover = true;
                    },
                    mouseleave: function () {
                        _this.isHover = false;
                    }
                }
            },[
                h('h'+this.level,[
                    this.title,
                    span
                ])
            ])
        },
        props: {
            level: {
                type: Number,
                default: 1
            },
            title: {
                type: String,
                required: true
            }
        },
        data: function(){
            return {
                isHover: false,
            }
        },
    })
    var vm = new Vue({
        el: '#app'
    });
</script>
```

除了把`template`替换成`render`，其他代码没啥变化，`render`函数传入一个参数，这个参数是用于创建虚拟节点的方法，我们调用它就可以创建虚拟节点。这个方法可以有三个参数，第一个是必选参数是，是html标签名，第二个是可选参数，是该节点的数据对象，用于绑定一些数据和监听器的东西，第三个也是可选参数，是子节点列表。

我们上面的render代码创建了一个a标签，并且绑定了href为`'#'+this.title`，用于点击后自动修改url，并且监听他的鼠标事件，去修改`isHover`，用于显示锚点图标，很多博客站的标题锚点功能都是这样的，包括本博客，就是鼠标进入标题时显示一个图标（本实例中是#号）。然后创建标题节点，作为a的子节点，这里可以访问到组件的作用域，所以`'h'+this.level`是可以的，然后给它两个子节点，一个是`this.title`，就是传入的标题，一个是`span`，这个`span`是在`render`函数最开始判断`isHover`动态创建的，因为`render`函数中不能用`v-if`这种命令，所以只能用js提供的结构体来满足功能。

