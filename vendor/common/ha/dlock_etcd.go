/*
 * author:        liujun
 * created:       2018-12-28 11:33
 * last modified: 2018-12-28 11:33
 * filename:      dlock.go
 * description:
 *	- 2018-12-28 基于etcd clientv3的分布式锁的基础版本
 */
package ha

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

var LockedError = errors.New("lock has been locked")
var AlreadyLockError = errors.New("already locked")
var ParamsInvalidError = errors.New("param invalid")

const MinTtl int = 3

type LeaseInfo struct {
	Lease   clientv3.Lease
	LeaseId clientv3.LeaseID
	Cancel  context.CancelFunc
}

func (obj LeaseInfo) Init(client *clientv3.Client, ttl int) error {
	if client == nil {
		return ParamsInvalidError
	}

	if ttl < MinTtl {
		ttl = MinTtl
	}

	var err error
	var ka_rsp <-chan *clientv3.LeaseKeepAliveResponse

	obj.Lease = clientv3.NewLease(client)
	lease_resp, grant_err := obj.Lease.Grant(context.TODO(), int64(ttl/3))
	ctx, cancel := context.WithCancel(context.TODO())
	obj.Cancel = cancel
	if grant_err != nil {
		err = errors.Wrap(grant_err, "grant fail")
		goto INIT_FAIL
	}

	obj.LeaseId = lease_resp.ID
	ka_rsp, err = obj.Lease.KeepAlive(ctx, obj.LeaseId)
	if err != nil {
		err = errors.Wrap(err, "set keep alive")
		goto INIT_FAIL
	}

	go listenLeaseChan(ka_rsp)

	return nil

INIT_FAIL:
	obj.Lease.Close()
	obj.Lease = nil

	return err
}

type Dlock struct {
	EndPoint    []string
	LockKeyName string
	NodeName    string
	Ttl         int
	client      *clientv3.Client
	Lease       *LeaseInfo
	kv          clientv3.KV
}

func (obj *Dlock) Init(endpoint []string, key string, ttl int, name string) error {
	obj.EndPoint = endpoint
	obj.LockKeyName = key
	obj.Ttl = ttl
	obj.NodeName = name

	var err error = nil
	obj.client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "初始化etcd-client失败")
	}

	obj.kv = clientv3.NewKV(obj.client)

	return err
}

func (obj *Dlock) Lock() error {
	if obj.Lease != nil {
		return nil
	}

	obj.Lease = &LeaseInfo{}
	obj.Lease.Init(obj.client, obj.Ttl)

	txn := obj.kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(obj.LockKeyName), "=", 0)).
		Then(clientv3.OpPut(obj.LockKeyName, obj.NodeName,
			clientv3.WithLease(obj.Lease.LeaseId))).
		Else(clientv3.OpGet(obj.LockKeyName))
	var resp *clientv3.TxnResponse
	var commit_err error
	if resp, commit_err = txn.Commit(); commit_err != nil {
		return errors.Wrap(commit_err, "commit fail")
	}

	if !resp.Succeeded {
		return LockedError
	}

	return nil
}

func (obj *Dlock) Unlock() error {
	if obj.Lease == nil {
		return nil
	}
	obj.Lease.Cancel()
	_, err := obj.client.Revoke(context.TODO(), obj.Lease.LeaseId)
	if err != nil {
		return errors.Wrap(err, "释放租约失败")
	}
	obj.Lease.Lease.Close()
	obj.Lease = nil

	return nil
}

func listenLeaseChan(leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse) {
	var (
		leaseKeepResp *clientv3.LeaseKeepAliveResponse
	)
	for {
		select {
		case leaseKeepResp = <-leaseRespChan:
			if leaseKeepResp == nil {
				goto END
			} else {
				// heart beat
				// do nothing
			}
		}
	}
END:
	fmt.Println("心跳退出")
}
