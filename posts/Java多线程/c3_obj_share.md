---
title: （二）数据可见性以及共享对象 —— 《Java并发编程实战》
date: 2019-07-26 19:13:30
tags: [java,thread]
categories: Java多线程
---

从上一篇文章中可以知道，同步代码块，也就是`synchronized`标识的代码块可以保证一个操作的原子性，即不可分割，保证了不会出现竞态条件，其他线程读取到的值不会错乱。

但是`synchronized`并不是只干了这一件事，他还保证了数据的可见性。

## 什么是可见性
可见性即能保证一个线程修改了一个数据之后，这个数据对另一个线程是可见的。

可见性说起来简单，但是却很难理解，因为它里面的很多东西都是违背我们的认知规律并且难以用示例说明的。

你可能会对可见性产生疑问，既然使用`synchronized`时只能有一个线程进行操作，所有线程都是顺序通过的，那怎么会有可见性问题？如果线程A先争抢到锁之后，把数据改了，那么当它释放锁，线程B得到锁的时候，这个数据已经改完了啊，这和单线程串行程序的道理是一样的，那么A改的数据对B，对其他后来线程肯定是可见的啊。为什么会把`synchronized`能保证可见性这件显而易见的事拿出来单说？

没那么简单，我们先看代码：

```java
package site.lilpig.tlearn02;

public class NoVisibility{
    private static boolean ready;
    private static int number;

    private static class ReaderThread extends Thread{
        public void run(){
            while (!ready){
                Thread.yield();
            }
            System.out.println(number);
        }
    }

    public static void main(String[] args){
        new ReaderThread().start();
        number = 42;
        ready = true;
    }
}
```

上面的代码看起来很容易判断结果，按照正常思路应该是结束循环并输出42，但是真的是这样吗？可能不是，它可能是死循环或者是输出0。为啥呢？分情况讨论。

造成死循环的原因是主线程修改后的`ready`对子线程来说并不可见。为什么会造成这个情况呢？这要从Java内存模型JMM说起。JMM有一个主内存区，每个线程操作数据时要把数据从主内存区复制到自己的工作内存区进行操作，JVM会选择适当的时机把工作内存和主内存区的数据进行同步，造成死循环的原因就是主线程的数据修改和`ReaderThread`中的数据不是同步的，可以通过同步代码块和`volatile`关键字修饰。`volatile`修饰的变量进行读写时会自动同步。

造成输出0的原因是代码的重排序，这个问题也可以称作“有序性”。

> 即使在单线程程序中，每条代码也不一定是顺序执行的

没错，我也被震惊到了，但是事实就是这样，即使感觉上这件事违背了常理。

代码执行的时候，内存的速度远比cpu的运算速度要慢，为了能发挥出cpu的最优性能，会将代码经过多层重排序，编译器会按照一定的规则将指令重排，指令到了cpu，cpu又会根据自己的规则对指令进行重排。比如下面的单线程代码：
```java
int a = 10;
int b = 20;
```
不要直接认为第一行先于第二行执行，这实际上是没准的事，不过重排序肯定会遵循一些规则，否则我们的代码逻辑就乱了套了。

重排序遵循`as-if-serial`语义，这名字直译过来是“好像是连续的...”，我们简单了解下这个语义。

对代码进行重排要满足数据依赖性，如果两个操作访问同一个变量，且这两个操作中有一个为写操作，此时这两个操作之间就存在数据依赖。如果存在数据依赖，那么重排就会影响代码的执行结果，`as-if-serial`语义保证了在单线程程序中重排不会影响代码执行结果。

但是多线程的话，程序的执行结果无法预判，所以编译器不好从语义上判断数据是否存在依赖性，所以`as-if-serial`在多线程程序中保证不了数据不会混乱，这下就出现了数据可见性问题。

上面的代码可能就是因为主线程的`ready`和`number`的赋值操作被重排了，`ready`先为`true`，然后`ReaderThread`读取到了，结束了循环并输出，这时主线程的`number`赋值还没执行，所以输出0。

