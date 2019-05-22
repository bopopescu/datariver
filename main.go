package main

import (
	"fmt"
	"os"
	"time"

	"datariver/boot"
)

func main() {
	err := boot.InitAll()
	if err != nil {
		fmt.Println("启动失败: ", err.Error())
		os.Exit(1)
	}

	for true {
		time.Sleep(time.Second * 10)
	}
}
