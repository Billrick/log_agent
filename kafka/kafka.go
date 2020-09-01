package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
	"time"
	"z.cn/logagentDemo/common"
)

var (
	client sarama.SyncProducer
	lock   sync.Mutex
)

func InitConnectKafka(addrs []string) (err error) {
	config := sarama.NewConfig()
	/*
		NoResponse RequiredAcks = 0  //不用响应acks响应
		WaitForLocal RequiredAcks = 1//需要主leader 返回ack
		WaitForAll RequiredAcks = -1//需要主leader->follower->leader-> 返回ack
	*/
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选取partition
	config.Producer.Return.Successes = true                   //成功后的消息在success channel返回
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("connect kafka err :",err)
		return err
	}
	go SendMsgToKafkaByJob()
	return
}

func SendMsgToKafkaByJob() (pid int32, offset int64, err error) {
	for {
		select {
			case logmsgdata := <-common.LogMsgJob:
				fmt.Printf("ready msg sennd type:%s msg:%s\n",logmsgdata.Topic,logmsgdata.Data)
				SendMsgToKafka(logmsgdata.Topic,logmsgdata.Data)
		default:
			time.Sleep(time.Second)
		}
	}
}

//发送消息
func SendMsgToKafka(topic, data string) (pid int32, offset int64, err error) {
	//构建消息
	msg := &sarama.ProducerMessage{
		Topic: topic, //消息分类
		Value: sarama.StringEncoder(data),
	}
	pid, offset, err = client.SendMessage(msg) //返回partition分区,和偏移量
	if err != nil {
		fmt.Print("send msg failed , err:", err)
		return
	}
	fmt.Printf("send msg success pid:%v offset:%v\n", pid, offset)
	return
}

//关闭生产者客户端
func Close() {
	lock.Lock()
	client.Close()
	lock.Unlock()
}
