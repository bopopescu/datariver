package canal

import (
	"reflect"
	"unsafe"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"

	"datariver/global"
)

type posSaver struct {
	pos   mysql.Position
	force bool
}

type eventHandler struct {
	c *SyncClient
}

func (h *eventHandler) OnRotate(e *replication.RotateEvent) error {
	pos := mysql.Position{
		string(e.NextLogName),
		uint32(e.Position),
	}

	h.c.syncCh <- posSaver{pos, true}

	return h.c.ctx.Err()
}

func (h *eventHandler) OnDDL(nextPos mysql.Position, _ *replication.QueryEvent) error {
	h.c.syncCh <- posSaver{nextPos, true}
	return h.c.ctx.Err()
}

func (h *eventHandler) OnXID(nextPos mysql.Position) error {
	h.c.syncCh <- posSaver{nextPos, false}
	return h.c.ctx.Err()
}

func (h *eventHandler) OnGTID(mysql.GTIDSet) error {
	// todo: gtid
	return nil
}
func (h *eventHandler) OnPosSynced(mysql.Position, mysql.GTIDSet, bool) error {
	return nil
}
func (h *eventHandler) OnTableChanged(schema string, table string) error {
	return nil
}

func hack_string(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

func (h *eventHandler) OnRow(e *canal.RowsEvent) error {
	rule, ok := h.c.getFilterInfo(e.Table.Schema, e.Table.Name)
	global.Logger.Debugf("%+v %+v.%+v", e.Action, e.Table.Schema, e.Table.Name)
	if !ok {
		global.Logger.Debugf("%+v.%+v filtered, continue", e.Table.Schema, e.Table.Name)
		return nil
	}
	data := EventData{}
	data.Action = EventAction(e.Action)
	data.Schema = e.Table.Schema
	data.Table = e.Table.Name
	var hit_index bool
	for index, col := range e.Table.Columns {
		if !hit_index && rule.KeyIndex != default_key_index && col.Name == rule.Key {
			rule.KeyIndex = index
			hit_index = true
		}
		data.Columns = append(data.Columns, col.Name)
		/*
			// do not need to process TYPE_TEXT
			if col.Type == 12 {
				for i, row := range e.Rows {
					if index < len(row) {
						if t, ok := e.Rows[i][index].([]byte); ok {
							e.Rows[i][index] = hack_string(t)
						}
					}
				}
			}
		*/
	}
	//should not come here
	if rule.KeyIndex >= len(data.Columns) || rule.KeyIndex == default_key_index {
		rule.KeyIndex = 0
	}
	data.Owner = *rule
	data.Rows = e.Rows
	global.Logger.Debugf("data:%+v", data)
	h.c.syncCh <- data
	return nil
}

func (h *eventHandler) String() string {
	return "BinlogEventHandler"
}
