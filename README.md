### 日志采集
大体架构
![image-架构图](https://github.com/Billrick/log_agent/blob/master/%E6%9E%B6%E6%9E%84%E5%9B%BE.png?raw=true)
- 1.logagent：这里应自己写的logagent代替ELK的Logstash进行日志收集
- 2.kafka：分布式消息队列,存储采集的日志信息
- 3.logTranfer：用于从Kafka中取出数据送入ES
- 4.ES：用于对日志建立索引并存储
- 5.etcd：用于系统配置管理，起到了个数据库的作用
- 6.kibana：用于对ES里的数据进行可视化展示，呈现的是一个web界面

