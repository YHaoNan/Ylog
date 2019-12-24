---
title: （五）任务执行 —— 《Java并发编程实战》
date: 2019-07-30 10:54:30
tags: [java,thread]
categories: Java多线程
---
在高并发应用中，如果把所有任务都放在一个线程中串行执行，那么会造成阻塞，导致我们的应用看起来奇慢无比，而且没有很好的利用CPU的资源。但如果我们给每个任务都创建一个线程，无节制的创建可能还会导致我们的资源被耗尽，然后发生一些不可预测的事情，并且每次请求都创建一个Thread，并在结束后销毁它这种做法在创建和销毁对象上浪费了很多资源，还会给Java的垃圾回收器带来一定的压力。

Java给我们提供了一些工具，可以通过它们来管理这些任务的执行流程，比如串行执行、并行执行，还可以设置一些限制，比如最大线程数等等，通过这些你可以灵活的控制你的任务。

我们先创建一个服务器模型，今天所有的代码都将在这个模型上修改。

```java
public class SyncServer {
    public static void main(String[] args) {
        try {
            ServerSocket server = new ServerSocket(80);
            while (true){
                Socket socket = server.accept();
                handleSocket(socket);
            }
        } catch (IOException e) {
            e.printStackTrace();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    private static void handleSocket(Socket socket) throws IOException, InterruptedException {
        socket.getOutputStream().write("HelloWorld!".getBytes());
        Thread.sleep(500);
        socket.close();
    }
}
```
这是一个单线程的简单的服务器，只要有连接过来，服务器就返回`HelloWorld`，然后等待500毫秒之后关闭，这个500毫秒我们用于模拟耗时操作，先在我们用nc来模拟连接：
```bash
nc localhost 80 & nc localhost 80
```
我们通过`&`让两条命令同步执行，结果发现两个输出还是间隔了一些，就是我们设置的`500`毫秒。

下面我们会用各种类型的任务执行器来对这个示例进行改进。

## Executor
```java
public interface Executor{
    void execute(Runnable command);
}
```
这是`Executor`接口的定义，Java中有很多`Executor`的实现类用于我们对我们的任务进行控制，下面我们使用`Executor`对之前的单线程服务器进行修改。

```java
public class AsyncServerWithFixedThreadPool {
    private static final Executor executor = Executors.newFixedThreadPool(1000);
    public static void main(String[] args) {
        try {
            ServerSocket server = new ServerSocket(80);
            while (true){
                Socket socket = server.accept();
                Runnable task = new Runnable() {
                    @Override
                    public void run() {
                        try {
                            handleSocket(socket);
                        } catch (Exception e) {
                            e.printStackTrace();
                        } 
                    }
                };
                executor.execute(task);
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private static void handleSocket(Socket socket) throws IOException, InterruptedException {
        socket.getOutputStream().write("HelloWorld!".getBytes());
        Thread.sleep(500);
        socket.close();
    }
}

```
上面的代码使用`Executor`实现了多线程服务器，用`Executors.newXXXXXThreadPool`创建了一个线程池，代码中用的是`newFixedThreadPool`，这个线程池可以指定一个最大线程数，如果线程池为空或者没有空闲线程，会为新请求创建一个线程，当然是没有到达最大线程数的情况，如果到达最大线程数则会产生阻塞，直到某个线程执行完毕后转换为空闲状态时执行。

每一个`Executor`实现都有自己对任务控制的逻辑，看看JVM中其他的线程池的控制机制吧：
* newFixedThreadPool  
    如上描述
* newCachedThreadPool  
    一个可缓存的线程池，它没有最大线程数，或者说它的最大线程数是`int`类型最大值，它可以为每个新建的线程设置一个超时时间，如果处于空闲状态的时间超过超时时间，这个线程将被销毁。如果在没有超过timeout的空闲时间内接收到了请求任务，则执行。当需求增加并没有多余线程时，增加线程池中线程数量。它的好处是回收机制保证了不会有很多空闲的线程。
* newSingleThreadExecutor  
    一个单线程的Executor，任务串行执行，但是能保证执行的顺序（FIFO、LIFO、优先级）
* newScheduledThreadPool  
    提供固定大小的线程池，并且以延时方式执行任务

## 停止执行
如果全用`Executor`来控制，那么停止貌似就成了一个问题，把所有线程全部停止可不像想象中那么简单，直接停止线程在Java中是不允许的，因为这可能造成数据丢失等问题，所以多线程中让线程停止其实可以看作是和线程进行协商的过程，如果你了解`interrupt`方法你应该理解了这段话。

`Executor`扩展了`ExecutorService`接口，提供了一些生命周期方法，其中一些用于停止的方法如下：
```java
public interface ExecutorService extends Executor{
    void shutdown();
    List<Runnable> shutdownNow();
    boolean isShutdown();
    boolean isTerminated();
    boolean awaitTermination(long timeout,TimeUnit unit) 
        throws InterruptedException; 
}
```
`ExecutorService`生命周期有三种状态，分别是运行、关闭和已终止，初始创建时是运行状态，`shutdown`是平缓关闭，即不接受新的任务，但是等待已经提交的任务完成，包括已经提交但未开始的任务。`shutdownNow`则粗暴的尝试取消所有运行中的任务，并且队列中没有开始执行的任务不会被启动。

