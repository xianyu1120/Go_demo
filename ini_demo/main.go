package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//ini 配置解析器
type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	UserName string `ini:"username"`
	Password string `ini:"password"`
}
type Redis struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
}
type Config struct {
	MysqlConfig `ini:"mysql"`
	Redis       `ini:"redis"`
}

//
func loadIni(fileName string, data interface{}) (err error) {
	//0. 参数的校验

	//0.1 传进来的data参数必须是指针类型（因为要在函数中对其赋值
	t := reflect.TypeOf(data) //传进来的结构体类型
	if t.Kind() != reflect.Ptr {
		//不是指针
		err = errors.New("data should be a pointer") //格式化输出之后返回一个error类型
		return
	}
	//0.2 传进来的data参数必须是结构体类型指针（因为配置文件中各种键值对需要赋值给结构体的字段）
	if t.Elem().Kind() != reflect.Struct { //0.1 保证了此时是指针，此时判断值类型是不是结构体
		err = errors.New("data parm should be a struct type")
		return
	}
	//1. 读文件得到字节型数据
	b, err := ioutil.ReadFile(fileName) //读取整个文件
	if err != nil {
		return
	}
	// string(b) //将文件内容由字节类型转换为字符串
	lineSlice := strings.Split(string(b), "\r\n")
	// fmt.Printf("%#v\n", lineSlice)
	//2. 一行一行得读数据

	var structName string              //结构体名称
	for idx, line := range lineSlice { //遍历读取处理每一行

		//如果是空行就跳过
		if len(line) == 0 {
			continue
		}
		//去掉字符串首尾空格
		line = strings.TrimSpace(line)
		//2.1 如果是注释就跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") { //以;或者#开头表示注释
			continue //跳过当前循环，开始下一次循环
		}
		//2.2 如果是[]开头表示是节（section）
		if strings.HasPrefix(line, "[") {
			//如果【】之间为空或者【】不完整
			if line[len(line)-1] != ']' { //不完整
				err = fmt.Errorf("line :%d systax errors", idx+1)
				return
			}
			//这一行[]去掉，拿到中间内容把首尾去掉
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 { //为空
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			//根据字符串sectionName去data里面根据反射找到对应结构体
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				if sectionName == field.Tag.Get("ini") {
					//说明找到了中间对应嵌套结构体，记录字段名
					structName = field.Name
					fmt.Printf("找到 %s对应的嵌套结构体 %s\n", sectionName, structName)
				}
			}
		} else {
			//行末为】
			if strings.HasSuffix(line, "]") {
				err = fmt.Errorf("line :%d systax errors", idx+1)
				return
			}
			//不为】
			//2.3如果不是[]开头就是=分割键值对
			// 1. 以等号分割这一行，左边为key,右边为value
			//丢弃不符合的字段
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") || strings.HasSuffix(line, "=") { //不存在等号,或者首字母为=,或者xxx=
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			index := strings.Index(line, "=")          //等号出现的下标
			key := strings.TrimSpace(line[:index])     //key为=之前的字段，字段去除空格
			value := strings.TrimSpace(line[index+1:]) //value 为=之后
			fmt.Println(key, value)
			//2. 根据structName去data里面把对应嵌套结构体给取出来
			v := reflect.ValueOf(data).Elem()   //传入结构体对应值
			sValue := v.FieldByName(structName) //得到嵌套结构体的结构体中的值信息
			sType := sValue.Type()              //嵌套结构体中结构体的类型

			fmt.Println(sType.Name(), sType, sValue)
			if sType.Kind() != reflect.Struct { //不是一个嵌套结构体
				err = fmt.Errorf("data 中的%s字段应该是一个结构体", structName)
				return
			}

			var fieldName string             //结构体中字段名
			var fileType reflect.StructField //字段名类型
			//3. 遍历嵌套结构体的每一个字段，判断tag是不是等于key
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i) //tag信息是存储在类型信息中的
				fileType = field        //结构体中字段
				if field.Tag.Get("ini") == key {
					//找到对应字段
					fieldName = field.Name //对应结构体字段名称
					break
				}
			}
			//4. 如果key=tag,给字段赋值
			//4.1 根据fieldName 去取出这个字段
			fileObj := sValue.FieldByName(fieldName) //根据名称得到结构体中的字段对象
			//4.2 对其赋值
			fmt.Printf("%T\n",fileObj)
			fmt.Println(fileObj, fileObj.Type())
			fmt.Println("=====")
			fmt.Printf("字段名称：%s, 字段类型：%v\n", fieldName, fileType.Type.Kind())
			
			switch fileType.Type.Kind() {
			case reflect.String:
				if fileObj.CanSet() {
					fileObj.SetString(value)
				}

			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16:
				var valueInt int64
				valueInt, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("line %d type error", idx+1)
					return
				}
				fileObj.SetInt(valueInt)
			}
		}
	}
	return
}
func main() {
	var c Config
	err := loadIni("./conf.ini", &c)
	if err != nil {
		fmt.Printf("load ini feiled,err:%v\n", err)
		return
	}

	fmt.Println(c)
}
