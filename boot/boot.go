package boot

import (
	"github.com/pkg/errors"
)

func InitAll() error {
	err := Init()
	if err != nil {
		return errors.Wrap(err, "启动配置失败")
	}

	err = InitLogger()
	if err != nil {
		return errors.Wrap(err, "启动日志失败")
	}
	err = StartGinServer()
	if err != nil {
		return errors.Wrap(err, "启动GIN服务失败")
	}

	err = StartDataSyncServer()
	if err != nil {
		return errors.Wrap(err, "启动数据同步服务失败")
	}

	return err
}
