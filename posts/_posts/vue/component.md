---
title: Vue中的组件
date: 2019-08-11 08:29:48
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

## 示例
```html
<div id="app">
    <button-counter></button-counter>
</div>
<script>
    Vue.component('button-counter',{
        data:function () {
            return {
                count: 0
            }
        },
        template: '<button @click="count++">You have clicked me {{count}} times.</button>'
    });
    var vm = new Vue({
        el: '#app',
        data: {
            contentFontSize: 1,
        }
    });
</script>
```

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e11.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e11.html"><button>在浏览器中运行</button></a>
</div>

使用`Vue.component`注册一个全局的组件，组件名是第一个参数，第二个参数是这个组件的模板、数据等信息，然后就可以在html中使用这个组件了。

## 局部组件
```html
<div id="app">
    <button-counter></button-counter>
</div>
<script>
    var ButtonCounter = {
        data:function () {
            return {
                count: 0
            }
        },
        template: '<button @click="count++">You have clicked me {{count}} times.</button>'
    };
    var vm = new Vue({
        el: '#app',
        data: {
            contentFontSize: 1,
        },
        components:{
            'button-counter': ButtonCounter
        }
    });
</script>
```
使用`components`可以指定依赖的组件，当你使用一些构建系统时可以使用`import`导入写好的组件并引用它。

## Props
`props`是组件和数据之间的接口，我们可以通过定义`props`来向组件中传递数据。

```html
<div id="app">
    <user-profile :username="user.username" :age="user.age" :gender="user.gender"></user-profile>
</div>

<script>

    Vue.component('user-profile',{
        props: ['username','age','gender'],
        template: '<div>{{username}} -- {{age}} -- {{gender}}</div>'
    });

    var vm = new Vue({
        el: '#app',
        data: {
            user: {
                username: 'LILPIG',
                age: 16,
                gender: 'male'
            }
        }
    })
</script>
```

这样有些啰嗦，不过Vue的`props`支持直接传入对象。

```html
<div id="app">
    <user-profile :user="user"></user-profile>
</div>

<script>

    Vue.component('user-profile',{
        props: ['user'],
        template: '<div>{{user.username}} -- {{user.age}} -- {{user.gender}}</div>'
    });

    var vm = new Vue({
        el: '#app',
        data: {
            user: {
                username: 'LILPIG',
                age: 16,
                gender: 'male'
            }
        }
    })
</script>
```

### 传入对象的所有属性
```html
<div id="app">
    <user-profile v-bind="user"></user-profile>
</div>

<script>

    Vue.component('user-profile',{
        props: ['username','age','gender'],
        template: '<div>{{username}} -- {{age}} -- {{gender}}</div>'
    });

    var vm = new Vue({
        el: '#app',
        data: {
            user: {
                username: 'LILPIG',
                age: 16,
                gender: 'male'
            }
        }
    })
</script>
```

## Props验证
可以验证Props的类型。

```html
<div id="app">
    <user-profile v-bind="user">

    </user-profile>
</div>
<script>
    Vue.component('user-profile',{
        props: {
            username: String,
            age: Number,
            gender: String
        },
        template: '<div>{{username}} -- {{age}} -- {{gender}}</div>'
    })
    var vm = new Vue({
        el: '#app',
        data: {
            user: {
                username: 'LILPIG',
                age: '16',
                gender: 'male'
            }
        }
    })
</script>
```
`props`中的age需要一个`Number`类型，但是传入的是一个字符串，打开控制台会看到这样的错误信息：
```text/plain
[Vue warn]: Invalid prop: type check failed for prop "age". Expected Number with value 16, got String with value "16".
```

除此之外，你还可以这样定义`props`的验证规则：

```javascript
Vue.component('my-component', {
  props: {
    // 基础的类型检查 (`null` 和 `undefined` 会通过任何类型验证)
    propA: Number,
    // 多个可能的类型
    propB: [String, Number],
    // 必填的字符串
    propC: {
      type: String,
      required: true
    },
    // 带有默认值的数字
    propD: {
      type: Number,
      default: 100
    },
    // 带有默认值的对象
    propE: {
      type: Object,
      // 对象或数组默认值必须从一个工厂函数获取
      default: function () {
        return { message: 'hello' }
      }
    },
    // 自定义验证函数
    propF: {
      validator: function (value) {
        // 这个值必须匹配下列字符串中的一个
        return ['success', 'warning', 'danger'].indexOf(value) !== -1
      }
    }
  }
})
```