`volatile`关键字修饰的变量除了会进行同步外，用它修饰的变量上面的操作不能被重排到下面，下面不能重排到上面，就像这样：
```java
int a;
int b;
volatile int x;
int y;
```
下面进行赋值操作：
```java
a=0;
b=1;
x=a;
y=b;
a=1;
```

上面代码`x`是`volatile`的，所以赋值操作中，第三行上面的两行，也就是对`a`、`b`赋值的操作不能重排到第三行下面，但内部可以重排，就是第一行和第二行可以重排，第三行下面的同样也是一样不能排到上面，但内部可以重排。

你可能跑去执行代码了但是没有出现问题，其实这些问题应当是个极小的概率事件吧，我们理解就好。虽然没有真正的把这个现象重现出来，但是并不代表它不存在，万一上帝就是那么偏爱你，你也没有办法，所以在多线程代码中我们要预防它的发生，否则可能会出现很大的问题。

对了，别忘了我之前说过

<span style="font-size:1.3em">synchronized会在一个线程退出后同步所有操作，所以能保证数据对其他线程的可见性</span>


## 最低安全性
Java可以保证就算可见性问题发生了，某个线程读到了失效的数据，但这个数据也是以前某个时刻某个线程设置过的，而不是一个随机值，如果做一些对数据要求不太精确的应用，比如网页访问时的计数器时，最低安全性还是能给我们一个保障的。

不过有一个例外，非原子的64位变量比如`long`和`double`无法保证最低安全性，JVM允许把它们的读取和写入操作分成两次，一次操作32位，所以如果发生了可见性问题，有可能读到的数据前32位是之前设置过的一个值的前32位，后32位是另一个之前设置过的值的后32位，这样就会出现问题。

解决办法是用`volatile`关键字声明它们或者用`synchronized`保护它们。

因为`volatile`不会产生阻塞，只是保证变量的可见性，所以它也保证不了操作的原子性。

## 发布对象
经常会有在多个线程中共享对象的场景，这时我们如何保证对象的可见性和有效性呢？

### 逸出
我们先看看逸出的概念，当发布一个不该被发布的对象时就称作这个对象逸出，可能会导致线程安全问题。

```java
public static List<Obj> lists;

public static init(){
    lists = new ArrayList();
}
```

上面我们通过最简单的公开静态域发布了一个`lists`对象，要注意发布一个对象时，对象中的所有共有属性和方法也会被发布，比如上面发布了`lists`，它其中的所有元素也都被发布了。

我们来看一个逸出的例子：
```java
private String[] state = {"A","B","C"};

public getState(){return state;}
```
上面把`state`设置成了私有，而通过`getState`方法其他线程可以随意的修改`state`中的内容，这时`state`已经逸出了它的私有作用域。

逸出会造成的危险不可估量，因为我们也不知道其他线程会对逸出的对象做什么，就算它们什么都不做，危险也是始终存在的，就像你的密码被发布到了公开网络上，虽然可能还没人对你的账号做手脚，但是你的账号已经不安全了。

```java
public class ThisEscape{
    public ThisEscape(EventSource event){
        source.registerListener(new EventListener(){
            @Override
            public void onEvent(Event e){
                doSomething(e);
            }
        })
    }
}
```

上面的代码会使`this`引用逸出，因为你在构造方法中就通过匿名内部类把`doSomething`传给了`source`控制，就是说当你把`doSomething`的控制权交给外部时，`ThisEscape`实例可能还没初始化完成，就算它位于构造方法的最后一行。

在构造函数启动一个线程时如果在线程内引用对象某个属性或方法的话，也会造成`this`引用逸出，解决办法就是不在构造函数内`start`它，而是提供一个`init`方法在对象构造完成后启动。

对于上面的`ThisEscape`可以这样修改：
```java
public class ThisEscape{
    private final EventListener listener;
    public ThisEscape(){
        listener = new EventListener(){
            @Override
            public void onEvent(Event e){
                doSomething(e);
            }
        }
    }

    public static ThisEscape newInstance(EventSource event){
        ThisEscape instance = new ThisEscape();
        event.registerListener(instance.listener);
        return instance;
    }
}
```

