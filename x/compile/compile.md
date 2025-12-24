
```bash
go build -gcflags "-S -N -l" -o /dev/null main.go 2> main.s

go build -gcflags "all=-S -N -l" -o /dev/null main.go 2> main.s
```