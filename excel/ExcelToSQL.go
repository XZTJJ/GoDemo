package excel

/**
 * 将excel的一行转成转成sql语句
 */
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// 组合操作
func CombinationOp() {
	defer ErrorHandler()
	//生成文件
	CombinationInsert()
}

// 插入语句集合
func CombinationInsert() {
	filePath := FilePathTip()
	allSQL := ToInsertSQL(filePath)
	fmt.Println(allSQL)
}

// 提示语
func FilePathTip() string {
	fmt.Print("请输入excel的决定路径，目前支持xlsx格式:")
	reader := bufio.NewReader(os.Stdin)
	filePath, err1 := reader.ReadString('\n')
	if err1 != nil {
		panic(err1)
	}
	return strings.TrimSpace(filePath)
}

// 提示语
func ColumsTip(sheetName string) ([]string, []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请指定" + sheetName + "表delete的colums,可以为空, 用,分隔名称:")
	delColumStr, err2 := reader.ReadString('\n')
	if err2 != nil {
		panic(err2)
	}
	orgDelColums := strings.Split(delColumStr, ",")
	var delColums []string
	for _, col := range orgDelColums {
		delColums = append(delColums, strings.TrimSpace(col))
	}
	fmt.Print("请指定" + sheetName + "表需要加上单引号的colums,可以为空, 用,分隔名称:")
	strColumStr, err3 := reader.ReadString('\n')
	if err3 != nil {
		panic(err2)
	}
	orgStrColums := strings.Split(strColumStr, ",")
	var strColums []string
	for _, col := range orgStrColums {
		strColums = append(strColums, strings.TrimSpace(col))
	}
	return delColums, strColums
}

// 生成对应的sql
func ToInsertSQL(filePath string) string {
	//打开文件
	f, err1 := excelize.OpenFile(filePath)
	if err1 != nil {
		panic(err1)
	}
	//关闭文件
	defer func() {
		if err2 := f.Close(); err2 != nil {
			panic(err2)
		}
	}()

	//拼接数据
	builder := strings.Builder{}
	//获取所有的sheet，遍历
	sheets := f.GetSheetList()
	for _, sheetName := range sheets {
		//指定表对应的SQL字段
		delCols, strCols := ColumsTip(sheetName)
		//获取指定的sheet的所有row
		oneSheetSQL := EachSheetHanlder(sheetName, f, delCols, strCols)
		builder.WriteString(sheetName)
		builder.WriteString("表相关的SQL:\n")
		builder.WriteString(oneSheetSQL)
		builder.WriteString("\n")
	}
	return builder.String()
}

func EachSheetHanlder(sheetName string, f *excelize.File, delCols []string, strCols []string) string {
	rows, err1 := f.GetRows(sheetName)
	if err1 != nil {
		panic(err1)
	}
	if len(rows) < 1 || len(rows[0]) < 1 {
		return ""
	}
	//标题和真实的数据
	titles := rows[0]
	datas := rows[1:]
	//生成delCols的索引map
	delColsIndexMap := make(map[string]int, len(delCols))
	for _, delCol := range delCols {
		for index, title := range titles {
			if delCol == title {
				delColsIndexMap[delCol] = index
			}
		}
	}
	//生成strCols的索引map
	strColsIndexMap := make(map[int]string, len(strCols))
	for _, strCol := range strCols {
		for index, title := range titles {
			if strCol == title {
				strColsIndexMap[index] = strCol
			}
		}
	}
	//拼接数据
	builder := strings.Builder{}
	//开始遍历数据
	for _, data := range datas {
		deleteSQL := GenerateDelete(sheetName, titles, data, delCols, delColsIndexMap, strColsIndexMap)
		insertSQL := GenerateInsert(sheetName, titles, data, strColsIndexMap)
		builder.WriteString(deleteSQL)
		builder.WriteString(insertSQL)
		builder.WriteString("\n")
	}
	return builder.String()
}

// 生成insert语句
func GenerateInsert(sheetName string, title []string, data []string, strColsIndexMap map[int]string) string {
	if len(title) < 1 || len(data) < 1 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("INSERT INTO `")
	builder.WriteString(sheetName)
	builder.WriteString("`(`")
	builder.WriteString(strings.Join(title, "`, `"))
	builder.WriteString("`) VALUES (")
	//遍历
	for index, value := range data {
		_, ok := strColsIndexMap[index]
		if ok {
			builder.WriteString("'")
			builder.WriteString(value)
			builder.WriteString("'")
		} else {
			builder.WriteString(value)
		}

		if index != len(data)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(");\n")
	return builder.String()
}

// 生成delete语句
func GenerateDelete(sheetName string, title []string, data []string, delCols []string, delColsIndexMap map[string]int, strColsIndexMap map[int]string) string {
	if len(delCols) < 1 || len(delColsIndexMap) < 1 || len(title) < 1 || len(data) < 1 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `")
	builder.WriteString(sheetName)
	builder.WriteString("` WHERE")
	//开始生成del字段
	for index, delCol := range delCols {
		//不为0的情况
		if index != 0 {
			builder.WriteString(" AND")
		}
		//
		builder.WriteString(" `")
		builder.WriteString(delCol)
		builder.WriteString("`")

		delIndex, ok1 := delColsIndexMap[delCol]
		if !ok1 {
			panic("del字段索引映射关系失败")
		}
		builder.WriteString("= ")
		_, ok2 := strColsIndexMap[delIndex]
		if ok2 {
			builder.WriteString("'")
			builder.WriteString(data[delIndex])
			builder.WriteString("'")
		} else {
			builder.WriteString(data[delIndex])
		}

	}
	builder.WriteString(";\n")
	return builder.String()
}

// 错误处理
func ErrorHandler() {
	if r := recover(); r != nil {
		fmt.Printf("错误：%s", r)
		time.Sleep(10 * time.Second)
	}
}
