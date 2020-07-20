package canal

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/siddontang/go-mysql/mysql"

	"common/ha"
	"datariver/lib/global"
)

func NewMetaInfo(group string) (*MetaInfo, error) {
	m := &MetaInfo{}
	key := fmt.Sprintf("/position/%s/%s", global.SERVERNAME, group)
	resp, err := global.GBusiEtcdClient.Api.Get(context.Background(), key)
	if err != nil {
		global.Logger.Error("get meta from etcd key:%v, err:%v", key, err)
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.New("meta信息为空")
	}

	err = json.Unmarshal([]byte(resp.Kvs[0].Value), m)
	m.Group = group

	return m, nil
}

type MetaInfo struct {
	sync.RWMutex
	//binlog filename
	Name string
	//binlog position
	Pos uint32
	//server group, check server.group == meta.group
	Group string
	//save time
	LastSaveTime time.Time
	//server self role
	MyRole ha.RoleType
}

func (p *MetaInfo) Save(pos mysql.Position) error {
	p.Lock()
	defer p.Unlock()

	if p.MyRole != ha.Main {
		return errors.New("从不可以写meta信息")
		return fmt.Errorf("subordinate server, cannot, save meta info")
	}
	p.Name, p.Pos, p.LastSaveTime = pos.Name, pos.Pos, time.Now()
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("/position/%s/%s", global.SERVERNAME, p.Group)
	if _, err = global.GBusiEtcdClient.Api.Put(context.Background(), key, string(d)); err != nil {
		return err
	}

	return nil
}

func (p *MetaInfo) Position() mysql.Position {
	p.Lock()
	defer p.Unlock()
	return mysql.Position{Name: p.Name, Pos: p.Pos}
}

func (p *MetaInfo) Close() {
	pos := p.Position()
	p.Save(pos)
	p.MyRole = ha.Subordinate
	return
}