### 线程封闭
保证线程安全的最有效的办法就是不共享数据。。。。就像保证你电脑不被黑客入侵的最有效的办法就是不用电脑。

#### ad-hoc线程封闭
ad-hoc线程封闭很好实现，你只需要记几（自己）控制记几不去用其他线程的数据就好了。（应该是这样理解）

但是这个方法还是不要用，因为你可能控制不了记几，并且因为你并没做一些语法层面的限制所以Java也控制不了记几了。

#### 栈封闭
栈封闭就是把对象都限制在局部变量中，其他线程无法访问。

```java
public int getPassCount(List<Student> students){
    int num = 0;
    List<Student> studentsCopy = new ArrayList();
    studentsCopy.addAll(students);
    for(Student s:studentsCopy){
        if(s.isPassed()){
            num++;
        }
    }
    return num;
}
```

上面的代码我们返回了一个基本类型，在Java中无法担心基本类型的引用被另一个线程获得，因为Java中的基本类型都是值传递，而且我们把`studentsCopy`封装在了方法中，其他线程无法修改。

#### ThreadLocal类
ThreadLocal是Java提供的一个类。

很多时候为了节省资源，我们会提供一个全局变量来缓存对象，如果对象之前被创建过我们就直接去拿那个全局变量，这在单线程程序中是安全的，但是在多线程程序中可能不安全。ThreadLocal可以轻易的给每个线程分配一个自己的全局变量，各个线程之间互不打扰。

```java
package site.lilpig.tlearn02;

public class ThreadLocalLearning {


    private static ThreadLocal<String> threadLocal = new ThreadLocal(){
        @Override
        protected Object initialValue() {
            return createObject();
        }
    };


    public static String getObject(){
        return threadLocal.get();
    }


    public static void main(String[] args) {
        new Thread(){
            @Override
            public void run() {
                System.out.println(getObject());
            }
        }.start();

        new Thread(){
            @Override
            public void run() {
                System.out.println(getObject());
            }
        }.start();

        new Thread(){
            @Override
            public void run() {
                System.out.println(getObject());
            }
        }.start();
    }



    /**
     * 模拟一个耗时创建对象的场景，假设这个对象是线程不安全的，比如获得数据库的Connection
     * @return
     */
    private static String createObject(){
        try {
            Thread.sleep(500);
            return new String(Thread.currentThread().getName());
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        return null;
    }
}

/**
* Output:
* Thread-0
* Thread-1
* Thread-2
*/
```

### 不变性
使用不可变对象依然可以保证数据的安全。不可变对象即创建后不能被修改的对象。

```java
@Immutable
static final class ImmutableObj{
    private final int x,y;
    public ImmutableObj(int x,int y){
        this.x = x;
        this.y = y;
    }

    public int sum(){
        try {
            Thread.sleep(100);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        return x+y;
    }
}
```
上面的类是一个不可变对象的类，`@Immutable`标记可以告诉其他开发人员该类是不可变的，类被定义成`final`保证了它不可被继承和改写，语义上也进一步告诉开发人员该类不可变。


该类在构造方法中传入x和y，并且声明为`final`，就算它们不是`final`也没有任何线程能修改它。

我们可以利用这个不可变对象工作。

```java
public static void main(String[] args) {
    class TaskThread extends Thread{
        @Override
        public void run() {
            super.run();
            Random random = new Random();
            int x = random.nextInt(10);
            int y = random.nextInt(10);
            int result = new ImmutableObj(x,y).sum();
            assertTrue(result == x+y,"The result of obj.sum() is "+result+", x is "+x +", y is "+y);

        }

        public void assertTrue(boolean condition,String msg){
            if (!condition)throw new RuntimeException(msg);
        }
    }

    for (int i=0;i<100;i++){
        new TaskThread().start();
    }
}
```

由于对象的不可变性，异常将永远不会被抛出。

