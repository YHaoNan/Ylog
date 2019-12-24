---
title: （三）对象的组合 —— 《Java并发编程实战》
date: 2019-07-28 11:19:30
tags: [java,thread]
categories: Java多线程
---

设计线程安全的代码要考虑的因素太多了，我们要考虑到每一个对象的线程安全性，并且要考虑到它们组合起来时的线程安全性，好像每一个操作都要做很多的思考并且要为之付出很多保证安全的代码。这一节就是一些利用对象组合的多线程编程技巧，在保证代码优雅的前提下还能兼顾到安全性。

## 封闭机制
我们可以通过把一个线程不安全的类封装到一个线程安全的类中并提供同步机制来保证线程的安全，比如这样做：
```java
@ThreadSafe
public class SafePersonSet {
    @GuardedBy("this")
    private final HashSet<String> personSet = new HashSet<>();
    public SafePersonSet(Collection<String> persons){
        personSet.addAll(persons);
    }

    public synchronized void add(String person){this.personSet.add(person);}

    public synchronized boolean contains(String person){return this.personSet.contains(person);}

}
```

`HashSet`并不是一个线程安全的类，在多线程中直接使用它会造成线程安全问题，我们根据我们的需求把`HashSet`封装在一个其他类里，并通过`synchronized`保证线程安全。


下面来看一个示例，检测车辆位置的类
```java
@ThreadSafe
public class MonitorVehicleTracker {
    
    @NotThreadSafe
    static class Point{
        public int x,y;
        public Point(Point point){
            this.x = point.x;
            this.y = point.y;
        }
        public Point(){ }
        
    }
    
    private final Map<String,Point> locations;
    
    public MonitorVehicleTracker(Map<String,Point> locations){
        this.locations = deepCopy(locations);
    }
    
    public synchronized Map<String,Point> getLocations(){
        return deepCopy(locations);
    }
    
    public synchronized Point getLocation(String id){
        Point point = locations.get(id);
        return point == null ? null : new Point(point);
    }
    
    public synchronized void setLocation(String id,int x, int y){
        Point point = locations.get(id);
        if (point == null)
            throw new IllegalArgumentException("No such ID: "+id);
        point.x = x;
        point.y = y;
    }
    
    
    private Map<String,Point> deepCopy(Map<String,Point> locations){
        Map<String,Point> result = new HashMap<>();
        for (String key:locations.keySet()){
            Point point = new Point(locations.get(key));
            result.put(key,point);
        }
        return Collections.unmodifiableMap(result);
    }

}
```

整个类维护了一个`Map<String,Point>`的表，`Point`不是线程安全的，所以当返回`Point`的时候我们采用了类似于快照的方式，就是生成一个新的，和原来的Point无关的对象来保证某个线程修改返回的`Point`对其他线程无法造成影响。而且，Map是非线程安全的，但是我们使用`synchronized`将所有操作保护了起来，这样一来`MonitorVehicleTracker`就是一个线程安全的类了。

不过在大量数据需要处理并且并发很高的时候，这个代码就会出问题，它会很慢，因为每次它都把整个Map复制一遍，并且重新创建每个`Point`对象，这可能会造成性能问题。

嘶。。Java中提供了那么多线程安全的集合类，我们为什么不直接用？假如我们用线程安全的`Map`并且把Point类也写成线程安全的，我们就不用每次都建立快照了。我们干脆把功能委托给`ConcurrentHashMap`来完成，这是一个Java提供的线程安全的HashMap，不过我们的Point类也得改，因为它是可变的，`ConcurrentHashMap`能保证线程安全，可没法保证拿出来的`Person`不被修改啊，所以我们要把`Point`改成不可变的。

