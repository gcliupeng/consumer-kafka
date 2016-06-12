## 介绍：  
本项目实现以下功能：    
将kafka中的消息取出来转到其他业务服务器
每个group的每个topic消费者可配  
能够保证在多个业务服务器消费相同topic的时候，消费的正确性  


## 配置  
配置文件conf/redirector.json.test
>{  
>    "Kafka": {  
>        "ZkHosts": ["127.0.0.1:2181"],  	//kafka和消费者的zk集群，消费者信息保存在zk上
>        "ZkPath": "/groups",  			//kafka在zk上的根目录
>        "KafkaAddrs": ["127.0.0.1:2181"], 	//kafka broker 地址 
>        "OffsetInterval": 50000000,  
>    },  
>    "Groups": [  
>        {  
>            "Group":"group1",       //Group可以用来指定业务方，可以有多个topic，  
>            "Topics":[
	     	{
		"Topic":"aaaa",         //Group和topic的组合要求唯一，  
>            	"EnsureConsume":false,  //保证消息必然被消费不丢失，如果失败会一直进行重试  
>            	"Access":["http://127.0.0.1:1111/aaaa"],  
>            	"Concurrency":2         //并发量  
>            	"Actions":["FieldAdd:act_id:2"]    //字段操作：
>            	"Transit":     //提取kafkadata中的有效数据key，转成自己的需要的key名称  
>            	{"id":"id"}
		}
		]	
>        }  
>    ],  
>    "Log": {  
>        "Output": "log",        //日志目录  
>        "Level": "DEBUG",       //日志级别  
>        "FileRotate": 5         //日志保留文件个数  
>    }  
}

## 运行  

wget http://mirror.bit.edu.cn/apache/zookeeper/stable/zookeeper-3.4.7.tar.gz   
先安装zookeeper-xxx/src/c  
configure && make && make install  

编译 make    
启动 nohup ./redirector --config=conf/redirector.json &



