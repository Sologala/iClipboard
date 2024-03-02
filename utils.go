package main

import (
	"os"
    "fmt"
)

func IsExistFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirectory(path string) error {
	if a, err := os.Stat(path); err != nil || !a.IsDir() {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func createFolderIfNotExists(folderPath string) error {
    if _, err := os.Stat(folderPath); os.IsNotExist(err) {
        err = os.MkdirAll(folderPath, 0755)
        if err != nil {
            return err
        }
        fmt.Printf("文件夹 %s 不存在，已成功创建！\n", folderPath)
    } else {
        fmt.Printf("文件夹 %s 已存在。\n", folderPath)
    }
    return nil
}


func CreateFileIfNotExist(file string) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
            os.Create(file)
            fmt.Printf("文件 %s 不存在，已成功创建！\n", file)
		}
	}
}
