package global

import (
	"go.etcd.io/etcd/clientv3"
)

var SERVERNAME string

type KeysApi struct {
	Addr []string
	Api  *clientv3.Client
}

var GConfigEtcdClient *KeysApi = nil
var GBusiEtcdClient *KeysApi = nil
