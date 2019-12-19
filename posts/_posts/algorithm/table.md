---
title: BST、2-3树、红黑树和散列函数实现符号表
date: 2019-06-29 9:56:24
tags: [数据结构]
categories: 算法和数据结构
---

## 有序符号表
有序符号表能保证键的顺序。

最简单的实现符号表的方式，我们可以用一个链表实现顺序存储，然后每次插入和删除时我们遍历链表，找到合适的位置。但是这样做的时间复杂度很高，需要线性复杂度，没法完成大任务，我们希望将插入和删除操作的时间复杂度限制在对数级别，所以应该另辟蹊径。

使用两个数组实现也可以，一个存键，一个存值，而且使用数组这种有序的数据结构，我们可以通过下标访问并不相邻的值，这样就可以使用二分算法来查找键，而不是遍历整个数组，这就实现了对数级别的查找功能。

因为是有序符号表，所以每一个键在数组中都是按顺序排列的，所以我们规定键必须是Comparable的，我们编写了一个rank函数，返回传入的key在符号表中应该摆放的位置，rank使用了二分查找，所以可以保证对数级别的查找速度。
```java
private int rank(Key key){
    int lo = 0 , hi = count-1;
    while (lo<=hi){
        int mid = lo+(hi-lo)/2;
        //比较
        int cmp = key.compareTo(keys[mid]);
        if (cmp == 0)return mid;
        if (cmp > 0) lo = mid+1;
        else hi = mid-1;
    }
    return lo;
}
```

以下是全部代码：
```java
public class BinarySearchST<Key extends Comparable<Key>,Value> implements Iterable<Key> {
    private Key[] keys;
    private Value[] values;
    private int count;

    public BinarySearchST(int size){
        keys = (Key[]) new Comparable[size];
        values = (Value[]) new Object[size];
    }


    public int size(){return count;}

    public Value get(Key key) {
        if (isEmpty())return null;
        int index = rank(key);
        //注意，index位置的key并不一定等于传入的key，所以还要判断是否相等
        if (index < count && keys[index].compareTo(key) == 0)
            return values[index];
        return null;
    }
    public boolean contains(Key key){
        return get(key)!=null;
    }

    public void put(Key key,Value value){
        if (isFull())throw new RuntimeException("Can't put element to a full map.");
        int index = rank(key);
        if (index<count && keys[index].compareTo(key)==0){
            values[index] = value;
            return;
        }
        //把index和后面的所有元素往后移动
        for (int i=count-1;i>=index;i--){
            keys[i+1] = keys[i];
            values[i+1] = values[i];
        }
        keys[index] = key;
        values[index] = value;
        count++;
    }

    public void delete(Key key){
        if (isEmpty())return;
        int index = rank(key);
        if (index<count && keys[index].compareTo(key) == 0){
            //把后面的所有元素往前移动
            for (int i=index+1;i<count;i++){
                keys[i-1] = keys[i];
                values[i-1] = values[i];
            }
            count--;
        }
    }

    private int rank(Key key){
        int lo = 0 , hi = count-1;
        while (lo<=hi){
            int mid = lo+(hi-lo)/2;
            int cmp = key.compareTo(keys[mid]);
            if (cmp == 0)return mid;
            if (cmp > 0) lo = mid+1;
            else hi = mid-1;
        }
        return lo;
    }



    //一些其他方法
    public Iterator<Key> keys(){
        return iterator();
    }

    private class Iter implements Iterator<Key>{

        private int cursor = 0;
        @Override
        public boolean hasNext() {
            return cursor<count;
        }

        @Override
        public Key next() {
            return keys[cursor++];
        }
    }

    @Override
    public Iterator<Key> iterator() {
        return new Iter();
    }

    public boolean isEmpty(){return count==0;}
    public boolean isFull(){return count==keys.length;}

    
}
```

这下就能保证对数级别的查找操作了，但是。。对于插入和删除操作好像。。。还是很慢，因为删除要移动rank返回的索引后面的所有元素，让他们都前进一个位置，插入新key操作则是要向后一个位置，虽然说符号表的插入新key和删除操作远没有查询用的多，但是这样的代码还是无法面对较大规模的任务。

