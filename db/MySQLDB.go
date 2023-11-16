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
	"time"
	"unicode"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed resource/*
var Fs embed.FS

// 组合操作
func CombinationOp() {
	defer ErrorHandler()
	ExampleStr()
	config := ParseConfigFile()
	connect := GetMySQLTCPConnect(*config.Db.User, *config.Db.Password, "information_schema",
		*config.Db.Host, *config.Db.Port)
	tableData := GetSpecificTableDefine(connect, *config.Db.DbName, *config.Db.TableName)
	columnsData := GetSpecificTableColumsDefine(connect, *config.Db.DbName, *config.Db.TableName)
	fmt.Println(tableData.String())
	for _, value := range *columnsData {
		fmt.Println(value.String())
	}
}

// 生成po的文档
func GeneratePO(tableData *map[string]interface{}, columnsData *[]InfoColumns) *string {
	return nil
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
func GetMySQLTCPConnect(user, password, dbname, host, port string) *gorm.DB {
	openUrl := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(openUrl))
	if err != nil {
		panic(err)
	}
	return db
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
		fmt.Printf("错误:%s", r)
		time.Sleep(10 * time.Second)
	}
}
