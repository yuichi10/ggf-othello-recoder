# ggf-othello-recoder
ggfファイルからオセロのデータを取り出す。

## Outputの形式 二種類

### board 情報
64この数字
```
1 2 1......1 2 
```

### 置いた場所
65の数字最後の数字はpassの時1
```
000....10000
```
pass は1がtrue 0がfalse

# Usage
```
ggf -f sample.ggf -w white -t white
```
-f ファイルを指定

-d ファイルの入ってるディレクトリを指定

-w 勝者を指定

-t ターンを指定(指定されなかった場合winnerと一緒)