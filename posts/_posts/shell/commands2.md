---
title: xargs和tr
date: 2019-09-14 08:11:27
tags: [shell]
categories: Shell
---

## xargs
xargs是一个把stdin转换成参数列表的工具。

比如有如下一个文件`args`：
```
arg1
arg2
arg3
```
你想把其中每一行都输出一次（而不是一次全部输出），那么xargs就派上了用场。

```bash
cat args | xargs -n 1 echo
```
xargs指定了一个`-n`参数，意思是每一个参数执行一次，上面的示例中xargs会把文件内容根据回车分割成三份，并把每一份交给echo命令执行一次。

最终的结果为：
```
arg1
arg2
arg3
```
虽然这看起来和直接cat没有区别，但是本质上它们是不一样的。

把上面的命令稍微修改一下：
```bash
cat args | xargs -n 2 echo
```
输出的结果是这样的：
```
arg1 arg2
arg3
```
显然，xargs一次给echo命令提交了两个参数。

### 指定分隔符
默认的分隔符可能是回车，可能是空格，可能是`\0`。但是很多时候我们需要自己指定。比如：

```bash
cat /etc/passwd | xargs -n 1
```
上面的代码通过xargs输出`/etc/passwd`的每一行，这个文件熟悉Linux的应该都知道，每一行是用`:`分割的。我们用xargs来分割它。

xargs用`-d`参数指定分隔符：
```bash
cat /etc/passwd | xargs -n 1 | xargs -n 1 -d :
```

### 指定参数位置
假如我们需要编写如下功能，输出每一个参数的前后加上`@`，也就是说xargs需要把参数放到中间，xargs提供了占位符功能实现这一需求。

```bash
cat args | xargs -I {} echo @ {} @
```
输出：
```
@ arg1 @
@ arg2 @
@ arg3 @
```

### 结合find
```bash
find . -type f -name '*.txt' -print0 | xargs -0 rm -f
```
如上命令删除当前目录下所有txt文件。`-print0`是用`\0`分割，而xargs的`-0`参数也是将分隔符设置成`\0`。一定不要设置简单的分隔符，或使用默认的，这样可能会误删文件。

```bash
find ~/source/gesture4book/app/src/main/ -type f -regex '.*\.\(java\|xml\)$' | xargs wc -l
```
如上命令统计某个项目下的代码行数。

### 循环和xargs
xargs不能执行多条语句，而循环可以，所以有的时候用循环也不失为一个好办法。

```bash
cat args | ( while read arg; do echo $arg; done )
```
## tr

tr用于文本转换，将一个字符串转换成另一个字符串。

### 大小写转换
```bash
cat args | tr 'a-z' 'A-Z'
```

输出:
```
ARG1
ARG2
ARG3
```
### ROT13
ROT13是一种简单的加密方法，通过把所有字母后移13位，比如a变成n。

```
echo "Yahoo~" | tr 'a-zA-Z' 'n-za-mN-ZA-M'
```
输出：
```
Lnubb~
```
### 删除指定字符集
比如我们的输入中有些无关紧要的数字：
```
echo 'Hello12,Wo1r34ld34!12'
```
我们想通过tr删除掉：
```bash
echo 'Hello12,Wo1r34ld34!12' | tr -d [0-9]
```
输出：
```
Hello,World!
```

除了`-d`删除，也可以用补集的方式：
```bash
echo 'Hello12,Wo1r34ld34!12' | tr -d -c [0-9]
```
这样的话就只会留下数字。

```
echo 'Hello12,Wo1r34ld34!12' | tr -c [0-9] ' '
```
这样的话会把所有不在集合1里的字符替换成空格。结果：
```bash
     12   1 34  34 12 
```

### 去除多余字符
刚刚我们利用补集把所有非数字变成了空格，但是这样参差不齐的空格有点丑，我们想只留下一个。

```bash
echo 'Hello12,Wo1r34ld34!12' | tr -c [0-9] ' ' | tr -s ' '
```
`-s`在遇到多个指定字符连续排列时，会只留下一个指定字符，删除多余的。常用来删除多余空格和换行符。

### 计算总和
假设我们有这样一个文本`nums`：
```bash
frist 1
second 2
third 3
```
我们想把所有的数字加起来：
```bash
cat nums | tr -d [a-z] | echo "total: $[ $( tr ' ' '+'  ) ]"
```

