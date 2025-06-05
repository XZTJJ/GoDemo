package dbV2

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/magiconair/properties"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义个全局变量的logger方便打印日志
var logger *logService

// 自动执行的init方法
func init() {
	logger = initLogServic()
}

// 方法的入口类
func FacadeFunc() {
	//只是休眠用,方便显示
	defer func() { time.Sleep(120 * time.Second) }()
	//引导方法
	propertiesConfig, err := bootStrapInfo()
	if err != nil {
		return
	}
	propertiesClass, err2 := parseConfig(propertiesConfig)
	if err2 != nil {
		return
	}
	//连接数据库配置信息
	mysqlConnect, err3 := getMySQLConnect(propertiesClass)
	if err3 != nil {
		return
	}
	//处理给定的所哟表
	allFiles, err4 := hanldeTalbles(propertiesClass, mysqlConnect)
	if err4 != nil {
		return
	}
	//生产zip包
	zipfileName, err5 := makeZipFile(allFiles)
	if err5 != nil {
		return
	}
	logger.logInfo("所需要的文件已经生在在当前目录下,文件名为:" + zipfileName)
}

// 生成zip包
func makeZipFile(allFiles *[]generateTemplateFile) (string, error) {
	//创建zip包
	zipFileName := "codeGenerate-" + string(time.Now().UnixMilli()) + ".zip"
	out, err := os.Create(zipFileName)
	if err != nil {
		logger.logError("创建"+zipFileName+"失败", err)
		return "", err
	}
	defer out.Close()
	//往zip包中逐一写入文件
	writer := zip.NewWriter(out)
	for _, file := range *allFiles {
		fileWriter, err2 := writer.Create(file.fullFileName)
		if err2 != nil {
			logger.logError("往"+zipFileName+"中添加"+file.fullFileName+"失败", err2)
			return "", err2
		}
		fileWriter.Write([]byte(file.content))
	}
	//生产zip包
	if err3 := writer.Close(); err3 != nil {
		logger.logError("生成"+zipFileName+"失败", err3)
		return "", err3
	}
	return zipFileName, nil
}

// 处理给定的表
func hanldeTalbles(p *propertiesClass, mysqlConnect *gorm.DB) (*[]generateTemplateFile, error) {
	//对每个表单独处理
	var allFiles []generateTemplateFile
	for _, table := range p.tables {
		logger.logInfo("开始处理" + table + "表")
		tableDefine := getTableDefine(mysqlConnect, p.dbName, table)
		if tableDefine == nil {
			logger.logInfo(table + "不存在的表,直接跳过该表")
			continue
		}
		//转变的类名
		className := sqlName2CodeName(table, p.tablePrefixes, p.tableSuffixes,
			p.tableStrategy, p.tableCamelSeparatorMap, p.tableSameSeparators)
		//处理给定表的字段
		logger.logInfo("开始处理" + table + "表的字段")
		//获取字段定义
		columnsMySQLs := getColumsDefine(mysqlConnect, p.dbName, table)
		//字段名转成代码的需要的信息，比如属性名等
		codeColumns := hanldeColumns(columnsMySQLs, p)
		logger.logInfo("开始渲染" + table + "表的所有模板")
		templateData := assembleTemplateData(p, table, *tableDefine.TableComment, className, codeColumns)
		fmt.Printf("渲染出来的数据为:%+v\n", *templateData)
		logger.logInfo("开始生产" + table + "表的所有模板文件")
		files, err := renderingTemplate(templateData, p, className)
		if err != nil {
			logger.logError("生成"+table+"文件错误", err)
			return &allFiles, err
		}
		allFiles = append(allFiles, (*files)...)
	}
	return &allFiles, nil
}

