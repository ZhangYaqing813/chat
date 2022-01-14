package loadConfig

import (
	messagetype "chat/Message_type"
	chatlog "chat/chatLog"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

var ServerCfg messagetype.ServerConfig
var ClientCfg messagetype.ClientConfig

func ServerFac(file string) *messagetype.ServerConfig {
	ServerLoadConfig(file, &ServerCfg)
	return &ServerCfg
}

func ClientFac(file string) *messagetype.ClientConfig {
	CCfg(file, &ClientCfg)
	return &ClientCfg
}

func ServerLoadConfig(file string, ServerCfg *messagetype.ServerConfig) {
	//获取ServerCfg 类型
	conType := reflect.TypeOf(ServerCfg)
	//获取ServerCfg 中的值
	conValue := reflect.ValueOf(ServerCfg)
	//定义需要的一些变量，存储临时数据
	//配置文件中的节点名
	var sectionName string
	//代码中定义的配置文件对应的结构体中的嵌套结构体名
	var conStructName string

	//1、读取文件
	//1.1 加载文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		chatlog.Std.Error(err)
	}
	//1.2 转换为字符串
	cfg := strings.Split(string(b), "\r\n")
	//2、解析
	fmt.Println(cfg)
	//2.1 解析读取到的配置文件
	for idx, line := range cfg {
		//2.2 判断配置文件中行是否存在不符合的语法的行
		// 去掉空行，注释行，以及"["开头不是"]"结尾，"="开头的行进行提示、过滤
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if line[0] == '[' && line[len(line)-1] != ']' || line[0] == '=' {
			chatlog.Std.Errorf("config.ini 第 %d 行格式错误", idx+1)
			return
		}
		//2.3 找到对应的节点，保存。
		// [] 表示节点的开始
		if len(line) > 0 && line[0] == '[' && line[len(line)-1] == ']' {
			sectionName = strings.Trim(line, "[]")
			chatlog.Std.Infof("config 找到节点%s 信息,并缓存至 sectionName 中", sectionName)
			//2.4  解析代码中对应的结构体中的节点信息
			//遍历结构体字段信息
			for i := 0; i < conType.Elem().NumField(); i++ {
				fieldName := conType.Elem().Field(i)
				//根据代码中配置结构体的tag 字段信息找到对应的节点信息
				if fieldName.Tag.Get("ini") == sectionName {
					conStructName = fieldName.Name
					chatlog.Std.Infof("找到与配置文件 %s 应的节点信息")
					break
				}
			}

		} else {
			//3、对应字段赋值
			//3.1 解析line 中的字段，获取key和value
			confLine := strings.Split(line, "=")
			//等号左边的是key
			key := strings.TrimSpace(confLine[0])
			//等号右边的是value
			value := strings.TrimSpace(confLine[1])

			//3.2 获取代码中配置struct 中节点的信息
			//获取sectionName 中的值信息
			subConfV := conValue.Elem().FieldByName(conStructName)
			// 获取subConfV 的类型信息
			subConfT := subConfV.Type()

			//3.2 遍历字段信息
			if subConfT.Kind() != reflect.Struct {
				chatlog.Std.Error("该配置节点不是不符合要求")
			}
			//缓存一下找到的每一个节点的字段信息
			var subConFieldName string
			var fieldName reflect.StructField
			for i := 0; i < subConfT.NumField(); i++ {
				fieldName = subConfT.Field(i)
				//根据ini tag 进行对比，找到字段后跳出
				if fieldName.Tag.Get("ini") == key {
					//找到字段后保存字段名
					subConFieldName = fieldName.Name
					chatlog.Std.Infof("找到对应字段%s", key)
					break
				}
			}
			//3.3 对应的字段进行赋值操作
			//获取到节点字段
			fieldobj := subConfV.FieldByName(subConFieldName)
			switch fieldobj.Type().Kind() {
			case reflect.String:
				fieldobj.SetString(value)
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					chatlog.Std.Errorf("%s 类型转换失败", value)
				}
				fieldobj.SetInt(v)
			}
		}
	}
}

// ClientCfg 配置文件解析

func CCfg(file string, clientcfg *messagetype.ClientConfig) {
	var sectionName string

	confT := reflect.TypeOf(clientcfg)
	confV := reflect.ValueOf(clientcfg)
	// 1、加载配置文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		chatlog.Std.Errorf("读取配置文件失败 %v", err)
	}
	clientCfg := strings.Split(string(b), "\r\n")

	//2、解析配置文件
	for index, line := range clientCfg {
		//处理注释空行
		if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		//处理错误表达式
		if line[0] == '[' && line[len(line)-1] != ']' || strings.HasPrefix(line, "=") {
			chatlog.Std.Errorf("%s 中第 %d 语法错误", file, index+1)
			return
		}

		// 查找节点信息
		if line[0] == '[' && line[len(line)-1] == ']' && len(line) > 0 {
			sectionName = strings.Trim(line, "[]")
			chatlog.Std.Info("找到节点信息", line)
			// 判断代码中传递的是否为结构体，
			if confV.Elem().String() == sectionName && confT.Elem().Kind() == reflect.Ptr && confV.Elem().Kind() == reflect.Struct {
				chatlog.Std.Info("配置文件对应")
			}
		} else {
			lcf := strings.Split(line, "=")
			key := strings.TrimSpace(lcf[0])
			value := strings.TrimSpace(lcf[1])
			var field reflect.StructField
			var fieldName string
			// 遍历结构体字段信息

			for i := 0; i < confV.Elem().NumField(); i++ {
				field = confT.Field(i)
				if field.Tag.Get("ini") == key {
					fieldName = field.Name
					chatlog.Std.Info("字段匹配成功")

					break
				}
			}
			fieldobj := confV.Elem().FieldByName(fieldName)
			switch field.Type.Kind() {
			case reflect.String:
				fieldobj.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					chatlog.Std.Errorf("%v 转换int64失败", value)
					return
				}
				fieldobj.SetInt(v)
			}
		}
	}
}
