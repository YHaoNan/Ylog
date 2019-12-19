---
title: Vue组件实例
date: 2019-08-15 10:29:48
tags: [前端,vue,js]
categories: Vue
---
看Vue官方文档，看完发现太过理论，到自己手里啥也写不出来，于是又看了本叫《Vue.js实战》的书，然后跟着写了俩组件的实例。

## 数字输入框

编写一个只能输入数字的输入框，并且提供加减按钮。

![图片](http://nsimg.cn-bj.ufileos.com/img-1565836997751.jpg)

这个需求乍一看用纯的html和js也不难实现，用vue是在是有些小题大做了，不过练习嘛。。

分析需求，我们需要让我们的组件内部和引用组件的地方进行双向数据绑定，所以要提供我们自己的`v-model`。我们打算把这个组件分成两个文件来实现。

- index.html  
    html模板文件
- input_number.js  
    组件文件

先看下html模板的代码，比较简单：
```html
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Title</title>
		<script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=0.5, maximum-scale=2.0, user-scalable=yes" />
	</head>
	<body>
		<div id="app">
			<input-number v-model="value" :max="10" :min="0"></input-number>
		</div>
		<script src="./input_number.js"></script>
		<script>
			var vm = new Vue({
				el: '#app',
				data: {
					value: 5
				}
			})
		</script>
	</body>
</html>
```

这个代码非常简单并且说明了我们的组件所提供的`props`参数，提供了三个`prop`，分别是用于数据绑定的`value`，用于限定输入最大值的`max`，用于限定输入最小值的`min`。

因为使用`v-model`，所以我们需要在组件中发射`input`事件告诉vue根节点更新`value`。这在上一篇笔记中有提到，详见[组件中的v-model](https://lilpig.site/post/vue-component#%E7%BB%84%E4%BB%B6%E4%B8%AD%E7%9A%84v-model)。

我们再开始编写`input_number.js`：
```javascript
Vue.component('input-number',{
	template: '<div class="input-number">\
				<input type="text" :value="currentValue" @input="handleChange">\
				<button @click="handleDown" :disabled="currentValue <= min">-</button>\
				<button @click="handleUp" :disabled="currentValue >= max">+</button>\
			</div>',
	props: {
		max: {
			type: Number,
			default: Infinity
		},
		min: {
			type: Number,
			default: -Infinity
		},
		value: {
			type: Number,
			default: 0
		}
	},
	data: function(){
		return {
			currentValue: this.value
		}
	}
    //...
})
```
我们先写一部分，把模板和`props`编写了，上面的代码我们把max和min分别设置成了正负无穷，并把`currentValue`初始化为从父控件中通过`v-model`传递进来的`value`，而且我们模板中的`input`绑定了`value`为`currentValue`，这是一个单向绑定，也就是说输入框中的值改变时`currentValue`并不会直接改变，而是提供了一个`@input="handleChange"`让输入时通过`handleChange`函数设置`currentValue`的值来完成双向数据绑定，这不是多此一举，因为我们做的是数字输入框，所以要在`handleChange`方法中验证输入数据的合法性。

`handleDown`和`handleUp`很好理解，就是加一减一的方法，这两个按钮通过`:disabled`来实现不满足条件下按钮不允许使用。

接下来看看这些方法都是怎么工作的：
```javascript
methods: {
    updateValue: function(val){
        if(val>this.max)val=this.max;
        if(val<this.min)val=this.min;
        this.currentValue = val;
    },
    handleUp: function(){
        if(this.currentValue<this.max)this.currentValue++;
    },
    handleDown: function(){
        if(this.currentValue>this.min)this.currentValue--;
    },
    handleChange: function(event){
        var val = event.target.value.trim();
        if((/^\d+$/).test(val)){
            val = Number(val);
            if(val>this.max)val=this.max;
            if(val<this.min)val=this.min;
            if(this.currentValue == val) event.target.value = val;
            this.currentValue = val;
        }else{
            event.target.value = this.currentValue;
        }
    },
},
watch: {
    currentValue: function(val){
        this.$emit('input',val);
    },
    value: function(val){
        this.updateValue(val);
    }
},
mounted: function(){
    this.updateValue(this.value);
}
```
首先这些方法大多是通过修改`currentValue`来完成的，并且我们监听了`currentValue`的修改事件，如果修改了就发出一个`input`事件给父级，这个事件会被父级的`v-model`接收，完成父级的数据双向绑定。

这些方法本身没什么好说的，都挺简单，和vue也没啥关系。就不说了。

在Vue这种基于数据驱动的js库中，中尽量不要手动操作dom，有什么需求请先考虑是否能用修改数据实现。

*PS:这个例子和原书中的有些区别，因为我觉得原书中的这个示例写的确实有很多问题，但是尊重原书，大部分还是按照原书的逻辑写的*


### 练习
1. 输入框聚焦时按键盘的上下键完成加一减一操作。  
    只需要在模板的输入框中添加如下内容：
    ```html
    <input type="text" @keyup.up="handleUp" @keyup.down="handleDown" :value="currentValue" @input="handleChange">
    ```
    详见Vue官方文档事件篇

2. 增加一个控制步数的`prop - step`，用于控制加减操作一次加减的值  
    略
## 标签页
下面写一个标签页功能，这个挺帅的

![图片](http://nsimg.cn-bj.ufileos.com/img-1565839315027.jpg)

没错，只是我的css写的有点丑而已，人家书里的示例挺好看的。。。

分析一下，首先我把这个分成两个组件，一个是上面的标签，一个是底下的面板。

我们将会在html模板中这样使用它：

```html
<div id="app">
    <tab>
        <pane label="FIRST" closable="true">
            FIRST PANE
            <p>
                <code>
                    var message = "Helo,Vue";<br>
                    console.log(message);
                </code>
            </p>
            <p>
                <code>
                    >>> Helo,Vue
                </code>
            </p>
        </pane>
        <pane label="SECOND">SECOND PANE</pane>
        <pane label="THRID" closable="true">THRID PANE</pane>
    </tab>
</div>
```

通过一个`tab`标签包裹住一些`pane`，`tab`会自动地获取每个`pane`的名字（label），然后自动生成对应的顶部标签。

在`pane`标签中必须指定一个`label`作为对应顶部标签的标题，closable是一个可选项，为`true`代表这个标签页可以关闭。

现在我们知道该如何使用这个组件了，那么下一步就是思考如何实现它。首先我们分成两个js，`tab.js`和`pane.js`。`tab.js`主要负责顶部的所有标签和设置当前显示的标签。因为所有的标签都在`tab`标签中，所以要提供一个插槽。

因为顶部标签有多少完全依赖于tab的子元素中有多少个pane，这是不定的，所以我们要用`v-for`来循环生成。
```javascript
Vue.component('tab',{
   template: '<div>' +
       '<div class="tab" v-for="(nav,index) in navList" @click="handleChange(index)" :class="{active: currentIndex == index}">{{nav.label}}</div>' +
       '<div class="pane-container"><slot></slot></div>'+
       '</div>',
});
```
如上是我们根据需求编写的组件模板，因为涉及到样式，所以还绑定了一些class。这其中引用了一些方法和变量，`navList`就是通过遍历子元素中的`pane`得来的标签列表，`currentIndex`用来记录当前页面的下标，`handleChange`用来更改页面。

我们先搞定这个`navList`，要搞定它就要遍历子元素，除非在做这种独立组件时，否则不推荐在vue中直接获取父子元素，这会让数据流向变得难以控制。

```javascript
methods: {
    getTabs: function () {
        return this.$children.filter(function (child) {
            return child.$options.name === "pane";
        })
    },
}
mounted: function () {
    var _this = this;
    this.getTabs().forEach(function (pane,index) {
        _this.navList.push({
            name: pane.name || index,
            label: pane.label,
        })
    });
    this.currentIndex = 0;
},
data: function () {
    return {
        navList: [],
        currentIndex: -1
    }
}
```
我们在挂载完毕时通过遍历`this.$children`中符合规则的元素来初始化`navList`和`currentIndex`。

接下来就要考虑修改当前页面，只有pane本身可以控制自己的显示和隐藏，所以我们打算在pane中提供一个`show`变量，并且作为`v-show`的条件。

```js
Vue.component('pane',{
    props:['name','label','closable'],
    template: '<div v-show="show"><slot></slot></div>',
    data: function () {
        return {
            show: false,
        }
    }
})
```

我们考虑通过修改`currentIndex`，并监听这个变量来修改当前页面：
```javascript
//...in methods
handleChange: function (i) {
    this.currentIndex = i;
},
update:function (cur) {
    this.getTabs().forEach(function (pane,index) {
        if(index == cur){
            pane.show = true;
        }else{
            pane.show = false;
        }
    })
}
//....
watch:{
    currentIndex: function (val) {
        this.update(val);
    }
},
```

这样就完事了，不过我们要难度加大，再给它添加一个关闭功能。

关闭功能，我们打算在pane中提供一个`alive`属性并且作为`v-if`的条件，并且在关闭时，tab会修改目标的`alive`属性。

最后我们放出完整代码：

index.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <script src="https://cdn.jsdelivr.net/npm/vue/dist/vue.js"></script>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=0.5, maximum-scale=2.0, user-scalable=yes" />
    <title>TAB BAR</title>
    <style>
        .tab{
            display: inline-block;
            border: #aaaaaa solid 1px;
            padding: 10px 20px;
            margin: 5px;
        }
        .active{
            border-bottom: deepskyblue solid 2px;
        }
        .pane-container{
            padding: 10px;
            margin: 0 5px;
            border: #aaa solid 1px;
        }
        code{
            font-family: "Fira Code";
            font-size: 0.8em;
        }
        .close{
            outline: none;
            border: none;
            background: rgba(0,0,0,0);
            color: red;
            font-size: 0.8em;
        }
    </style>
</head>
<body>
    <div id="app">
        <tab>
            <pane label="FIRST" closable="true">
                FIRST PANE
                <p>
                    <code>
                        var message = "Helo,Vue";<br>
                        console.log(message);
                    </code>
                </p>
                <p>
                    <code>
                        >>> Helo,Vue
                    </code>
                </p>
            </pane>
            <pane label="SECOND">
                SECOND PANE
                <p>
                    <code>
                        var message = "Helo,Vue";<br>
                        console.log(message);
                    </code>
                </p>
                <p>
                    <code>
                        >>> Helo,Vue
                    </code>
                </p>
            </pane>
            <pane label="THRID" closable="true">THRID PANE</pane>
        </tab>
    </div>

    <script src="./pane.js"></script>
    <script src="./tabs.js"></script>
    <script>
        var vm = new Vue({
            el: '#app'
        });
    </script>
</body>
</html>
```

tab.js

```javascript
Vue.component('tab',{
   template: '<div>' +
       '<div class="tab" v-for="(nav,index) in navList" @click="handleChange(index)" :class="{active: currentIndex == index}">{{nav.label}} <button class="close" @click.stop="closePane(index)" v-if="nav.closable">X</button></div>' +
       '<div class="pane-container"><slot></slot></div>'+
       '</div>',
   methods: {
      getTabs: function () {
         return this.$children.filter(function (child) {
            return child.$options.name === "pane" && child.alive;
         })
      },
      handleChange: function (i) {
         this.currentIndex = i;
      },
      closePane: function (i) {
         this.navList.splice(i,1);
         this.getTabs()[i].alive = false;
         this.update(0);
      },
      update:function (cur) {
         this.getTabs().forEach(function (pane,index) {
            if(index == cur){
               pane.show = true;
            }else{
               pane.show = false;
            }
         })
      }
   },
   watch:{
      currentIndex: function (val) {
         this.update(val);
      }
   },
   mounted: function () {
      var _this = this;
      this.getTabs().forEach(function (pane,index) {
         _this.navList.push({
            name: pane.name || index,
            label: pane.label,
            closable: pane.closable
         })
         console.log(_this.navList[index].closable)
      });
      this.currentIndex = 0;
   },
   data: function () {
      return {
         navList: [],
         currentIndex: -1
      }
   }
});
```

pane.js

```javascript
Vue.component('pane',{
    props:['name','label','closable'],
    template: '<div v-if="alive" v-show="show"><slot></slot></div>',
    data: function () {
        return {
            show: false,
            alive: true
        }
    }
})
```