// 开始渲染模板
func renderingTemplate(templateData *map[string]any, p *propertiesClass, className string) (*[]generateTemplateFile, error) {
	var files []generateTemplateFile
	//遍历所有的模板进行处理
	for key, value := range p.templates {
		var err error
		var tmpl1 *template.Template
		switch key {
		//处理po的模板
		case PO_TEMPLATE_DEFAULT:
			//使用使用默认模板
			if value.templateStr == EMPTY_STRING {
				tmpl1, err = template.New(key + "Template").Parse(showDefaultPOTemplate())
			} else {
				tmpl1, err = template.ParseFiles(value.templateStr)
			}
			//处理mybatisplusJava的模板
		case MYBATISPLUSJAVA_TEMPLATE_DEFAULT:
			//使用使用默认模板
			if value.templateStr == EMPTY_STRING {
				tmpl1, err = template.New(key + "Template").Parse(showDefaultMybatisplusJavaTemplate())
			} else {
				tmpl1, err = template.ParseFiles(value.templateStr)
			}
			//处理mybatisplusXml的模板
		case MYBATISPLUSXML_TEMPLATE_DEFAULT:
			//使用使用默认模板
			if value.templateStr == EMPTY_STRING {
				tmpl1, err = template.New(key + "Template").Parse(showDefaultMybatisplusXmlTemplate())
			} else {
				tmpl1, err = template.ParseFiles(value.templateStr)
			}
			//默认其他都是从模板文件加载
		default:
			tmpl1, err = template.ParseFiles(value.templateStr)
		}
		//判断是否存在异常
		if err != nil {
			logger.logError(key+"类型模板加载失败", err)
			return &files, err
		}
		//保存渲染出来的模板字节，转成字符
		var tmplBytes bytes.Buffer
		//渲染模板
		err = tmpl1.Execute(&tmplBytes, *templateData)
		if err != nil {
			logger.logError(key+"类型模板处理失败", err)
			return &files, err
		}
		//将内容转成字符
		fileContent := tmplBytes.String()
		//转换成全路径类名
		fullFileName := strings.ReplaceAll(value.packageStr+".", ".", string(os.PathSeparator))
		fullFileName = fullFileName + string(os.PathSeparator) + value.classNamePrefix +
			className + value.classNameSuffix + value.fileTypeStr
		file := generateTemplateFile{fullFileName: fullFileName, content: fileContent}
		files = append(files, file)
	}
	return &files, nil
}

// 准备模板需要的所有数据
func assembleTemplateData(p *propertiesClass, tableName, tableComment, className string, codeColumns *[]columnInfos) *map[string]any {
	dataMap := make(map[string]any)
	dataMap["author"] = p.authorStr
	dataMap["version"] = p.versionStr
	dataMap["desc"] = p.descStr
	dataMap["datetime"] = p.datetimeStr
	//模板相关信息
	for key, value := range p.templates {
		dataMap[key+"Package"] = value.packageStr
		dataMap[key+"Template"] = value.templateStr
		dataMap[key+"ClassNamePrefix"] = value.classNamePrefix
		dataMap[key+"ClassNameSuffix"] = value.classNameSuffix
		dataMap[key+"FileType"] = value.fileTypeStr
	}
	dataMap["tableComment"] = tableComment
	dataMap["tableSQLName"] = tableName
	dataMap["tableJavaName"] = className
	dataMap["columnInfos"] = *codeColumns
	return &dataMap
}

// 根据字段结构处理表的字段
func hanldeColumns(cols *[]columnsMySQL, p *propertiesClass) *[]columnInfos {
	var columns []columnInfos
	for _, col := range *cols {
		comment := *col.ColumnComment
		sqlName := *col.ColumnName
		sqlDateType := *col.DataType
		codeDateType := p.sqlCodeType[sqlDateType]
		codeFieldName := sqlName2CodeName(sqlName, p.fieldPrefixes, p.fieldSuffixes, p.fieldStrategy,
			p.fieldCamelSeparatorMap, p.fieldSameSeparators)
		column := columnInfos{columnComment: comment, columnSQLName: sqlName,
			columnJavaName: codeFieldName, columnSQLType: sqlDateType, columnJavaType: codeDateType}
		columns = append(columns, column)
	}
	return &columns
}

// 获取字段定义
func getColumsDefine(connect *gorm.DB, dbName, tableName string) *[]columnsMySQL {
	var results []columnsMySQL
	connect.Where(" table_schema = ? AND table_name = ? ", dbName, tableName).Order("ordinal_position").Find(&results)
	return &results
}

// 将数据库名称或者字段名称按照规则转变成代码名称
func sqlName2CodeName(sqlName string, prefiexArray []string, suffixesArray []string, strategy string,
	camelSeparatorSet map[rune]struct{}, sameSeparatorArray []string) string {
	var codeName = sqlName
	//替换前缀
	for _, prefiex := range prefiexArray {
		if prefiex == EMPTY_STRING {
			continue
		}
		codeName = strings.TrimPrefix(codeName, prefiex)
	}
	//替换后缀
	for _, suffix := range suffixesArray {
		if suffix == EMPTY_STRING {
			continue
		}
		codeName = strings.TrimSuffix(codeName, suffix)
	}
	//其他数据转换规则
	switch strategy {
	// 驼峰写法的转换规则
	case STRATEGY_CAMEL:
		var resultName []rune
		//首字符大写,默认为true
		nextUpper := true
		for _, r := range codeName {
			//是否包含某个分隔符
			_, ok := camelSeparatorSet[r]
			if ok {
				nextUpper = true
				continue
			}
			//非分隔符的处理逻辑
			if nextUpper {
				resultName = append(resultName, unicode.ToUpper(r))
				nextUpper = false
			} else {
				resultName = append(resultName, r)
			}
		}
		codeName = string(resultName)
	// 原始名称，不做任何修改，只是去掉分隔符
	case STRATEGY_SAME:
		for _, sameSeparator := range sameSeparatorArray {
			if sameSeparator == EMPTY_STRING {
				continue
			}
			codeName = strings.ReplaceAll(codeName, sameSeparator, EMPTY_STRING)
		}
	}
	return codeName
}

