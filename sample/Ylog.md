# Ylog
Ylog是一个由Go编写的轻量级静态博客生成器。

Ylog只是一个Go语言练手项目，如果你想用于自己的博客，Ylog肯定不比Hexo这些成熟的静态博客系统。它没有海量的主题和插件供你使用。

## Ylog VS Hexo
Hexo给我的个人感觉就是对非前端开发者用户不太友好，而且有点乱套。可能是太注重扩展性带来的反作用，反正我第一次配置Hexo博客的时候废了很大力气。当然也不排除是我那时候实在太菜了。

Ylog从设计之初就避开这种复杂性，对于非前端用户，甚至非开发者用户，他们只需要配置一套结构清晰的配置文件就好了。还会提供VSCode插件，可以可视化编写博客。

## Ylog的输出
Ylog输出一套静态的html博客源代码文件，不包含任何CSS。有关页面排版的内容都需要前端去做。这也是舍弃一些可配置性来换可用性了。

## 插件
Ylog中一切皆插件。甚至连将博客Markdown文件编译成HTML代码都是由插件来完成的。

编译插件是Ylog的内置插件，它永远是插件列表中第一个插件。它针对每个待构建的HTML创建一个HTML标签树，后面的所有用户插件都可以对这个标签树进行操作。

Ylog在编译Markdown文件时会按顺序调用用户配置文件里指定的所有插件。

Ylog的最后两个插件是主题插件和静态资源插件，它们总被最后调用，用于把主题要注入的JS和CSS文件还有引用的静态资源打包进最终的博客源码包中。

## 第二类插件
第二类插件是工具类插件。它们不参与整个编译过程，一般是用于在编译后对整个源码包做一些操作。比如push到github，比如生成sitemap。

这类插件通过Ylog命令行调用，所以也叫命令行插件。

编译功能本身也是一个第二类插件。
## 项目结构
```
.
├── config.ymal
├── _out
├── posts
├── pages
└── themes
```

`config.ymal`是项目配置文件。

`_out`是打包出来的项目。

`posts`是用户编写的markdown文件。

`pages`是用户创建的页面文件。

`themes`是项目依赖的主题。之所以叫`themes`是因为这点上Ylog和Hexo的做法是一样的，可以在themes中存放多个主题，然后在`config.yaml`中切换。

## 约束
posts文件夹中的任意深度的文件夹都会被编译到out中，并且只剩下最外层的文件夹。这个最外层的文件夹作为里面博客文件的分类。

```
posts
    Java
        Spring
            spring-note01.md
        reflect.md
```

比如如上文件夹结构会编译成：
```
out
    Java
        reflect.html
    Spring
        spring-note01.html
```

pages中的任意深度的文件夹都会直接编译到`out`文件夹里。如：
```
pages
    test1
        testpage1.md
```
会被编译成
```
out
    testpage1.html
```

所有的markdown文件都必须有一份文件头，该文件头和Hexo一样，只是多了一些约束。
```
---
title: 标题
date: 创建时间
categories: 分类
---
```
这里要说一下，上面三个是必须的字段，而且必须Markdown文件第一行开始就是文件头，否则会报错。

使用`Ylog new post 分类/标题`会自动生成文件并写入文件头。

插件可以使用文件头与用户交互。

categories文件头如果按照Ylog的文件夹分类标准的话其实可以不用填，但是我之前写进去了，然后也懒得改了，请保持该字段和所在文件夹一致。


## 命令行
```
Usage: 
    Ylog <command>

```

`command`是一个第二类Plugin的名字，第二类Plugin是编译好的，在Ylog安装目录下的二进制文件。Ylog本身只有这一个命令，最主要是靠Plugin扩展，一切皆Plugin。

下面来看看自带的Plugin。

```
Ylog build [Path] -- 编译Ylog项目
Ylog serve [Path Port] -- 在本地运行一个HTTP服务器运行所选项目
Ylog gitpush [Path] -- 推送Ylog项目到git仓库
Ylog install [SourceFile] -- 编译安装一个插件（SourceFile为插件的Go语言源码）
Ylog install [PluginName] -- 编译安装一个插件（PluginName为插件名，这只适用于官方代码仓库中包含的Plugin）
```


