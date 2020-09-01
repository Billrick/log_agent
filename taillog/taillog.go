package taillog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"time"
	"z.cn/logagentDemo/common"
)
//管理器
var tailMgr *TailTaskMgr

type TailTaskMgr struct {
	tails []*TailTask
	initData map[string]*TailTask
}

type TailTask struct {
	Path,Topic string
	instance *tail.Tail
	ctx context.Context
	cancelFunc context.CancelFunc
}
//初始化任务
func NewTailTask(path,topic string) (task *TailTask,err error){
	ctx,cancelFunc := context.WithCancel(context.Background())
	task = &TailTask{
		Path: path,
		Topic: topic,
		ctx:ctx,
		cancelFunc:cancelFunc,
	}
	tailobj ,err := task.InitTailf(path)
	if err != nil {
		fmt.Println("init log file failed , err :",err)
		return nil,err
	}
	task.instance = tailobj
	go task.run() //启动读取日志的任务
	//把初始化的日志读取器放入管理器中
	tailMgr.tails = append(tailMgr.tails,task)
	tailMgr.initData[topic+"_"+path] = task
	return task,nil
}

func (t *TailTask)run(){
	for  {
		select {
		case <- t.ctx.Done(): //收到cancelFunc调用过后的通道数据时 退出
			fmt.Printf("%s_%s stop run() \n",t.Topic,t.Path)
			return
		//读取数据 放入jobchannel中
		case line := <- t.instance.Lines:
			common.LogMsgJob<-common.NewLogMsgData(t.Topic,line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 初始化日志文件读取器
func (t *TailTask)InitTailf(path string) (tailobj *tail.Tail,err error) {
	tailobj, err = tail.TailFile(path, tail.Config{
		ReOpen:    true, //切文件后 重新打开文件
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪开始读
		MustExist: false,                                //log文件是否存在报错
		Poll:      true, //不断拉取日志
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return nil,err
	}
	return tailobj,nil
}

func StartReadLog(logSettings []*common.LogMsgSetting) error{
	tailMgr = &TailTaskMgr{
		initData: make(map[string]*TailTask,0),
	}
	for _, logmsg := range logSettings {
		//3. 读取日志信息 ,向job通道发送消息
		_, err := NewTailTask(logmsg.Path, logmsg.Topic)
		if err != nil {
			fmt.Println("init log file failed , err :", err)
			return err
		}
	}
	return nil
}

func OpTails(newConf string,opType int32){
	if opType == 0 {//新增
		var newConfigs []*common.LogMsgSetting
		err := json.Unmarshal([]byte(newConf),&newConfigs)
		if err != nil {
			fmt.Println("unmarshal json failed , err ",err)
			return
		}
		for _,newconf := range newConfigs {
			newconfKey := fmt.Sprintf("%s_%s",newconf.Topic,newconf.Path)
			if _,b := tailMgr.initData[newconfKey];!b{//配置不存在的时候，新启动启动器
				NewTailTask(newconf.Path,newconf.Topic)
				fmt.Println("new task :", newconfKey)
			}
		}
		//记录 需要被删除的读取器
		for _,tsk := range tailMgr.tails {
			initDataKey := fmt.Sprintf("%s_%s",tsk.Topic,tsk.Path)
			var isDelete bool = true
			for _,newconf := range newConfigs{
				newconfKey := fmt.Sprintf("%s_%s",newconf.Topic,newconf.Path)
				if initDataKey == newconfKey{
					isDelete =false
					continue

				}
			}
			if isDelete{//不存在删除
				//remove(tailMgr.tails,tsk)
				delete(tailMgr.initData, initDataKey)
				fmt.Println("删除后map:",tailMgr.initData)
				tsk.cancelFunc() //会像context.done 发送信号通知该退出了
			}
		}

	}else if opType ==1 {//删除

	}
}