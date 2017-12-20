# ggf-othello-recoder
ggfファイルからオセロのデータを取り出す。

## Outputの形式 二種類

### board 情報
64この数字
```
1 2 1......1 2 
```

### 置いた場所
3つの数字
```
Y X pass
4 6  2
```
pass は1がtrue 2がfalse

# Usage
```
ggf -f sample.ggf -w white -t white
```
-f ファイルを指定

-w 勝者を指定

-t ターンを指定(指定されなかった場合winnerと一緒)