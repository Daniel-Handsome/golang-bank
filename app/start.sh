# /bin/sh
## https://stackoverflow.com/questions/73706572/reference-shell-script-variable-in-dockerfile
## docker 如何用變數 , 其實就是因為docker就是環境變數 機上sh本來就可以印出環境變數了

##若指令传回值不等于0，则立即退出shell 簡單來說看變數有沒有載入
set -e

## $@ 會將 sh a,b,c 印出a,b,c這樣
## exec是會去執行
## 簡單來說 獲取參數並執行
exec "$@"