这时候可以考虑用链式存储，链式存储没有数组的限制，用链表来思考，如果在链表中间删除一个元素的话，只需要将他的前后结点相连即可，而数组中删除一个元素则需要移动元素。

但是，我们说的可不是用链表，而是用树。

## 二叉查找树

和链表不同，二叉查找树包含多个指向子节点的链接，一个是左节点，一个是右节点。二叉查找树中一个元素的左节点小于它，右节点大于它，这很像快速排序中的基数。

查找一个元素的时候，只需要从根节点开始对比，如果key相同就返回，如果key比当前节点小，就去左子树中找，如果大则去右子树，这样就保证了对数级别的查找操作。当然，二叉查找树最坏情况下还是需要线性时间，也就是当整棵树排成一个链表时，也就是都只有右节点或都只有左节点，但是这种情况基本不可能出现。

插入操作有两种可能，一是该键存在，找到它并充值它的值即可。一种是该键不存在，这时需要在合适位置添加一个节点。

删除操作比较麻烦，我们需要考虑四种情况：
1. 删除的节点只有左节点，那么需要用它的左节点顶替他的位置
2. 删除的节点只有右节点，那么需要用它的右节点顶替他的位置
3. 删除的节点左右节点都有，需要用右子节点中的最小节点顶替它的位置
4. 左右子节点都没有，直接删

关于第三点，为什么要用右子节点中最小节点顶替呢？假设要删的节点是x，在它的父节点的左侧，那么x的子节点肯定比x的父节点小，否则它就在x的父节点的右面了，反之亦然。而x的右节点中最小的子节点肯定大于x的左子节点，所以，综上所述，这个数可以胜任。
```java
public class BST<Key extends Comparable<Key>,Value>{
    class Node{
        private Node left;
        private Node right;
        private Key key;
        private Value value;
        private int N;

        public Node(Key key, Value value, int n) {
            this.key = key;
            this.value = value;
            N = n;
        }
    }

    private Node root;

    public BST(){ }

    public void put(Key key,Value value){
        root = put(key,value,root);
    }

    private Node put(Key key,Value value,Node node){
        if (node == null) return new Node(key,value,1);
        int cmp = key.compareTo(node.key);
        if (cmp > 0) node.right = put(key,value,node.right);
        else if (cmp < 0) node.left = put(key,value,node.left);
        else {
            node.value = value;
        }
        node.N = size(node.left)+size(node.right)+1;
        return node;
    }


    public Value get(Key key){
        return get(key,root);
    }
    private Value get(Key key,Node node){
        if (node == null) return null;
        int cmp = key.compareTo(node.key);
        if (cmp > 0) return get(key,node.right);
        else if (cmp < 0) return get(key,node.left);
        else {
            return node.value;
        }
    }

    public Node min(){return min(root);}
    private Node min(Node node){
        if (node.left == null) return node.right;
        return min(node.left);
    }

    public void deleteMin(){
        root = deleteMin(root);
    }
    private Node deleteMin(Node node){
        if (node.left == null) return node.right;
        node.left = deleteMin(node.left);
        node.N = size(node.left)+size(node.right)+1;
        return node;
    }

    public void delete(Key key){root = delete(key,root);}

    public Node delete(Key key,Node node){
        if (node == null) return null;
        int cmp = key.compareTo(node.key);
        if (cmp > 0) node.right = delete(key,node.right);
        else if (cmp < 0) node.left = delete(key,node.left);
        else {
            //没有左节点
            if (node.left==null) return node.right;
            //没有右节点
            if (node.right==null) return node.left;
            Node t = node;
            //找到右节点中最小的
            node = min(t.right);
            //删掉这个最小的
            node.right = deleteMin(t.right);
            //顶替之前的元素的位置
            node.left = t.left;
        }
        node.N = size(node.left)+size(node.right)+1;
        return node;
    }
    public int size(){return size(root);}
    private int size(Node node){
        return node!=null?node.N:0;
    }
}
```

