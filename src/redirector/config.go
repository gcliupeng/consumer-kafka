package main

import (
	log "dlog"
)

type GroupConfig struct {
	Group         string
	Topics		  []TopicConfig
}

type TopicConfig struct {
	Topic 			string
	EnsureConsume	bool
	Access			[]string
	Concurrency		int
	Actions			[]string
	Transit			map[string]string
	RetryTime		int
	ErrNo			[]int
}

// type LogConfig struct {
// 	type 
// 	Output     string
// 	Level      string
// 	FileRotate int
// }

type KafkaConfig struct {
	ZKHosts      	[]string
	ZKPath		 	string
	KafkaAddrs		[]string
	OffsetInterval  int64
}

type Config struct {
	Kafka     KafkaConfig
	Groups 	  []GroupConfig
	Log       log.LogConfig
}
