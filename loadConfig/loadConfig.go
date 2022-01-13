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

func ClientFac() *messagetype.ClientConfig {
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

	//3、对应字段赋值

}

func loadini(configFile string, cfg *config) {
	conType := reflect.TypeOf(cfg)
	conVaule := reflect.ValueOf(cfg)
	var confStructName string
	var sectionName string

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("read config file failed ", err)
	}
	// 读取配置文件，并且跟"\r\n" 换行进行字符串分割
	conf := strings.Split(string(b), "\r\n")
	//fmt.Println(conf)
	// 循环处理读取的配置文件
	for idx, line := range conf {
		//判断是否符合相应的要求
		// 空行，注释去掉
		if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		//判断格式是否满足
		if strings.HasPrefix(line, "=") || line[0] == '[' && line[len(line)-1] != ']' {
			//提示哪一行不符合规范
			err := fmt.Errorf("config file line %d is wrong ", idx+1)
			fmt.Println(err)
		}
		// 获取配置的节点信息,并保存到sectionName
		if line[0] == '[' && line[len(line)-1] == ']' {

			sectionName = strings.Trim(line, "[]")
			fmt.Println("找到配置文件 ", sectionName)
			//遍历配置结构体，查询是否有对应的配置文件
			for i := 0; i < conType.Elem().NumField(); i++ {
				fieldName := conType.Elem().Field(i)
				if fieldName.Tag.Get("ini") == sectionName {
					confStructName = fieldName.Name
					fmt.Println("在结构体中查询到对应的配置信息", confStructName)
					break
				}
			}

		} else {
			//切割行变换成 切片，条件“=”
			line := strings.Split(line, "=")
			// 获取kv
			key := strings.TrimSpace(line[0])
			vaule := strings.TrimSpace(line[1])
			// 根据上文拿到的代码中的配置
			subConV := conVaule.Elem().FieldByName(confStructName)
			subConT := subConV.Type()
			var subConFieldName string
			var subConFieldType reflect.StructField

			for i := 0; i < subConV.NumField(); i++ {
				field := subConT.Field(i)
				subConFieldType = field
				if field.Tag.Get("ini") == key {
					subConFieldName = field.Name
					fmt.Println(subConFieldName)
					break
				}
			}
			fieldObj := subConV.FieldByName(subConFieldName)
			switch subConFieldType.Type.Kind() {
			case reflect.String:
				fieldObj.SetString(vaule)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, _ := strconv.ParseInt(vaule, 10, 64)
				fieldObj.SetInt(v)
			}
		}
	}
}
