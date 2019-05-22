package config

import ()

var GConfig ServerConfig = ServerConfig{}

type TopicInfo struct {
	// table name
	Table string
	// to kafka topic name, default name is Schema name
	Topic string
	// kafka partion key name, default table name
	Key string
}

type Source struct {
	Schema string
	Tables []TopicInfo
}

type DBConfig struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DSN    string
}

type SourceConfig struct {
	//db conn: root:password@tcp(localhost:3306)
	MysqlConn string
	//parse from MysqlConn
	DBConfig struct {
		User   string
		Passwd string
		Net    string
		Addr   string
		DSN    string
	} `json:"-"`
	//db and table list
	Sources []Source
	//if sources config is empty and SyncAll is true,broker sync all, else do nothing
	SyncAll bool
	//if binlog position not update during RebootSeconds, server will restart. RebootSeconds < 0, never restart
	RebootSeconds int64
}

type BrokerConfig struct {
	//group name from cmd line
	Group   string `json:"-"`
	LogDir  string
	LogFile string
	//etcd url, from cmd line
	EtcdAddr []string
	//binlog push to kafka
	KafkaAddr []string
	//admin port
	RPCListen string
}

type ServerConfig struct {
	//etcd: /v2/keys/config/binlogbroker
	BrokerConfig BrokerConfig
	//etcd: /v2/keys/config/binlogbroker:${group}, group name from the cmd line
	SourceConfig SourceConfig
}
