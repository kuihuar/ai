包的依赖关系依次对包进行初始化，在这个过程中，每个包的 init 函数会自动执行
init 函数没有参数，也没有返回值，其函数签名为 func init() {...}。
一个包内可以定义多个 init 函数，它们会按照在文件里的出现顺序依次执行
执行顺序固定
包的初始化顺序为：首先初始化包级别的变量，然后按照 init 函数在文件里的出现顺序依次执行 init 函数。如果存在多个包，会先初始化依赖的包，再初始化当前包。

### 使用场景
初始化全局变量
可以利用 init 函数对全局变量进行复杂的初始化操作。

注册操作
在一些框架中，经常会使用 init 函数进行组件的注册。例如，在数据库驱动里，会在 init 函数中注册驱动。
```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "user:password@/dbname")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
    // 使用数据库连接
}
```
```go
func init() {
	if driverName != "" {
		sql.Register(driverName, &MySQLDriver{})
	}
}
```
github.com/go-sql-driver/mysql 包的 init 函数会注册 MySQL 驱动，这样在 sql.Open 时就能使用该驱动。

配置加载
可以在 init 函数里加载配置文件，保证在程序启动时配置信息已经准备好
```go
package main

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "io/ioutil"
)

type Config struct {
    Port int `yaml:"port"`
}

var config Config

func init() {
    data, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        panic(err.Error())
    }
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        panic(err.Error())
    }
}

func main() {
    fmt.Println("配置的端口号为:", config.Port)
}
```


