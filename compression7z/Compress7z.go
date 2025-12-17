package compression7z

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// 判断是否中文
func isChinese(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// 中文 → 拼音（首字母大写）
func chineseToPinyinTitle(s string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal

	var b strings.Builder

	for _, r := range s {
		if isChinese(r) {
			py := pinyin.Pinyin(string(r), args)
			if len(py) > 0 && len(py[0]) > 0 {
				b.WriteString(strings.Title(py[0][0]))
			}
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// 调用 7z 压缩单个文件
func compressFile(
	software7zPath string,
	format string,
	password string,
	srcFile string,
	dstFile string,
) error {

	args := []string{
		"a",
		"-t" + format,
		dstFile,
		srcFile,
		"-y",
	}

	if password != "" {
		args = append(args,
			"-p"+password,
			"-mhe=on",
		)
	}

	cmd := exec.Command(software7zPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("压缩失败: %w\n%s", err, string(out))
	}
	return nil
}

// 开始压缩文件
func StartCompress() {
	//从命令行中读取数据
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入待压缩的目录,换行完成输入")
	srcDir, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入异常，请重新输入: %w\n", err)
		return
	}
	srcDir = strings.TrimSpace(srcDir)

	fmt.Println("请输入压缩后存放文件的目录,换行完成输入")
	outDir, err2 := reader.ReadString('\n')
	if err2 != nil {
		fmt.Printf("读取输入异常，请重新输入: %w\n", err2)
		return
	}
	outDir = strings.TrimSpace(outDir)

	fmt.Println("请输入压缩格式,换行完成输入,常见的格式为(7z,zip)")
	format, err3 := reader.ReadString('\n')
	if err3 != nil {
		fmt.Printf("读取输入异常，请重新输入: %w\n", err3)
		return
	}
	format = strings.TrimSpace(format)

	fmt.Println("请输入密码,换行完成输入")
	password, err4 := reader.ReadString('\n')
	if err4 != nil {
		fmt.Printf("读取输入异常，请重新输入: %w\n", err4)
		return
	}
	password = strings.TrimSpace(password)

	fmt.Println("请输入7z软件的可执行文件的全路径(一般为7z.exe,7z.sh),换行完成输入")
	software7zPath, err4 := reader.ReadString('\n')
	if err4 != nil {
		fmt.Printf("读取输入异常，请重新输入: %w\n", err4)
		return
	}
	software7zPath = strings.TrimSpace(software7zPath)

	//创建目录
	_ = os.MkdirAll(outDir, 0755)
	//读取目录，并且进行单独压缩
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}

	//来时遍历下面的所有文件
	for _, entry := range entries {
		//目录跳过
		if entry.IsDir() {
			continue
		}

		//获取文件名，扩展名,以及不包含扩展名的文件名
		oldName := entry.Name()
		ext := filepath.Ext(oldName)
		base := strings.TrimSuffix(oldName, ext)
		//创建压缩后的文件名，中文转拼音(首字母大写)
		newBase := chineseToPinyinTitle(base)
		archiveName := newBase + "." + format
		//拼接全路劲
		srcPath := filepath.Join(srcDir, oldName)
		dstPath := filepath.Join(outDir, archiveName)

		fmt.Printf("压缩：%s -> %s\n", oldName, archiveName)
		//开始压缩
		err := compressFile(software7zPath, format, password, srcPath, dstPath)
		if err != nil {
			fmt.Println(err)
		}
	}
}
