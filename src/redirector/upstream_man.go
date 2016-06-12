package main

import (
	//"fmt"
	"time"
	//"strings"
	log "dlog"
	"github.com/gcliupeng/consumer/consumergroup"
	"action"
)

type GroupConsumer struct {
	kafkaConfig        *KafkaConfig
	group 			   string
	topics     		   []TopicConfig
	consumer 		   *consumergroup.ConsumerGroup
	topicsConsumers	   map[string]*TopicConsumer
	actionCenter       *action.ActionCenter
}

func NewGroupConsumer(kafka *KafkaConfig, group string,topics []TopicConfig, actionCenter *action.ActionCenter) *GroupConsumer {
	return &GroupConsumer{
		kafkaConfig:        kafka,
		group:     			group,
		topics:				topics,
		actionCenter:       actionCenter,
	}
}

func (this *GroupConsumer) Run() error {
	//每一个group产生一个consumer
	consumergroupConf:=consumergroup.NewConfig()
	topics:=make([]string,100)
	for _,topic:=range this.topics{
		topics=append(topics,topic.Topic)
	}
	consumergroupConf.Zookeeper.Chroot=this.kafkaConfig.ZKPath
	consumer,err := consumergroup.JoinConsumerGroup(this.group,topics,this.kafkaConfig.ZKHosts,consumergroupConf)
	if err!=nil{
		return err
	}
	this.consumer=consumer
	this.topicsConsumers=make(map[string]*TopicConsumer)
	for i, _ := range this.topics {
		if err := this.newTopicConsumer(this.topics[i], this.actionCenter); err != nil {
			return err
		}
	}
	go func () {
		for k, v := range this.topicsConsumers {
				log.Warningf("group:%v topic:%v, done:%v",this.group, k, v.Count())
			}
			time.Sleep(time.Second * 5)
		}()
	return nil
}

func (this *GroupConsumer) Stop() error {
	// for _, v := range this.groupConsumer {
	// 	v.Stop()
	// }
	err:=this.consumer.Close()
	return err
}

func (this *GroupConsumer) newTopicConsumer(conf TopicConfig, actionCenter *action.ActionCenter) error {

	topicsConsumer, err := NewTopicConsumer(this.kafkaConfig,this.group, conf, this.consumer, actionCenter)
	if err != nil {
		return err
	}
	this.topicsConsumers[conf.Topic] = topicsConsumer
	return nil
}
