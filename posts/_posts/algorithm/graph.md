---
title: 图简单实现以及深度优先搜索和广度优先搜索
date: 2019-07-07 17:01:19
tags: [图,算法,数据结构]
categories: [数据结构,算法]
---

## 写在最前
博客最近都没怎么更新，天天练车，也没啥时间打理。但是我一直都在学习吼～～

## 图
图是啥？我们不妨想地图。

地图上都有啥？各种地名，把地名连接起来的线路。

图也是这样的，它由各种节点和连接节点的边组成。图可以干很多事，比如社交网络的朋友关系、计算机网络、地图等等。

这篇文章中我们介绍下无向图，即边中没有任何方向信息，只是负责连接节点。


## 树
树就是一个图，满足如下条件的图就是一棵树：
1. 图的边是节点数-1条且不含有环
2. 每条边是连通的
3. 图是连通的(可以由一个节点到所有节点)，但删除任意一条边都会让它不连通
4. 添加任何一条边都会产生一个环
5. 任意一对顶点之间只存在一条简单路径

## 无向图
我们先抽象一个无向图的API：
public class|Graph|<span></span>
-:|:-|:-
<span></span>|Graph(int V)|创建一个含有V个顶点但不含有边的图
int|V()|顶点数
int|E()|边数
void|addEdge(int v,int w)|添加一条边，连接v-w
Iterable\<Integer>|adj(int v)|和v相邻的所有顶点
String|toString()|对象的字符串表示
int|*degree(Graph G,int v)*|计算v的度数
int|*maxDegree(Graph G)*|计算图中顶点的最大度数
double|*avgDegree(Graph G)*|计算图中顶点的平均度数
int|*numberOfSelfLoops(Graph G)*|计算自环个数

下面有几个静态方法，我们会在以后的所有图API中用到。

## 邻接表数组表示图
我们提供一个数组，数组中每一个项目是一个链表，代表与当前下标（每个下标就代表相应的元素）相邻的节点。

每当调用`addEdge`去添加一条边的时候，我们需要把`v`放到`w`的邻接列表里，把`w`放到`v`的邻接列表里。

我们得出了这样一个图模型(我把V和E方法省略了，直接公开属性)：
```java
/**
 * @Description: description...
 * @author: LILPIG
 * @date 19-7-7
 * God bless my code...
 */
package graph;

import java.util.LinkedList;

public class Graph {


    public int E;
    public final int V;
    private final LinkedList<Integer> adj[];


    public Graph(int V){
        this.V = V;
        this.adj = (LinkedList<Integer>[]) new LinkedList[V];
        for (int i=0;i<adj.length;i++){
            this.adj[i] = new LinkedList<>();
        }
    }

    public void addEdge(int v,int w){
        adj[v].add(w);
        adj[w].add(v);
        E++;
    }


    /**
     * 和某个节点相邻的节点
     * @param v
     * @return
     */
    public Iterable<Integer> adj(int v){
        return adj[v];
    }

    /**
     * 计算图的度数
     * @param G
     * @param v
     * @return
     */
    public static int degree(Graph G,int v){
        int degree = 0;
        for (int w:G.adj(v))degree++;
        return degree;
    }

    /**
     * 获取最大的度数
     * @param G
     * @return
     */
    public static int maxDegree(Graph G){
        int max = 0;
        for (int i=0;i<G.V;i++){
            max = Math.max(max,degree(G,i));
        }
        return max;
    }

    /**
     * 获取平均度数
     * @param G
     * @return
     */
    public static double avgDegreeSlow(Graph G){
        float sDeg = 0;
        for (int i=0;i<G.V;i++)
            sDeg += degree(G,i);
        return sDeg / G.V;
    }

    /**
     * 获取平均度数
     * @param G
     * @return
     */
    public static double avgDegree(Graph G){
        return 2.0 * G.E / G.V;
    }

    /**
     * 获取自环个数
     * @param G
     * @return
     */
    public static int numberOfSelfLoops(Graph G){
        int c = 0;
        for (int i=0;i<G.V;i++){
            for (int num:G.adj(i))
                if (num == i) c++;
        }
        return c/2;
    }

    public String toString() {
        StringBuilder s = new StringBuilder();
        s.append(V + " vertices, " + E + " edges " + "\n");
        for (int v = 0; v < V; v++) {
            s.append(v + ": ");
            for (int w : adj[v]) {
                s.append(w + " ");
            }
            s.append("\n");
        }
        return s.toString();
    }
}

```

## 搜索算法
我们需要一个搜索算法搜索两个节点是否连通，和某个节点连通的节点总数，两个节点之间的一条路径。

所以我们抽象如下API：

public interface|Search|<span></span>
-:|:-|:-
int|count|与某个节点连通的节点数
boolean|hasPathTo(int v)|是否能连通到这个节点
Iterable\<Integer>|pathTo(int v)|到v的路径

## 深度优先搜索
```java
/**
 * @Description: description...
 * @author: LILPIG
 * @date 19-7-7
 * God bless my code...
 */
package graph;


/**
 * 深度优先搜索
 */
public class DepthFirstSearch implements Search{


    private boolean marked[];
    private int edgeTo[];
    private final int s;
    private int count;
    public DepthFirstSearch(Graph G, int s){
        marked = new boolean[G.V];
        edgeTo = new int[G.V];
        this.s = s;
        dfs(G,s);
    }

    private void dfs(Graph G,int v){
        marked[v] = true;
        count++;
        for (int i:G.adj(v)){
            if (!marked[i]) {edgeTo[i]=v;dfs(G,i);}
        }
    }
    @Override
    public int count(){
        return count;
    }

    @Override
    public boolean hasPathTo(int v) {
        return marked[v];
    }

    @Override
    public Iterable<Integer> pathTo(int v) {
        if (!hasPathTo(v))return null;
        Stack<Integer> stack = new Stack<>();
        for (int x=v;x!=s;x=edgeTo[x])
            stack.push(x);
        stack.push(s);
        return stack;
    }
}

```
## 广度优先搜索
```java
/**
 * @Description: description...
 * @author: LILPIG
 * @date 19-7-7
 * God bless my code...
 */
package graph;


public class BreadthFirstSearch implements Search{
    private boolean[] marked;
    private int[] edgeTo;
    private final int s;

    public BreadthFirstSearch(Graph G,int s) {
        this.s = s;
        marked = new boolean[G.V];
        edgeTo = new int[G.V];
        bfs(G,s);
    }

    private void bfs(Graph G,int s) {
        Queue<Integer> queue=new LinkedBlockingQueue<>();
        marked[s] = true;
        queue.add(s);
        while (!queue.isEmpty()){
            int c = queue.poll();
            for (int x:G.adj(c)){
                if (!marked[x]){edgeTo[x]=c;bfs(G,x);queue.add(x);}
            }
        }

    }

    @Override
    public int count() {
        return 0;
    }

    @Override
    public boolean hasPathTo(int v) {
        return marked[v];
    }

    @Override
    public Iterable<Integer> pathTo(int v) {
        if (!hasPathTo(v))return null;
        Stack<Integer> stack = new Stack<>();
        for (int x=v;x!=s;x=edgeTo[x])
            stack.push(x);
        stack.push(s);
        return stack;
    }
}

```

内容来自《算法 第四版》