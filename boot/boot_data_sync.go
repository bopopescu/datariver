package boot

import (
	"context"
	"fmt"

	"common/ha"
	"datariver/app/service"
	"datariver/config"
	"datariver/global"

	"github.com/pkg/errors"
)

func get_key() string {
	return fmt.Sprintf("/lock/%s/%s", global.SERVERNAME, config.GConfig.BrokerConfig.Group)
}

func StartDataSyncServer() error {
	ha, err := ha.NewHaWrapper(config.GConfig.BrokerConfig.EtcdAddr, get_key(), 3,
		"no-use", &service.DataSyncServer{})
	if err != nil {
		return errors.Wrap(err, "启动同步服务失败")
	}
	go ha.Run(context.Background())

	return nil
}
