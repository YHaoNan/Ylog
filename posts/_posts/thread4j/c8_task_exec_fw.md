---
title: （七）任务执行框架 —— 《Java并发编程实战》
date: 2019-08-03 19:13:08
tags: [java,thread]
categories: Java多线程
---

之前我们多少接触过一些`Executor`方面的内容，知道了它是Java提供的任务执行框架，能很方便的管理线程，对任务进行管理控制。

在实际开发中，线程池设置的大小需要仔细斟酌或者经过程序计算动态修改，因为线程池过大或者过小都会有问题。如果线程池过大，会造成大部分系统资源被占用但却是闲置的，等于占着茅坑不拉屎，而线程池如果过小则会导致响应缓慢，大部分系统资源没有被利用上。

如果任务中存在依赖关系，比如一个任务需要另一个任务的结果，那么这时如果它们还在一个线程池中的话，在线程吃过小或任务过多的时候可能会造成死锁，我们称之为“饥饿死锁”。看个例子：
```java
public class StarvationDeadLock {
    private static final ExecutorService executor = Executors.newFixedThreadPool(1);

    public static void main(String[] args) {
        Runnable task = new Runnable() {
            @Override
            public void run() {
                System.out.println("doSomething...");
                Future future = executor.submit(this);
                try {
                    future.get();
                } catch (InterruptedException e) {
                    e.printStackTrace();
                } catch (ExecutionException e) {
                    e.printStackTrace();
                }finally {
                    System.out.println("Over...");
                }
            }
        };

        executor.submit(task);
    }
}
```
这个代码就会造成饥饿死锁，为了缩小问题范围，我们把线程池大小设置成了1，并在任务中再一次提交自己并等待获取值。因为线程池只有一个，被当前任务占用了，第二个任务等着第一个任务结束后运行，而第一个任务又因为用`future.get`获取第二个任务的结果而阻塞，所以这个时候就会造成死锁，你会发现`finally`块中的`Over...`永远不会被打印，并且程序永远不会结束运行。

而且就算不是完全阻塞，只是线程中的某个操作比较耗时，那当线程池过小的时候响应还是会慢。不过对于这种问题我们可以使用超时技术，Java平台中的阻塞方法一般都提供了两个版本，一个是不带超时的，一个是带超时的版本，比如我们可以把上面的代码改进：
```java
public class StarvationDeadLock {
    private static final ExecutorService executor = Executors.newFixedThreadPool(1);

    public static void main(String[] args) {
        Runnable task = new Runnable() {
            @Override
            public void run() {
                System.out.println("doSomething...");
                Future future = executor.submit(this);
                try {
                    future.get(5,TimeUnit.SECONDS);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                } catch (ExecutionException e) {
                    e.printStackTrace();
                } catch (TimeoutException e) {
                    System.out.println("Can not get the result of the sub task...");
                } finally {
                    System.out.println("Over...");
                }
            }
        };

        executor.submit(task);
    }
}
```

我们给`future.get`设置了超时，如果五秒获取不到子任务的结果就放弃，这时会抛出一个`TimeoutException`可以捕获到它做一些操作，由于这个示例都是用的一个`Runnable`实例所以会一直死循环执行任务。

## 设置线程池大小
根据处理器个数，内存大小，数据库连接，还有执行任务的类型确定线程池大小，这个在原书中做了详细介绍，我不写上来了。

除了显式的设置线程池大小，还有隐式原因可能让你间接的限制了线程池中的线程数量，比如你的JDBC连接数量只有10个，任务中需要依赖JDBC连接，那么无论你的线程池多大都是没用的，线程池中永远只能有10个（或更少）任务。
## 配置ThreadPoolExecutor
Thread参数最多的构造器如下：
```java
public ThreadPoolExecutor(int corePoolSize,
                        int maximumPoolSize,
                        long keepAliveTime,
                        TimeUnit unit,
                        BlockingQueue<Runnable> workQueue,
                        ThreadFactory threadFactory,
                        RejectedExecutionHandler handler)
```
它们的作用如下：

名字|作用
:-:|:-:
corePoolSize|核心线程数，就是池子中一般情况下的线程数量，线程会保持在这个数量，少了就创建，多余的会删掉
maximumPoolSize|最大线程数，到了最大线程数后其他的任务会阻塞等待
keepAliveTime|线程存活时间，核心线程数之外的线程的存活时间，到时间没任务使用会删除
unit|时间单位
workQueue|工作队列，线程池中线程的队列
threadFactory|线程工厂，线程池创建新线程时会用到工厂方法
handler|拒绝执行处理器，当线程池满了后的拒绝处理器


`Executors`提供了一些方法创建线程池，比如`newFixedThreadPool`，`newCachedThreadPool`，前者把`corePoolSize`和`maximumPoolSize`全设置成传入的`nThreads`，`keepAliveTime`设置成0，这样不会超时。而后者把`corePoolSize`设置成0，`maximumPoolSize`设置成`Integer.MAX_VALUE`，`keepAliveTime`为60秒，也就是一个线程都不留，只要过了超时时间就不用了，然后来新任务如果线程池满了就会创建新的线程。

其实如果业务不是很复杂，我们直接利用`ThreadPoolExecutor`的构造器就可以设计出符合自己需求的线程池，就像官方提供的那些线程池一样，但是很多情况下我们不能直接用，我们需要更复杂的逻辑，比如在任务进入和执行完毕后写入日志，添加一个未捕获异常，那么我们就要自己继承`ThreadPoolExecutor`了，下面我们就一起深入了解下它的各种用法。

