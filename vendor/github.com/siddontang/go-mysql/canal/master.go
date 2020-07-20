package canal

import (
	"sync"

	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
)

type mainInfo struct {
	sync.RWMutex

	pos mysql.Position

	gset mysql.GTIDSet

	timestamp uint32
}

func (m *mainInfo) Update(pos mysql.Position) {
	log.Debugf("update main position %s", pos)

	m.Lock()
	m.pos = pos
	m.Unlock()
}

func (m *mainInfo) UpdateTimestamp(ts uint32) {
	log.Debugf("update main timestamp %d", ts)

	m.Lock()
	m.timestamp = ts
	m.Unlock()
}

func (m *mainInfo) UpdateGTIDSet(gset mysql.GTIDSet) {
	log.Debugf("update main gtid set %s", gset)

	m.Lock()
	m.gset = gset
	m.Unlock()
}

func (m *mainInfo) Position() mysql.Position {
	m.RLock()
	defer m.RUnlock()

	return m.pos
}

func (m *mainInfo) Timestamp() uint32 {
	m.RLock()
	defer m.RUnlock()

	return m.timestamp
}

func (m *mainInfo) GTIDSet() mysql.GTIDSet {
	m.RLock()
	defer m.RUnlock()

	if m.gset == nil {
		return nil
	}
	return m.gset.Clone()
}
