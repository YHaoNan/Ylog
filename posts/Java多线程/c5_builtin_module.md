---
title: （四）基础构建模块 —— 《Java并发编程实战》
date: 2019-07-29 18:19:30
tags: [java,thread]
categories: Java多线程
---
多线程程序往往是很难控制的，因为我们要保证每一个对象状态的安全，本篇介绍一些Java中内置的用来构建并发程序的模块，让我们更快更优雅的实现并发程序。

## `concurrent`并发数据结构
我们先来看看一个具有同步机制的线程安全的数据结构——`Vector`。

先看下JDK官方文档中对它的描述，出自JDK1.8：
> Unlike the new collection implementations, {@code Vector} is synchronized.  If a thread-safe implementation is not needed, it is recommended to use.
>
> 不同于一个新的集合实现，Vector是同步的。如果不需要一个线程安全实现（指的应该是你不需要主动实现线程安全），推荐使用它。

`Vector`继承了`AbstractList`并且使用一个数组来管理元素，它通过为每个方法加上`synchronized`关键字保证每个操作的同步。

也就是说它就是一个同步的List，不过我们使用它时要注意的是，它只能保证单个操作的同步，如果你的代码中有对`Vector`的操作依赖上面几行对`Vector`操作的结果时，需要自行加锁，假设我们对它做迭代操作：
```java
public class TestVector {

    private final static Vector<String> vector = new Vector<>();

    public static void main(String[] args) {
        for (int i=0;i<vector.size();i++){
            System.out.println(vector.get(i));
        }
    }
}
```
如上是一个单线程环境下的代码，不过它在多线程环境下会出现问题，`main`方法中的循环判断条件依据`vector.size()`的返回结果，当`vector.size()`返回后当前线程的锁被释放，这时可能有其他线程在当前线程调用`vector.get(i)`时抢先获得锁，如果它执行了一些删除操作，我们循环中的代码就可能抛出`ArrayIndexOutOfBoundsException`。

解决办法是当进行遍历（或任何其他非原子的操作）时，请自行加锁：
```java
synchronized (vector){
    for (int i=0;i<vector.size();i++){
        System.out.println(vector.get(i));
    }
}
```
如果你用迭代器操作，也是会出现异常的，但是不是`ArrayIndexOutOfBoundsException`，迭代器会在每次迭代时检测整个数组是否被修改，如果被修改则会抛出一个`ConcurrentModifactionException`。

所以，只要你想在多线程环境下遍历`Vector`，不管你用什么办法，你都必须得获得锁，如果`Vector`是一个非常大的集合，而且遍历时对每一个元素的处理很耗时，这时就会造成阻塞，其他线程被都在外面无法进入，对于高并发的程序来说这是不好的。

**而且要注意，当你输出一个vector时，或者把它连接到字符串中时，也会隐式的进行迭代，因为这时调用了它默认的toString方法**

好了，介绍完`Vector`的缺点现在我们开始并发容器的主题。。。

### ConcurrentHashMap
这是一个并发的HashMap，Java的`concurrent`包中有很多为了实现高并发程序设计的数据结构，它就是其中之一，这些数据结构并不是通过为每一个方法加锁来确保安全性，而是使用更细粒度的锁，保证线程安全性的同时提高并发性能。

在`ConcurrentHashMap`中多个线程可以并发访问它，而不是一个线程访问另一个线程就必须等待，因为在高并发情况下容器内的元素在很短时间内会经常的改变，所以一定程度的损失是可以容忍的，在迭代操作中，它能保证遍历当前已有的元素并且尽量把所有的修改操作反映到容器。

实现高并发性能的代价就是一些在高并发应用中不那么重要的方法被减弱了，例如`size()`、`isEmpty()`，因为在高并发环境下容器的内容频繁改变，所以这些方法的返回值在计算时可能已经过期了，不过仔细想想这些方法在高并发应用中我们确实不太常用。

`ConcurrentHashMap`实现了一些常用的复合操作的原子版，比如“没有则添加”、“相等则移除”等。

### CopyOnWriteArrayList
从名字分析，就是在写入时复制的列表，它可以替代同步List。