## 任务队列
上面我们看到了我们需要在`ThreadPoolExecutor`的构造器中设置任务队列，默认的`newFixedThreadPool`使用的是一个无界的`LinkedBlockingQueue`，而`newCachedThreadPool`使用的是`SynchronousQueue`，它其实不是一个队列，它通常作为生产者和消费者之间直接传递工作的一个伪队列，它会把任务直接交给消费者，这适用于没有限制的，大的线程池。

有时为了避免资源消耗，会使用有界的阻塞队列，有界的就会有满的时候，我们成为“饱和”，队列如果满了其他任务就要阻塞，等待空闲位置，这时候`RejectedExecutionHandler`就发挥作用了，它是当发生饱和时的拒绝处理器，有如下几个处理器：

名称|策略
:-:|:-:
AbortPolicy|直接拒绝任务并抛出异常
CallerRunsPolicy|用调用者线程执行任务
DiscardPolicy|无声的直接拒绝任务，不会抛出异常
DiscardOldestPolicy|无声的拒绝最老的任务

下面的代码用`CallerRunsPolicy`作为拒绝处理器，并限制阻塞队列最大为10：
```java
public class ExecutorTest {
    private static final ThreadPoolExecutor executor =
            new ThreadPoolExecutor(10,10,
                    0L, TimeUnit.MILLISECONDS,
                    new LinkedBlockingQueue<Runnable>(10));
    public static void main(String[] args) throws ExecutionException, InterruptedException {
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                try {
                    Thread.sleep(2000);
                    System.out.println("TASK RUNNING..."+Thread.currentThread().getName());
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        };


        for (int i=0;i<2000;i++){
            Thread.sleep(100);
            executor.submit(runnable);
        }
    }
}

```
你会时不时的发现控制台会打印一个`TASK RUNNING...main`，这代表任务被退回到了主线程执行。

## 线程工厂
`ThreadFactory`是一个接口，实现这个接口的`newThread`方法返回一个你自己的线程，你可以在这里定制你自己的线程对象，线程池会通过你的工厂来创建新线程，比如：

```java
public class MyThreadFactory implements ThreadFactory {

    static class MyThread extends Thread{
        private static final AtomicInteger created = new AtomicInteger(0);
        public MyThread(Runnable runnable,String name){
            super(runnable,name+created.incrementAndGet());
            setUncaughtExceptionHandler(new UncaughtExceptionHandler() {
                @Override
                public void uncaughtException(Thread t, Throwable e) {
                    System.out.println(getName()+" found a uncaught exception..."+e.getMessage());
                }
            });
        }

        public static int getCreated() {
            return created.get();
        }
    }

    @Override
    public Thread newThread(Runnable r) {
        return new MyThread(r,"MyThread");
    }
}
```
上面的示例为线程设置了自己的名字，并且设置了一个未捕获异常的处理器，所有未捕获的异常都会被这个处理器捕获，你可以在这里写入日志或做一些别的操作。

## 扩展ThreadPoolExecutor
上面的一通操作你其实已经可以很灵活的定制线程池了，但是很多时候我们要在任务开始和结束记录日志，或者做一些别的操作，这时我们可以继承`ThreadPoolExecutor`。

`ThreadPoolExecutor`提供了`beforeExecute`和`afterExecute`方法帮助我们完成这个想法。

下面我们通过一个示例结合上面学到的所有东西来制作一个自动关闭的线程池，这个线程池只接受指定数量的线程，多余的不会接了，并且当任务都执行完，它还会自动关闭。

```java
public class ExecutorTest {
    static class AutoCloseThreadPoolExecutor extends ThreadPoolExecutor{
        private final AtomicInteger doneTaskCount;
        private final AtomicInteger canAddTaskCount;

        public AutoCloseThreadPoolExecutor(int nThreads, long keepAliveTime, TimeUnit unit, ThreadFactory threadFactory) {
            super(nThreads, nThreads, keepAliveTime, unit, new LinkedBlockingQueue(nThreads), threadFactory);
            doneTaskCount = new AtomicInteger(0);
            canAddTaskCount = new AtomicInteger(nThreads);
            setRejectedExecutionHandler(new AbortPolicy());

        }

        @Override
        public Future<?> submit(Runnable task) {
            if (canAddTaskCount.decrementAndGet()==0){
                throw new IllegalStateException("Can not submit a task because the pool is fulled!");
            }
            return super.submit(task);
        }

        @Override
        protected void afterExecute(Runnable r, Throwable t) {
            if (canAddTaskCount.incrementAndGet()==this.getMaximumPoolSize()-1){
                shutdown();
            }
            super.afterExecute(r, t);
        }
    }
    private static final ThreadPoolExecutor executor =
            new AutoCloseThreadPoolExecutor(10,
                    0L, TimeUnit.MILLISECONDS, new MyThreadFactory());
    public static void main(String[] args) throws ExecutionException, InterruptedException {
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                try {
                    Thread.sleep(2000);
                    System.out.println("TASK RUNNING..."+Thread.currentThread().getName());
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        };


        for (int i=0;i<2000;i++){
            Thread.sleep(100);
            executor.submit(runnable);
        }
    }
}
```
我们通过修改它的`submit`方法和`afterExecute`实现了限制提交和关闭功能。

