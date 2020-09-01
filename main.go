package main

import (
	"fmt"
	"sync"
	"z.cn/logagentDemo/conf"
	"z.cn/logagentDemo/etcd"
	"z.cn/logagentDemo/kafka"
	"z.cn/logagentDemo/taillog"
)

func main() {
	appconf, err := conf.InitConf()
	if err != nil {
		fmt.Println("init appconf failed , err :", err)
		return
	}
	fmt.Println(appconf)
	//1.初始化etcd 获取配置信息
	etcd.InitEtcdConfig(appconf.EtcdConf.Addrs)
	//conf.PutConf("/logagentSetting", `[{"topic":"notify","path":"./logs/notify.log"},{"topic":"notify2","path":"./logs/notify2.log"},{"topic":"notify3","path":"./logs/notify3.log"}]`)
	//1.1 获取配置信息
	logSettings, err := etcd.GetConf(appconf.EtcdConf.Key)


	if err != nil {
		fmt.Println("get conf failed , err :", err)
		return
	}
	//2. 初始化kafka连接,启动后台任务  读取任务通道中的数据发送数据到kafka
	kafka.InitConnectKafka(appconf.KafkaConf.Addrs)
	//3. 开始向读取日志项 放入任务通道中
	taillog.StartReadLog(logSettings)

	//3.1 监控配置，获取新的配置信息
	go etcd.WatchConfChange(appconf.EtcdConf.Key)
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()
}
