package db

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// 定义常量
const (
	//相关包名后缀
	PO_PACKEAGE_SUFFIX           = ".api.po"
	DTO_PACKEAGE_SUFFIX          = ".api.dto"
	MAPPER_PACKEAGE_SUFFIX       = ".mapper"
	CONTROLLER_PACKEAGE_SUFFIX   = ".controller"
	SERVICE_PACKEAGE_SUFFIX      = ".service"
	SERVICE_IMPL_PACKEAGE_SUFFIX = ".service.impl"
	//相关类名后缀
	PO_CLASS_SUFFIX           = "PO"
	DTO_CLASS_SUFFIX          = "DTO"
	MAPPER_CLASS_SUFFIX       = "Mapper"
	CONTROLLER_CLASS_SUFFIX   = "Controller"
	SERVICE_CLASS_SUFFIX      = "Service"
	SERVICE_IMPL_CLASS_SUFFIX = "ServiceImpl"
	//时间格式化
	TIME_FORMAT = "2006-01-02 15:04:06"
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
	var buffer bytes.Buffer
	buffer.WriteString("{")
	if i.TableCatalog != nil {
		buffer.WriteString(fmt.Sprintf("TableCatalog:\"%s\",", *i.TableCatalog))
	}
	if i.TableSchema != nil {
		buffer.WriteString(fmt.Sprintf("TableSchema:\"%s\",", *i.TableSchema))
	}
	if i.TableNameStr != nil {
		buffer.WriteString(fmt.Sprintf("TableNameStr:\"%s\",", *i.TableNameStr))
	}
	if i.TableType != nil {
		buffer.WriteString(fmt.Sprintf("TableType:\"%s\",", *i.TableType))
	}
	if i.Engine != nil {
		buffer.WriteString(fmt.Sprintf("Engine:\"%s\",", *i.Engine))
	}
	if i.Version != nil {
		buffer.WriteString(fmt.Sprintf("Version:\"%d\",", *i.Version))
	}
	if i.RowFormat != nil {
		buffer.WriteString(fmt.Sprintf("RowFormat:\"%s\",", *i.RowFormat))
	}
	if i.TableRows != nil {
		buffer.WriteString(fmt.Sprintf("TableRows:\"%d\",", *i.TableRows))
	}
	if i.AvgRowLength != nil {
		buffer.WriteString(fmt.Sprintf("AvgRowLength:\"%d\",", *i.AvgRowLength))
	}
	if i.DataLength != nil {
		buffer.WriteString(fmt.Sprintf("DataLength:\"%d\",", *i.DataLength))
	}
	if i.MaxDataLength != nil {
		buffer.WriteString(fmt.Sprintf("MaxDataLength:\"%d\",", *i.MaxDataLength))
	}
	if i.IndexLength != nil {
		buffer.WriteString(fmt.Sprintf("IndexLength:\"%d\",", *i.IndexLength))
	}
	if i.DataFree != nil {
		buffer.WriteString(fmt.Sprintf("DataFree:\"%d\",", *i.DataFree))
	}
	if i.AutoIncrement != nil {
		buffer.WriteString(fmt.Sprintf("AutoIncrement:\"%d\",", *i.AutoIncrement))
	}
	if i.CreateTime != nil {
		buffer.WriteString(fmt.Sprintf("CreateTime:\"%s\",", *i.CreateTime))
	}
	if i.UpdateTime != nil {
		buffer.WriteString(fmt.Sprintf("UpdateTime:\"%s\",", *i.UpdateTime))
	}
	if i.CheckTime != nil {
		buffer.WriteString(fmt.Sprintf("CheckTime:\"%s\",", *i.CheckTime))
	}
	if i.TableCollation != nil {
		buffer.WriteString(fmt.Sprintf("TableCollation:\"%s\",", *i.TableCollation))
	}
	if i.Checksum != nil {
		buffer.WriteString(fmt.Sprintf("Checksum:\"%s\",", *i.Checksum))
	}
	if i.CreateOptions != nil {
		buffer.WriteString(fmt.Sprintf("CreateOptions:\"%s\",", *i.CreateOptions))
	}
	if i.TableComment != nil {
		buffer.WriteString(fmt.Sprintf("TableComment:\"%s\",", *i.TableComment))
	}
	tempBytes := buffer.Bytes()
	if tempBytes[len(tempBytes)-1] == ',' {
		tempBytes = tempBytes[:len(tempBytes)-1]
	}
	tempBytes = append(tempBytes, '}')
	return string(tempBytes)
}