`ExecutorService`属于关闭状态时提交的任务将被拒绝，可能会抛出一个未检查的`RejectedExecutionException`，当所有任务都完成后，`ExecutorService`将进入终止状态。

`awaitTermination`方法等待`ExecutorService`达到终止状态，`isTerminated`返回是否终止。

下一篇会详细介绍关闭和取消任务。

## Callable和Future
`Executor`使用`Runnable`有一些局限性，`Runnable`的接口决定了它不能有返回值，不能向上抛出异常，许多任务都需要返回结果，所以`Runnable`并不能适合所有的场景，`Callable`可以代替它，它是这样的：
```java
public interface Callable<V> {
    V call() throws Exception;
}
```
它有返回结果和异常声明。

`Future`表示一个任务的生命周期，它有几个状态：未开始、运行中
完成，它除了可以获取任务执行的结果外还可以判断是否已经完成、是否已经取消，还可以手动取消任务等。

`Future`的`get`方法用于获取任务结果，如果任务处于完成状态，它会立即返回或抛出异常，分别对应成功或失败（异常或取消），如果任务运行中则会一直阻塞到完成。

任务如果被取消则会抛出`CancellationException`，如果是任务运行时出现了异常则会抛出`ExecutionException`，不管抛出的是什么异常都会封装成这个，可以通过`getCause`获得原始异常。

`ExecutorService`的所有`submit`方法都返回一个`Future`，你需要传入`Runnable`或`Callable`代表任务，可以通过返回的`Future`获取任务结果。

```java
public class ImageDownloader {
    private String[] imgs;
    private ExecutorService service;
    private List<Future<String>> futures;
    public ImageDownloader(String[] imgs){
        this.imgs = imgs;
        this.service = Executors.newFixedThreadPool(imgs.length);
        futures = Collections.synchronizedList(new ArrayList());
    }

    public void startDownload() throws ExecutionException, InterruptedException {
        for (String img:imgs){
            Callable<String> callable = new Callable<String>() {
                @Override
                public String call() throws Exception {
                    return DownloadUtils.download(img);
                }
            };
            futures.add(service.submit(callable));
        }
        for (int i=0;i<futures.size();i++){
            System.out.println(futures.get(i).get());
        }
        System.exit(0);
    }

    public static void main(String[] args) throws ExecutionException, InterruptedException {

        new ImageDownloader(图片列表).startDownload();

    }
}
```
上面代码用`Callable`和`Future`实现了一个图片下载器，`Future`会返回图片的保存路径，当然`DownloadUtils`我没放出来，你需要自己实现。

这个虽然在`ExecutorService`中所有图片是并发下载的，但是我们后面获取下载结果的操作是串行执行的，如果我们下载完图片还需要进行渲染到界面上，那么假设第一个图片下载用了10秒，其他的图片都只用了0.1秒，实际上用户还是要等第一个图片下载完才能看到后面的图片，尽管后面的图片早就下载完毕，这是个问题。我们可以用`CompletionService`解决。

## CompletionService
`ExecutorCompletionService`是`CompletionService`的一个实现类，它使用阻塞队列把每个`submit`的任务执行结果添加进去，然后我们可以直接使用阻塞队列的`take`方法，每当有一个任务完成，我们就可以同时取出结果。这是一个生产者消费者模式，我们就是消费者角色。

```java
public class ImageDownloader {
    private String[] imgs;
    private ExecutorService service;
    public ImageDownloader(String[] imgs){
        this.imgs = imgs;
        this.service = Executors.newFixedThreadPool(imgs.length);
    }

    public void startDownload() throws ExecutionException, InterruptedException {
        ExecutorCompletionService<String> completionService = new ExecutorCompletionService<>(service);
        for (String img:imgs){
            Callable<String> callable = new Callable<String>() {
                @Override
                public String call() throws Exception {
                    return DownloadUtils.download(img);
                }
            };
            completionService.submit(callable);
        }
        for (int i=0;i<imgs.length;i++){
            System.out.println(completionService.take());
        }
        System.exit(0);
    }

    public static void main(String[] args) throws ExecutionException, InterruptedException {

        new ImageDownloader(new String[]{"http://nsimg.cn-bj.ufileos.com/img-1562809043876.jpg"
                ,"http://nsimg.cn-bj.ufileos.com/img-1563418446349.jpg","http://nsimg.cn-bj.ufileos.com/img-1563750980266.jpg",
                "http://nsimg.cn-bj.ufileos.com/img-1563751115703.jpg","http://nsimg.cn-bj.ufileos.com/img-1564023283075.jpg"})
        .startDownload();

    }
}
```

`Future`和`Executor`还有很多方法，比如定时、调用全部什么的，不多介绍了，去看看API吧！