如果你自己去看这个类的源码，你就会发现它确实是每次add操作都会把整个数组复制一遍，复制成一个新的数组，这也说明了它不适合做一个元素会经常被改动的容器，因为操作的时间复杂度和空间复杂度都很高，也很耗资源。它适合做一些监听器列表等操作，只需要在初始化时把所有监听器设置进去，当事件发生时，遍历所有的监听器发送事件。

### BlockingQueue
阻塞队列是一个多线程环境下特别特别特别特别特别特别好用的FIFO队列数据结构，当`put`时如果队列满就会阻塞，等待何时队列腾出空位再加入，`take`操作时如果队列为空也会阻塞，等待何时队列中有内容时再弹出。

反正我看到它的介绍时一下就想到了多线程中经典的生产者消费者模式，不知道你想到没有，哈哈。这个队列让我们可以方便的实现生产者消费者模式，而不用我们自己实现锁，自己操作`wait`和`notify`。下面是使用阻塞队列实现生产者消费者的一个示例：

```java
public class Factory {
    private static final BlockingQueue<String> productQueue =
            new ArrayBlockingQueue(20);

    public static void main(String[] args) {
        class Worker extends Thread{
            public Worker(String name){super(name);}
            @Override
            public void run() {
                super.run();
                Random random = new Random();
                while (true){
                    try {
                        Thread.sleep(random.nextInt(5)*1000);
                        String product = String.valueOf(System.nanoTime());
                        productQueue.add(product);
                        System.out.println("Worker "+Thread.currentThread().getName()
                                +" put "+product+" to the queue.");
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            }
        }

        class Consumer extends Thread{
            public Consumer(String name){super(name);}
            @Override
            public void run() {
                super.run();
                Random random = new Random();
                while (true){
                    try {
                        Thread.sleep(random.nextInt(5)*1000);
                        String product = productQueue.take();
                        System.out.println("Consumer "+Thread.currentThread().getName()
                                +" take "+product+" from the queue.");
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            }
        }

        for (int i=0;i<5;i++){
            new Worker("W"+i).start();
        }
        for (int i=0;i<5;i++){
            new Consumer("C"+i).start();
        }

    }
}
```
上面代码用了阻塞队列的一个实现`ArrayBlockingQueue`，它还有一些其他实现，比如按照优先级排序的`PriorityBlockingQueue`和同步队列`SynchronousQueue`（不维护空间，直接把产品交给消费者）

### 双端队列与密取
Deque的阻塞实现是`BlockingDeque`，是JDK6中加入的双端队列数据结构，可以从两端加入或弹出。

密取模式则是生产者消费者的升级版，传统的生产者消费者模式因为共用一个队列，经常会出现生产者跟不上消费者或消费者跟不上生产者的情况，然后速度比较快的一方会经常阻塞，效率不高。密取模式则各有各的队列，当消费者没有可消费的东西时就偷偷的去其他消费者那取，当生产者队列满时则悄悄的去其他生产者那里生产。

密取使用双端队列可以更好的减少竞争，密取操作从尾部获取，而正常操作从头部。

### 同步工具类
同步工具类即控制多个线程等待执行状态的工具类，所以BlockingQueue也可以作为一个同步工具类，Java中提供了一些线程的同步工具类。

先介绍一种东西叫闭锁，它像一个大门，当条件未满足时这扇大门关着，一个线程都不能通过，而条件满足时这扇大门敞开，所有线程都可以通过。
#### CountDownLatch
`CountDownLatch`可以作为一个闭锁，它提供一个计数器，在构造时传入，调用`await`方法会产生阻塞，只有当计数器为0时阻塞才会结束，调用`countDown`方法可以使计数器减一。下面是一个用例：


```java

public class TestHarness {
    public long timeTask(int nThreads,Runnable task) throws InterruptedException {
        CountDownLatch startGate = new CountDownLatch(1);
        CountDownLatch endGate = new CountDownLatch(nThreads);
        for (int i=0;i<nThreads;i++){
            new Thread(){
                @Override
                public void run() {
                    super.run();
                    try {
                        startGate.await();
                        task.run();
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }finally {
                        endGate.countDown();
                    }
                }
            }.start();
        }

        long startTime = System.nanoTime();
        startGate.countDown();
        endGate.await();
        long endTime = System.nanoTime();
        return endTime - startTime;
    }

    public static void main(String[] args) throws InterruptedException {
        Runnable task  = new Runnable() {
            @Override
            public void run() {
                try {
                    Thread.sleep(1000);
                    System.out.println(Thread.currentThread().getName());
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        };
        long timePast = new TestHarness().timeTask(1000,task);
        System.out.println("All done , time : "+timePast);
    }
}

```
使用两个`CountDownLatch`第一个用于控制1000个测试线程一起开始，第二个用于确保所有测试线程执行完成。