// 绑定表名
func (i *InfoColumns) TableName() string {
	return "COLUMNS"
}

// 打印对象
func (i *InfoColumns) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	if i.TableCatalog != nil {
		buffer.WriteString(fmt.Sprintf("TableCatalog:\"%s\",", *i.TableCatalog))
	}
	if i.TableSchema != nil {
		buffer.WriteString(fmt.Sprintf("TableSchema:\"%s\",", *i.TableSchema))
	}
	if i.TableNameStr != nil {
		buffer.WriteString(fmt.Sprintf("TableNameStr:\"%s\",", *i.TableNameStr))
	}
	if i.ColumnName != nil {
		buffer.WriteString(fmt.Sprintf("ColumnName:\"%s\",", *i.ColumnName))
	}
	if i.OrdinalPosition != nil {
		buffer.WriteString(fmt.Sprintf("OrdinalPosition:\"%d\",", *i.OrdinalPosition))
	}
	if i.ColumnDefault != nil {
		buffer.WriteString(fmt.Sprintf("ColumnDefault:\"%s\",", *i.ColumnDefault))
	}
	if i.IsNullable != nil {
		buffer.WriteString(fmt.Sprintf("IsNullable:\"%s\",", *i.IsNullable))
	}
	if i.DataType != nil {
		buffer.WriteString(fmt.Sprintf("DataType:\"%s\",", *i.DataType))
	}
	if i.CharacterMaximumLength != nil {
		buffer.WriteString(fmt.Sprintf("CharacterMaximumLength:\"%d\",", *i.CharacterMaximumLength))
	}
	if i.CharacterOctetLength != nil {
		buffer.WriteString(fmt.Sprintf("CharacterOctetLength:\"%d\",", *i.CharacterOctetLength))
	}
	if i.NumericPrecision != nil {
		buffer.WriteString(fmt.Sprintf("NumericPrecision:\"%d\",", *i.NumericPrecision))
	}
	if i.NumericScale != nil {
		buffer.WriteString(fmt.Sprintf("NumericScale:\"%d\",", *i.NumericScale))
	}
	if i.DatetimePrecision != nil {
		buffer.WriteString(fmt.Sprintf("DatetimePrecision:\"%d\",", *i.DatetimePrecision))
	}
	if i.CharacterSetName != nil {
		buffer.WriteString(fmt.Sprintf("CharacterSetName:\"%s\",", *i.CharacterSetName))
	}
	if i.CollationName != nil {
		buffer.WriteString(fmt.Sprintf("CollationName:\"%s\",", *i.CollationName))
	}
	if i.ColumnType != nil {
		buffer.WriteString(fmt.Sprintf("ColumnType:\"%s\",", *i.ColumnType))
	}
	if i.ColumnKey != nil {
		buffer.WriteString(fmt.Sprintf("ColumnKey:\"%s\",", *i.ColumnKey))
	}
	if i.Extra != nil {
		buffer.WriteString(fmt.Sprintf("Extra:\"%s\",", *i.Extra))
	}
	if i.Privileges != nil {
		buffer.WriteString(fmt.Sprintf("Privileges:\"%s\",", *i.Privileges))
	}
	if i.ColumnComment != nil {
		buffer.WriteString(fmt.Sprintf("ColumnComment:\"%s\",", *i.ColumnComment))
	}
	if i.GenerationExpression != nil {
		buffer.WriteString(fmt.Sprintf("GenerationExpression:\"%s\",", *i.GenerationExpression))
	}
	if i.SrsId != nil {
		buffer.WriteString(fmt.Sprintf("SrsId:\"%s\",", *i.SrsId))
	}
	tempBytes := buffer.Bytes()
	if tempBytes[len(tempBytes)-1] == ',' {
		tempBytes = tempBytes[:len(tempBytes)-1]
	}
	tempBytes = append(tempBytes, '}')
	return string(tempBytes)
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
	Summary     *SummaryConfig     `json:"summary"`
	TypeMapping *map[string]string `json:"typeMapping"`
}

