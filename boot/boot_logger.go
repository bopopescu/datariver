package boot

import (
	"fmt"
	"math"

	"github.com/op/go-logging"
	"github.com/pkg/errors"

	"common/futil"
	"datariver/config"
	"datariver/global"
)

var log_format string = "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"

func InitLogger() error {
	fpath := config.GConfig.BrokerConfig.LogDir + "/" + config.GConfig.BrokerConfig.LogFile
	fmt.Println("路径: ", fpath)
	fp, err := futil.NewFileLogWriter(fpath, false, math.MaxInt64)
	if err != nil {
		return errors.Wrap(err, "日志文件创建错误")
	}
	blog := logging.NewLogBackend(fp, "", 0)
	format := logging.MustStringFormatter(log_format)

	formatter := logging.NewBackendFormatter(blog, format)
	level := logging.AddModuleLevel(formatter)
	level.SetLevel(logging.INFO, "")

	logging.SetBackend(level)
	l := logging.MustGetLogger(global.SERVERNAME)
	global.Logger = &global.GoLoggingLogger{
		Backend: level,
		Logger:  l,
	}
	//global.Logger.ExtraCalldepth += 1
	global.Logger.Infof("初始化日志成功")

	return nil
}
