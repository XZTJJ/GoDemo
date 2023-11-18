package db

/**
 * 生成Java对应的pojo，service等
 * 引入对应的包
 * `go mod init  xxxx`，xxxx为module的名称
 *  go get gorm.io/gorm
 *  go get gorm.io/driver/mysql
 */
import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed resource/*
var Fs embed.FS

const TimeFormat = "2023-11-17 12:47:00"

// 组合操作
func CombinationOp() {
	defer ErrorHandler()
	ExampleStr()
	config := ParseConfigFile()
	connect := GetMySQLTCPConnect(config.Db, "information_schema")
	tableData := GetSpecificTableDefine(connect, *config.Db.DbName, *config.Db.TableName)
	columnsData := GetSpecificTableColumsDefine(connect, *config.Db.DbName, *config.Db.TableName)
	fileStr := GenerateJavaString(tableData, columnsData, config)
	for _, value := range fileStr {
		fmt.Println(value.FilePathName)
		fmt.Println(value.content)
		fmt.Printf("\n\n\n\n\n")
	}
}

// 生成Java的文件内容
func GenerateJavaString(tableData *InfoTable, columnsData *[]InfoColumns, config *JsonConfig) []TemplateJavaFile {
	//生成对应的po模板对象
	p := JavaPoTemplate{}
	p.FillPoTemplate(tableData, columnsData, config)
	//生成PoJavaClass对象
	po := PoJavaClass{}
	po.FillPoJavaClass(&p)
	poFile := ParseTemplate(po)
	//生成DtoJavaClass对象
	dto := DtoJavaClass{}
	dto.FillDtoJavaClass(&p)
	dtoFile := ParseTemplate(dto)
	//生成mapperJavaClass对象
	mapperJava := MapperJavaClass{}
	mapperJava.FillMapperJavaClass(&p, &po)
	mapperJavaFile := ParseTemplate(mapperJava)
	//生成mapperXml对象
	mapperXml := MapperXmlFile{}
	mapperXml.FillMapperXmlFile(&p, &po, &mapperJava)
	mapperXmlFile := ParseTemplate(mapperXml)
	//生成controller对象
	controller := ControllerJavaClass{}
	controller.FillControllerJavaClass(&p)
	controllerFile := ParseTemplate(controller)
	//生成service对象
	service := ServiceJavaClass{}
	service.FillServiceJavaClass(&p, &po)
	serviceFile := ParseTemplate(service)
	//生成serviceImpl对象
	serviceImpl := ServiceImplJavaClass{}
	serviceImpl.FillServiceImplJavaClass(&p, &po, &mapperJava, &service)
	serviceImplFile := ParseTemplate(serviceImpl)
	return []TemplateJavaFile{poFile, dtoFile, mapperJavaFile, mapperXmlFile, controllerFile, serviceFile, serviceImplFile}
}

// 解析模板并且返回字符串
func ParseTemplate(v interface{}) TemplateJavaFile {
	tFile := TemplateJavaFile{}
	resourceTemplateFile := ""
	switch t := v.(type) {
	case PoJavaClass:
		tFile.FilePathName = PoJavaClass(t).FilePathName
		resourceTemplateFile = "resource/poTemplate.txt"
	case DtoJavaClass:
		tFile.FilePathName = DtoJavaClass(t).FilePathName
		resourceTemplateFile = "resource/dtoTemplates.txt"
	case MapperJavaClass:
		tFile.FilePathName = MapperJavaClass(t).FilePathName
		resourceTemplateFile = "resource/mapperJava.txt"
	case MapperXmlFile:
		tFile.FilePathName = MapperXmlFile(t).FilePathName
		resourceTemplateFile = "resource/mapperXml.txt"
	case ControllerJavaClass:
		tFile.FilePathName = ControllerJavaClass(t).FilePathName
		resourceTemplateFile = "resource/controller.txt"
	case ServiceJavaClass:
		tFile.FilePathName = ServiceJavaClass(t).FilePathName
		resourceTemplateFile = "resource/service.txt"
	case ServiceImplJavaClass:
		tFile.FilePathName = ServiceImplJavaClass(t).FilePathName
		resourceTemplateFile = "resource/serviceImpl.txt"
	default:
		panic("没有配置对应类型解析")
	}
	var tmplBytes bytes.Buffer
	t1, err1 := template.ParseFS(Fs, resourceTemplateFile)
	if err1 != nil {
		panic(err1)
	}
	err2 := t1.Execute(&tmplBytes, v)
	if err2 != nil {
		panic(err2)
	}
	tFile.content = tmplBytes.String()
	return tFile
}

// 读取要定义的表结构
func GetSpecificTableDefine(connect *gorm.DB, dbName, tableName string) *InfoTable {
	var results InfoTable
	connect.Select("table_name, table_comment, create_time").
		Where(" table_schema = ? AND table_name = ? ", dbName, tableName).Take(&results)
	return &results
}

// 读取要定义的表结构
func GetSpecificTableColumsDefine(connect *gorm.DB, dbName, tableName string) *[]InfoColumns {
	var results []InfoColumns
	connect.Where(" table_schema = ? AND table_name = ? ", dbName, tableName).Order("ordinal_position").Find(&results)
	return &results
}

// 获取数据库连接
func GetMySQLTCPConnect(db *DbConfig, dbname string) *gorm.DB {
	var user, password, host, port string = "", "", "", "80"
	if db.User != nil {
		user = *db.User
	}
	if db.Password != nil {
		password = *db.Password
	}
	if db.Host != nil {
		host = *db.Host
	}
	if db.Port != nil {
		port = *db.Port
	}
	if user == "" || password == "" || host == "" || port == "" {
		panic("数据库配置错误，账号，密码，host，端口不能为空")
	}
	openUrl := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	connect, err := gorm.Open(mysql.Open(openUrl))
	if err != nil {
		panic(err)
	}
	return connect
}

// 模板示例
func ExampleStr() {
	//提供选项
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("是否需要显示json文件配置示例[Y/N]：")
	name, err1 := reader.ReadString('\n')
	if err1 != nil {
		panic(err1)
	}
	name = strings.TrimSpace(name)
	if name == "N" || name == "n" {
		return
	}
	content, err2 := Fs.ReadFile("resource/configjson.txt")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(content))
}

// 用于解析配置文件
func ParseConfigFile() *JsonConfig {
	//输入文件位置,去掉换行符
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入json配置文件路径，回车键结束，模板参考前面示例：")
	name, err1 := reader.ReadString('\n')
	if err1 != nil {
		panic(err1)
	}
	name = strings.TrimSpace(name)
	//开始解析文件
	const BufferSize = 1024
	file, err2 := os.Open(name)
	if err2 != nil {
		panic(err2)
	}
	defer file.Close()
	buffer := make([]byte, BufferSize)
	for {
		_, err3 := file.Read(buffer)
		if err3 != nil {
			if err3 != io.EOF {
				panic(err3)
			}
			break
		}
	}
	//去掉不可见字符
	needTrimChar := []byte{'\r', '\t', '\v', '\f', '\n', ' ', 0x85, 0xA0}
	for _, value := range buffer {
		if !unicode.IsPrint(rune(value)) {
			needTrimChar = append(needTrimChar, value)
		}
	}
	for _, value := range needTrimChar {
		buffer = bytes.Replace(buffer, []byte{value}, []byte(""), -1)
	}
	//转成map
	var jsonConfig JsonConfig
	err4 := json.Unmarshal(buffer, &jsonConfig)
	if err4 != nil {
		panic(err4)
	}
	return &jsonConfig
}

// 错误处理
func ErrorHandler() {
	if r := recover(); r != nil {
		fmt.Printf("错误：%s", r)
		time.Sleep(10 * time.Second)
	}
}

// 下划线转驼峰写法
func UderscoreToUpperCamelCase(s string) []byte {
	strBytes := []byte(s)
	i, isSplit := 0, false
	for _, value := range strBytes {
		if value == '_' || value == '-' {
			isSplit = true
		} else {
			//a对应的ASCII编码为97,A的编码为65，z为122，Z为90
			if isSplit && value >= 97 && value <= 122 {
				value = value - 32
			}
			strBytes[i] = value
			i, isSplit = i+1, false
		}
	}
	return strBytes[:i]
}

// 生成文件名
func ModifyFileName(fullPathClassName string, fileSuffix string) string {
	path := strings.ReplaceAll(fullPathClassName, ".", string(os.PathSeparator))
	return path + fileSuffix
}