// 获取数据库中表定义相关信息
func getTableDefine(mySQLConnect *gorm.DB, dbName, tableName string) *tableMySQL {
	var tableDefine tableMySQL
	mySQLConnect.Select("table_name, table_comment, create_time").
		Where(" table_schema = ? AND table_name = ? ", dbName, tableName).Find(&tableDefine)
	return &tableDefine
}

// 初始化数据库连接
func getMySQLConnect(p *propertiesClass) (*gorm.DB, error) {
	//连接mysql,RL编码处理特殊字符
	openUrl := url.QueryEscape(p.dbUser) + ":" + url.QueryEscape(p.dbPassword) + "@tcp(" + p.dbHost + ":" + p.dbPort + ")/" +
		url.QueryEscape(DB_INFORMATION_SCHEMA) + "?charset=utf8mb4&parseTime=True&loc=Local"
	connect, err := gorm.Open(mysql.Open(openUrl))
	if err != nil {
		logger.logError("数据库连接失败，请检查配置信息", err)
		return nil, err
	}
	logger.logInfo("数据库连接成功")
	return connect, nil
}

// 根据properties文件的配置项转成struct
func parseConfig(p *properties.Properties) (*propertiesClass, error) {
	var propertiesClass = propertiesClass{}
	var err error
	//数据库信息解析
	propertiesClass.dbHost, err = getValueFromProperties(p, BD_HOST, "数据库地址没有配置,请配置数据库信息")
	if err != nil {
		return &propertiesClass, err
	}
	propertiesClass.dbPort, err = getValueFromProperties(p, BD_PORT, "数据库端口没有配置,请配置数据库信息")
	if err != nil {
		return &propertiesClass, err
	}
	propertiesClass.dbUser, err = getValueFromProperties(p, BD_USER, "数据库用户没有配置,请配置数据库信息")
	if err != nil {
		return &propertiesClass, err
	}
	propertiesClass.dbPassword, err = getValueFromProperties(p, BD_PASSWORD, "数据库密码没有配置,请配置数据库信息")
	if err != nil {
		return &propertiesClass, err
	}
	propertiesClass.dbName, err = getValueFromProperties(p, BD_NAME, "没有指定数据库名称,请配置数据库名称")
	if err != nil {
		return &propertiesClass, err
	}
	var tables string
	tables, err = getValueFromProperties(p, HD_TABLES, "没有指定表名称,请配置表名称")
	if err != nil {
		return &propertiesClass, err
	}
	propertiesClass.tables = strings.Split(tables, SEPARATOR_DEFAULT)
	//表的转换规则
	propertiesClass.tablePrefixes = strings.Split(p.GetString(HD_IG_T_PREFIX, EMPTY_STRING), SEPARATOR_DEFAULT)
	propertiesClass.tableSuffixes = strings.Split(p.GetString(HD_IG_T_SUFFIX, EMPTY_STRING), SEPARATOR_DEFAULT)
	propertiesClass.tableStrategy = p.GetString(HD_T_STRATEGY, T_STRATEGY_DEFAULT)
	propertiesClass.tableSameSeparators = strings.Split(p.GetString(HD_F_SAME_SEPARATOR, EMPTY_STRING), SEPARATOR_DEFAULT)
	//驼峰写法的特殊处理逻辑，将分割符转变成map
	tableCamelSeparatorSet := make(map[rune]struct{})
	for _, camelSeparator := range strings.Split(p.GetString(HD_T_CAMEL_SEPARATOR, T_CAMEL_SEPARATOR_DEFAULT), SEPARATOR_DEFAULT) {
		if camelSeparator == EMPTY_STRING { // 过滤空字符串
			continue
		}
		//拆分成每个字符
		runes := []rune(camelSeparator)
		if len(runes) > 1 {
			return &propertiesClass, errors.New("分隔符只支持单个字符")
		}
		for _, runeStr := range runes {
			tableCamelSeparatorSet[runeStr] = struct{}{}
		}
	}
	propertiesClass.tableCamelSeparatorMap = tableCamelSeparatorSet

	//字段的转换规则
	propertiesClass.fieldPrefixes = strings.Split(p.GetString(HD_IG_F_PREIFX, EMPTY_STRING), SEPARATOR_DEFAULT)
	propertiesClass.fieldSuffixes = strings.Split(p.GetString(HD_IG_F_SUFFIX, EMPTY_STRING), SEPARATOR_DEFAULT)
	propertiesClass.fieldStrategy = p.GetString(HD_F_STRATEGY, F_STRATEGY_DEFAULT)
	propertiesClass.fieldSameSeparators = strings.Split(p.GetString(HD_F_SAME_SEPARATOR, EMPTY_STRING), SEPARATOR_DEFAULT)
	//驼峰写法的特殊处理逻辑，将分割符转变成map
	fieldCamelSeparatorSet := make(map[rune]struct{})
	for _, camelSeparator := range strings.Split(p.GetString(HD_F_CAMEL_SEPARATOR, F_CAMEL_SEPARATOR_DEFAULT), SEPARATOR_DEFAULT) {
		if camelSeparator == EMPTY_STRING { // 过滤空字符串
			continue
		}
		//拆分成每个字符
		runes := []rune(camelSeparator)
		if len(runes) > 1 {
			return &propertiesClass, errors.New("分隔符只支持单个字符")
		}
		for _, runeStr := range runes {
			fieldCamelSeparatorSet[runeStr] = struct{}{}
		}
	}
	propertiesClass.fieldCamelSeparatorMap = fieldCamelSeparatorSet

	//公共创作者信息
	propertiesClass.authorStr = p.GetString(J_C_AUTHOR, EMPTY_STRING)
	propertiesClass.versionStr = p.GetString(J_C_VERSION, EMPTY_STRING)
	propertiesClass.descStr = p.GetString(J_C_DESC, EMPTY_STRING)
	propertiesClass.datetimeStr = p.GetString(J_C_DATETIME, EMPTY_STRING)

	//获取所有的Key
	configMap := p.Map()
	//自定义模板信息,遍历
	ctemplateMap := make(map[string]customeTemplate)
	ctemplateMap[PO_TEMPLATE_DEFAULT] = customeTemplate{classNameSuffix: "PO", fileTypeStr: ".java"}
	ctemplateMap[MYBATISPLUSJAVA_TEMPLATE_DEFAULT] = customeTemplate{classNameSuffix: "Mapper", fileTypeStr: ".java"}
	ctemplateMap[MYBATISPLUSXML_TEMPLATE_DEFAULT] = customeTemplate{classNameSuffix: "Mapper", fileTypeStr: ".xml"}
	for key, value := range configMap {
		//非java.template.开头的直接跳过
		if !strings.HasPrefix(key, J_TEMPLATE_TYPE_PREFIX) {
			continue
		}
		//模板文件类型
		typeStr := strings.TrimPrefix(key, J_TEMPLATE_TYPE_PREFIX)
		switch {
		case strings.HasSuffix(key, J_TEMPLATE_TYPE_SUFFIX_PACKGE):
			//处理模板的package信息
			typeStr = strings.TrimSuffix(typeStr, J_TEMPLATE_TYPE_PREFIX)
			ctemplate, ok := ctemplateMap[typeStr]
			if !ok {
				ctemplateMap[typeStr] = customeTemplate{packageStr: value}
			} else {
				ctemplate.packageStr = value
			}
		case strings.HasSuffix(key, J_TEMPLATE_TYPE_SUFFIX_TEMPLATE):
			//处理模板的template信息
			typeStr = strings.TrimSuffix(typeStr, J_TEMPLATE_TYPE_SUFFIX_TEMPLATE)
			ctemplate, ok := ctemplateMap[typeStr]
			if !ok {
				ctemplateMap[typeStr] = customeTemplate{templateStr: value}
			} else {
				ctemplate.templateStr = value
			}
		case strings.HasSuffix(key, J_TEMPLATE_TYPE_SUFFIX_CLASSNAMEPREFIX):
			//处理模板的classNamePrefix信息
			typeStr = strings.TrimSuffix(typeStr, J_TEMPLATE_TYPE_SUFFIX_CLASSNAMEPREFIX)
			ctemplate, ok := ctemplateMap[typeStr]
			if !ok {
				ctemplateMap[typeStr] = customeTemplate{classNamePrefix: value}
			} else {
				ctemplate.classNamePrefix = value
			}
		case strings.HasSuffix(key, J_TEMPLATE_TYPE_SUFFIX_CLASSNAMESUFFIX):
			//处理模板的classNameSuffix信息
			typeStr = strings.TrimSuffix(typeStr, J_TEMPLATE_TYPE_SUFFIX_CLASSNAMESUFFIX)
			ctemplate, ok := ctemplateMap[typeStr]
			if !ok {
				ctemplateMap[typeStr] = customeTemplate{classNameSuffix: value}
			} else {
				ctemplate.classNameSuffix = value
			}
		case strings.HasSuffix(key, J_TEMPLATE_TYPE_SUFFIX_FILETYPE):
			//处理模板的fileType信息
			typeStr = strings.TrimSuffix(typeStr, J_TEMPLATE_TYPE_SUFFIX_FILETYPE)
			ctemplate, ok := ctemplateMap[typeStr]
			if !ok {
				ctemplateMap[typeStr] = customeTemplate{fileTypeStr: value}
			} else {
				ctemplate.fileTypeStr = value
			}
		}
	}
	//检查模板的必填项
	for key, value := range ctemplateMap {
		if value.templateStr == "" && key != PO_TEMPLATE_DEFAULT && key != MYBATISPLUSJAVA_TEMPLATE_DEFAULT && key != MYBATISPLUSXML_TEMPLATE_DEFAULT {
			return &propertiesClass, errors.New(key + "类型的模板文件不存在，请检查")
		}
	}
	propertiesClass.templates = ctemplateMap

	//初始化SQL类型默认值
	sqlCodeType := make(map[string]string)
	defualtTypeArrays := []string{"bigint,Long", "bit,Boolean", "char,String",
		"date,LocalDate", "datetime,LocalDateTime", "decimal,BigDecimal",
		"double,Double", "float,Float", "int,Integer", "integer,Integer", "longtext,String",
		"mediumint,Integer", "mediumtext,String", "smallint,Integer", "text,String", "timestamp,LocalDateTime",
		"tinyint,Byte", "tinytext,String", "varchar,String"}
	for _, value := range defualtTypeArrays {
		parts := strings.Split(value, SEPARATOR_DEFAULT)
		sqlCodeType[parts[0]] = parts[1]
	}
	for key, value := range configMap {
		//跳过配置
		if !strings.HasPrefix(key, S_TYPE_PREFIX) || !strings.HasSuffix(key, S_TYPE_SUFFIX) {
			continue
		}
		sqlType := strings.TrimSuffix(strings.TrimPrefix(key, S_TYPE_PREFIX), S_TYPE_SUFFIX)
		sqlCodeType[sqlType] = value
	}
	propertiesClass.sqlCodeType = sqlCodeType
	logger.logInfo("提取出来的配置项为:" + fmt.Sprintf("%+v", propertiesClass))
	return &propertiesClass, nil
}

