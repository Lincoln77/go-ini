package goini

// 在这里，配置你想从ini文件中读出信息的结构体，并把它汇总到Config结构体下
// 如果ini文件中的一节信息为：
// 		[mysql]
//		address=192.168.190.130
// 		port=3306
// 		username=root
// 		password=123456
// 那么我们需要定义的结构体如下：
// 		type MysqlConfig struct {
// 			Address  string `ini:"address"`
// 			Port     int    `ini:"port"`
// 			Username string `ini:"username"`
// 			Password string `ini:"password"`
// 		}
//
// 		type Config struct {
// 			MysqlConfig `ini:"mysql"`
// 			RedisConfig `ini:"redis"`
// 		}
//
// 要注意的是结构体里字段的标签一定要为ini的信息相对应
// 以下是我测试时定义的3个结构体：

// MysqlConfig MySQL配置结构体
type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

// RedisConfig 结构体
type RedisConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
	Test     bool   `ini:"test"`
}

// Config 配置文件结构体
type Config struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}
