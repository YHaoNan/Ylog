---
title: 儿的生日，爹的苦日
date: 2019-10-28 12:45:27
categories: Others
---

<script>
    window.onload = function(){
        
        var playerContainer = document.getElementsByClassName("player-container");
        var audios = document.getElementsByClassName("song");
        console.log(audios);
        console.log(playerContainer);
        for(let ii=0;ii<audios.length;ii++){
            audios[ii].load();
            playerContainer[ii].onmouseenter = function(){
                audios[ii].play();
            }
            playerContainer[ii].onclick = function(){
                audios[ii].play();
            }

            playerContainer[ii].onmouseleave = function(){
                audios[ii].pause();
            }
        }
    }
</script>

<style>
  @keyframes hover-anim{
    30%{
        transform: scale(1.05);
        filter: grayscale(70%)
    }
    60%{
        transform: scale(1.08);
        filter: grayscale(40%)
    }
    100%{
        transform: scale(1.1);
        filter: grayscale(20%)
    }
}
.song{
    display: none;
}
.player-body{
    width: 260px;
    height: 70px;
    background-size: 100% auto;
    background-position-y: center;
    border-radius: 5px;
    box-shadow: 0px 30px 60px rgba(0,0,0,0.12),
                0 20px 20px rgba(95,23,101,0.2);
}
.player-container{
    margin:40px;
    width: 260px;
    filter: grayscale(80%);
}
.player-container:hover{
    animation: hover-anim .3s ease-in forwards
}

.player-album{
    box-sizing:unset;
    width: 15px;
    height: 15px;
    border: 30px solid #16244d;
    border-radius: 100%;
    display: inline-block;
    margin-top: -30px;
    margin-left: 10px;
}

.player-name-and-bar{
    width: 220px;
    display:inline-block;
    height: 45px;
    box-shadow: 0px 30px 60px rgba(0,0,0,0.12),
                0 20px 20px rgba(95,23,101,0.2);
    line-height: 40px;
    padding-left: 80px;
    background: #eee;
    margin-left: 20px;
    margin-right: 20px;
    border-top-left-radius: 5px;
    border-top-right-radius: 5px;
}

.player-songname{
    font-weight: bold;
    line-height: 0px;
    font-family: sans-serif;
    font-size: 0.8em;
}  
.player-description{
    font-family: sans-serif;
    font-family: sans-serif;
    font-size: 0.8em;
    line-height: 1;
}
.imgdes{
    font-size: 0.8em;
    text-align: center;
    text-decoration: underline;
}

    
</style>

## 序