// 从properties中获取key的value
func getValueFromProperties(p *properties.Properties, key string, tip string) (string, error) {
	if val, ok := p.Get(key); ok {
		return val, nil
	} else {
		err := errors.New(tip)
		logger.logError(tip, err)
		return "", err
	}
}

// 交互引导方法
func bootStrapInfo() (*properties.Properties, error) {
	reader := bufio.NewReader(os.Stdin)
	helpInfo := `1.查看示例的properties配置文件请输入P,
2.查看示例模板文件请输入T,
3.直接读取当前目录下面的configApp.properties请出入Y，
4.读取其他目录下的配置文件请输入配置文件绝对路径
按回车键完成输入`

	//支持循环提示
	for {
		fmt.Println(helpInfo)
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.logError("读取输入异常，请重新输入", err)
			return nil, err
		}
		input = strings.TrimSpace(input)
		//指引判断
		if input == "" {
			fmt.Println("输入为空，请重新输入")
		} else if input == "P" || input == "p" {
			fmt.Printf("##################################以下是配置文件示例################################：\n")
			fmt.Println(showDefaultProperties())
		} else if input == "T" || input == "t" {
			fmt.Printf("#################################以下是模板可使用的变量以及使用方法示例#################################\n")
			fmt.Println(showDefaultVariable())
			fmt.Printf("#################################以下是po模板示例#################################\n")
			fmt.Println(showDefaultPOTemplate())
			fmt.Printf("#################################以下是mybatisplusJava模板示例#################################\n")
			fmt.Println(showDefaultMybatisplusJavaTemplate())
			fmt.Printf("#################################以下是mybatisplusXml模板示例#################################\n")
			fmt.Println(showDefaultMybatisplusXmlTemplate())
		} else if input == "Y" || input == "y" {
			p, err := loadProperties(DEFAULT_FILE_NAME)
			if err != nil {
				logger.logError("读取"+DEFAULT_FILE_NAME+"配置文件失败", err)
				return nil, err
			}
			logger.logInfo(DEFAULT_FILE_NAME + "properties文件读取成功")
			return p, nil
		} else {
			p, err := loadProperties(input)
			if err != nil {
				logger.logError("读取"+input+"配置文件失败", err)
				return nil, err
			}
			logger.logInfo(DEFAULT_FILE_NAME + "properties文件读取成功")
			return p, nil
		}
	}
}