// JsonConfig的格式化方式
func (j *JsonConfig) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	if j.Db != nil {
		buffer.WriteString("Db{")
		if j.Db.DbName != nil {
			buffer.WriteString(fmt.Sprintf("DbName:\"%s\",", *j.Db.DbName))
		}
		if j.Db.Host != nil {
			buffer.WriteString(fmt.Sprintf("Host:\"%s\",", *j.Db.Host))
		}
		if j.Db.Password != nil {
			buffer.WriteString(fmt.Sprintf("Password:\"%s\",", *j.Db.Password))
		}
		if j.Db.Port != nil {
			buffer.WriteString(fmt.Sprintf("Port:\"%s\",", *j.Db.Port))
		}
		if j.Db.TableName != nil {
			buffer.WriteString(fmt.Sprintf("TableName:\"%s\",", *j.Db.TableName))
		}
		if j.Db.User != nil {
			buffer.WriteString(fmt.Sprintf("User:\"%s\",", *j.Db.User))
		}
		tempBytes := buffer.Bytes()
		if tempBytes[len(tempBytes)-1] == ',' {
			tempBytes = tempBytes[:len(tempBytes)-1]
			buffer = *bytes.NewBuffer(tempBytes)
		}
		buffer.WriteString("},")
	}
	if j.Summary != nil {
		buffer.WriteString("Summary:{")
		if j.Summary.Author != nil {
			buffer.WriteString(fmt.Sprintf("Author:\"%s\",", *j.Summary.Author))
		}
		if j.Summary.Email != nil {
			buffer.WriteString(fmt.Sprintf("Email:\"%s\",", *j.Summary.Email))
		}
		if j.Summary.PackageName != nil {
			buffer.WriteString(fmt.Sprintf("PackageName:\"%s\",", *j.Summary.PackageName))
		}
		if j.Summary.TablePrefix != nil {
			buffer.WriteString(fmt.Sprintf("TablePrefix:\"%s\",", *j.Summary.TablePrefix))
		}
		tempBytes := buffer.Bytes()
		if tempBytes[len(tempBytes)-1] == ',' {
			tempBytes = tempBytes[:len(tempBytes)-1]
			buffer = *bytes.NewBuffer(tempBytes)
		}
		buffer.WriteString("},")
	}
	if j.TypeMapping != nil {
		buffer.WriteString(" TypeMapping:{")
		for key, value := range *j.TypeMapping {
			buffer.WriteString(fmt.Sprintf("%s:\"%s\",", key, value))
		}
		tempBytes := buffer.Bytes()
		if tempBytes[len(tempBytes)-1] == ',' {
			tempBytes = tempBytes[:len(tempBytes)-1]
			buffer = *bytes.NewBuffer(tempBytes)
		}
		buffer.WriteString("}")
	}
	buffer.WriteString("}")
	return buffer.String()
}

// 字段模板
type JavaPoField struct {
	ColumnComments string
	ColumnName     string
	FieldName      string
	DataType       string
}

// PO对象模板
type JavaPoTemplate struct {
	HasDate          bool
	HasBigDecimal    bool
	HasLocalDate     bool
	HasLocalDateTime bool
	PackageName      string
	TableComments    string
	Author           string
	Datetime         string
	TableName        string
	ClassName        string
	PoFields         []JavaPoField
	FilePathName     string
}

