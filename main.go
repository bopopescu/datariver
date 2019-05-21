package main

import (
	"fmt"
	"os"

	"datariver/boot"
	"datariver/canal"
)

func main() {
	err := boot.InitAll()
	if err != nil {
		fmt.Println("启动失败: ", err.Error())
		os.Exit(1)
	}

	_, err = canal.NewMetaInfo("test")
	fmt.Println(err)
}
