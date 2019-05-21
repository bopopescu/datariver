package boot

import (
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

type KeysApi struct {
	Addr []string
	Api  *clientv3.Client
}

var GConfigEtcdClient *KeysApi = nil
var GBusiEtcdClient *KeysApi = nil

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

func init_etcd_client(addr []string) (*KeysApi, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "创建etcd-client失败")
	}

	return &KeysApi{
		Addr: addr,
		Api:  cli,
	}, nil
}

func RelaseEtcdClient() {
	if GConfigEtcdClient != nil {
		GConfigEtcdClient.Api.Close()
	}
	if GBusiEtcdClient != nil {
		GBusiEtcdClient.Api.Close()
	}
}

func init_config_etcd_client(addr ...string) error {
	var err error = nil
	GConfigEtcdClient, err = init_etcd_client(addr)
	return err
}

func init_busi_etcd_client(addr ...string) error {
	var err error = nil
	GBusiEtcdClient, err = init_etcd_client(addr)
	return err
}
