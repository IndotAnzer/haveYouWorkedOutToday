package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 定义需要忽略的目录名（常见的非源码目录）
var ignoreDirs = map[string]bool{
	".git":         true,
	".svn":         true,
	"node_modules": true,
	"dist":         true,
	"build":        true,
	"vendor":       true,
	".idea":        true,
	".vscode":      true,
	"tmp":          true,
	"temp":         true,
	"logs":         true,
	"cache":        true,
}

// 定义需要忽略的文件后缀（可根据需求调整）
var ignoreExts = map[string]bool{
	".exe":    true,
	".dll":    true,
	".so":     true,
	".a":      true,
	".o":      true,
	".zip":    true,
	".tar.gz": true,
	".rar":    true,
	".7z":     true,
	".jpg":    true,
	".png":    true,
	".gif":    true,
	".mp4":    true,
	".mp3":    true,
	".pdf":    true,
	".docx":   true,
	".xlsx":   true,
}

func main() {
	// 解析命令行参数
	targetDir := flag.String("dir", ".", "要提取源码的目标目录（默认当前目录）")
	outputFile := flag.String("output", "", "输出文件路径（为空则输出到控制台）")
	flag.Parse()

	// 验证目标目录是否存在
	if _, err := os.Stat(*targetDir); os.IsNotExist(err) {
		fmt.Printf("错误：目标目录 %s 不存在\n", *targetDir)
		os.Exit(1)
	}

	// 收集所有文件的源码内容
	var codeBuilder strings.Builder
	err := filepath.WalkDir(*targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("警告：无法访问路径 %s - %v\n", path, err)
			return nil // 跳过错误路径，继续遍历
		}

		// 跳过目录（尤其是忽略列表中的目录）
		if d.IsDir() {
			dirName := filepath.Base(path)
			if ignoreDirs[dirName] {
				return fs.SkipDir // 跳过整个忽略的目录
			}
			return nil
		}

		// 跳过忽略的文件后缀
		fileExt := filepath.Ext(path)
		if ignoreExts[fileExt] {
			return nil
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("警告：无法读取文件 %s - %v\n", path, err)
			return nil
		}

		// 跳过空文件
		if len(content) == 0 {
			return nil
		}

		// 整理文件路径（转为相对路径，更易读）
		relPath, err := filepath.Rel(*targetDir, path)
		if err != nil {
			relPath = path // 转换失败则使用绝对路径
		}

		// 按格式写入内容（分隔线 + 文件路径 + 源码内容）
		codeBuilder.WriteString("========================================\n")
		codeBuilder.WriteString(fmt.Sprintf("文件路径：%s\n", relPath))
		codeBuilder.WriteString("----------------------------------------\n")
		codeBuilder.Write(content)
		codeBuilder.WriteString("\n\n") // 空行分隔不同文件

		return nil
	})

	if err != nil {
		fmt.Printf("遍历目录时出错：%v\n", err)
		os.Exit(1)
	}

	// 输出结果（控制台或文件）
	if *outputFile != "" {
		// 写入到指定文件
		err = os.WriteFile(*outputFile, []byte(codeBuilder.String()), 0644)
		if err != nil {
			fmt.Printf("错误：无法写入输出文件 %s - %v\n", *outputFile, err)
			os.Exit(1)
		}
		fmt.Printf("成功！源码已提取到文件：%s\n", *outputFile)
	} else {
		// 输出到控制台
		fmt.Println("======= 项目源码提取结果 =======")
		fmt.Println(codeBuilder.String())
	}
}
