/*
 * author:        liujun
 * created:       2018-12-28 11:31
 * last modified: 2018-12-28 11:31
 * filename:      ha.go
 * description:
 *	- 2018-12-28 ha的基础版本
 */
package ha

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Server interface {
	SwitchToMaster() error
	SwitchToSlave() error
	IsOK() bool
}

type RoleType int

const (
	Master RoleType = 1
	Slave  RoleType = 2
)

type HaWrapper struct {
	Lock           *Dlock
	Role           RoleType
	ServerInstance Server
	Stop           chan struct{}
}

func (obj HaWrapper) Run(ctx context.Context) {
	// init role
	if lerr := obj.Lock.Lock(); lerr != nil {
		obj.ServerInstance.SwitchToMaster()
		obj.Role = Master
	}

	interval := obj.Lock.Ttl / 3

	tick := time.NewTicker(time.Second * time.Duration(interval))
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if obj.Role != Master {
				if lerr := obj.Lock.Lock(); lerr != nil {
					obj.ServerInstance.SwitchToMaster()
					obj.Role = Master
				}
			} else {
				// todo: check server is ok
			}
		case <-ctx.Done():
			// todo: show why?
			if obj.Role == Master {
				obj.ServerInstance.SwitchToSlave()
				obj.Role = Slave
			}
		}
	}
}

func NewHaWrapper(endpoint []string, key string, ttl int, name string,
	ser Server) (*HaWrapper, error) {
	ha := HaWrapper{}
	ha.Lock = &Dlock{}
	init_err := ha.Lock.Init(endpoint, key, ttl, name)
	if init_err != nil {
		return nil, errors.Wrap(init_err, "init dlock error")
	}

	ha.ServerInstance = ser

	return &ha, nil
}