// 常量类，对应配置文件
const (
	//配置文件
	DEFAULT_FILE_NAME = "configApp.properties"
	//数据库
	BD_HOST               = "database.host"
	BD_PORT               = "database.port"
	BD_NAME               = "database.dbName"
	BD_USER               = "database.user"
	BD_PASSWORD           = "database.password"
	DB_INFORMATION_SCHEMA = "information_schema"
	//表明和字段名处理信息
	HD_TABLES            = "handle.toCode.tables"
	HD_IG_T_PREFIX       = "handle.ignore.tablePrefix"
	HD_IG_T_SUFFIX       = "handle.ignore.tableSuffix"
	HD_T_STRATEGY        = "handle.classNameStrategy"
	HD_T_CAMEL_SEPARATOR = "handle.classNameStrategy.camel.separator"
	HD_T_SAME_SEPARATOR  = "handle.classNameStrategy.same.separator"
	HD_IG_F_PREIFX       = "handle.ignore.fieldPrefix"
	HD_IG_F_SUFFIX       = "handle.ignore.fieldSuffix"
	HD_F_STRATEGY        = "handle.fieldNameStrategy"
	HD_F_CAMEL_SEPARATOR = "handle.fieldNameStrategy.camel.separator"
	HD_F_SAME_SEPARATOR  = "handle.fieldNameStrategy.same.separator"
	//创作者公共信息
	J_C_AUTHOR   = "java.common.author"
	J_C_VERSION  = "java.common.version"
	J_C_DESC     = "java.common.desc"
	J_C_DATETIME = "java.common.datetime"
	//自定义模板类型
	J_TEMPLATE_TYPE_PREFIX                 = "java.template."
	J_TEMPLATE_TYPE_SUFFIX_PACKGE          = ".package"
	J_TEMPLATE_TYPE_SUFFIX_TEMPLATE        = ".template"
	J_TEMPLATE_TYPE_SUFFIX_CLASSNAMEPREFIX = ".classNamePrefix"
	J_TEMPLATE_TYPE_SUFFIX_CLASSNAMESUFFIX = ".classNameSuffix"
	J_TEMPLATE_TYPE_SUFFIX_FILETYPE        = ".fileType"
	//自定义类型
	S_TYPE_PREFIX = "sql.type."
	S_TYPE_SUFFIX = ".toJava"
	//一些默认值
	T_STRATEGY_DEFAULT        = "camel"
	T_CAMEL_SEPARATOR_DEFAULT = "-,_"
	EMPTY_STRING              = ""
	F_STRATEGY_DEFAULT        = "camel"
	F_CAMEL_SEPARATOR_DEFAULT = "-,_"
	SEPARATOR_DEFAULT         = ","
	//转换规则的两个枚举类
	STRATEGY_CAMEL = "camel"
	STRATEGY_SAME  = "same"
	//自定义模板类别名称
	PO_TEMPLATE_DEFAULT              = "po"
	MYBATISPLUSJAVA_TEMPLATE_DEFAULT = "mybatisplusJava"
	MYBATISPLUSXML_TEMPLATE_DEFAULT  = "mybatisplusXml"
)

