package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"sync"
	log "dlog"

	"action"
	"utils"
)
var 	wg     sync.WaitGroup
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	configFile := flag.String("config", "conf/redirector.json", "The path of config file")
	flag.Parse()

	config := new(Config)
	if err := utils.LoadConfig(*configFile, config); err != nil {
		fmt.Printf("load config err %v", err)
		return
	}

	//log init
	err :=log.Init(config.Log)
	if err != nil {
		fmt.Println("init log failed")
		return
	}
	// b, err := log.NewFileBackend(config.Log.Output) //log文件目录
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// log.SetLogging(config.Log.Level, b) //只输出大于等于INFO的log
	// log.Rotate(config.Log.FileRotate, 1024*1024*1024*10)

	//action center init
	actionCenter := action.NewActionCenter()
	groupConsumers:=make(map[string]*GroupConsumer)

	//upstream manager init
	for _,group:=range config.Groups{
		wg.Add(1)
		groupConsumers[group.Group] = NewGroupConsumer(&config.Kafka, group.Group, group.Topics, actionCenter)
		err:=groupConsumers[group.Group].Run()
		if err!=nil {
			log.Close()
			os.Exit(0)
		}
	}
	//signal
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
			<-signalChan
			// fmt.Printf("%v\n",ss)
			// fmt.Printf("SHUT DOWN!!!")
			for _,group:=range config.Groups{
				groupConsumers[group.Group].Stop()
				wg.Add(-1)
			}
			log.Close()
	}()

	wg.Wait()

	return
}
