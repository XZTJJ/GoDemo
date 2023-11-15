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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 表结构
type InfoTables string

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
func (i InfoTables) TableName() string {
	return "TABLES"
}

// 绑定表名
func (i *InfoColumns) TableName() string {
	return "COLUMNS"
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
	var tip string
	tip = "目前只是支持MySQL数据库,db表示数据库相关信息，author表示作者信息，typeMapping表示字段类型和java类型的映射关系,下面给出配置文件的json示例："
	fmt.Println(tip)
	mapConfig := make(map[string]map[string]string)
	mapConfig["db"] = map[string]string{"user": "admin", "password": "123456", "host": "127.0.0.1", "port": "3306", "dbName": "demoDatabase", "tableName": "demoXiaohaizi"}
	mapConfig["author"] = map[string]string{"packageName": "com.zhc.test", "author": "xiaohaizi", "email": "12345@qq.com", "tablePrefix": "fk_"}
	mapConfig["typeMapping"] = map[string]string{"tinyint": "Byte", "smallint": "Integer", "mediumint": "Integer", "int": "Integer", "integer": "Integer", "bigint": "Long", "float": "Float", "double": "Double", "decimal": "BigDecimal", "bit": "Boolean", "char": "String", "varchar": "String", "tinytext": "String", "text": "String", "mediumtext": "String", "longtext": "String", "date": "LocalDate", "datetime": "LocalDateTime", "timestamp": "LocalDateTime"}
	filebytes, err := json.MarshalIndent(mapConfig, "", "    ")
	if err != nil {
		fmt.Println("系统格式化json失败，直接退出", err)
		os.Exit(1)
	}
	fileStr := string(filebytes)
	fmt.Println(fileStr)
}

// 用于解析配置文件
func ParseConfigFile() *map[string]map[string]string {
	//输入文件位置,去掉换行符
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入json配置文件路径，回车键结束，模板参考前面示例：")
	name, err1 := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	name = strings.TrimSuffix(name, "\r")
	if err1 != nil {
		fmt.Println("输入有误", err1)
		os.Exit(1)
	}
	//开始解析文件
	const BufferSize = 1024 * 1024
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
	//解析成map形式
	//buffer = bytes.ReplaceAll(buffer, []byte{'\r', '\n'}, []byte{})
	//fmt.Println(buffer)
	mapTemp := make(map[string]map[string]string)
	err4 := json.Unmarshal(buffer, &mapTemp)
	if err4 != nil {
		fmt.Println("json反序列化失败", err4)
		os.Exit(1)
	}
	fmt.Println(mapTemp)
	return nil
}