## 自定义事件
```html
<div id="app">
    <custom-component @custom-event="onCustomEvent">Custom Component</custom-component>
</div>

<script>
    Vue.component('custom-component',{
        template: "<div><slot></slot><button @click=\"sendEvent\">\n" +
            "            Send Event\n" +
            "        </button></div>",
        methods: {
            sendEvent: function(){
                this.$emit("custom-event");
            }
        }
    })
    var vm = new Vue({
        el: '#app',
        methods: {
            onCustomEvent: function(){
                alert("On Custom Event");
            }
        }
    })
</script>
```

可以像如上一样自定义一个事件，利用`$emit`函数和`v-on`。

## 插槽
在上面的示例中，`custom-component`的模板中有一对`slot`标签，里面没有任何内容。如果你删掉这个标签，那么你会发现我们引用这个组件时在其中写的文字`Custom Component`没了。

Vue在渲染组件时会把组件中的内容放到插槽中，如果组件模板没有定义插槽那么将会舍弃内容。

插槽中可以提供一个默认的备用值，在组件引用时没有填写内容时显示。

### 具名插槽
```html
<div id="app">
    <page>
        <template v-slot:header>
            <h3>Header</h3>
        </template>

        <p>
            Hello,World
        </p>

        <template v-slot:footer>
            <strong>Footer</strong>
        </template>
    </page>
</div>
<script>
    Vue.component('page',{
        template: '<div class="container">\n' +
            '    <header>\n' +
            '        <slot name="header"></slot>\n' +
            '    </header>\n' +
            '    <main>\n' +
            '        <slot></slot>\n' +
            '    </main>\n' +
            '    <footer>\n' +
            '        <slot name="footer"></slot>\n' +
            '    </footer>\n' +
            '</div>',
    })
    var vm = new Vue({
        el: '#app'
    })
</script>
```
可以给插槽取个名字，然后使用`template`标签来包裹插入的内容。默认的名字是`default`。

插槽访问数据是很必要的一件事。

```html
<div id="app">
    <custom-component :user="user">
        <template v-slot:default="slotProps">
            {{slotProps.user.name}} -- {{slotProps.user.age}}
        </template>
    </custom-component>
</div>
<script>
    Vue.component('custom-component',{
        props:['user'],
        template: "<div><slot :user='user'>{{user.name}}</slot></div>"
    })
    var vm = new Vue({
        el: '#app',
        data: {
            user: {
                name: "LILPIG",
                age: 16,
                gender: "male"
            }
        }
    })
</script>
```

## is

is属性可以指定当前的模板使用哪个组件，这个可以用来制作标签页功能。

```html
<div id="app" class="demo">
    <button
            v-for="tab in tabs"
            v-bind:key="tab"
            v-bind:class="['tab-button', { active: currentTab === tab }]"
            v-on:click="currentTab = tab"
    >{{ tab }}</button>

    <component
            v-bind:is="currentTabComponent"
            class="tab"
    ></component>
</div>

<script>

    Vue.component('tab-home',{
        template: "<div>Home Page</div>"
    })

    Vue.component('tab-post',{
        template: "<div>Post Page</div>"
    })

    var vm = new Vue({
        el: '#app',
        data: {
            currentTab: 'Home',
            tabs: ['Home','Post']
        },
        computed:{
            currentTabComponent: function () {
                return 'tab-' + this.currentTab.toLowerCase();
            }
        }
    })
</script>
```

<div id="exm">
    <iframe id="exm" src="http://host.lilpig.site/vue-exm/e16.html"></iframe>
    <a href="http://host.lilpig.site/vue-exm/e16.html"><button>在浏览器中运行</button></a>
</div>

每次vue都是重新渲染这个组件，这会导致丢失换页前用户的操作，比如选中checkbox，所以分页功能需要保存页面状态。可以通过添加`keep-alive`标签实现缓存功能。

```html
<keep-alive>
    <component
            v-bind:is="currentTabComponent"
            class="tab"
    ></component>
</keep-alive>
```
## 组件中的v-model
我们自己编写的组件也可以使用`v-model`，`v-model`中的变量会动态传入给组件`props`中的`value`，这是固定的。这是组件的父节点到组件的单向绑定，而从组件到父节点的绑定需要在组件中用`$emit`发射`input`事件。

如下代码可以将`message`和`my-input`组件进行双向绑定，无论哪边的数据更改了都会响应。
```html
<div id="app">
	<p>{{message}}</p>
	<my-input v-model="message"></my-input>
</div>

<script>
	Vue.component('my-input',{
		props: ['value'],
		template: '<div><input :value="value" @input="emitEvent"/>\
				</div>',
		methods: {
			emitEvent: function(event){
				this.$emit('input',event.target.value);
			}
		}
	})
	var vm = new Vue({
		el: '#app',
		data: {
			message: 'Hello,World!'
		}
	})
</script>
```