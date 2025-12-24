### 单元测试三原则

- 首先是 Automatic（自动化）原则
- 接着是 Independent（独立性）原则。
- 最后是 Repeatable（可重复）原则。
 ```shell
 https://github.com/stretchr/testify
 https://github.com/smartystreets/goconvey
 "github.com/agiledragon/gomonkey/v2"
 "github.com/golang/mock/gomock"
 https://github.com/bytedance/mockey
https://github.com/micvbang/go-mocky
 go test -coverprofile=coverage.out

 table-driven
 
PASS
coverage: 45.5% of statements
ok      server-go/21    0.569s

go tool cover -html=coverage.out -o coverage.html
 ```