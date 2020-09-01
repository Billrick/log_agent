package common

var LogMsgJob chan *LogMsgData

type LogMsgData struct {
	Topic , Data string
}

func NewLogMsgData(topic,data string) *LogMsgData{
	return &LogMsgData{Topic: topic ,Data: data}
}
//配置信息
type LogMsgSetting struct {
	Topic string `json:"topic"`
	Path string	`json:"path"`
}