#### 习题
1. 为二叉查找树添加一个方法`height()`来计算树的高度。
    ```java
    public int height() {
        return height(root);
    }
    private int height(Node x) {
        int size = 1;
        if (node==null)return 0;
        if (node.right != null) size = Math.max(1+height(node.right),size);
        if (node.left != null) size = Math.max(1+height(node.left),size);
        return size;
    }
    ```

2. 为二叉查找树实现非递归的`get`和`put`方法。
    ```java
    public void put(Key key,Value value){
        Node node = new Node(key,value,1);
        if (root == null){root=node;return;}
        Node cur = root,parent = null;
        while (cur!=null){
            parent = cur;
            int cmp = key.compareTo(cur.key);
            if (cmp > 0) cur=cur.right;
            else if (cmp < 0) cur=cur.left;
            else {cur.value = value;return;}
        }
        int cmp = key.compareTo(parent.key);
        if (cmp < 0) parent.left  = node;
        else parent.right = node;
    }
    ```
    ```java
    public Value get(Key key){
        Node cur = root;
        while (cur!=null){
            int cmp = key.compareTo(cur.key);
            if (cmp > 0) cur=cur.right;
            else if (cmp < 0) cur=cur.left;
            else {return cur.value;}
        }
        return null;
    }
    ```
3. 修改二叉查找树的实现，将最近访问的节点Node保存在一个变量中，这样`get()`或`put()`再次访问同一个键的时候就只需要常数时间了。
    ```java
    private Node lastUse = null;
    ```
    ```java
    //get
    public Value get(Key key){
        if (lastUse!=null && lastUse.key.compareTo(key)==0) return lastUse.value;
        return get(key,root);
    }
    private Value get(Key key,Node node){
        if (node == null) return null;
        int cmp = key.compareTo(node.key);
        if (cmp > 0) return get(key,node.right);
        else if (cmp < 0) return get(key,node.left);
        else {
            lastUse = node;
            return node.value;
        }
    }
    ```
    ```java
    //put
    public void put(Key key,Value value) {
        if (lastUse!=null && lastUse.key.compareTo(key)==0) {lastUse.value = value;return;}
        root = put(key,value,root);
    }

    private Node put(Key key,Value value,Node node){
        if (node == null) {lastUse = new Node(key,value,1);return lastUse;}
        int cmp = key.compareTo(node.key);
        if (cmp > 0) node.right = put(key,value,node.right);
        else if (cmp < 0) node.left = put(key,value,node.left);
        else {
            lastUse = node;
            node.value = value;
        }
        node.N = size(node.left)+size(node.right)+1;
        return node;
    }
    ```

## 平衡查找树
虽然二叉查找树基本能实现对数级别的插入、删除和查找操作，但是最坏情况下它仍是线性级别的，因为根据输入键的不同，二叉查找树无法保证树的平衡，比如：

向一个以字符为键的二叉查找树中依次插入`A,B,C,D,E,F,G...Z`，则这个树是这样的

![img](http://nsimg.cn-bj.ufileos.com/img-1561880516848.png)

前面也提到过，这样的一颗树其实和链表没什么区别了，如果这时再去查找`Z`，是需要线性时间的。所以我们需要研究一种平衡的查找树来完成这个工作。

平衡查找树一般情况下会在插入和删除时为保证树的完美平衡而做一些变换操作，用这些变换操作换来完美平衡是值得的，因为首先这些变换操作不会太慢，其次是查找操作完全比插入（插入新节点）和删除操作用的多。

平衡查找树有很多种实现，2-3查找树、红黑树、AVL算法实现的平衡查找树等等，我们主要研究2-3查找树和红黑树。

### 2-3树
和传统的二叉查找树不同，二叉查找树中的节点都有一个键值和两个子节点（左节点、右节点），我们把它称作2节点，而2-3树中还有一种含有两个键值三个子节点的3节点。

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-anatomy.png)

3节点左侧的键大于右侧的键，从左侧子节点开始，后面的树中每个节点的键都小于左侧键（后统称为左册子节点小于左键），中间子节点的键在左右键中间（大于左键小于右键），右侧的子节点的键大于右侧键。

