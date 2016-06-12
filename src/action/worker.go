package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"errors"
	log "dlog"
)

var ExpireErr error = errors.New("try too many times")

type SrcData struct {
	No   int64                  `json:"no"`
	Type int64                  `json:"type"`
	Ct   interface{}            `json:"ct"`
	Data map[string]interface{} `json:"data"`
	TryTime int               `json:"tryTime"`
}

type Worker struct {
	actions []func(arg map[string]interface{}) error
	Servs   []string
	Transit map[string]string
	ReTryTime int
	ErrNo	[]int  
}

func NewWorker(t map[string]string, s []string, actStrs []string,
	 actCenter *ActionCenter,retryTime int,errNo []int) (*Worker, error) {
	this := &Worker{
		actions: make([]func(arg map[string]interface{}) error, 0),
		Servs:   s,
		Transit: t,
		ReTryTime: retryTime,
		ErrNo:	 errNo,
	}
	for _, j := range actStrs {
		tmp := actCenter.GetAction(j)
		this.actions = append(this.actions, tmp)
	}
	return this, nil
}

func (this *Worker) Do(in []byte) error {
	d, err := this.ParseData(in)
	if err != nil {
		return err
	}
	d["WORKER"] = this
	for _, j := range this.actions {
		err = j(d)
		if err!=nil{
			return err
		}
	}
	return err
}

func (this *Worker) ParseData(in []byte) (map[string]interface{}, error) {
	var data SrcData
	decoder := json.NewDecoder(bytes.NewReader(in))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	ctJson := data.Ct.(json.Number)
	if ctTime, err := ctJson.Int64(); err == nil {
		log.Warningf("kafka duration %v | now:%v, inkafka:%v", (time.Now().Unix() - ctTime), time.Now().Unix(), data.Ct)
	}
	src := data.Data
	if src == nil {
		return nil, fmt.Errorf("there is no key data; %s", in)
	}
	src["ct"] = data.Ct
	src["kafkaMsgType"]=data.Type
	return src, nil
}
func (this *Worker) AddTryTime(in []byte) (out []byte,err error) {
	var data SrcData
	decoder := json.NewDecoder(bytes.NewReader(in))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	tryTime:=data.TryTime
	data.TryTime=tryTime+1
	if(data.TryTime > this.ReTryTime){
		return nil,ExpireErr
	}
	out, err = json.Marshal(data)
	return out, err
}