// 生成PoTemplate对象的方法
func (po *JavaPoTemplate) FillPoTemplate(tableData *InfoTable, columnsData *[]InfoColumns, config *JsonConfig) {
	if config.Summary != nil && config.Summary.PackageName != nil {
		po.PackageName = *config.Summary.PackageName
	}
	if tableData.TableComment != nil {
		po.TableComments = *tableData.TableComment
	}
	if config.Summary != nil && config.Summary.Author != nil {
		po.Author = *config.Summary.Author
	}
	if tableData.CreateTime != nil {
		timeStr := (*tableData.CreateTime).Format(TIME_FORMAT)
		po.Datetime = timeStr
	}
	if config.Db != nil && config.Db.TableName != nil {
		po.TableName = *config.Db.TableName
		className := UderscoreToUpperCamelCase(po.TableName)
		if className[0] >= 97 && className[0] <= 122 {
			className[0] = className[0] - 32
		}
		po.ClassName = string(className)
	}
	if columnsData != nil {
		for _, column := range *columnsData {
			var field JavaPoField
			if column.ColumnComment != nil {
				field.ColumnComments = *column.ColumnComment
			}
			if column.ColumnName != nil {
				field.ColumnName = *column.ColumnName
				field.FieldName = string(UderscoreToUpperCamelCase(*column.ColumnName))
			}
			if column.DataType != nil {
				dateType := *column.DataType
				field.DataType = (*config.TypeMapping)[dateType]
				switch strings.ToLower(field.DataType) {
				case "date":
					po.HasDate = true
				case "localdate":
					po.HasLocalDate = true
				case "localdatetime":
					po.HasLocalDateTime = true
				case "bigdecimal":
					po.HasBigDecimal = true
				}
			}
			po.PoFields = append(po.PoFields, field)
		}
	}
}

// PO的java对象
type PoJavaClass JavaPoTemplate

// 构建PO的java对象
func (po *PoJavaClass) FillPoJavaClass(p *JavaPoTemplate) {
	po.HasDate = p.HasDate
	po.HasBigDecimal = p.HasBigDecimal
	po.HasLocalDate = p.HasLocalDate
	po.HasLocalDateTime = p.HasLocalDateTime
	po.PackageName = p.PackageName + PO_PACKEAGE_SUFFIX
	po.TableName = p.TableName
	po.TableComments = p.TableComments
	po.Author = p.Author
	po.Datetime = p.Datetime
	po.ClassName = p.ClassName + PO_CLASS_SUFFIX
	po.PoFields = p.PoFields
	po.FilePathName = ModifyFileName(po.PackageName+"."+po.ClassName, ".java")
}

// DTO的java对象
type DtoJavaClass JavaPoTemplate

// 构建DTO的java对象
func (dto *DtoJavaClass) FillDtoJavaClass(p *JavaPoTemplate) {
	dto.HasDate = p.HasDate
	dto.HasBigDecimal = p.HasBigDecimal
	dto.HasLocalDate = p.HasLocalDate
	dto.HasLocalDateTime = p.HasLocalDateTime
	dto.PackageName = p.PackageName + DTO_PACKEAGE_SUFFIX
	dto.TableName = p.TableName
	dto.TableComments = p.TableComments
	dto.Author = p.Author
	dto.Datetime = p.Datetime
	dto.ClassName = p.ClassName + DTO_CLASS_SUFFIX
	dto.PoFields = p.PoFields
	dto.FilePathName = ModifyFileName(dto.PackageName+"."+dto.ClassName, ".java")
}

// mapper的java的接口对象
type MapperJavaClass struct {
	PackageName   string
	ClassName     string
	TableComments string
	TableName     string
	Author        string
	Datetime      string
	FilePathName  string
	PoPackageName string
	PoClassName   string
}

// 生成 mapper的java对象
func (i *MapperJavaClass) FillMapperJavaClass(p *JavaPoTemplate, po *PoJavaClass) {
	i.PackageName = p.PackageName + MAPPER_PACKEAGE_SUFFIX
	i.ClassName = p.ClassName + MAPPER_CLASS_SUFFIX
	i.TableComments = p.TableComments
	i.TableName = p.TableName
	i.Author = p.Author
	i.Datetime = p.Datetime
	i.FilePathName = ModifyFileName(i.PackageName+"."+i.ClassName, ".java")
	i.PoPackageName = po.PackageName
	i.PoClassName = po.ClassName
}

