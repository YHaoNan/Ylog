---
title: VSCode 扩展开发(十) TreeView
date: 2019-06-19 9:35:54
tags: [vscode,editor]
categories: VSCode
---
## 学习目标
了解TreeView的用法

## 开篇寒暄
TreeView？？听起来挺屌的，应该是个什么树形控件。

Ohh对，就是一个树形控件，它就是我们最常用的侧边栏。

![img](http://nsimg.cn-bj.ufileos.com/img-1560908995830.png)

其实我们很多时候都要和侧边栏打交道，VSCode提供了几个默认的侧边栏，比如文件管理、搜索、版本控制、debug、插件市场等。

我们可以通过本篇文章的学习，学会自己开发一个侧边栏工具。

Ohh，不过别高兴得太早，按照国际惯例，本节也不会写一个非常有用的插件，有用的在实战部分呐，到时候我们写几个上班偷懒的插件，比如在侧边栏里刷知乎，听音乐～～

## 配置
首先要在`package.json`的`contributes`节点里配置它一些子节点，如下即是。

`viewContainers.activitybar`节点就是侧边栏中的按钮，在这里配置它的标题，图标和id等信息。

```json
"viewsContainers": {
    "activitybar": [
        {
            "id": "my-tree",
            "title": "My TreeView",
            "icon": "res/tree.svg"
        }
    ]
}
```

而`views`节点则绑定了该按钮弹出的视图：
```json
"views": {
    "my-tree":[
        {
            "id": "treeItems",
            "name": "Tree Items"
        }
    ]
}
```

然后，要在`activationEvents`节点中配置激活事件，我们要在`treeItems`触发时激活。

```json
"activationEvents": [
    "onView:treeItems"
]
```

## TreeDataProvider
TreeDataProvider用于提供TreeView展示的数据，它使用泛型设计，需要一个实体类代表树的一个节点，这个实体类继承自`vscode.TreeItem`：
```typescript
export class MyTreeItem extends vscode.TreeItem {
    //label是该条数据显示的文本，collapsibleState是该节点的状态
    //vscode.TreeItemCollapsibleState是个枚举类，有三个值
    //None表示它不能折叠，也就是没有子节点
    //Collapsed表示它可以折叠但是是未打开状态
    //Expanded表示它可以折叠但是是已打开状态
    constructor(public readonly label: string, public readonly collapsibleState: vscode.TreeItemCollapsibleState) {
        super(label, collapsibleState);
    }

    
    iconPath = {
        light: path.join(__filename, '..', '..', 'res', 'light', 'item.svg'),
        dark: path.join(__filename, '..', '..', 'res', 'dark', 'item.svg')
    }

}
```

这个实体还有很多属性，比如`description`、`tooltip`等，这个查一下API就好了。

有了节点类，我们去创建TreeDataProvider的实现类：
```typescript
export class MyTreeDataProvider implements vscode.TreeDataProvider<MyTreeItem>{
    onDidChangeTreeData?: vscode.Event<MyTreeItem | null | undefined> | undefined;    getTreeItem(element: MyTreeItem): vscode.TreeItem | Thenable<vscode.TreeItem> {
        throw new Error("Method not implemented.");
    }
    getChildren(element?: MyTreeItem | undefined): vscode.ProviderResult<MyTreeItem[]> {
        throw new Error("Method not implemented.");
    }

    
}
```
onDidChangeTreeData是当树的数据改变后的回调，getChildren是获取它的子节点时的回调，也就是展开时的回调。

知道了这些，我们写一个迷你的文件管理器，最后要实现这个效果：

![效果](https://s2.ax1x.com/2019/06/19/VOk0E9.gif)

虽然有点丑哈，但是无伤大雅，分析一下怎么实现。

首先点击它会获取当前工作空间的路径，然后读取子文件，如果是文件夹，则能展开继续获取子文件，如果是文件，则不能。

工作空间路径，我们在外部传入，可以使用`vscode.workspace.rootPath`获取。

import:
```typescript
import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
```

我们先抽象实体类：
```typescript
//两个辅助方法 获取图标
function lightRes(filename: string){
    return path.join(__filename, '..', '..', 'res', 'light', filename); 
}

function darkRes(filename: string){
    return path.join(__filename, '..', '..', 'res', 'dark', filename); 
}


export class MyTreeItem extends vscode.TreeItem {
    //提供一个parent，当它是null时则代表它是根节点
    constructor(public readonly label: string, public readonly collapsibleState: vscode.TreeItemCollapsibleState,public parent: MyTreeItem | null) {
        super(label, collapsibleState);
    }
    //提供一个path的get方法用于获取路径，如果是根节点则就是当前的label，否则还要加上父节点的
    get path(): string {
        return this.parent?path.join(this.parent.path,this.label):this.label;
    }

    //对文件和文件夹进行分类 如果collapsibleState是None则代表它是文件
    iconPath = {
        light: this.collapsibleState == vscode.TreeItemCollapsibleState.None ? lightRes('file.svg') :  lightRes('dir.svg'),
        dark: this.collapsibleState == vscode.TreeItemCollapsibleState.None ?  darkRes('file.svg') :  darkRes('dir.svg')
    }

}
```

然后去写Provider的实现类：
```typescript
export class MyTreeDataProvider implements vscode.TreeDataProvider<MyTreeItem>{
    //让外部传入根路径
    constructor(public rootDir: string){}
    
    //当发生改变后进行刷新UI 但是我不太明白这两行 明白的大佬麻烦告诉我下 嘿嘿～
    private _onDidChangeTreeData: vscode.EventEmitter<MyTreeItem | undefined> = new vscode.EventEmitter<MyTreeItem | undefined>();
    readonly onDidChangeTreeData: vscode.Event<MyTreeItem | undefined> = this._onDidChangeTreeData.event;

    
    getTreeItem(element: MyTreeItem): vscode.TreeItem {
		return element;
	}

    //关键，获取子节点
    getChildren(element?: MyTreeItem | undefined): vscode.ProviderResult<MyTreeItem[]> {
        //当element不存在，则代表它是根节点，这时我们把this.rootDir设进去，并且让它默认展开
        if(!element){
            return Promise.resolve([new MyTreeItem(this.rootDir,vscode.TreeItemCollapsibleState.Expanded,null)]);
        }else{
            let items: MyTreeItem[] = [];
            //否则遍历文件夹
            let subfiles = fs.readdirSync(element.path);
            subfiles.forEach(file=>{
                let stat = fs.statSync(path.join(element.path,file));
                //如果是目录则设置它可以展开，否则设置它不可展开
                items.push(new MyTreeItem(file,(stat.isDirectory()?vscode.TreeItemCollapsibleState.Collapsed:vscode.TreeItemCollapsibleState.None),element))
            })
            return Promise.resolve(items);
        }
    }
}
```

然后在`extension.ts`中注册它：
```typescript
context.subscriptions.push(
    vscode.window.registerTreeDataProvider('treeItems',new MyTreeDataProvider(vscode.workspace.rootPath as string))
);
```
完活。运行试试吧。

## 添加功能
现在它好像啥功能都没有呢，只能展开关闭，这是很蛋疼的一件事。我们给它添加一些功能，比如删除文件，打开文件。

在`menus`节点中我们可以设置关于这些项目的一些操作。

```json
"menus": {
    "view/item/context": [
        {
            "command": "tree.delete",
            "when": "view == treeItems && viewItem == fileItems"
        },
        {
            "command": "tree.delete",
            "when": "view == treeItems && viewItem == dirItems"
        },
        {
            "command": "tree.open",
            "group": "inline",
            "when": "view == treeItems && viewItem == fileItems"
        }
    ]
}
```
`view/item/context`设置每一个项目的菜单，我们提供两个菜单，分别是打开文件和删除文件。

`when`中我们设置了当view是treeItems并且viewItem的值等于特定值时才出现，viewItem我们一会在代码里设置。

menus节点中关于侧栏的其他设置请走[这边](https://code.visualstudio.com/api/extension-guides/tree-view#view-actions)

注册命令：
```json
"commands":[
    {
        "command": "tree.delete",
        "title": "Delete"
    },
    {
        "command": "tree.open",
        "title": "Open"
    }
]
```
当然可以设置图标，我这里没设置。

在再Provider中添加几个方法，供实现功能：
```typescript
delete(element: MyTreeItem){
    let path = element.path;
    let stat = fs.statSync(path);
    if(stat.isDirectory()){
        exec('rm -rf "'+path+'"')
    }else{
        fs.unlinkSync(path);
    }
    this.refresh();
}

open(element: MyTreeItem){
    vscode.window.showTextDocument(vscode.Uri.file(element.path))
}

//刷新
refresh(): void {
    this._onDidChangeTreeData.fire();
}
```

最后在TreeItem类里加入这样一行代码：
```typescript
contextValue = this.collapsibleState == vscode.TreeItemCollapsibleState.None ? "fileItems" : "dirItems";
```
这就是设置刚刚package.json中viewItem的取值。

然后就是extension.ts
```typescript
const provider = new MyTreeDataProvider(vscode.workspace.rootPath as string);
vscode.window.registerTreeDataProvider('treeItems',provider)
vscode.commands.registerCommand('tree.delete',(node: MyTreeItem) => provider.delete(node))
vscode.commands.registerCommand('tree.open',(node: MyTreeItem)=> provider.open(node))
```

![img](http://nsimg.cn-bj.ufileos.com/img-1560992418573.png)

运行，可以达到我们预期的效果，但是我们发现，使用一个open按钮打开文档不太符合我们的使用习惯，我们习惯点击项目打开文件，所以还要继续扩展功能。其实TreeItem有一个command属性，就是用于设置点击该项目时触发的命令，我们可以在TreeItem的实体类中或者provider的`getTreeItem`扩展它。

```typescript
getTreeItem(element: MyTreeItem): vscode.TreeItem {
    element.command = element.collapsibleState == vscode.TreeItemCollapsibleState.None? {
        command: 'tree.open',
        arguments: [element],
        title: 'Open'
    } : void 0;
    return element;
}
```

现在，点击项目的时候也会打开文件了。

**[点此返回目录](/post/vscode-ext1#目录)**

---

## 参考
[tree-view](https://code.visualstudio.com/api/extension-guides/tree-view#)