---
title: 归并排序、快速排序
date: 2019-06-26 19:19:18
tags: [排序,算法]
categories: 算法和数据结构
---

会写点业务逻辑，但是算法太渣了，之前看过《算法 第四版》，因为高考，一直没弄这些东西，最近刷leetcode发现我连个归并排序自己都写不出来了，打算重新看看算法这些东西吧，毕竟真的不想只当一个copy、paste的平庸之辈～～

## 归并排序
归并排序的思想就是分治，把问题分成一堆小问题，依次解决，最后整个大问题就解决了。


如果让你把两个有序的数组`arr1`和`arr2`合并为一个有序的数组`arr3`，你应该会这么写：
```java
int[] arr3 = new int[arr1.length+arr2.length];
int i=0,j=0,k=0;

while (i<arr1.length&&j<arr2.length)
    arr3[k++] = arr1[i]<arr2[j]?arr1[i++]:arr2[j++];

//处理剩余
while (i<arr1.length)
    arr3[k++]=arr1[i++];
while (secondStart<=hi)
    arr3[k++]=arr2[j++];
```

归并排序利用分治，每次把数组分两等分，再把每一份分两等分，一直这样分下去，直到最后每份只剩下一个元素，一个元素我们始终可以认为是有序的，然后再把这些小的有序数组依次排成一个大的有序数组，这种操作非常适合递归完成。

但是每次分成两等分需要耗费许多不必要的空间，所以我们采用原地归并，利用三个变量`lo`、`mid`、`hi`表示最低下标、中间和最大下标，然后把lo-mid看成一份，mid-hi看成一份。
```java
import java.util.Arrays;

public class MergeSort {
    
    public void sort(int[] array){
        divide(array,0,array.length-1);
    }

    //分割数组为两等份(实际不一定是等份)，然后递归再分，lo==hi，也就是分成一个元素
    private void divide(int[] array,int lo,int hi){
        if (lo == hi)return;
        int mid = (lo+hi)/2;
        //继续分左边
        divide(array,lo,mid);
        //继续分右边
        divide(array,mid+1,hi);
        merge(array,lo,mid,hi);
    }

    //合并数组
    private void merge(int[] array,int lo,int mid,int hi){
        int firstStart = lo;
        int secondStart = mid+1;
        int tempI=0;
        int[] temp = new int[hi-lo+1];
        while (firstStart<=mid&&secondStart<=hi){
            temp[tempI++] = array[firstStart]<array[secondStart]?array[firstStart++]:array[secondStart++];
        }
        while (firstStart<=mid){
            temp[tempI++]=array[firstStart++];
        }
        while (secondStart<=hi){
            temp[tempI++]=array[secondStart++];
        }

        tempI=0;
        while (lo<=hi){
            array[lo++]=temp[tempI++];
        }
    }


    public static void main(String[] args) {
        int[] array = {54,41,563,2,1,47,21,341,2,6,1,5,3,2,1,4,5,5,5,5,6,3,1,47,5,1};
        new MergeSort().sort(array);
        System.out.println(Arrays.toString(array));
    }
}
```

## 快速排序
快速排序和归并排序一样，也用到了分治思想，它的时间复杂度和归并排序一样，但是归并排序需要一个暂存数组，在空间复杂度上稍高，而快速排序解决了这个问题，并且它的原理和代码都是很简洁的。

快速排序的思想是在数组中找一个基数，然后遍历整个数组，把比他大的排在它右边，比它小的排在它左边，这样就保证了这个元素的位置是正确的。再以这个正确的位置分割，将左右两边的子数组分别再做此操作，最后所有的元素都会找到自己该有的位置。

```java
import java.util.Arrays;

public class QuickSort {
    public void sort(int[] arr){
        divide(arr,0,arr.length-1);
    }

    private void divide(int[] arr,int lo,int hi){
        if (lo>=hi)return;
        //找中心位置
        int j = partition(arr,lo,hi);
        //再次对左右两边重复操作
        divide(arr,lo,j-1);
        divide(arr,j+1,hi);
    }

    
    private int partition(int[] arr,int lo,int hi){
        //把第一个元素当做基数
        int base=arr[lo];
        //开始做排队操作
        while (lo < hi){
            //如果右侧的元素比基数大或等于基数，hi--，如果比基数小，就放在左边(放到arr[lo])，因为最开始的元素被作为基数，所以从hi这边开始遍历
            while (lo<hi && arr[hi]>=base) hi--; 
            arr[lo]=arr[hi];
            //如果左侧的元素比基数小或等于基数，lo++，如果比基数大，则放在右边(放到arr[hi])
            while (lo<hi && arr[lo]<=base) lo++;
            arr[hi]=arr[lo];
        }
        //把基数归位
        arr[lo] = base;
        return lo;
    }

    public static void main(String[] args) {
        int[] array = {54,41,563,2,1,47,21,341,2,6,1,5,3,2,1,4,5,5,5,5,6,3,1,47,5,1};
        new QuickSort().sort(array);
        System.out.println(Arrays.toString(array));
    }
}
```