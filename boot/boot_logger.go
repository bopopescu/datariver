package boot

import (
	"github.com/gogf/gf/g/os/glog"
	"github.com/pkg/errors"

	"datariver/config"
	"datariver/global"
)

/*
import (
	//"math"

	"github.com/op/go-logging"
	"github.com/gogf/gf/g/os/glog"
	//"github.com/pkg/errors"

	"datariver/config"
	"datariver/global"
)

var log_format string = ""

func InitLogger(path string) error {
	fpath := config.GConfig.BrokerConfig.LogDir + "/" + config.GConfig.BrokerConfig.LogFile
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
	global.Logger = logging.MustGetLogger(global.SERVERNAME)
	global.Logger.ExtraCalldepth += 1

	return nil
}
*/

func InitLogger() error {
	logger := glog.New()
	err := logger.SetPath(config.GConfig.BrokerConfig.LogDir)
	if err != nil {
		return errors.Wrap(err, "设置日志目录")
	}
	logger.SetStdPrint(true)

	global.Logger = logger

	return nil
}