// mapper的xml的接口对象
type MapperXmlFile struct {
	MapperJavaClassPackageName string
	MapperJavaClassClassName   string
	ClassName                  string
	FilePathName               string
	PoPackageName              string
	PoClassName                string
	PoFields                   []JavaPoField
}

// 生成 mappper 的xml对象
func (m *MapperXmlFile) FillMapperXmlFile(p *JavaPoTemplate, po *PoJavaClass, i *MapperJavaClass) {
	m.MapperJavaClassPackageName = i.PackageName
	m.MapperJavaClassClassName = i.ClassName
	m.ClassName = p.ClassName
	m.FilePathName = ModifyFileName(m.MapperJavaClassPackageName+"."+m.MapperJavaClassClassName, ".xml")
	m.PoPackageName = po.PackageName
	m.PoClassName = po.ClassName
	m.PoFields = p.PoFields
}

// controller的Java对象
type ControllerJavaClass struct {
	PackageName    string
	ClassName      string
	TableComments  string
	TableName      string
	Author         string
	Datetime       string
	RequestMapping string
	FilePathName   string
}

// 生成controller的Java对象
func (c *ControllerJavaClass) FillControllerJavaClass(p *JavaPoTemplate) {
	c.PackageName = p.PackageName + CONTROLLER_PACKEAGE_SUFFIX
	c.ClassName = p.ClassName + CONTROLLER_CLASS_SUFFIX
	c.TableComments = p.TableComments
	c.TableName = p.TableName
	c.Author = p.Author
	c.Datetime = p.Datetime
	c.RequestMapping = p.ClassName
	c.FilePathName = ModifyFileName(c.PackageName+"."+c.ClassName, ".java")
}

// service的Java对象
type ServiceJavaClass struct {
	PackageName   string
	ClassName     string
	TableComments string
	TableName     string
	Author        string
	Datetime      string
	PoPackageName string
	PoClassName   string
	FilePathName  string
}

// 生成service的Java对象
func (s *ServiceJavaClass) FillServiceJavaClass(p *JavaPoTemplate, po *PoJavaClass) {
	s.PackageName = p.PackageName + SERVICE_PACKEAGE_SUFFIX
	s.ClassName = p.ClassName + SERVICE_CLASS_SUFFIX
	s.TableComments = p.TableComments
	s.TableName = p.TableName
	s.Author = p.Author
	s.Datetime = p.Datetime
	s.PoPackageName = po.PackageName
	s.PoClassName = po.ClassName
	s.FilePathName = ModifyFileName(s.PackageName+"."+s.ClassName, ".java")
}

// serviceImpl的Java对象
type ServiceImplJavaClass struct {
	PackageName                string
	ClassName                  string
	TableComments              string
	TableName                  string
	Author                     string
	Datetime                   string
	PoPackageName              string
	PoClassName                string
	MapperJavaClassPackageName string
	MapperJavaClassClassName   string
	ServicePackageName         string
	ServiceClassName           string
	FilePathName               string
}

// 生成serviceImpl的Java对象
func (si *ServiceImplJavaClass) FillServiceImplJavaClass(p *JavaPoTemplate, po *PoJavaClass, i *MapperJavaClass, s *ServiceJavaClass) {
	si.PackageName = p.PackageName + SERVICE_IMPL_PACKEAGE_SUFFIX
	si.ClassName = p.ClassName + SERVICE_IMPL_CLASS_SUFFIX
	si.TableComments = p.TableComments
	si.TableName = p.TableName
	si.Author = p.Author
	si.Datetime = p.Datetime
	si.PoPackageName = po.PackageName
	si.PoClassName = po.ClassName
	si.MapperJavaClassPackageName = i.PackageName
	si.MapperJavaClassClassName = i.ClassName
	si.ServicePackageName = s.PackageName
	si.ServiceClassName = s.ClassName
	si.FilePathName = ModifyFileName(si.PackageName+"."+si.ClassName, ".java")
}

// 保存内容用的dto
type TemplateJavaFile struct {
	FilePathName string
	content      string
}