```java
@ThreadSafe
public class MonitorVehicleTracker2 {

    @Immutable
    static
    class Point{
        public final int x,y;
        public Point(int x,int y){
            this.x = x;
            this.y = y;
        }
    }

    private final ConcurrentHashMap<String,Point> locations;
    private final Map<String,Point> unmodifiableMap;
    public MonitorVehicleTracker2(Map<String,Point> locations){
        this.locations = new ConcurrentHashMap<>(locations);
        this.unmodifiableMap = Collections.unmodifiableMap(locations);
    }



    public Map<String,Point> getLocations(){
        return unmodifiableMap;
    }

    public Point getLocation(String id){
        return locations.get(id);
    }

    public void setLocation(String id,int x, int y){
        if (locations.replace(id,new Point(x,y))==null)
            throw new RuntimeException("No such id : "+id);
    }

}
```
我们甚至没在这里使用任何显式的同步代码，因为我们把同步委托给其他类了。而且这还有个好处就是当有一个线程通过`getLocations`获取到所有车位置的集合后，其他线程重新调用`setLocation`更新`locations`时，对之前的线程也是可见的。


我们还可以把`Point`设计成线程安全的并且允许修改`x,y`：
```java
@ThreadSafe
public class MonitorVehicleTracker3{

    @ThreadSafe
    static
    class Point{
        private int x,y;
        public Point(int x,int y){
            this.x = x;
            this.y = y;
        }
        public Point(Point point){
            this(point.get());
        }
        public Point(int[] point){
            this.x = point[0];
            this.y = point[1];
        }

        public synchronized int[] get(){return new int[]{x,y};}
        public synchronized void set(int x,int y){this.x=x;this.y=y;}
    }

    private final ConcurrentHashMap<String,Point> locations;
    private final Map<String,Point> unmodifiableMap;
    public MonitorVehicleTracker3(Map<String,Point> locations){
        this.locations = new ConcurrentHashMap<>(locations);
        this.unmodifiableMap = Collections.unmodifiableMap(locations);
    }



    public Map<String,Point> getLocations(){
        return unmodifiableMap;
    }

    public Point getLocation(String id){
        return locations.get(id);
    }

    public void setLocation(String id,int x, int y){
        if (!locations.containsKey(id))
            throw new IllegalArgumentException("No such ID :" +id);
        locations.get(id).set(x,y);
    }

}
```
它依然是线程安全的，而且客户端可以通过`Point`对象修改x，y了。这是书上写的代码，我对`setLocation`是否是安全的表示怀疑。。。因为它调用了两次`locations`的方法，还调用了`Point`的set方法，我觉得如果不用`synchronized`不能保证原子性啊，如果有了解的请告诉我。

## 扩展线程安全方法
就像上面我对`setLocation`方法的疑问一样，比如我们需要给`Vector`类添加一个`putIfAbsent`方法，该方法检测如果元素存在就不添加，不存在才添加，并且保证这个操作是线程安全的，我们可以通过继承的方法来写。

```java
public class BetterVector<E> extends Vector<E> {
    public synchronized boolean putIfAbsent(E e){
        boolean absent = !contains(e);
        if (absent)
            add(e);
        return absent;
    }
}
```

这不能满足所有的场景，比如我们在使用`Collections.synchronizedList`方法返回一个同步的List时就无法通过子类的形式扩展，这时我们可以通过一个助手类的方式扩展：
```java

public class ConcurrentListHelper<E> {
    public List<E> list = Collections.synchronizedList(new ArrayList<>());
    ...
    public boolean putIfAbsent(E e){
        synchronized(list){
            boolean absent = !list.contains(e);
            if (absent)
                list.add(e);
            return absent;
        }
    }
}
```

注意这里我们是针对`list`加锁，如果对`this`加锁那么对`list`的其他操作与`putIfAbsent`不会同步。

还可以通过委托扩展：

```java
public class ImproveList<E> implements List<E> {
    private final List<E> list = new ArrayList<>();

    public synchronized boolean putIfAbsent(E e){
        boolean absent = !contains(e);
        if (absent)
            add(e);
        return absent;
    }


    @Override
    public int size() {
        return list.size();
    }

    //..省略其他接口规定的方法
}
```
