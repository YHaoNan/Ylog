---
title: cat和find
date: 2019-09-13 12:58:27
tags: [shell]
categories: Shell
---

平时开发用Ubuntu，SHELL应该是接触的最多的东西，但是一直也没系统的学习SHELL，就是靠着平时遇到需求就百度的积累。但是这样不太行，很多东西还是不会用，而且在SHELL中总会有更简单快速的解决办法，所以最近在读《Linux Shell脚本攻略》一书，顺便记个笔记。

## cat
cat用于拼接文件或基本输入，应该说是非常常用的命令了。

对于cat，以前我只会`cat xxx`来查看文件内容，其实cat的功能还是不少的。

```bash
cat file1 file2 file3
```
如上命令可以把三个文件的内容拼接，默认cat会发送到stdout，所以会显示到终端中，可以用重定向符号重定向到另一个文件。比如：
```bash
cat emp1 emp2 emp3 > allemp
```
当然也可以简化成这样：
```bash
cat emp* > allemp
```

cat还可以从标准输入中读取：
```bash
ls | cat -n
```
如上命令使用管道符，将ls的输出作为下一条命令的stdin，cat可以读取这个stdin，使用`-n`为每个项目加上行号。

如果想将stdin和其他文件的数据拼接在一起，可以用`-`代表stdin的数据：
```bash
ls | cat - file1.txt
```

使用`-s`可以去除多余的空白行：
```
//tmp.txt
line1

line2


line3



line4

```
```bash
cat tmp.txt
```
执行如上命令可以将多行空行压缩成一行。（注意，不是删除所有空行）

> 注意，cat命令的目的是拼接，所以绝不会修改你的文件，cat规定如果输入和输出是一个文件的话，则会报错，如果使用管道符号重定向输出会清空文件。

## find
find命令也是很常用的，但是其实它有很多功能我们都没发挥出来。

### 查找
find可以使用`-a`和`-o`进行逻辑与和或操作（或者用`-and`和`-or`），比如：
```bash
find . \( -name '*.txt' -o -name '*.md'\)
```
如上命令查找当前目录下所有的`txt`和`md`文件。
`-path`可以指定匹配的路径，因为find命令是递归查找的，会查找子文件夹，所以可以通过`-path`指定一个路径，在指定的路径下的文件才会被匹配。比如：
```bash
find / -path '*/node_modules/*' -name -type d axios
```

该命令查看所有`node_modules`下的`axios`文件夹，运行之后会从根目录开始查找所有nodejs项目中引用的axios库。

`-regex`可以使用正则表达式，这样的话就可以避免麻烦的`-a`、`-o`等选项。如：
```bash
find . -regex '.*\.\(py\|sh\)$'
```
如上命令查找当前目录下所有`py`和`sh`文件。

`-iregex`忽略大小写。

`!`选项可以启动排除模式，会显示所有不符合规则的文件。

```bash
find . ! -name '*.txt'
```
如上命令查找当前目录下所有非`txt`文件

Linux系统提供链接，可以用`ln`命令创建，这和windows的快捷方式很像。find命令默认不会跟随链接，`-L`可以启动跟随链接。

可以通过`-maxdepth`和`-mindepth`选项可以限制find遍历目录的深度，以免因为指向自身的链接造成死循环。

### 类型
`-type`可以指定查找的类型，它接受如下类型：

类型|标识符
:-:|:-:
普通文件|f
符号链接|l
目录|d
字符设备|c
块设备|b
套接字|s
FIFO|p

### 根据时间
find还可以根据文件的时间戳进行搜索。

`-atime`可以指定用户最近一次访问文件的时间，`-mtime`指定用户最近一次修改的时间，`-ctime`可以指定文件元数据（权限等）最后一次改变的时间。

这三个参数以天为基准，并用`-`和`+`分别表示小于和大于，比如：
```bash
find . -type f -atime -7 -name '*.md'
```
以上命令查看最近七天内被访问过的md文件。
```bash
find . -type f -atime 7 -name '*.md'
```
以上命令查看恰好七天前被访问过的md文件。
```bash
find . -type f -atime +7 -name '*.md'
```
以上命令查看最近访问时间超过七天的md文件。

还有更精确的以分钟为单位计时的选项：`-amin`、`-mmin`、`-cmin`。

`-newer`可以找出比某个文件修改时间更近的所有文件：
```bash
find . -type f -newer file.txt
```
### 根据文件大小
find命令还可以基于文件大小进行搜索
```bash
find . -type f -size +2k
```
上面命令查找当前目录下所有大于2k的文件，还可以用`-`和什么都不用。

还支持如下文件单位：
* `b` 块（512字节）
* `c` 字节
* `w` 字（2字节）
* `k` KB
* `M` MB
* `G` GB

### 根据文件权限
`-perm`可以根据权限来匹配文件。

```bash
find . -type f -perm 644
```
匹配所有权限为664的文件

```bash
find . -type f -user lilpig
```
匹配所有所有权为lilpig的文件。

### 执行操作
find命令还可以对匹配的文件执行一些操作，如：
```bash
find . -type f -name "*.swp" -delete
```
删除所有扩展名为`swp`的文件
```bash
find . -type f -name '*.sh' -exec chmod u+x {}
```
对于当前文件夹下的所有sh文件添加执行权限。

### 跳过指定目录
```bash
find project -name '.git' -prune -type f
```
使用`-prune`可以指定跳过的目录，上面命令跳过`.git`目录。
