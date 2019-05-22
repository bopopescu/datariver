package service

import (
	"fmt"

	"datariver/canal"
	"datariver/config"
)

type DataSyncServer struct {
	Broker *canal.SyncClient
}

func (p *DataSyncServer) SwitchToMaster() error {
	// todo: 启动同步逻辑
	var err error
	p.Broker, err = canal.NewSyncClient(&config.GConfig)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		p.Broker.Start()
	}
	return nil
}

func (p *DataSyncServer) SwitchToSlave() error {
	// todo: 启动释放同步逻辑
	if p.Broker != nil {
		p.Broker.Close()
		p.Broker = nil
	}
	return nil
}

func (p *DataSyncServer) IsOK() bool {
	// todo: 状态检查
	return true
}
