package boot

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"

	"datariver/config"
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

	flag.StringVar(&config.GConfig.BrokerConfig.Group, "g", "",
		"this parameter is must, group name like mysql instance, ex: db_mall")
	flag.StringVar(&g_conf_etcd, "etcd", "", "etcd addr")
	flag.StringVar(&g_conf_file, "c", "", "json file config")

	flag.Parse()

	if config.GConfig.BrokerConfig.Group == "" {
		return errors.New("必须提供broker group")
	}

	err := parseConfig()
	if err != nil {
		return errors.Wrap(err, "解析配置失败")
	}

	if len(config.GConfig.BrokerConfig.EtcdAddr) != 0 {
		if err := init_busi_etcd_client(config.GConfig.BrokerConfig.EtcdAddr...); err != nil {
			return errors.Wrap(err, "初始化业务etcd-client失败")
		}
	}

	return nil
}

func getConfig(config_file string, config *config.ServerConfig) error {
	file, err := os.Open(config_file)
	if err != nil {
		return errors.Wrap(err, "打开文件错误")
	}
	defer file.Close()

	config_str, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "读文件异常")
	}
	err = json.Unmarshal(config_str, config)
	if err != nil {
		err = errors.Wrap(err, "json文件格式错误")
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
	return json.Unmarshal([]byte(rsp.Kvs[0].Value), cfg)

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
	err = getConfig(g_conf_file, &config.GConfig)
	if err != nil {
		return errors.Wrap(err, "解析配置失败")
	}
	err = parseDSN(config.GConfig.SourceConfig.MysqlConn, &config.GConfig.SourceConfig)
	if err != nil {
		err = errors.Wrap(err, "解析DSN配置失败")
	}
	return err

LoadEtcdConfig:
	key := global.SERVERNAME
	for loop := true; loop; loop = false {
		if err = GetConfigFromEtcd(global.SERVERNAME, &config.GConfig.BrokerConfig); err != nil {
			err = errors.Wrap(err, "从etcd读取broker配置异常")
			break
		}
		key = fmt.Sprintf("%v:%v", global.SERVERNAME, config.GConfig.BrokerConfig.Group)
		if err = GetConfigFromEtcd(key, &config.GConfig.SourceConfig); err != nil {
			err = errors.Wrap(err, "从etcd读取source配置异常")
			break
		}

		if config.GConfig.SourceConfig.MysqlConn == "" {
			err = errors.New("MySQL配置为空")
			break
		}
		err = parseDSN(config.GConfig.SourceConfig.MysqlConn, &config.GConfig.SourceConfig)
		if err != nil {
			err = errors.Wrap(err, "解析DSN异常")
		}
	}

	return err
}

//admin_root:bjfmg1nSsynKggb@tcp(timeline:3306) -> DSN
func parseDSN(dsn string, cfg *config.SourceConfig) error {
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