![preface](https://s2.ax1x.com/2019/11/03/KOb1Yj.jpg)

<div class="imgdes">图为“儿砸儿砸”在小米录音机里的声波频谱</div>

## 昨天和今天

时常还能想到初中的时候，一个农村娃和一个黑不出溜的城里娃就莫名其妙的被分到了一组。那时候农村娃才一米四几，和同学走在一起就像一个碗和一根筷子在做游戏。

某一天，一米四的农村娃和城里娃传纸条玩

某一天，农村娃长到了一米六

某一天，城里娃运动会跑了个第一

某一天，农村娃体育测试不合格

某一天，考试

某一天，农村娃求城里娃帮他不及格的试卷签字

某一天，农村娃长到了一米七

某一天，城里娃运动会摔了一跤，所有人都跑过去看她，农村娃没有过去

某一天，农村娃长到了一米七五

某一天，中考

农村娃考了一个理所应当的成绩，并且理所应当的被所有高中抛弃

城里娃跑到了很远很远的地方读书

他最后见到城里娃那天，城里娃在他的毕业证上写了几个字

然后农村娃就再也没见过城里娃

三年多了...

他有三年没见到城里娃了...他...

有些话，她在他旁边时候他没说，她在他毕业证上签字的时候他没说，她走的那天他没说，最后就...

今天他要大胆说出来，对，就是今天。今天不说，难道还要等到十一月十二号？还是一年以后？或者两年？二十年？三十年？

他想说——老子能从一米四长到一米八怎么你就不他娘的长呢？？？？？？

（手动道歉，求不杀👀）


## 谢谢你突然出现

我也记不清是哪天了，你突然跟我说了句没头没尾的话。

你说：“我数学老师的儿子长得好像你。”

然后，两个有两年没联系的人就开始扯起了皮。从那以后的每天，都会扯那么几句。用你的话说叫撸嘴。

一切看起来都很好，只是你不知道的是，那段时间，是我这18年来最难以度过的一段。

颓废、无助、天昏地暗。。。甚至感觉自己活着无意义。

你一定想不到一个大老爷们儿每天可以哭十遍，他甚至看熊出没都会哭。

我不能让我爸妈知道，更不能和别人说。毕竟。。谁愿意去接你从心里扒拉出来的苦水呢。

那段正好临近高考，如果不是亲自尝试的话，我甚至永远都不知道“化悲愤为动力”这句话只是个天大的屁，那些好端端的人放出来转身就走了，他们品不出也不会去品这话有多臭。我根本无法专注学习。

我希望有个人能陪我说说话，同时又不希望那样，我怕我一不小心就把这份不快带给任何人。

万幸，你出现了，而且我忍住了，没把它带给你一点。

有一天我跟你说，只要每天和你骂几句就很高兴了，那是真话，那不是犯贱。。。如果把你扔进深海冰窟里，你就能感受到那些话里蕴含的温度，是它们把我从刺骨的冰窟里打捞出来。谢谢宁。


## 想和你分享快乐
想把所有开心的事告诉喜欢的人，至于讨厌的人，最好让他们下课回寝室时候喝冰露矿泉水不小心被淹死。。。

所以啊，想和你互相交换开心事，这样就是双倍快乐。比蜜雪冰城的大快乐还快乐。

每天都快快乐乐的，那样多好。


* 11月1号你傻笑了14次
* 11月2号你傻笑了13次
* 11月3号你傻笑了7次
* 11月4号你傻笑了10次
* 11月5号你傻笑了2次 宁怎么今天智商在线了？
* 11月6号你傻笑了9次
* 11月7号你傻笑了28次 可能是全年智商最低的一天并扬言要包养我
* 11月7号你傻笑了27次
* 11月9号你傻笑了40次 ？？？ 
* 11月10号你傻笑了45次 WTF

虽然很傻，但是爸爸不嫌弃的

## 庆幸我和你说话不用动脑

总有几个人，我和他说话可以不经思考，因为我知道，我就是说什么他都不会介意的。

就算他介意，我大不了把他弄死呗。

## 希望那答案唯一

我印象最大的一件事就是封面上那件事，我问你，将来想过什么样的生活，你说快乐安逸就好。

是啊，有啥比这更好啊。

生活温饱，无忧无虑，工作内容恰好是爱好，转头就是你最爱的人，他们也同样安逸，同样快乐......

我想应该每个人都想过这种生活吧，但是，慢慢地，有人的方向就会偏离。他们会活成自己不喜欢的样子，他们会活不出自己的人生。

“这不现实，你太幼稚了，这可不是乌托邦”成了他们没能力一直追求自己喜欢的东西的借口，或者是他们根本不知道自己啥样才会快乐。

现在如果再问你，答案会变吗？

如果你还是我认识的你，那我想，你永远都不会变的。

不变，嗯，不变......


## 明天

明天，不敢想。

反正，早晚我们会忙起来。你可能会忙着学业，忙着工作；可能会忙着照顾家庭，忙着照顾上天赐给你的小天使。我可能也在千里之外忙着同样的事。我可能每年只有今天才会喊你一声儿砸，你也可能只有四月份的某个早晨会喊我声父亲。我们不会再有时间每天一起吹牛逼，我们会忙到不知不觉就灭火了，那艘船也可以还给你的小宝贝儿了。我们可能一辈子都不会再见面，但是不管怎样，我希望你在不开心的时候永远能第一个来找我。

那个女人可以夺走“小头儿子”旁边的那艘船，但她永远拿不走我肚子里那艘～拿不走的

## 生日快乐
嗯哼哼，生日快乐，相信我吧，就是天塌下来我也不会忘的。

你一直走，我一定在你身后(好几把土)

祝你心想事成，万事如意...祝你早日成为彭于晏的女人...祝你早日拥有妹妹一样的妈妈...

PS：我不想我的儿子十一月十一号过生日，我想你天天过生日，当然是不用我花钱的生日


请宁包容...
啊哈哈哈哈

## 关于礼物
其实，，，怎么说呢。。。

那相框本就是给你买的。。。

我。。。我这么嘴硬怎么好意思说那是给你的呢。。。

其实早在我过生日的时候我就在构想怎么给你过生日了。。。

我想的特别好。。。

但是。。。

我有些懒惰。。。

以至于这份礼物相对于构想之中的不断缩水。。。

ver1版本是，我买个单片机，一个屏幕，加一个半透明的镜子，做一个智能的镜子。它除了是一块镜子外，还可以显示一些其他信息，比如天气、备忘录、TODOLIST等。并且你可以每天像问你手机里的智能助手一样问他话。

但是，别的我都可以搞，就是这个镜子。。。它的木制边框我该怎么做。。。。并且成本有些高我觉得宁有点不配(不不不对不起)

ver2版本和你收到的礼物差不多，只是构想中效果比现在的要好很多。

别说了，一切都是因为我懒😂

以下是和ver2设想中不一致的地方：

- 少了一份手写的祝福(也就是本文)
- 少了一份我自己读的祝福(还是本文)
- 少了一份精心设计过的页面(所以你看到的是本页)
- 少了几份歌词图片
- 10寸画框缩水成6寸相框
- 多了一份妹妹的照片

## 惊吓
> 不知道是惊喜还是惊吓，最近感冒了，而且在寝室瓦也不敢大声逼逼，嘤嘤嘤。。。
>
> 微信/QQ系浏览器可能没得声音，请复制下面网址并使用其他浏览器，推荐电脑端使用Chrome浏览器
>
> http://lilpig.site/post/other-sonbir/




<audio controls="true" src="http://nsimg.cn-bj.ufileos.com/lt.mp3" class="song" id="song-lt"></audio>


<audio controls="true" src="http://nsimg.cn-bj.ufileos.com/xsdzt.mp3" class="song" id="song-q"></audio>


<audio controls="true" src="http://nsimg.cn-bj.ufileos.com/0e0b_0f08_015d_9a5e1cddac8bb492ec8a6717f674725f.mp3" id="song-wxnl" class="song"></audio>



### 骆驼

<div class="player-container">
   <div class="player-name-and-bar">
       <p class="player-songname">骆驼</p>
       <p class="player-description">全世界都是沙漠</p>
   </div>
   <div class="player-body" style="background-image:url(https://s2.ax1x.com/2019/11/03/KOYM38.jpg)">
        <div class="player-album" style="border-color:#4c3851"></div>
   </div>
</div>


> 其实，这世界就是沙漠，我们是骆驼  
> 沙漠大的看不着边，太大了就显得空旷，显得寂寞  
> 骆驼们孤身在沙漠里来来回回，漫无目的，相遇后扭头别离  
> 不管怎样，总会有一只骆驼穿越千里，最后，目的是你


### 消失的昨天

<div class="player-container">
   <div class="player-name-and-bar">
       <p class="player-songname">消失的昨天</p>
       <p class="player-description">挡住日月，才能...</p>
   </div>
   <div class="player-body" style="background-image:url(https://s2.ax1x.com/2019/11/10/Mn2lvV.md.jpg)">
        <div class="player-album" style="border-color:#1e1b13"></div>
   </div>
</div>

> 如果昨天的美好突然消失，也要安安稳稳的过好今天，谁能保证明天它不会再次出现呢？


### 电鳗想你了

<div class="player-container">
   <div class="player-name-and-bar">
       <p class="player-songname">我想你了</p>
       <p class="player-description">最近老做同一个梦</p>
   </div>
   <div class="player-body" style="background-image:url(https://s2.ax1x.com/2019/11/10/Mn2SNd.md.jpg)">
        <div class="player-album"></div>
   </div>
</div>

> 咱也加个电吧。
