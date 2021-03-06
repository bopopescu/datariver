package boot

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"datariver/lib/global"
)

var g_conf_file, g_conf_etcd string

func init_server_nanme() {
	if i := strings.LastIndex(os.Args[0], "/"); i >= 0 {
		i++
		global.SERVERNAME = os.Args[0][i:]
	}
}

func Init() error {
	init_server_nanme()

	flag.StringVar(&global.GConfig.BrokerConfig.Group, "g", "",
		"this parameter is must, group name like mysql instance, ex: db_mall")
	flag.StringVar(&g_conf_etcd, "etcd", "", "etcd addr")
	flag.StringVar(&g_conf_file, "c", "", "json file config")

	flag.Parse()

	if global.GConfig.BrokerConfig.Group == "" {
		return errors.New("必须提供broker group")
	}

	err := parseConfig()
	if err != nil {
		return errors.Wrap(err, "解析配置失败")
	}

	if len(global.GConfig.BrokerConfig.EtcdAddr) != 0 {
		if err := init_busi_etcd_client(global.GConfig.BrokerConfig.EtcdAddr...); err != nil {
			return errors.Wrap(err, "初始化业务etcd-client失败")
		}
	}

	return nil
}

func getConfig(config_file string, config *global.ServerConfig) error {
	file, err := os.Open(config_file)
	if err != nil {
		return errors.Wrap(err, "打开文件错误")
	}
	defer file.Close()

	config_str, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "读文件异常")
	}
	/*
		err = json.Unmarshal(config_str, config)
		if err != nil {
			err = errors.Wrap(err, "json文件格式错误")
		}
	*/
	_, err = toml.Decode(string(config_str), config)
	if err != nil {
		err = errors.Wrapf(err, "Decode文件[%s]格式错误", config_file)
	}
	return err
}

func GetConfigFromEtcd(keyName string, cfg interface{}) error {
	return getServerConfig(global.GConfigEtcdClient, keyName, cfg)
}

func getServerConfig(c *global.KeysApi, servername string, cfg interface{}) error {
	var environment = os.Getenv("GOENV")
	if environment == "" {
		environment = "online"
	} else {
		environment = strings.ToLower(environment)
	}

	key_name := fmt.Sprintf("/config/%s/%s", servername, environment)
	rsp, err := c.Api.Get(context.Background(), key_name)
	if err != nil {
		return errors.Wrap(err, "读取ETCD错误")
	}

	if len(rsp.Kvs) == 0 {
		return errors.New("etcd配置为空")
	}
	_, err = toml.Decode(string(rsp.Kvs[0].Value), cfg)
	return err
}

func parseConfig() error {
	var err error = nil
	if g_conf_etcd != "" {
		if err := init_config_etcd_client(g_conf_etcd); err != nil {
			return errors.Wrap(err, "初始化etcd-client失败")
		}
		goto LoadEtcdConfig
	} else if g_conf_file != "" {
		goto LoadConfig
	} else {
		fmt.Fprintf(os.Stderr, "No configuration source\n")
		return errors.New("没有配置etcd地址")
	}

LoadConfig:
	//init config
	err = getConfig(g_conf_file, &global.GConfig)
	if err != nil {
		return errors.Wrap(err, "解析配置失败")
	}
	err = parseDSN(global.GConfig.SourceConfig.MysqlConn, &global.GConfig.SourceConfig)
	if err != nil {
		err = errors.Wrap(err, "解析DSN配置失败")
	}
	return err

LoadEtcdConfig:
	key := global.SERVERNAME
	for loop := true; loop; loop = false {
		if err = GetConfigFromEtcd(global.SERVERNAME, &global.GConfig.BrokerConfig); err != nil {
			err = errors.Wrap(err, "从etcd读取broker配置异常")
			break
		}
		key = fmt.Sprintf("%v:%v", global.SERVERNAME, global.GConfig.BrokerConfig.Group)
		if err = GetConfigFromEtcd(key, &global.GConfig.SourceConfig); err != nil {
			err = errors.Wrap(err, "从etcd读取source配置异常")
			break
		}

		if global.GConfig.SourceConfig.MysqlConn == "" {
			err = errors.New("MySQL配置为空")
			break
		}
		err = parseDSN(global.GConfig.SourceConfig.MysqlConn, &global.GConfig.SourceConfig)
		if err != nil {
			err = errors.Wrap(err, "解析DSN异常")
		}
	}

	return err
}

//admin_root:bjfmg1nSsynKggb@tcp(timeline:3306) -> DSN
func parseDSN(dsn string, cfg *global.SourceConfig) error {
	var err error = nil
	var left string = dsn
	for loop := true; loop; loop = false {
		var i int
		//user
		if i = strings.Index(left, ":"); i < 0 {
			err = errors.New("Invalid DSN: can not find user")
			break
		}
		cfg.DBConfig.User = left[0:i]
		i++
		left = left[i:]

		//password
		if i = strings.Index(left, "@"); i < 0 {
			err = errors.New("Invalid DSN: can not find passord")
			break
		}
		cfg.DBConfig.Passwd = left[0:i]
		i++
		left = left[i:]

		//addr
		if i = strings.Index(left, "("); i < 0 {
			err = errors.New("Invalid DSN: can not find addr")
			break
		}
		i++
		left = left[i:]

		if i = strings.Index(left, ")"); i < 0 {
			err = errors.New("Invalid DSN: can not find addr")
			break
		}
		cfg.DBConfig.Addr = left[0:i]
	}

	fmt.Printf("parseDSN %+v, %+v, err:%+v", dsn, *cfg, err)

	return err
}