一颗完美平衡的2-3查找树的所有空连接到根节点的距离是相同的，也就是说命中查找的最大比较次数相同，如上图所示。

#### 查找
查找元素其实和二叉查找树没啥区别，只是多了中间节点，直接看图应该就能懂了。

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-search.png)

#### 插入
插入比较难，因为要考虑树的平衡，所以需要考虑两种情况。如果查找结束于一个2节点，那么就直接把该2节点转换为3节点即可（转换为3节点就可以不增加树的高度）。

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-insert2.png)

但是如果结束于一个3节点则比较麻烦，为保证树的平衡，所以需要往上分解，如何分解？我们先引入4节点的概念，不过不要被这么多花里胡哨的节点搞怕了，你也不用想这么复杂的东西如何用代码实现，我们不会用代码实现2-3树，而是用另一种简单的方式（红黑树）实现它。

4节点包含三个键，左侧最小，其次是中间，然后是右边，它包含四个子节点，分别小于左侧键，在左侧和第二个键之间，在第二和第三个键之间，在第三和第四个键之间。

我们现在假设有一个3节点，我们想往这个3节点下插入一个键，那么会经过如下步骤：

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-insert3a.png)

首先把要插入的键和这个3节点合并，并把键放在合适的位置，上图中S大于A和E，所以放在右边，然后抽出中间的节点，假设这个节点的引用叫t，这个4节点变成了三个普通的2节点，然后把左键放到t的左子节点，把右键放到t的右子节点。这就完成了分解。

但是实际的情况并没这么简单，这个代码也没这么容易实现，我们还要考虑，如果这个被插入的3节点有一个2节点作为父节点，那么就需要如下步骤：

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-insert3b.png)

插入合并为4节点后，抽出的那个中间节点t则会和这上面的2节点合并成一个3节点，并且左键放到这个3节点的中间，右键放到这个3节点的右边。

再考虑一个问题，如果向3节点中插入新节点，而父节点也是一个3节点，那该怎么办？答案是递归向上分解。

就是不断组成4节点，向上分解，直到遇到了一个2节点，这个过程就可以结束了。

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-insert3c.png)

如果到了根节点还是3节点，那么需要继续向上分解，根节点分解后树的高度才增加1。