// 读取配置文件，判断是否为空
func loadProperties(fileName string) (*properties.Properties, error) {
	if fileExists(fileName) {
		// 读取 .properties 文件
		p, err := properties.LoadFile(fileName, properties.UTF8)
		if err != nil {
			logger.logError("读取"+fileName+"文件异常，请确定当前配置项是否正确", err)
			return nil, err
		}
		return p, nil
	} else {
		err2 := errors.New("查找文件失败")
		logger.logError("当前目录不存在"+fileName+"文件请确认", err2)
		return nil, err2
	}
}

// 判断文件是否存在
func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		logger.logError(fileName+"查找失败", err)
	}
	return !os.IsNotExist(err)
}

/*************************配置文件和模板示例开始********************************/
//默认配置文件
func showDefaultProperties() string {
	propertiesStr := `
#配置信息采用properties格式,默认读取执行文件目录下面的configGoPro.properties执行文件信息

#MySQL数据库的地址，必填
database.host=
#MySQL数据库的port，必填
database.port=
#MySQL数据库的名称，必填
database.dbName=
#MySQL数据库的用户名，必填
database.user=
#MySQL数据库的账号密码，必填
database.password=

#需要生成java代码的数据库表明,通过,分割多个表名，必填
handle.toCode.tables=
#生成java代码需要忽略调的表名的前缀,按照指定的顺序处理前缀,通过,分割多个前缀
handle.ignore.tablePrefix=
#生成java代码需要忽略调的表名的后缀,按照指定的顺序处理后缀,通过,分割多个前缀
handle.ignore.tableSuffix=
#生成的java的类名策略: camel(驼峰形式), same(保持表名原样不处理)
handle.classNameStrategy=camel
#生成的java的类名如果采用camel(驼峰形式),需要处理的分隔符，通过,分割多个前缀，这些符号会被忽略，并且将后面字段改成大写
handle.classNameStrategy.camel.separator=-,_
#生成的java的类名如果采用same(保持表名原样不处理),需要处理的分隔符，通过,分割多个前缀，这些符号会被忽略
handle.classNameStrategy.same.separator=
#生成java代码需要忽略调的字段的前缀,按照指定的顺序处理前缀,通过,分割多个前缀
handle.ignore.fieldPrefix=
#生成java代码需要忽略调的字段的后缀,按照指定的顺序处理后缀,通过,分割多个前缀
handle.ignore.fieldSuffix=
#生成的java的字段策略: camel(驼峰形式), same(保持表名原样不处理)
handle.fieldNameStrategy=camel
#生成的java的字段如果采用camel(驼峰形式),需要处理的分隔符，通过,分割多个前缀，这些符号会被忽略，并且将后面字段改成大写
handle.fieldNameStrategy.camel.separator=-,_
#生成的java的字段如果采用same(保持表名原样不处理),需要处理的分隔符，通过,分割多个前缀，这些符号会被忽略
handle.fieldNameStrategy.same.separator=


#公共属性配置,作者
java.common.author=
#公共属性配置,版本号
java.common.version=
#公共属性配置,描述
java.common.desc=
#公共属性配置,日期
java.common.datetime=

#需要生产哪些类，支持自定义，只需要按照下面格式配置就行，支持属性覆盖，
#   java.template.xx.package=
#   java.template.xx.template=
#   java.template.xx.classNamePrefix=
#   java.template.xx.classNameSuffix=
#   java.template.xx.fileType=
#xx表示类型，默认会生产po,mybatisplusJava，mybatisplusXml这3种类型的文件,
# package表示包名, template表示模板文件的绝对路径，classNamePrefix表示class类的前缀，
# classNameSuffix表示class类的后缀，fileType表示文件类型
# 下面是默认的3种类型的2个默认属性
# java.template.po.classNameSuffix=PO     java.template.po.fileType=.java
# java.template.mybatisplusJava.classNameSuffix=Mapper     java.template.mybatisplusJava.fileType=.java
# java.template.mybatisplusXml.classNameSuffix=Mapper     java.template.mybatisplusXml.fileType=.xml
# 

#需要将Mysql的类型映射成java的那种类型,支持自定义，只需要按照下面格式配置就行，支持属性覆盖
#  sql.type.xx1.toJava=xx2
#xx1表示数据库类型，xx2表示java类型
#默认会将会进行如下转换
#   sqltype    -   javatype
#   bigint  -   Long
#   bit -   Boolean
#   char    -   String
#   date    -   LocalDate
#   datetime    -   LocalDateTime
#   decimal -   BigDecimal
#   double  -   Double
#   float   -   Float
#   int -   Integer
#   integer -   Integer
#   longtext    -   String
#   mediumint   -   Integer
#   mediumtext  -   String
#   smallint    -   Integer
#   text    -   String
#   timestamp   -   LocalDateTime
#   tinyint -   Byte
#   tinytext    -   String
#   varchar -   String	
	`
	return propertiesStr
}

