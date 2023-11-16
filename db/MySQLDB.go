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
	"unicode"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed resource/*
var Fs embed.FS

// 表结构，通是绑定表名，使用指针的原因，是为了处理空的默认值操作
type InfoColumns struct {
	TableCatalog           *string `gorm:"column:TABLE_CATALOG"`
	TableSchema            *string `gorm:"column:TABLE_SCHEMA"`
	TableNameStr           *string `gorm:"column:TABLE_NAME"`
	ColumnName             *string `gorm:"column:COLUMN_NAME"`
	OrdinalPosition        *int    `gorm:"column:ORDINAL_POSITION"`
	ColumnDefault          *string `gorm:"column:COLUMN_DEFAULT"`
	IsNullable             *string `gorm:"column:IS_NULLABLE"`
	DataType               *string `gorm:"column:DATA_TYPE"`
	CharacterMaximumLength *int    `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   *int    `gorm:"column:CHARACTER_OCTET_LENGTH"`
	NumericPrecision       *int    `gorm:"column:NUMERIC_PRECISION"`
	NumericScale           *int    `gorm:"column:NUMERIC_SCALE"`
	DatetimePrecision      *int    `gorm:"column:DATETIME_PRECISION"`
	CharacterSetName       *string `gorm:"column:CHARACTER_SET_NAME"`
	CollationName          *string `gorm:"column:COLLATION_NAME"`
	ColumnType             *string `gorm:"column:COLUMN_TYPE"`
	ColumnKey              *string `gorm:"column:COLUMN_KEY"`
	Extra                  *string `gorm:"column:EXTRA"`
	Privileges             *string `gorm:"column:PRIVILEGES"`
	ColumnComment          *string `gorm:"column:COLUMN_COMMENT"`
	GenerationExpression   *string `gorm:"column:GENERATION_EXPRESSION"`
	SrsId                  *string `gorm:"column:SRS_ID"`
}

// 绑定表名
func (i *InfoColumns) TableName() string {
	return "COLUMNS"
}

// 组合操作
func CombinationOp() {
	ExampleStr()
	configMap := ParseConfigFile()
	connect := GetMySQLTCPConnect((*configMap)["db"]["user"], (*configMap)["db"]["password"],
		"information_schema", (*configMap)["db"]["host"], (*configMap)["db"]["port"])
	tableData := GetSpecificTableDefine(connect, (*configMap)["db"]["dbName"], (*configMap)["db"]["tableName"])
	columnsData := GetSpecificTableColumsDefine(connect, (*configMap)["db"]["dbName"], (*configMap)["db"]["tableName"])
	fmt.Println(*tableData)
	for _, value := range *columnsData {
		fmt.Println(*(value.ColumnName), *(value.DataType), *(value.ColumnComment), *(value.ColumnKey), *(value.Extra))
	}
}

// 生成po的文档
func GeneratePO(tableData *map[string]interface{}, columnsData *[]InfoColumns) *string {
	return nil
}

// 读取要定义的表结构
func GetSpecificTableDefine(connect *gorm.DB, dbName, tableName string) *map[string]interface{} {
	var results map[string]interface{}
	connect.Table("TABLES").Select("table_name as tableName, table_comment as tableComment, create_time as createTime").
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
		fmt.Println("数据库连接失败", err)
		os.Exit(1)
	}
	return db
}

// 模板示例
func ExampleStr() {
	//提供选项
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n是否需要显示json文件配置示例[Y/N]：")
	name, err1 := reader.ReadString('\n')
	if err1 != nil {
		fmt.Println("是否需要显示json文件配置示例输出失败", err1)
		os.Exit(1)
	}
	name = strings.TrimSpace(name)
	if name == "N" || name == "n" {
		return
	}
	content, err2 := Fs.ReadFile("resource/configjson.txt")
	if err2 != nil {
		fmt.Println("系统读取json模板文件失败，直接退出", err2)
		os.Exit(1)
	}
	fmt.Println(string(content))
}

// 用于解析配置文件
func ParseConfigFile() *map[string]map[string]string {
	//输入文件位置,去掉换行符
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入json配置文件路径，回车键结束，模板参考前面示例：")
	name, err1 := reader.ReadString('\n')
	if err1 != nil {
		fmt.Println("文件路径输入有误", err1)
		os.Exit(1)
	}
	name = strings.TrimSpace(name)
	//开始解析文件
	const BufferSize = 1024
	file, err2 := os.Open(name)
	if err2 != nil {
		fmt.Println("文件打开失败", err2)
		os.Exit(1)
	}
	defer file.Close()
	buffer := make([]byte, BufferSize)
	for {
		_, err3 := file.Read(buffer)
		if err3 != nil {
			if err3 != io.EOF {
				fmt.Println("文件读取失败", err3)
				os.Exit(1)
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
	mapTemp := make(map[string]map[string]string)
	err4 := json.Unmarshal(buffer, &mapTemp)
	if err4 != nil {
		fmt.Println("json反序列化失败", err4)
		os.Exit(1)
	}
	return &mapTemp
}
