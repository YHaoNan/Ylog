---
title: （一）基本概念 —— 《Java并发编程实战》
date: 2019-07-24 14:13:30
tags: [java,thread]
categories: Java多线程
---

一直对多线程的很多东西很模糊，以前看过一本《Java多线程设计模式》，讲了很多多线程应用程序的设计模式，感觉受益匪浅，但是当时正在准备高考，看了之后没时间实现出来，也没时间记录，现在基本全都忘了。

最近在读《Java并发编程实战》，顺便把读到的记下来，转换成自己的东西。

## 基本概念
多线程的概念、线程和进程的区别这些内容基本上你随便打开一篇多线程相关的文章或者是书籍时都会说，我就不准备在这个地方废话了。


### 线程安全  
线程安全就是说一个类（或者方法？）被多个线程同时访问时，它会正常的工作，不会出错。这个错不一定是显式的异，更多是其中的数据出现错乱。

比如如下代码：
```java
int incr(){return ++counter;}
```
该方法将类中的一个变量`counter`递增返回，变量递增要进行三个操作，先读取`counter`的值，再加一，最后再设置给`counter`，尽管这三个操作对你并不可见，但是Java中它的确是这样的。

这时如果线程A、B同时访问这个方法时，我们希望`incr`方法能正确的递增，但是操作顺序很有可能是这样的
```
counter为0
A读取counter为0
B读取counter为0
A递增 counter为1
B递增 counter为1
A设置counter为1并返回
B设置counter为1并返回
```

它并没有返回正确的结果，这就是线程不安全的例子。
### 原子操作  
所谓原子操作，字面意思就是不可分割的操作，也就是一次完成的操作，上面的递增就不是一个原子操作，它分为三次才能操作完。

### 竞态条件  
竞态条件简单的说就是返回的数据已经丧失了有效性，比如如下的单例模式代码就会有产生竞态：
```java
package site.lilpig.tlearn01;

public class SingletonRace {
    private static SingletonRace instance = null;
    
    public static SingletonRace getInstance(){
        if (instance!=null)
            instance = new SingletonRace();
        return instance;
    }
}
```
假如A、B两个线程共同调用`getInstance`
```
A进入getInstance
B进入getInstance
A判断instance为null
B判断instance为null
A创建instance并返回
B创建instance并返回
```

A和B返回了两个不一样的instance，这在单例模式中是绝对不被允许的，如果后面的线程不会和AB产生竞态的话，那么它们用的应该都是B创建的instance，那么如果线程A给instance设置了一些状态，在整个应用中这些状态都是没用的，因为其他线程用到的instance和A不是同一个。

### 锁  
锁是针对一个对象的，默认情况下的对象就是该对象本身。用锁锁住的操作都可以视为一个原子操作，当一个线程进入时会持有锁，这时其他线程无法进入此操作，当线程操作完毕退出时会释放锁，其他线程开始争抢这个锁。

例如修改之前计数器的代码：
```java
int synchronized incr(){return ++counter;}
```
这样就可以把本来是三步的操作当作一个原子操作，因为操作是不可分割的，所以不可能会产生线程不安全的问题。

### 死锁  
A持有锁1并想获得锁2，B持有锁2并想获得锁1，这个时候两个线程会永远的等待下去。
### 重入  
试想一下，如果有一个这样的代码会出现啥情况
```java
int synchronized incrWithTimes(int times){
    ++counter;
    return times==1?counter:incrWithTimes(times-1);
}
```
如果直接看的话，我们可能认为它会发生死锁，因为在return语句递归调用`incrWithTimes`的时候，线程持有的锁还没有被释放，显然也不能进入这个方法，就是说线程一直等着自己释放锁。

但是它能正确运行，这就说明Java获取锁操作是基于线程的而不是方法，就是说同一个线程可以在自己持有锁的状态下继续获取锁，Java会提供一个计数器，每次获取锁加1，直到所有操作完成，计数器变成0，这时锁被释放。

这个现象叫做重入。

### 活跃性与性能  
如果直接对所有操作加锁，那么肯定会出现性能问题，从加锁的关键字就可以看出，加锁后所有的线程的操作都是串行的，也就是同步的，这完全就发挥不出多线程异步的优势，所以要有针对性的加锁。

假设我们有这样的代码：
```java
package site.lilpig.tlearn01;

import java.util.Random;

public class TaskService {
    private int counter = 0;
    //模拟任务
    public synchronized void doTask(int num){
        //随机选一个200或2000来模拟现实任务中极快和极慢的情况
        int delayMS = new Random().nextBoolean()?200:2000;
        try {
            Thread.sleep(delayMS);
        } catch (InterruptedException e) {
            throw new RuntimeException("Catch an exception...");
        }
        //计数
        counter++;
        System.out.println("Counter:"+counter+",Input:"+num+",Result:"+num*num);
    }

    public static void main(String[] args) {
        TaskService service = new TaskService();
        for (int i=0;i<20;i++){
            lunchANewTask(service,i);
        }
    }

    public static void lunchANewTask(TaskService service,int param){
        new Thread(new Runnable() {
            @Override
            public void run() {
                service.doTask(param);
            }
        }).start();
    }
}

```
这个代码会有一个问题，你会发现把整个方法都加上一个synchronized后，整个操作变成了串行的，所有线程一个一个的通过（因为是抢占锁所以不一定是按顺序），当一个线程调用doTask时，如果随即延时是200毫秒还可以忍受，等200毫秒下一个线程就会执行，而如果是2000毫秒，下一个线程会等好久也得不到执行的机会。如果此时我们设计的是一个web后端应用，如果我们面对的是成千上万的用户，而不是上面示例中的20次调用，那么整个服务器都会阻塞，可能用户今天发起的请求，我们隔了几天才会处理，因为它前面还可能有若干个请求在等待处理...这是很可怕的。

分析上面的代码，其实我们只有计数的时候调用了方法作用域外部的变量，也只有`counter`变量会引发线程不安全问题，所以我们只需要针对`counter`自增操作加锁就可以，这个操作很快，加锁也不会引起阻塞。而上面用`Thread.sleep`模拟的执行任务因为没有调用什么其他变量，不会引起线程安全问题，这个操作根本不用加锁。

所以我们修改代码：
```java
public void doTask(int num){
    int delayMS = new Random().nextBoolean()?200:2000;
    try {
        Thread.sleep(delayMS);
    } catch (InterruptedException e) {
        throw new RuntimeException("Catch an exception...");
    }
    synchronized (this) {
        counter++;
        System.out.println("Counter:"+counter+",Input:"+num+",Result:"+num*num);
    }

}
```
只在用到外部变量的时候加锁，再运行时发现速度很快，上面的模拟任务处理操作也都是同时执行的。

所以千万不要对外部的大型IO操作或者什么耗时操作加锁，你需要根据实际情况进行加锁，而不是上来就加整个方法。