// 默认变量
func showDefaultVariable() string {
	variables := `
author：作者，对应配置java.common.author
version：版本号，对应配置java.common.version
desc：描述，对应配置java.common.desc
datetime:日期，对应配置java.common.datetime
xxPackage: 类型的包名，xx要改成对应的类型，对应配置java.template.xx.package
xxTemplate: 类型的模板文件，xx要改成对应的类型，对应配置java.template.xx.template
xxClassNamePrefix: 类型的前缀，xx要改成对应的类型，对应配置java.template.xx.classNamePrefix
xxClassNameSuffix: 类型的后缀，xx要改成对应的类型，对应配置java.template.xx.classNameSuffix
xxFileType: 类型的文件后缀，xx要改成对应的类型，对应配置java.template.xx.fileType
tableComment：数据库中表的描述说明
tableSQLName: 数据库中表名
tableJavaName: 数据库中表名转变的代码名称
columnInfos: 数据库中表的字段信息容器，遍历用，不能单独使用
columnComment: 数据库中表的字段的注解，在columnInfos的遍历中使用
columnSQLName：数据库中表的字段的名称，在columnInfos的遍历中使用
columnJavaName：sql中的字段名转变的代码名称，在columnInfos的遍历中使用
columnSQLType：数据库中表的字段的类型，在columnInfos的遍历中使用
columnJavaType：sql类型映射的代码类型，在columnInfos的遍历中使用

通过一个遍历columnInfos类介绍变量使用
{{range $index, $value := .columnInfos}}
  字段名：{{$value.columnSQLName}},字段注解:{{$value.columnComment}}
{{end}}
	`
	return variables
}

// 默认po模板
func showDefaultPOTemplate() string {
	pOTemplate := `
package {{.poPackage}};

import lombok.Data;

/**
 * {{.tableName}} 对应的Po类型,
 * {{.tableComment}}
 *
 * @author {{.author}}
 * @version {{.version}}
 * @desc {{.desc}}
 * @date {{.datetime}}
 */  
public class {{.poClassNamePrefix}}{{.tableJavaName}}{{.poClassNameSuffix}} implements Serializable {
	private static final long serialVersionUID = 1L;
	{{range $index, $value := .columnInfos}}
	/**
	 * {{$value.columnComment}}
	 */
	private {{$value.columnJavaType}} {{$value.columnJavaName}};
	{{end}}
}
	`
	return pOTemplate
}

