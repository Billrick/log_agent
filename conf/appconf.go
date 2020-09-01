package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"z.cn/logagentDemo/common"
)

type AppConf struct {
	KafkaConf `ini:"kafka"`  //对应conf.ini的小节 [*]
	EtcdConf `ini:"etcd"`
}
type KafkaConf struct {
	Addrs []string `ini:"addrs"`
	BufSize int `ini:"bufSize"`
}
type EtcdConf struct {
	Addrs []string `ini:"addrs"`
	Key string `ini:"key"`
}

func setJobMaxSize(maxSize int){
	common.LogMsgJob = make(chan *common.LogMsgData,maxSize)
}

func InitConf()  (appconf AppConf,err error){
	err = ini.MapTo(&appconf,"./conf/conf.ini")
	if err != nil {
		fmt.Println("init appconf failed , err :",err)
	}
	setJobMaxSize(appconf.BufSize)
	return
}