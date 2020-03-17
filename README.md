# go-ini
简单的读取ini配置文件，将配置信息保存在结构体中。
Simple read ini file, the configuration information is stored in the struct.

### 使用前，你需要根据你的需求在`settings.go`中定义相关结构体用来保存你读取的信息：
在这里，配置你想从ini文件中读出信息的结构体，并把它汇总到Config结构体下。
如果ini文件中的一节信息为：
```
[mysql]
address=192.168.190.130
port=3306
username=root
password=123456
```
那么我们需要定义的结构体如下：
```go
type MysqlConfig struct {
    Address  string `ini:"address"`
    Port     int    `ini:"port"`
    Username string `ini:"username"`
    Password string `ini:"password"`
}

type Config struct {
    MysqlConfig `ini:"mysql"`
}
```
#### 要注意的是结构体里字段的标签一定要为ini的信息相对应

## 以下代码可以帮助你快速的使用：
```go
func main() {
	var cfg goini.Config
	err := goini.LoadIni("./conf.ini", &cfg)
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	fmt.Println(cfg)
}
```

`goini.LoadIni`接收两个参数，第一个是`ini`文件的路径，第二个传递`goini.Config`类型的指针。
`goini.LoadIni`返回错误的信息,将ini中的配置信息保存在`goini.Config`中。