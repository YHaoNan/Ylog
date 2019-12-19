---
title: 优先队列
date: 2019-06-27 13:41:24
tags: [数据结构]
categories: 算法和数据结构
---
优先队列是一种数据结构，不同于队列的先进先出，优先队列是最大的元素先出。这种数据结构应用场景很广泛，在很多有优先级需求的场景可以用到。

使用二叉堆或者树可以很好的实现优先队列并把插入和取出操作的时间复杂度限制在对数级别。

如果使用数组编写二叉堆，那么我们为了方便计算不用第0个空间，用第一个空间表示最大的数，也就是二叉堆的根，用k表示二叉堆中第k个节点，那么`k*2`和`k*2+1`就是它们的子节点，而k的父节点就是`k/2`。

插入操作直接把元素插入到二叉堆最后一个位置，然后再用`swim`方法让它向上游，如果父节点小于它那么它将和父节点换位置，如果父节点大于它那么该处则是它的正确位置。

对于删除元素，因为优先队列删除的永远是最大的，也就是二叉堆的根，那么我们可以把二叉堆的最后一个元素放到首位，之后用`sink`方法让它下沉，如果子节点大于它那么他将和子节点交换位置。
```java
import java.util.Scanner;

public class MaxPQ {
    private int[] pq;
    private int length;

    public MaxPQ(int max){
        pq = new int[max+1];
        length = 1;
    }

    public int max(){
        return pq[1];
    }

    public int delMax(){
        if (isEmpty()) {System.out.println("Empty");return -1;}
        int temp = pq[1];
        pq[1]=pq[--length];
        pq[length] = 0;
        sink(1);
        return temp;
    }

    public void insert(int a){
        if (isFull()){ System.out.println("Full");return; }
        pq[length] = a;
        swim(length++);
    }

    public boolean isFull(){return length==pq.length;}
    public boolean isEmpty(){return length==1;}

    public int size(){return length-1;}

    //下沉
    private void sink(int i){
        while (2*i<=length){
            int j = 2*i;
            if (pq[j]<pq[j+1])j++;
            if (pq[i]<pq[j]){int t=pq[i];pq[i]=pq[j];pq[j]=t;}else break;
            i=j;
        }
    }
    //上浮
    private void swim(int i){
        while (i>1 && pq[i/2]<pq[i]){
            int tmp = pq[i];pq[i]=pq[i/2];pq[i/2]=tmp;
            i/=2;
        }
    }

    public static void main(String[] args) {
        MaxPQ maxPQ = new MaxPQ(100);
        Scanner scanner = new Scanner(System.in);
        while (true){
            switch (scanner.nextLine()){
                case "i":
                    maxPQ.insert(scanner.nextInt());
                    break;
                case "d":
                    System.out.println(maxPQ.delMax());
                    break;
                case "m":
                    System.out.println(maxPQ.max());
                    break;
                case "l":
                    System.out.println(maxPQ.size());
                    break;
            }
        }

    }
}

```