---
title: （八）活跃性问题 —— 《Java并发编程实战》
date: 2019-08-09 11:52:08
tags: [java,thread]
categories: Java多线程
---
活跃性问题的原因有很多，如死锁、活锁等等。出现活跃性问题后，我们的系统可能无法继续正常工作，或者响应性变差，大量资源被浪费，所以避免发生活跃性问题在多线程开发中还是很重要的。

## 死锁
造成死锁的原因是两个线程互相等待获得对方的锁而双方又都不释放锁。例如线程A持有锁L，并想要获得锁M，与此同时，线程B持有锁M，并想要获得锁L，两个线程都在等待着对方释放锁，但是Java的内置锁可没有发现死锁并自动释放的机制，所以两个线程就会僵持不下，造成无限等待。

就像我们小时候和同班的plmm打闹一样，你抓着她的书包挑逗人家，人家抓着你的肉想要抢回自己的书包，你太疼了于是你说：“松手！”  
她拧着你肉的手又向右旋转了一些，说：“你先放！”  
你疼的龇牙咧嘴，但是你还要装作云淡风轻，因为不能在plmm面前丢脸，于是你说：“我不，你先放！”

你俩的行为就是我们说的死锁。你是线程A，plmm的书包是锁L，plmm是线程B，你的肉是锁M。你持有plmm的书包并想获得plmm手捏着的你的肉的所有权，plmm捏着你的肉，并想要夺回书包。

*（PS:我是不知道现在的小孩还这样不了，反正我小时候干过这事，嗯没错就是这样）*

看一个死锁的例子：
```java
public class DeadLockExample {
    private final Object lockA = new Object();
    private final Object lockB = new Object();

    public static void main(String[] args) throws InterruptedException {

        DeadLockExample example = new DeadLockExample();

        Runnable task1 = ()->{
            synchronized (example.lockA){
                synchronized (example.lockB){
                    System.out.println("T1...");
                }
            }
        };

        Runnable task2 = ()->{
            synchronized (example.lockB){
                synchronized (example.lockA){
                    System.out.println("T2...");
                }
            }
        };

        for (int i=0;i<10000;i++){
            Thread t1 = new Thread(task1);
            Thread t2 = new Thread(task2);
            t1.start();t2.start();
            t1.join();t2.join();
        }
    }
}
```

因为获得锁的过程很快，不一定每次都出现死锁，所以我这里用了个循环一直重复10000次，其实在线程获得锁的过程中加一个`Thread.sleep`也可以，但是我总觉得`try-catch`那个中断异常很丑。

造成死锁的原因主要是获得锁的顺序不一致导致的，如果两个线程都获取锁的顺序都是`A-B`就不会发生死锁。

### 动态死锁
思考一下下面的代码是否会发生死锁：
```java
public void transferMoney(Account from,Account to,int money){
    synchronized (from){
        synchronized (to){
            if (from.balance>=money){
                from.balance -= money;
                to.balance += money;
            }
        }
    }
}
```
看起来是不会，因为获取锁的顺序都是从`from`到`to`，但是。。。`from`和`to`可以随便传啊，第一次把用户1当作`from`把用户2当`to`，第二次把用户2当`from`用户1当`to`，这样就可能造成死锁。

可以调整两个账户获得锁的顺序，从而避免死锁。比如使用hashcode来避免死锁。
```java
private final Object tieLock = new Object();
public void transferMoney(Account from,Account to,int money){
    int fromHash = System.identityHashCode(from);
    int toHash = System.identityHashCode(to);
    if (fromHash>toHash){
        synchronized (from){
            synchronized (to){
                if (from.balance>=money){
                    from.balance -= money;
                    to.balance += money;
                    System.out.println(from+" to "+to);
                }
            }
        }
    }else if (fromHash<toHash){
        synchronized (to){
            synchronized (from){
                if (from.balance>=money){
                    from.balance -= money;
                    to.balance += money;
                    System.out.println(from+" to "+to);
                }
            }
        }
    }else {
        synchronized (tieLock){
            synchronized (from){
                synchronized (to){
                    if (from.balance>=money){
                        from.balance -= money;
                        to.balance += money;
                        System.out.println(from+" to "+to);
                    }
                }
            }
        }
    }
}
```
虽然增加了很多额外代码，但是避免了死锁发生的可能性。关于最后一个else块是因为hashcode可能相等，这时`from`和`to`其实没办法分辨，先锁哪个都不行，所以外面再加一层锁，这样就可以避免因为hashcode冲突而造成的死锁了。

如果你的用户类有唯一的ID的话那就更好了。

### 协作对象之间的死锁
上面的两种发生死锁的原因都很明显，下面我们要看的这个例子就不那么明显了，这个死锁发生在两个协同工作的对象之间。

```java
public class CooperativeObjectDeadLock {

    static class Taxi {
        private volatile int x, y;
        private final Dispatcher dispatcher;

        Taxi(int x, int y, Dispatcher dispatcher) {
            this.x = x;
            this.y = y;
            this.dispatcher = dispatcher;
        }

        public synchronized void setLocation(int x, int y) {
            this.x = x;
            this.y = y;
            dispatcher.notifyMove(this);
        }

        public synchronized int[] getLocation() {
            return new int[]{x, y};
        }
    }

    static class Dispatcher {
        private Set<Taxi> taxis = new LinkedHashSet<>();

        public synchronized void notifyMove(Taxi taxi) {
            taxis.add(taxi);
        }

        public synchronized void drawImage() {
            for (Taxi taxi : taxis) {
                int[] loc = taxi.getLocation();
                //dosomething
            }
        }
    }
}
```

上面示例中的两个类如果协同工作可能产生死锁，因为`Taxi`的`setLocation`方法需要获得两个锁，分别是`this`和`dispatcher`，而`dispatcher`的`drawImage`同样需要获得两个锁，是`this`和`taxi`，这就造成了获取锁的顺序不同。

解决办法就是控制锁的粒度，让它更精细，比如：
```java
public void setLocation(int x, int y) {
    synchronized (this){
        this.x = x;
        this.y = y;
    }
    dispatcher.notifyMove(this);
}
```
```java
public void drawImage() {
    for (Taxi taxi : taxis) {
        int[] loc = taxi.getLocation();
        synchronized (this){
            //dosomething
        }
    }
}
```
但是锁的粒度更细了，不在锁中的操作就越多，无法保证同步的操作也就越多。这就需要开发人员在更高的活跃性和更高的安全性之间做权衡。

## 死锁的避免与诊断
### 定时锁
不使用内置锁，使用Java提供的`Lock`类中的`tryLock`，该方法可以提供一个时限，如果指定时限内没有获得锁就放弃并抛出异常。

### 线程转储
略

## 饥饿
除了死锁，饥饿也是引发活跃性问题的一个原因。

饥饿就是线程无法访问它需要访问的资源，比如CPU的时钟周期

## 耗时操作
我们之前有说过在耗时操作里不要获取锁，因为其它的线程必须等待耗时任务结束，在结束之前它们没法获得锁。

## 活锁
活锁和死锁不一样，它不能造成程序假死，但是他会消耗无用的资源。

当一个任务执行出现问题，它可能会被回滚到任务队列尾部等待继续执行，一旦这个任务是个不可能成功的任务，就是每次他都会出现问题，那么他就一直会在任务队列中存在。

解决办法是添加一个任务失败的重试次数。