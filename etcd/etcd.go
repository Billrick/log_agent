package etcd

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"z.cn/logagentDemo/common"
	"z.cn/logagentDemo/taillog"
)

var (client *clientv3.Client)
//初始化读取etcd配置
func InitEtcdConfig(endpoints []string) (err error){
	conf := clientv3.Config{
		Endpoints: endpoints,
	}
	client,err = clientv3.New(conf)
	if err != nil {
		fmt.Println("connect etcd failed ,err :",err)
		return
	}
	return
}

func WatchConfChange(key string){
	ConfChan := client.Watch(context.Background(),key)
	for  {
		select {
		case info:=<-ConfChan:
			fmt.Printf("linsenter conf changed , info : %#v \n",info)
			for _,ifo := range info.Events{
				fmt.Printf("conf change type:%s newvalue:%v \n",ifo.Type,string(ifo.Kv.Key)+":"+string(ifo.Kv.Value))
				//处理操作
				taillog.OpTails(string(ifo.Kv.Value),int32(ifo.Type))//
				//if ifo.Type == clientv3.EventTypePut{// 新增或修改
				//
				//}else if ifo.Type == clientv3.EventTypeDelete {
				//
				//}
			}
		default:
			time.Sleep(time.Millisecond*500)
		}
	}
}

//取配置
func GetConf(key string) (logmsgs []*common.LogMsgSetting,err error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	res,err := client.Get(ctx,key)
	if err != nil {
		return
	}
	cancel()
	logmsgs = make([]*common.LogMsgSetting,0)
	for _, ev := range res.Kvs {
		fmt.Printf("etcd get key:%s value:%s \n",string(ev.Key),string(ev.Value))
		err = json.Unmarshal(ev.Value,&logmsgs)
	}
	return
}
//添加配置
func PutConf(key,conf string) (err error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_ ,err = client.Put(ctx,key,conf)
	if err != nil {
		return
	}
	cancel()
	return
}