// 默认mybatisplusJava模板
func showDefaultMybatisplusJavaTemplate() string {
	mybatisplusJavaTemplate := `
package {{.mybatisplusJavaPackage}};

import {{.poPackage}}.{{.poClassNamePrefix}}{{.tableJavaName}}{{.poClassNameSuffix}};
import com.baomidou.mybatisplus.core.mapper.BaseMapper;

/**
 * {{.tableName}} 对应的mybatisplus操作类，
 * {{.tableComment}}
 *
 * @author {{.author}}
 * @version {{.version}}
 * @desc {{.desc}}
 * @date {{.datetime}}
 */
public interface {{.mybatisplusClassNamePrefix}}{{.tableJavaName}}{{.mybatisplusClassNameSuffix}} extends BaseMapper<{{.poClassNamePrefix}}{{.tableJavaName}}{{.poClassNameSuffix}}> {
	
}
	`
	return mybatisplusJavaTemplate
}

// 默认mybatisplusJava模板
func showDefaultMybatisplusXmlTemplate() string {
	mybatisplusXmlTemplate := `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">

<mapper namespace="{{.mybatisplusJavaPackage}}.{{.mybatisplusClassNamePrefix}}{{.tableJavaName}}{{.mybatisplusClassNameSuffix}}">

	<!-- 可根据自己的需求，是否要使用 -->
	<resultMap type="{{.poPackage}}.{{.poClassNamePrefix}}{{.tableJavaName}}{{.poClassNameSuffix}}" id="{{.poPackage}}.{{.poClassNamePrefix}}{{.tableJavaName}}{{.poClassNameSuffix}}Map">
	{{range $index, $value := .columnInfos}}
		<result property="{{$value.columnJavaName}}" column="{{$value.columnSQLName}}"/>
	{{end}}
	</resultMap>

</mapper>
`
	return mybatisplusXmlTemplate
}

/*************************配置文件和模板示例结束********************************/

/*************************构建模板用到的结构体开始*********************************/
// 字段整个相关的结构体
type columnInfos struct {
	columnComment  string
	columnSQLName  string
	columnJavaName string
	columnSQLType  string
	columnJavaType string
}

// 根据模板生成文件持有的信息
type generateTemplateFile struct {
	fullFileName string
	content      string
}

// 对应的配置文件的解析结构体
type propertiesClass struct {
	dbHost                 string
	dbPort                 string
	dbUser                 string
	dbPassword             string
	dbName                 string
	tables                 []string
	tablePrefixes          []string
	tableSuffixes          []string
	tableStrategy          string
	tableCamelSeparatorMap map[rune]struct{}
	tableSameSeparators    []string
	authorStr              string
	versionStr             string
	descStr                string
	datetimeStr            string
	fieldPrefixes          []string
	fieldSuffixes          []string
	fieldStrategy          string
	fieldCamelSeparatorMap map[rune]struct{}
	fieldSameSeparators    []string
	templates              map[string]customeTemplate
	sqlCodeType            map[string]string
}

// 自定义的模板类型
type customeTemplate struct {
	packageStr      string
	templateStr     string
	classNamePrefix string
	classNameSuffix string
	fileTypeStr     string
}

/*************************构建模板用到的结构体结束*********************************/

/*************************MySQL对应的表和字段映射的结构体开始*******************/
// 表的表结构()，通是绑定对应表，使用指针的原因，是为了处理空的默认值操作
type tableMySQL struct {
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

// 指定tableMySQL结构体对应的表名，gorm对应语法规则
func (i *tableMySQL) TableName() string {
	return "TABLES"
}

// 字段的表结构()，通是绑定对应字段，使用指针的原因，是为了处理空的默认值操作
type columnsMySQL struct {
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

// 指定columnsMySQL结构体对应的表明，gorm对应语法规则
func (i *columnsMySQL) TableName() string {
	return "COLUMNS"
}

/*************************MySQL对应的表和字段映射的结构体结束*******************/

/***************日志结构体定义和相关方法开始************************/
// 定义日志结构体
type logService struct {
	logger *zap.Logger
}

// 定义日志结构体初始化方法
func initLogServic() *logService {
	logger, err := zap.NewProduction()
	//清理资源
	defer logger.Sync()
	if err != nil {
		fmt.Println("初始化日志系统失败")
		panic(err)
	}
	return &logService{logger: logger}
}

// 定义 LogService的info日志
func (l logService) logInfo(msg string) {
	l.logger.Info(msg)
}

// 定义 LogService的error日志
func (l logService) logError(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

/***************日志结构体定义和相关方法结束************************/