![img](https://algs4.cs.princeton.edu/33balanced/images/23tree-split.png)

我猜现在你能想到这种做法为啥能够保证树的平衡了，如果没想到就多在脑袋里走几遍这个过程。

（上面的说法是根据图示说的，实际的情况比这个复杂，分解过程可能出现6种情况呢）

树一旦平衡了，它就可以把操作的时间复杂度限制在对数级别了，说一个激动人心的事，加入我们有10亿个节点的一颗2-3树，那么它的高度也就在19-30之间，也就是说最多需要30次访问即可在10亿个数据中找出一个数据。不得不说上面的2-3树很精妙，不过很难实现，因为可能发生的情况太多，处理这些情况所需要的花销有时候可能还不如直接用二叉查找树。

下面，我们使用在二叉树中设置链接颜色的方式来用二叉树实现一个2-3树，这种方式叫红黑树。

### 红黑树
红黑树用节点间链接的颜色来表示2节点或3节点，红色表示链接的是一个3节点，黑色表示链接的是一个2节点。以下是一个红黑树和2-3树的对比：

![img](https://algs4.cs.princeton.edu/33balanced/images/redblack-1-1.png)

可以清晰的看出红黑树是如何表示2-3树的。

一颗红黑树满足如下条件：
1. 只红链接都是左链接
2. 不会出现连续的两个红链接（这里把红链接连接的两个节点看做一个3节点就知道为什么了）
3. 到最底层的每一个节点的每一条路径中经过的黑色节点的数量相同

给出节点的定义：
```java
private static final boolean RED = true;
private static final boolean BLACK = false;

private class Node{
    Key key;
    Value value;
    Node left;
    Node right;
    int N;
    boolean color;


    public Node(Key key, Value value, int n, boolean color) {
        this.key = key;
        this.value = value;
        N = n;
        this.color = color;
    }
}
```

在2-3树中，向上分解需要借助4节点，在红黑树中也有4节点，如果一个节点的左链接和右链接都是红的，你可以把它看做一个4节点，当然和2-3树中4节点会被分解一样，红黑树中的4节点也会被分解成三个2节点（当然，中间节点最后会向上组成3节点）。

而且在红黑树这种数据结构中，把4节点分解成3个2节点非常简单，只需要把链接变成黑的就好了：
```java
public void flipColors(Node h){
    //父节点变为红色，相当于2-3树中的中间节点向上组成3节点
    h.color = RED;
    h.left.color = BLACK;
    h.right.color = BLACK;
}
```
如下图：

![img](https://algs4.cs.princeton.edu/33balanced/images/color-flip.png)

为了方便实现，我们决定，只要插入一个节点，我们就把它以红色的形式插到树底，但是有如下几个问题需要解决：

插入的元素是插入到右节点的（插入的键比较大），如果父节点的左节点不是红的，那么就需要通过旋转把这个向右的红链接转成向左的

![img](https://algs4.cs.princeton.edu/33balanced/images/redblack-left-rotate.png)

插入的元素是插入到右节点的（插入的键比较大），但是父节点的左节点也是红的，那么如上文所说，这时就构成了4节点，需要通过`flipColors`修正颜色。

![img](https://algs4.cs.princeton.edu/33balanced/images/color-flip.png)

插入的元素是插入到左节点的（插入的键比较小），并且父节点是红的，这下就构成了两个连续的左节点，可以通过两次操作解决它（向右旋转变成4节点再分解）。

以下是三种可能性的示意图：

![img](https://img-blog.csdn.net/20171028113830847?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvbGVvbmxpdTE5OTU=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

接下来给出红黑树的实现：
```java
import utils.ConsoleTest;

/**
 * @Description: 红黑树的Java实现
 * @author: LILPIG
 * @date 2019/6/30 16:52
 * God bless my code...
 */
public class RedBlackBST <Key extends Comparable<Key> , Value>{

    private static final boolean RED = true;
    private static final boolean BLACK = false;

    private class Node{
        Key key;
        Value value;
        Node left;
        Node right;
        int N;
        boolean color;


        public Node(Key key, Value value, int n, boolean color) {
            this.key = key;
            this.value = value;
            N = n;
            this.color = color;
        }
    }

    private Node root = null;


    public void put(Key key,Value value){root = put(key,value,root);root.color=BLACK;}
    private Node put(Key key,Value value,Node node){
        if (node == null) {return new Node(key,value,1,RED);}
        int cmp = key.compareTo(node.key);
        if (cmp > 0) node.right = put(key,value,node.right);
        else if (cmp < 0) node.left = put(key,value,node.left);
        else {
            node.value = value;
        }
        if (isRed(node.right) && !isRed(node.left)) node = rotateLeft(node);
        if (isRed(node.left) && isRed(node.left.left)) node = rotateRight(node);
        if (isRed(node.left) && isRed(node.right)) flipColors(node);

        node.N = size(node.left)+size(node.right)+1;
        return node;
    }

    public Value get(Key key){
        return get(key,root);
    }
    private Value get(Key key, Node node){
        if (node == null) return null;
        int cmp = key.compareTo(node.key);
        if (cmp > 0) return get(key,node.right);
        else if (cmp < 0) return get(key,node.left);
        else {
            return node.value;
        }
    }

    //左旋转
    private Node rotateLeft(Node h){
        Node x = h.right;
        h.right = x.left;
        x.left = h;
        x.color = h.color;
        h.color = RED;
        x.N = h.N;
        h.N = 1 + size(h.left) + size(h.right);
        return x;
    }
    //右旋转
    private Node rotateRight(Node h){
        Node x = h.left;
        h.left = x.right;
        x.right = h;
        x.color = h.color;
        h.color = RED;
        x.N = h.N;
        h.N = 1 + size(h.left) + size(h.right);
        return x;
    }
    //改变颜色
    private void flipColors(Node h){
        h.color = RED;
        h.left.color = BLACK;
        h.right.color = BLACK;
    }

    public int size(){return size(root);}
    private int size(Node node){
        return node!=null?node.N:0;
    }
    private boolean isRed(Node node){return node!=null?node.color:false;}

}
```
### 删除
删除操作是红黑树中最麻烦的操作。
#### 删除最小键   
如果删除的键是一个3节点，那么非常容易实现，直接拿掉就好了，如果是2节点就不能直接拿掉了，这样会产生一个空链接，那样的话就不是一个平衡的树了。

我们把所有经过的节点都组合成3节点或暂时的4节点，这样删除操作就不难了。删除后要把4节点分解回去。

删除操作的分解过程会出现三种情况：
1. 当前及节点左子节点不是2节点，不动
2. 当前节点是一个2节点，当前节点的右兄弟节点是一个非2节点，那就需要从父节点中取出一个键放到当前节点，因为递归调用所以父节点肯定是一个非2节点，这个操作不难实现，只需要父节点进行一次左旋转即可完成，但是因为右兄弟节点是一个非2节点，那么从父节点左旋转可能会让兄弟节点跟着跑（讲真的我不知道是不是这个原因，跟着跑好像也没啥事），就先让右兄弟节点进行一次右旋转，旋转到父节点中，再把父节点进行左旋转。
3. 当前节点是一个2节点，右兄弟节点也是一个2节点，这时需要把左兄弟，父节点中的最小键和右兄弟合并成一个4节点就好了，在红黑树中，这不需要做什么旋转操作，只需要一次颜色修正即可完成。

但是由于红黑树的根节点肯定是一个2节点，所以根节点要特殊处理，如果根节点的两个子节点都是2节点的话，我们将它们合并为一个4节点，否则按照上面的两点仍然可以得出正确的结果。

如下图所示：

![img](https://img-blog.csdn.net/20171028114059951?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvbGVvbmxpdTE5OTU=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

有了这些，我们可以实现删除最小键这个操作了。
```java
public void deleteMin(){
    if (!isRed(root.left) && !isRed(root.right))//这个即判断根节点的两个子节点是不是2节点，如果是则把root的颜色改为红色，这个操作会在后面的moveRedLeft方法调用中使之和左右子节点构成4节点
        root.color = RED;
    root = deleteMin(root);
    //因为根节点始终是黑的，递归调用完之后修复
    root.color = BLACK;
}
private Node deleteMin(Node h) {
    //为空代表是最小了
    if (h.left == null) return null;
    //这个操作可以确定左子节点是个2节点（如果h.left.left是红的，则h.left是一个3节点的右键）
    if (!isRed(h.left) && !isRed(h.left.left))
        //从自己身上或者右子节点借键给左子节点
        h= moveRedLeft(h);
    //递归向左调用
    h.left= deleteMin(h.left);
    //最后修复树为一个标准2-3树（过程中可能生成4节点和向右的红链接）
    return balance(h);
}
//向左移动，也就是我们说的借给左子树
private Node moveRedLeft(Node h){
    //修复颜色，这个方法把当前设置成黑的，左右节点设置成红的，这样就生成了一个4节点，如果右子节点是一个2节点，那么该操作会将左子节点，当前节点和右子节点合并为一个4节点
    flipColorsMove(h);
    if (isRed(h.right.left)){//如果h的左侧子节点的最小的兄弟节点是红的，也就是不是2节点（因为经过颜色修复后h是一个4节点，所以h.right.left就是h.left的最小兄弟节点）
        //向右旋转h.right的目的是把最小的兄弟节点排到最前
        h.right = rotateRight(h.right);
        //向左旋转实际上就是把节点接过来
        h = rotateLeft(h);
        //修复颜色
        flipColors(h);
    }
    return h;
}
```

#### 删除最大键
略
#### 删除
略
## 散列表
利用各种树来实现的符号表已经很快了，但是。。。还可以更快吗？？？可以可以，散列表使用散列函数在最优状态的时间复杂度为O(1)，也就是一下就能找到。

但是它也有缺点，就是键无序，而且有时会浪费很多空间。

略。。。

## 参考
- 《算法 第4版 - Robert Sedgewick / Kevin Wayne》
- [3.3 Balanced Search Trees](https://algs4.cs.princeton.edu/33balanced/)