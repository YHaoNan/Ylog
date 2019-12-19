---
title: Git常用命令
date: 2019-06-21 17:45:27
tags: [git]
categories: Others
---
```
git init //初始化git仓库
git add [文件...] //添加到仓库
git commit [文件...] -m "注释" //提交一个版本
git status //查看状态
git diff [文件] //查看变更
git diff HEAD -- readme.txt //查看某个版本中的文件与当前文件的不同
git log [--pretty=oneline] //查看提交记录
git reset --hard HEAD //重置到当前版本
git reset --hard HEAD^ //重置到上一个版本
git reset --hard HEAD^^ //重置到上上一个版本
git reset --hard HEAD~100 //重置到第前100个版本
git reset --hard 1094a //重置到版本id为1094a开头的那个版本
git reflog //查看操作记录
git checkout -- readme.txt //丢弃工作区的更改
git reset HEAD readme.txt //丢弃暂存区的修改（从某一个版本中的某个文件覆盖）
git rm test.txt //从版本库中删除文件 需要commit
```
## 说明
上面内容全部来自[廖雪峰 Git教程](https://www.liaoxuefeng.com/wiki/896043488029600)，为了阅读方便放到了这里。