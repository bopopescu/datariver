package main

import (
	"fmt"
	"os"

	"datariver/boot"
)

func main() {
	err := boot.Init()
	if err != nil {
		fmt.Println("初始化失败: ", err.Error())
		os.Exit(1)
	}

}