#### FutureTask
`FutureTask`同样可以作为一个闭锁使用，它表示一个异步任务，当任务正在运行时调用`get`获取结果会产生阻塞，当任务运行完成调用`get`则会立即返回。

下面是我们用`FutureTask`处理耗时任务的一个示例：
```java
public class TestFuture {
    private final FutureTask<String> task = new FutureTask(new Callable<String>(){
        @Override
        public String call() throws Exception {
            //模拟耗时任务
            Thread.sleep(1000);
            return randomString(32);
        }
    });

    private final Thread thread = new Thread(task);
    public void start(){thread.start();}

    public String getResult() throws InterruptedException {
        try {
            return task.get();
        } catch (InterruptedException e) {
            throw e;
        } catch (ExecutionException e) {
            Throwable cause = e.getCause();
            if (cause instanceof RuntimeException){
                throw (RuntimeException) cause;
            }else if (cause instanceof Error){
                throw (Error) cause;
            }else {
                throw new IllegalStateException("Not unchecked exception: ",cause);
            }
        }
    }

    public static void main(String[] args) throws InterruptedException {
        TestFuture future = new TestFuture();
        future.start();
        System.out.println(future.getResult());
    }

    private static String randomString(int len) {
        StringBuffer buffer = new StringBuffer();
        Random random = new Random();
        for (int i=0;i<len;i++)
            buffer.append((char)(random.nextInt(26)+97));
        return buffer.toString();
    }
}
```
主线程中的`future.getResult`会产生阻塞，直到`Callable`的`call`方法执行完毕返回，然后把结果返回给主线程。

`FutureTask`执行任务时的错误处理是个问题，无论发生什么异常，都会被封装到`ExecutionException`中，所以我们要根据不同的异常类型向父级继续抛出。


#### Semaphore信号量
`Semaphore`维护一组固定数量的许可，这个数量可以在构造时传入，每当调用`acquire`时将颁发一个许可，当没有多余的许可可以颁发时`acquire`进入阻塞状态，调用`release`方法可以释放许可，下面用`Semaphore`实现了一个类似阻塞队列的东西，当然他不是阻塞队列。


```java
public class TestSemaphore {
    private final Set<String> set;
    private final Semaphore semaphore;

    public TestSemaphore(int bound){
        set = Collections.synchronizedSet(new HashSet<>());
        semaphore = new Semaphore(bound);
    }

    public boolean add(String e) throws InterruptedException {
        semaphore.acquire();
        boolean wasAdded = false;
        try {
            wasAdded = set.add(e);
            return wasAdded;
        } finally {
            if (!wasAdded)
                semaphore.release();
        }
    }

    public boolean remove(String e){
        boolean wasRemoved = set.remove(e);
        if (wasRemoved){
            semaphore.release();
        }
        return wasRemoved;
    }

    public static void main(String[] args) throws InterruptedException {
        TestSemaphore mySet = new TestSemaphore(5);
        mySet.add("ADFADF");
        mySet.add("ADFavADF");
        mySet.add("AvasdbDFADF");
        mySet.add("ADasdfFADF");
        mySet.add("ADFasbbADF");
        mySet.add("ADFababbADF");
    }
}
```

Set并不知道关于边界的任何内容，边界是我们的信号量在通过许可控制。

### 构建结果缓存
下面将用我们学到的这些模块来构建一个结果缓存。

```java
interface Computable<A,V>{
    V compute(A a) throws InterruptedException;
}
```