声明一个不可变的属性时把他设置为`final`是一个好习惯，其他开发人员能很快的明白你的意图，尽管`final`关键字并不能限制对对象属性的修改。

我们还可以用不可变对象做一个缓存功能：
```java
@Immutable
static final class SumResultCache{
    private final int x,y,result;
    public SumResultCache(int x,int y,int result){
        this.x = x;
        this.y = y;
        this.result = result;
    }

    public Integer get(int x,int y){
        if (this.x == x && this.y == y)return result;
        else return null;
    }
}

static volatile SumResultCache cache = new SumResultCache(0,0,0);
public static void main(String[] args) {
    class TaskThread extends Thread{

        @Override
        public void run() {
            super.run();
            Random random = new Random();
            int x = random.nextInt(10);
            int y = random.nextInt(10);
            Integer result = cache.get(x,y);


            if (result==null){
                result = x+y;
                //假设计算需要经过很长时间，是个耗时操作，所以需要缓存
                try {
                    Thread.sleep(2000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                cache = new SumResultCache(x,y,result);
            }

            assertTrue(result == x+y,"The result of obj.sum() is "+result+", x is "+x +", y is "+y);

        }
        public void assertTrue(boolean condition,String msg){
            if (!condition)throw new RuntimeException(msg);
        }
    }

    for (int i=0;i<10000;i++){
        new TaskThread().start();
    }
}
```
上面我们对`cache`只有一次`get`操作，之后就判断`result`是不是空，如果是就计算，并重新建立缓存。为了避免出现可见性问题，我们要加上`volatile`关键字确保我们对`cache`的修改对其他线程可见。

### 安全发布
#### 不安全的发布
很奇怪哈，明明是在说对象的发布，我们却一直在讨论怎么保证对象不被发不出去，虽然不把对象发不出去能解决很多问题，但是如果我们的确需要在多线程间共享数据，那就没办法了。

来看一个很有意思的例子：
```java
public class Holder{
    private int n;
    public Holder(int n){this.n = n};
    public void assertSanity(){
        if(n!=n)
            throw new RuntimeException("n is not equals to n...");
    }
}
```

最开始看到这个例子，我无法理解`n!=n`这件事，我认为这根本不可能发生，我去问了一个乐于助人的老师，他告诉我是因为重排造成的可见性问题，在多线程中，也许`this.n=n`还没有执行完毕，对象就返回了，所以`assertSanity`中读取的第一个`n`可能是默认值0，这时`this.n=n`执行了，读到的第二个`n`是正确的n，就发生了这个异常。

当然你不要试图运行，发生这件事的概率太小了，几乎无法重现，我的那个老师也说至今没找出重现的办法。

这里要推荐一下老师的课：[老葛课堂 -- 网易云课堂](https://study.163.com/provider/400000000447005/index.htm)

嘶，昨天向老师咨询这个多线程问题的时候我不知道老师有多线程的课，我和他说我是看《Java并发编程实战》这本书产生的疑问，他只是给我解答，对自己的课只字未提，真的是...今天看到他有个收费的多线程的课我莫名有点小感动，不过才5块，在这个知识付费的年代，几千都不一定学到啥，反正我是去支持了。


Java能保证不可变对象的可见性，所以可以通过把`n`声明成final来解决这个问题。

#### 如何安全发布对象
- 在静态域中初始化对象  
    ```java
    public static final Holder holder = new Holder(10);
    ```

    静态域将在JVM初始化类的时候完成初始化，所以是安全的，在单例模式中，常常应用这个技巧。

- 将对象设置为`volatile`或原子对象中  
    Java提供了一系列原子操作的对象，可以保证对象操作的原子性，volatile可以保证可见性。

- 将对象的引用保存到某个正确构造的final类型域中
- 将对象的引用保存到由锁保护的域中



安全发布只能保证发布那一刻的可见性，如果发布的是可变对象，每次访问时还要做同步操作。


---

## 参考
* [多线程之指令重排序](https://www.cnblogs.com/gtaxmjld/p/5274779.html)
