---
title: GO语言笔记和习题
date: 2019-12-07 14:52:27
tags: [Golang]
categories: golang
---

开始学Go语言了，今天把GO语言的语法中我认为和其他语言差异很大的值得记录的地方记录下来。并记录了一些官方习题。

## 命名返回值
Go语言中的函数返回值可以是命名的，它的好处就是可以直接作为文档使用，并且看起来会很简洁。

```go
package main

import "fmt"
func split(sum int) (x, y int){
	x = sum * 4 / 9
	y = sum - x
	return
}
func main(){
	x,y := split(17)
	fmt.Println(x,y)
}
```

可以看到命名返回值可以直接使用`return`语句返回，只要函数体中有对应变量就好。

不过官方并不推荐将直接返回语句应用在长函数中，因为会影响代码可读性。

## if前置语句
在`for`循环中可以有一个初始化语句，但是大部分语言都没提供一个if语句的初始化语句，然而有时候这是必要的。Go语言允许在if中存在一个初始化语句，如：

```go
if conn,err := connect() ; err != nil {
    //Do something
}
```

上面的代码中，如果`connect()`方法返回的错误为空，就执行if里面的内容。这样做除了简洁外，还可以把`conn`和`err`的作用域限制在if语句内。

## 结构体指针
假如我们有一个指向结构体的指针，可以通过`(*p).X`访问其中的字段，相当于先针对指针取值，再取字段。

在C++中可以通过`->`运算符取字段，go中为了避免麻烦，直接允许通过指针访问结构体中的字段：`p.X`。它会被编译成`(*p).X`的形式。

## 结构体方法
Go没有类，只有结构体，它为结构体定义方法也和其他语言不太相同。

它把结构体看做方法的接收者，你需要在定义的方法签名中指定接收者，接收者在方法签名中类似于一个特殊的参数。

```go
package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
}
```
上面的`Abs`方法既是为`Vertex`结构体定义的方法，`Vertex`为接收者，可以看到Go中需要把接收者放在方法名前。在方法中可以引用这个接收者，上例中的`v`类似于其他面向对象编程语言中的`this`。

Go官方给出的方法定义是：方法只是个带接收者参数的函数。

那岂不就是说我们就可以在任意位置给结构体扩展方法了？就像这样：
```go
func (s string) times(c int) (result string){ //cannot define new methods on non-local type string
    result = ""
    for i := 0; i<c ; i++ {
        result = result + s
    }

    return
}

func main(){
    fmt.Println("a".times(20)) //"a".times undefined (type string has no field or method times)
}
```


显然不是的，运行代码后编译不会通过，因为Go只允许为包内的结构体定义方法。而string的定义显然不在我们的包里。

我们可以通过自定义类型来完成需求，我们只需要这样修改代码：
```go
type MyString string

func (s MyString) times(c int) (result MyString){
	result = ""
	for i := 0; i<c ; i++ {
		result = result + s
	}

	return
}

func main(){
	fmt.Println(MyString("a").times(20))
}
```
## 指针接收者
你可以定义方法的接收者为指针，这样可以修改指向者的值，因为在普通的值接收者（上例）中，Go会先对原始的Vertex对象创建副本，然后再传入方法，也就是说我们操作的都是副本。所以当你的方法需要修改原始对象的值时就使用指针接收者。

```go
import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func main() {
	v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())
}
```
指针接收者接受一个指针，而我们调用的是`v.Scale(10)`，它会被编译为`(&v).Scale(10)`

相反，对于普通函数来说，接受一个值作为参数的函数，你传入的参数既能是值又能是指。



