---
title: VSCode在Ubuntu上交换Capslock和Esc
date: 2019-09-22 16:35:54
tags: [vscode,editor]
categories: VSCode
---

## 起源
之前入了VIM的坑，真的特别香。

用VIM的应该都交换了Capslock和Esc，这是系统级别的交换，看似和VSCode不发生关系，不过在Ubuntu下的VSCode上确实是不管用的，其他编辑器没出现这种情况。

百度了一番，发现了这个issue：[vscode-issues-23991](https://github.com/microsoft/vscode/issues/23991#issuecomment-292336504)

在设置里加入：`"keyboard.dispatch": "keyCode"`即可。
