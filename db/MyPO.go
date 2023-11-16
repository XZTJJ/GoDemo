package db

import (
	"bytes"
	"fmt"
	"time"
)

// 表的表结构，通是绑定对应表，使用指针的原因，是为了处理空的默认值操作
type InfoTable struct {
	TableCatalog   *string    `gorm:"column:TABLE_CATALOG"`
	TableSchema    *string    `gorm:"column:TABLE_SCHEMA"`
	TableNameStr   *string    `gorm:"column:TABLE_NAME"`
	TableType      *string    `gorm:"column:TABLE_TYPE"`
	Engine         *string    `gorm:"column:ENGINE"`
	Version        *int       `gorm:"column:VERSION"`
	RowFormat      *string    `gorm:"column:ROW_FORMAT"`
	TableRows      *int64     `gorm:"column:TABLE_ROWS"`
	AvgRowLength   *int64     `gorm:"column:AVG_ROW_LENGTH"`
	DataLength     *int64     `gorm:"column:DATA_LENGTH"`
	MaxDataLength  *int64     `gorm:"column:MAX_DATA_LENGTH"`
	IndexLength    *int64     `gorm:"column:INDEX_LENGTH"`
	DataFree       *int64     `gorm:"column:DATA_FREE"`
	AutoIncrement  *int64     `gorm:"column:AUTO_INCREMENT"`
	CreateTime     *time.Time `gorm:"column:CREATE_TIME"`
	UpdateTime     *time.Time `gorm:"column:UPDATE_TIME"`
	CheckTime      *time.Time `gorm:"column:CHECK_TIME"`
	TableCollation *string    `gorm:"column:TABLE_COLLATION"`
	Checksum       *string    `gorm:"column:CHECKSUM"`
	CreateOptions  *string    `gorm:"column:CREATE_OPTIONS"`
	TableComment   *string    `gorm:"column:TABLE_COMMENT"`
}

// 字段的表结构，通是绑定对应字段，使用指针的原因，是为了处理空的默认值操作
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
func (i *InfoTable) TableName() string {
	return "TABLES"
}

// 打印对象
func (i *InfoTable) String() string {
	format := `TableCatalog:%v, TableSchema:%v, TableNameStr:%v, TableType:%v, Engine:%v, Version:%v, 
RowFormat:%v, TableRows:%v, AvgRowLength:%v, DataLength:%v, MaxDataLength:%v, IndexLength:%v, DataFree:%v, 
AutoIncrement:%v, CreateTime:%v, UpdateTime:%v, CheckTime:%v, TableCollation:%v, Checksum:%v, CreateOptions:%v, TableComment:%v`
	str := fmt.Sprintf(format, *i.TableCatalog, *i.TableSchema, *i.TableNameStr, *i.TableType, *i.Engine, *i.Version,
		*i.RowFormat, *i.TableRows, *i.AvgRowLength, *i.DataLength, *i.MaxDataLength, *i.IndexLength, *i.DataFree,
		*i.AutoIncrement, *i.CreateTime, *i.UpdateTime, *i.CheckTime, *i.TableCollation, *i.Checksum,
		*i.CreateOptions, *i.TableComment)
	return str
}

// 绑定表名
func (i *InfoColumns) TableName() string {
	return "COLUMNS"
}

// 打印对象
func (i *InfoColumns) String() string {
	format := `"TableCatalog:%s, TableSchema:%s, TableNameStr:%s, ColumnName:%s, OrdinalPosition:%d, 
ColumnDefault:%s, IsNullable:%s, DataType:%s, CharacterMaximumLength:%d, CharacterOctetLength:%d, 
NumericPrecision:%d, NumericScale:%d, DatetimePrecision:%d, CharacterSetName:%s, CollationName:%s, 
ColumnType:%s, ColumnKey:%s, Extra:%s, Privileges:%s, ColumnComment:%s, GenerationExpression:%s, SrsId:%s"`
	str := fmt.Sprintf(format, *i.TableCatalog, *i.TableSchema, *i.TableNameStr, *i.ColumnName, *i.OrdinalPosition,
		*i.ColumnDefault, *i.IsNullable, *i.DataType, *i.CharacterMaximumLength, *i.CharacterOctetLength,
		*i.NumericPrecision, *i.NumericScale, *i.DatetimePrecision, *i.CharacterSetName, *i.CollationName,
		*i.ColumnType, *i.ColumnKey, *i.Extra, *i.Privileges, *i.ColumnComment, *i.GenerationExpression, *i.SrsId)
	return str
}

// 数据库的配置
type DbConfig struct {
	DbName    *string `json:"dbName"`
	Host      *string `json:"host"`
	Password  *string `json:"password"`
	Port      *string `json:"port"`
	TableName *string `json:"tableName"`
	User      *string `json:"user"`
}

// 作者信息配置
type SummaryConfig struct {
	Author      *string `json:"author"`
	Email       *string `json:"email"`
	PackageName *string `json:"packageName"`
	TablePrefix *string `json:"tablePrefix"`
}

// json配置文件配置
type JsonConfig struct {
	Db          *DbConfig          `json:"db"`
	Summary     *SummaryConfig     `json:"author"`
	TypeMapping *map[string]string `json:"typeMapping"`
}

// JsonConfig的格式化方式
func (j *JsonConfig) String() string {
	db := fmt.Sprintf("Db:{DbName:%s,Host:%s,Password:%s,Port:%s,TableName:%s,User:%s}", *j.Db.DbName,
		*j.Db.Host, *j.Db.Password, *j.Db.Port, *j.Db.TableName, *j.Db.User)
	summary := fmt.Sprintf("Summary:{Author:%s,Email:%s,PackageName:%s,TablePrefix:%s}", *j.Summary.Author,
		*j.Summary.Email, *j.Summary.PackageName, *j.Summary.TablePrefix)
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(db)
	buffer.WriteString(",")
	buffer.WriteString(summary)
	buffer.WriteString(",TypeMapping:{")
	for key, value := range *j.TypeMapping {
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(value)
		buffer.WriteString(",")
	}
	temp := []byte(buffer.String())
	temp = temp[:len(temp)-1]
	temp = append(temp, '}')
	temp = append(temp, '}')
	return string(temp)
}
