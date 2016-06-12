package main

import (
	"action"
	"sync/atomic"
	"time"
	"fmt"
	log "dlog"
	"github.com/gcliupeng/consumer/consumergroup"
	"github.com/Shopify/sarama"
)

const CONCURRENT_POOL = 10

/*
获得消息进行并行化处理
*/
type TopicConsumer struct {
	kafkaConfig        *KafkaConfig
	group 		  		string
	conf          		TopicConfig
	worker        		*action.Worker
	count         		int64
	messages 	  		<-chan *sarama.ConsumerMessage
	consumer           *consumergroup.ConsumerGroup
}

func NewTopicConsumer(kafkaConfig *KafkaConfig, group string, conf TopicConfig, consumer *consumergroup.ConsumerGroup,
	 actionCenter *action.ActionCenter) (*TopicConsumer, error) {
	w, err := action.NewWorker(conf.Transit, conf.Access, conf.Actions, actionCenter,conf.RetryTime,conf.ErrNo)
	if err != nil {
		return nil, err
	}

	u := &TopicConsumer{
		kafkaConfig:   kafkaConfig,
		group:		   group,
		conf:          conf,
		messages:	   consumer.Messages(conf.Topic),
		worker:        w,
		count:         0,
		consumer:	   consumer,
	}
	for i := 0; i < conf.Concurrency; i++ {
		go u.consume(i)
	}
	return u, nil
}

func (this *TopicConsumer) consume(i int) {
	defer func () {
		if err := recover(); err != nil {
		return
		//fmt.Println(err)    //这里的err其实就是panic传入的内容，55
		//fmt.Printf("recover")
	}
	}()
	
	for {
		// fmt.Printf("%s\n",this.conf.Group)
		// return;
		msg := <- this.messages
		beginTime := time.Now().UnixNano()
		//log.Debugf("upstream %p routine %v get msg %v", this, i, msg.Offset)
		fmt.Printf("upstream %p routine %v get msg %v\n", this, i, msg.Offset)
		for {
			if err := this.worker.Do(msg.Value); err != nil {
				log.Errorf("group:%v topic:%v, err:%v content:%s",this.group,this.conf.Topic, err, msg.Value)
				if this.conf.EnsureConsume == false {
					break
				} else {
					if err == action.ActionHttpErr{
						log.Debugf("http status error ,topic is %s ,group is %s ", msg.Topic,this.group)
						//do nothing ,try again
						}else if err == action.ActionErrorCodeErr{
							//开始重写消息
							errRes := this.ProduceMsg(msg,err)
							if(errRes == nil){
//								log.Debugf("rewrite msg ok,topic is %s ,group is %s", msg.Topic,this.conf.Group)
							//	log.Debugf("rewrite msg ok,topic is %s ,group is %s ", msg.Topic,this.conf.Group)
							}else{
//								log.Errorf("rewrite msg faile,topic is %s ,group is %s,error is %v", msg.Topic,this.conf.Group,errRes)
							//	log.Debugf("rewrite msg faile,topic is %s ,group is %s,error is %v ", msg.Topic,this.conf.Group,errRes)
							}
							break;
						}else{
							break;
						}
				}
			} else {
				break
			}
			time.Sleep(time.Second * 10)
		}
		this.consumer.CommitUpto(msg)
		atomic.AddInt64(&this.count, 1)
		endTime := time.Now().UnixNano()
		log.Debugf("deal over msg %v", msg.Offset)
		log.Infof("deal msg ok group:%v, topic:%v, partition:%v, offset :%v duration:%v ms",
			this.group, this.conf.Topic, msg.Partition, msg.Offset, (endTime-beginTime)/1000000)
	}
}

func (this *TopicConsumer) Count() int64 {
	ret := this.count
	atomic.StoreInt64(&this.count, 0)
	return ret
}

func (this *TopicConsumer) ProduceMsg(msg *sarama.ConsumerMessage,errType error)(err error) {
	producer, err :=sarama.NewSyncProducer(this.kafkaConfig.KafkaAddrs, nil)
	if err != nil {
    	return err
	}
	defer func() {
    		if err := producer.Close(); err != nil {
     	   		return
    		}
	}()
	//kafka 重写三次
	for i:=0;i<3;i++ {
		if( errType == action.ActionHttpErr){
			msg := &sarama.ProducerMessage{Topic: msg.Topic, Value: sarama.ByteEncoder(msg.Value)}
			_, _, err := producer.SendMessage(msg)
			if err==nil{
				// fmt.Printf("> message sent to partition %d at offset %d\n", partition, offset)
				return err
			}
		}else{
			newMsg,err:=this.worker.AddTryTime(msg.Value)
			if(err==action.ExpireErr){
//				log.Errorf("retry too many time, topic is %s ,group is %s,RetryTime is %v", msg.Topic,this.conf.Group,this.worker.ReTryTime)
				log.Debugf("retry too many time, topic is %s ,group is %s,RetryTime is %v ", msg.Topic,this.group,this.worker.ReTryTime)
				return nil
			}
			msg := &sarama.ProducerMessage{Topic: msg.Topic, Value: sarama.ByteEncoder(newMsg)}
			_, _, err = producer.SendMessage(msg)
			if err==nil{
				return nil
			}
		}	
	}
	return err
}
