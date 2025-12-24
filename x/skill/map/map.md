### 命令行查看GC
```shell

GODEBUG=gctrace=1 go run main.go

go build -o main main.go

GODEBUG=gctrace=1 ./main
```


### 缓存大规模数据时

1. 缓存大规模数据时，为了避免GC开销,key-value不能包含指针类型且key-value的大小不能超过128字节
2. 可以使用开源库，比如allegro/bigcache, coocood/freecache
