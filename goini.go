package goini

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// ini配置文件解析器


// LoadIni 读取ini配置文件
func LoadIni(fileName string, data *Config) (err error) {
	// 1. 参数校验
	// 1.1 传进来的data参数必须是指针类型（因为需要在函数中对其赋值
	t := reflect.TypeOf(data)
	// fmt.Println(t, t.Kind())
	if t.Kind() != reflect.Ptr {
		err = errors.New("data param should be a pointer") // 创建一个错误
		return
	}

	// 1.2 传进来的data参数必须是结构体类型的指针
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data param should be a struct pointer")
		return
	}

	// 2. 读文件得到字节类型的数据
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	lineSlice := strings.Split(string(b), "\r\n")
	// fmt.Printf("%#v\n", lineSlice)

	// 3. 一行一行得到数据
	var structName string

	for idx, line := range lineSlice {
		line = strings.TrimSpace(line) // 去掉字符串首尾空格
		if len(line) == 0 {
			continue // 跳过空行
		}

		// 3.1 如果是注释就跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 3.2 如果是[开头的就表示是节(section)
		if strings.HasPrefix(line, "[") {
			if line[0] != '[' || line[len(line)-1] != ']' {
				err = fmt.Errorf("line:%d syntax error1", idx+1)
				return
			}
			// 取[]中间的内容把并首尾的空格去掉
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			// fmt.Println(sectionName)

			if len(sectionName) == 0 {
				err = fmt.Errorf("line:%d syntax error2", idx+1)
				return
			}
			// 根据字符串sectionName去data里面根据反射找到对应的结构体
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				if sectionName == field.Tag.Get("ini") {
					// 说明找到了对应的嵌套结构体，把字段名记下来
					structName = field.Name
					fmt.Printf("找到%s对应的嵌套结构体%s\n", sectionName, structName)
				}
			}
		} else {
			// 3.3 如果不是[开头就是=分割的键值对
			// 1. 以=分割这一行，等号左边是key，等号右边是value
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line:%d syntax error3", idx+1)
				return
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])
			// 2. 根据structName 去 data 里把对应的嵌套结构体取出来
			v := reflect.ValueOf(data)
			sValue := v.Elem().FieldByName(structName) // 拿到嵌套结构体的值信息
			sType := sValue.Type()                     // 拿到嵌套结构体的类型信息

			if sType.Kind() != reflect.Struct {
				err = fmt.Errorf("data中的%s字段应该是一个结构体", structName)
				return
			}
			// 3. 遍历结构体的每一个字段，判断tag是不是等于key
			var fieldName string
			var fieldType reflect.StructField
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i) // tag信息存储在类型信息中
				fieldType = field
				if field.Tag.Get("ini") == key {
					// 找到对应的字段
					fieldName = field.Name
					break
				}
			}
			// 4. 找到对应字段（key = tag） 给这个字段赋值
			// 4.1 根据fieldName 取出这个字段
			if len(fieldName) == 0 {
				// 在结构体中找不到对应的字符
				continue
			}
			// 4.2 对其赋值
			fieldObj := sValue.FieldByName(fieldName)
			fmt.Println(fieldName, fieldType.Type.Kind())
			switch fieldType.Type.Kind() {
			// 字符串
			case reflect.String:
				fieldObj.SetString(value)
			// 整型
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var valueInt int64
				valueInt, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("line:%d value type error4", idx+1)
					return
				}
				fieldObj.SetInt(valueInt)
			// BOOL型
			case reflect.Bool:
				var valueBool bool
				valueBool, err = strconv.ParseBool(value)
				if err != nil {
					err = fmt.Errorf("line:%d value type error5", idx+1)
					return
				}
				fieldObj.SetBool(valueBool)
			// 浮点型
			case reflect.Float32, reflect.Float64:
				var valueFloat float64
				valueFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					err = fmt.Errorf("line:%d value type error6", idx+1)
					return
				}
				fieldObj.SetFloat(valueFloat)
			}
		}
	}
	return
}
