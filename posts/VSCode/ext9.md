---
title: VSCode 扩展开发(九) Editor操作
date: 2019-06-19 7:35:54
tags: [vscode,editor]
categories: VSCode
---
## 学习目标
使用VSCode API操作编辑器，了解常用的API。

其他API请去官方文档。

## 获取Editor
VSCode中的编辑器是一个TextEditor对象，主要通过以下方法获取：
1. vscode.commands.registerTextEditorCommand
2. vscode.window.activeTextEditor

第一个方法适合与命令搭配使用的情景，我们之前用到过，用这个方法注册命令我们需要传入一个回调函数，这个函数有一个TextEditor对象的参数，VSCode会帮你把当前的TextEditor传入，我们在里面就可以直接使用它了。

第二个方法没什么可说的，就是获取当前的TextEditor。

```typescript
context.subscriptions.push(
    vscode.commands.registerTextEditorCommand('lilpig.command1',(editor)=>{
        if(editor == vscode.window.activeTextEditor){
            console.log(true);
        }
    })
)
```
运行上面代码，控制台会输出true。

也就是说，通过这两个方法获得的Editor都是当前编辑器，只是使用场景不同而已。

## Position和Range
在之前一篇文章中简要提过。

可以把编辑器中的一段文本想象成二维的面，有行有列，通过第几行第几列我们可以精准的定位到一个字符。

`vscode.Position`是位置对象，它是一个点，它可以定位到一个准确的字符。它有如下属性：
* line  
    所在行
* character  
    所在列，也就是字符在该行的位置

这两个属性从0开始。

而`vscode.Range`对象可以表示一个二维的范围，就是在编辑器的文本中取出一段。它由两个Position对象构成，一个是开始位置，一个是结束位置。

它有如下属性和方法：
* `constructor(start: Position, end: Position)`
* `constructor(startLine: number, startCharacter: number, endLine: number, endCharacter: number)`
* `start`  
    开始位置
* `end`  
    结束位置
* `contains(positionOrRange: Position | Range): boolean`  
    范围中是否包含range或position


Range还有很多集合操作，如交集并集等，不一一介绍。
## 编辑器API

属性
* ### document
    > 文档对象，用它来访问当前编辑器中的所有文本，在provider里也很常见，通常作为provider的第一个参数

    * `document.getText(range?: vscode.Range)`  
        根据Range对象获取指定范围文本，如果是undefined，则获取全部
    * `document.getWordRangeAtPosition(position: vscode.Position,regex?: RegExp)`  
        根据Position对象获取范围单词范围，后面可以跟单词规则的正则，该方法常用来获取当前的一个单词。
    * lineAt  
        定位到行

* ### selection & selections
    > 前面我们用过这个东西，就是选中的段落，是Range的子类，没什么好说的

    需要注意的是，无论是否选中，都存在一个selection，只不过未选中时它是空的（isEmpty()==true）。

方法
* `insertSnippets(snippet: SnippetString, location?: Position | Range | Position[] | Range[], options?: { undoStopBefore: boolean; undoStopAfter: boolean; }): Thenable<boolean>`
    > 向编辑器中插入Snippets，第二个参数不指定则在当前位置插入

* `edit(callback: (editBuilder: TextEditorEdit) => void, options?: { undoStopBefore: boolean; undoStopAfter: boolean; }): Thenable<boolean>`
    > 操作编辑器，传入一个带有`TextEditorEdit`的回调

## TextEditorEdit
前面我们见过一个类似的，就是在做自动完成功能时，有一个`additionalTextEdits`，这是一个`TextEdit`对象，他们都差不多，只不过`TextEditorEdit`是针对`TextEditor`的。

这个对象主要用于对编辑器做一些常规操作，插入，替换，删除等，有如下方法：
* `replace(location: Position | Range | Selection, value: string): void`  
    > 将给定范围内的文本替换成value
* `insert(location: Position, value: string): void;`   
    > 将value插入到指定位置
* `delete(location: Range | Selection): void;`   
    > 删除指定范围
* `setEndOfLine(endOfLine: EndOfLine): void;`  
    > 设置换行符 有LF和CRLF两种，不常用

**[点此返回目录](/post/vscode-ext1#目录)**

---

## 参考
* [TextEditor](https://code.visualstudio.com/api/references/vscode-api#TextEditor)