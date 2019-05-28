package canal

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"

	dconfig "datariver/config"
	"datariver/global"
)

import ()

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags|log.Lshortfile)
}

var defaultQueue KafKaQueue

type KafKaQueue struct {
	Addr     []string
	Producer sarama.SyncProducer
}

func InitQueue(addr []string) error {
	global.Logger.Infof("Begin Queue Init addr:%+v", addr)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 2
	config.Version = sarama.V2_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForLocal
	//config.Producer.Flush.MaxMessages = 1
	config.Producer.Return.Successes = true
	config.ClientID = global.SERVERNAME + "-" + dconfig.GConfig.BrokerConfig.Group

	defaultQueue = KafKaQueue{Addr: addr}
	var err error
	if defaultQueue.Producer, err =
		sarama.NewSyncProducer(defaultQueue.Addr, config); err != nil {
		global.Logger.Errorf("Queue Init Failed addr:%+v, err:%+v", addr, err)
	}
	global.Logger.Info("Queue Init End addr:%+v, err:%+v", addr, err)
	return err
}

func PushBack(data *EventData) error {
	global.Logger.Debug("push to kafka [%v %v.%v] begin", data.Action, data.Schema, data.Table)
	msg := genKafkaMsg(data)
	var err = retrySendMessages(data, msg)
	if err != nil {
		global.Logger.Error("push to kafka err:%v, data:%+v", err, *data)
	} else {
		global.Logger.Debug("push to kafka [%v %v.%v] success", data.Action, data.Schema, data.Table)
	}
	return err
}

func genKafkaMsg(data *EventData) []*sarama.ProducerMessage {
	msg_list := make([]*sarama.ProducerMessage, 0)
	if data.Owner.KeyIndex >= 0 {
		dispatch_msg := make(map[TopicInfo]*EventData)
		for _, row := range data.Rows {
			data.Owner.Key = fmt.Sprintf("%v", row[data.Owner.KeyIndex])
			if _, ok := dispatch_msg[data.Owner]; !ok {
				dispatch_msg[data.Owner] = &EventData{
					Action: data.Action, Schema: data.Schema, Table: data.Table,
					Columns: data.Columns, Rows: make([][]interface{}, 0),
				}
			}
			dispatch_msg[data.Owner].Rows = append(dispatch_msg[data.Owner].Rows, row)
		}
		for topic, data := range dispatch_msg {
			msg := &sarama.ProducerMessage{
				Topic: topic.Topic, Key: sarama.StringEncoder(topic.Key), Value: data}
			msg_list = append(msg_list, msg)
		}
	} else {
		msg := &sarama.ProducerMessage{
			Topic: data.Owner.Topic, Key: sarama.StringEncoder(data.Owner.Key), Value: data}
		msg_list = append(msg_list, msg)
	}
	return msg_list
}

func retrySendMessages(data *EventData, msgs []*sarama.ProducerMessage) error {
	var err error
	for loop := true; loop; loop = false {
		if err = defaultQueue.Producer.SendMessages(msgs); err == nil {
			break
		}
		global.Logger.Info("begin retry push msg, last err:%+v ,data:%+v", err, *data)
		deadline := time.After(2 * time.Minute)
		interval := 5 * time.Second
		tick := time.NewTicker(interval)
		defer tick.Stop()
	InnerLoop:
		for {
			select {
			case <-tick.C:
				global.Logger.Info("retry push msg, last err:%+v ,data:%+v", err, *data)
				if err = defaultQueue.Producer.SendMessages(msgs); err == nil {
					break InnerLoop
				}
			case <-deadline:
				global.Logger.Error("retry timeout with err:%v, data:%+v", err, *data)
				break InnerLoop
			}
		}
	}
	return err
}

func CloseProducer() error {
	var err error
	if defaultQueue.Producer != nil {
		err = defaultQueue.Producer.Close()
		global.Logger.Error("Close Producer whit err:%+v", err)
	}
	return err
}