```java
public class Cache<A,V> implements Computable<A,V> {
    private Computable<A,V> c;
    private final Map<A,V> cache;
    public Cache(Computable<A,V> computable){
        c = computable;
        cache = new HashMap<>();
    }

    @Override
    public synchronized V compute(A a) throws InterruptedException {
        V v = cache.get(a);
        if (v==null){
            v = c.compute(a);
            cache.put(a,v);
        }
        return v;
    }
}
```
这个示例很简单了，直接用一个`HashMap`来构建缓存，并且为了保证线程安全，我们只能使用锁把`compute`方法锁住。接下来是调用示例：
```java
public class CacheTest {
    public static void main(String[] args) {
        Cache<String,Integer> converter = new Cache<>(new Computable<String, Integer>() {
            @Override
            public Integer compute(String s) throws InterruptedException {
                //模拟耗时
                Thread.sleep(1000);
                return Integer.parseInt(s);
            }
        });

        try {
            Integer result = converter.compute("200");
            System.out.println(result);
            result = converter.compute("200");
            System.out.println(result);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
}
```
在单线程环境中，它是没有问题的，第二次调用`compute`处理200会直接返回`HashMap`中的缓存，但是在多线程中还是有问题，因为计算是一个耗时任务，直接用`synchronized`包裹起来会造成大量的计算请求被阻塞在外面，之前也说过，这种耗时操作能不用`synchronized`尽量还是不用，我们要用更细粒度的操作来保证性能。

我们尝试用`ConcurrentHashMap`来完成这个操作：
```java
public class Cache2<A,V> implements Computable<A,V> {
    private Computable<A,V> c;
    private final Map<A,V> cache;
    public Cache2(Computable<A,V> computable){
        c = computable;
        cache = new ConcurrentHashMap<>();
    }

    @Override
    public V compute(A a) throws InterruptedException {
        V v = cache.get(a);
        if (v==null){
            v = c.compute(a);
            cache.put(a,v);
        }
        return v;
    }
}

```
我们要做的很简单，只需要把`cache`的实现类换成`ConcurrentHashMap`并把方法上的`synchronized`去掉就好了，由于是并发操作所以它能保证这个应用是线程安全的，但是因为去掉了`synchronized`导致`compute`方法的原子性丢失，所以当有两个线程同时计算同一个对象时，线程1执行`compute`时发现缓存里没有，然后去计算，因为计算是耗时任务，不能马上计算完放进缓存，所以线程2也发现了缓存里没有这个，然后又去计算，这会导致多了很多不必要的计算，我们的目的是让线程2发现已经有其他线程计算这个问题了，我只需要等待它计算完毕即可。

`FutureTask`正好能解决我们的问题：
```java
public class Cache3<A,V> implements Computable<A,V> {
    private Computable<A,V> c;
    private final Map<A, Future<V>> cache;
    public Cache3(Computable<A,V> computable){
        c = computable;
        cache = new ConcurrentHashMap<>();
    }

    @Override
    public V compute(A a) throws InterruptedException {
        Future<V> future = cache.get(a);
        if (future == null){
            FutureTask<V> ft = new FutureTask(new Callable<V>(){
                @Override
                public V call() throws Exception {
                    return c.compute(a);
                }
            });
            future = ft;
            cache.put(a,future);
            ft.run();
        }
        try {
            return future.get();
        } catch (ExecutionException e) {
            throw new RuntimeException("ERR");
        }
    }
}
```
上面的代码使用`Future`重写，让我们的Map中保存`Future`，由于它的特性，计算过程中会产生阻塞，其他线程就会等它计算完毕，而计算完毕后就会直接返回。最后我们对错误的处理选择直接简单的抛出了`RuntimeException`，在实际开发中不要这么做。

不过还是有问题，由于`if`中的`put`操作还是非原子的，所以还是会可能出现两个线程同时执行同一个任务的尴尬场面。使用`ConcurrentHashMap`的`putIfAbsent`方法来避免这个问题：
```java
public class Cache4<A,V> implements Computable<A,V> {
    private Computable<A,V> c;
    private final ConcurrentHashMap<A, Future<V>> cache;
    public Cache4(Computable<A,V> computable){
        c = computable;
        cache = new ConcurrentHashMap<>();
    }

    @Override
    public V compute(A a) throws InterruptedException {
        Future<V> future = cache.get(a);
        if (future == null){
            FutureTask<V> ft = new FutureTask(new Callable<V>(){
                @Override
                public V call() throws Exception {
                    return c.compute(a);
                }
            });
            future = cache.putIfAbsent(a,ft);
            if (future == null){future = ft;ft.run();}
        }
        try {
            return future.get();
        } catch (ExecutionException e) {
            throw new RuntimeException("ERR");
        }
    }
}
```
