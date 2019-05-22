package boot

import (
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"

	"datariver/global"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

func init_etcd_client(addr []string) (*global.KeysApi, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "创建etcd-client失败")
	}

	return &global.KeysApi{
		Addr: addr,
		Api:  cli,
	}, nil
}

func RelaseEtcdClient() {
	if global.GConfigEtcdClient != nil {
		global.GConfigEtcdClient.Api.Close()
	}
	if global.GBusiEtcdClient != nil {
		global.GBusiEtcdClient.Api.Close()
	}
}

func init_config_etcd_client(addr ...string) error {
	var err error = nil
	global.GConfigEtcdClient, err = init_etcd_client(addr)
	return err
}

func init_busi_etcd_client(addr ...string) error {
	var err error = nil
	global.GBusiEtcdClient, err = init_etcd_client(addr)
	return err
}