## 切片练习
[在线Playground](https://tour.go-zh.org/moretypes/18)


```
实现 Pic。它应当返回一个长度为 dy 的切片，其中每个元素是一个长度为 dx，元素类型为 uint8 的切片。当你运行此程序时，它会将每个整数解释为灰度值（好吧，其实是蓝度值）并显示它所对应的图像。

图像的选择由你来定。几个有趣的函数包括 (x+y)/2, x*y, x^y, x*log(y) 和 x%(y+1)。

（提示：需要使用循环来分配 [][]uint8 中的每个 []uint8；请使用 uint8(intValue) 在类型之间转换；你可能会用到 math 包中的函数。）
```
```go
package main

import "golang.org/x/tour/pic"
func Pic(dx, dy int) [][]uint8 {
	outer := make([][]uint8,dy)
	for y := range outer {
		inter := make([]uint8,dx)
		for x := range inter {
			inter[x] = uint8(x%(y+1))
		} 
		outer[y] = inter
	}
	return outer
}

func main() {
	pic.Show(Pic)
}
```

## 映射练习
[在线Playground](https://tour.go-zh.org/moretypes/23)

```
实现 WordCount。它应当返回一个映射，其中包含字符串 s 中每个“单词”的个数。函数 wc.Test 会对此函数执行一系列测试用例，并输出成功还是失败。

你会发现 strings.Fields 很有帮助。
```

```go

package main
import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	var mapToResult map[string]int
	mapToResult = make(map[string]int)
	for _,keyword := range strings.Fields(s) {
		if curV,ok := mapToResult[keyword]; ok {
			mapToResult[keyword] = curV + 1
		}else {
			mapToResult[keyword] = 1
		}
	}
	return mapToResult
}

func main() {
	wc.Test(WordCount)
}
```

## 闭包练习（斐波那契闭包）

[在线Playground](https://tour.go-zh.org/moretypes/26)

```
让我们用函数做些好玩的事情。

实现一个 fibonacci 函数，它返回一个函数（闭包），该闭包返回一个斐波纳契数列 `(0, 1, 1, 2, 3, 5, ...)`。
```

```go
package main

import "fmt"

// 返回一个“返回int的函数”
func fibonacci() func() int {
	last,curr := 0 ,1
	return func() (result int){
		result = last + curr
		last = curr
		curr = result
		return
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

## Stringer练习
[在线Playground](https://tour.go-zh.org/methods/18)

Stringer是一个接口，用于其它类型的字符串表示，它像Java中的`toString`。
```
type Stringer interface {
	String() string
}
```

```
通过让 IPAddr 类型实现 fmt.Stringer 来打印点号分隔的地址。

例如，IPAddr{1, 2, 3, 4} 应当打印为 "1.2.3.4"

```

```go
package main

import "fmt"

type IPAddr [4]byte


// TODO: 给 IPAddr 添加一个 "String() string" 方法

func (addr IPAddr) String() string{
	return fmt.Sprintf("%d.%d.%d.%d",addr[0],addr[1],addr[2],addr[3])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}
```


## 练习 Reader
[在线Playground](https://tour.go-zh.org/methods/22)

```
实现一个 Reader 类型，它产生一个 ASCII 字符 'A' 的无限流。

Reader接口有如下方法
func (T) Read(b []byte) (n int, err error)
```

```go
package main

import "golang.org/x/tour/reader"

type MyReader struct{}

// TODO: 给 MyReader 添加一个 Read([]byte) (int, error) 方法
func (reader MyReader) Read(buf []byte) (int,error){
	n := 0
	for i := range buf {
		buf[i] = 'A'
		n++
	}
	return n ,nil
}


func main() {
	reader.Validate(MyReader{})
}
```

## rot13Reader 练习
[在线Playground](https://tour.go-zh.org/methods/23)

```
练习：rot13Reader
有种常见的模式是一个 io.Reader 包装另一个 io.Reader，然后通过某种方式修改其数据流。

例如，gzip.NewReader 函数接受一个 io.Reader（已压缩的数据流）并返回一个同样实现了 io.Reader 的 *gzip.Reader（解压后的数据流）。

编写一个实现了 io.Reader 并从另一个 io.Reader 中读取数据的 rot13Reader，通过应用 rot13 代换密码对数据流进行修改。

rot13Reader 类型已经提供。实现 Read 方法以满足 io.Reader。
```

```go

package main

import (
"io"
"os"
"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (reader rot13Reader) Read(buf []byte) (int,error){
	n , err := reader.r.Read(buf)
	for i,b := range buf {
		switch {
		case b >= 'A' && b < 'N' || b >= 'a' && b < 'n':
			buf[i] = b + 13
		case b >= 'N' && b <= 'Z' || b >= 'n' && b <= 'z':
			buf[i] = b - 13
		}
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
```
