package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"io"
	"fmt"
	"io/ioutil"
	log "dlog"
)

var ActionHttpErr error = errors.New("http status wrong")
var ActionErrorCodeErr error = errors.New("the return error code wrong")

func HttpAccess(arg map[string]interface{}) error {
	servs := arg["WORKER"].(*Worker).Servs
	if servs == nil {
		log.Errorf("cannot find servs")
		return nil
	}
	defer func () {
		delete(arg, "WORKER")
	}()
	input, err := json.Marshal(arg)
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}
	v := url.Values{}
	v.Set("params", string(input))
	for _,ser := range servs {
		res, err := http.PostForm(ser, v)
		if err != nil {
			return ActionHttpErr
		}
		err= ErrNoProcess(arg,res.Body,ser)
		if err !=nil{
			return err
		}
	}
	
	return nil
}

func ErrNoProcess(arg map[string]interface{},Body io.ReadCloser,ser string) error {
	if (len(arg["WORKER"].(*Worker).ErrNo)>0){
    		//业务返回码判断
    		body, err := ioutil.ReadAll(Body)
    		if err != nil {
        		Body.Close()
        		return ActionErrorCodeErr
    		}
    		var dat map[string]int
    		if err = json.Unmarshal(body, &dat); err == nil {
				u := dat["errno"]
				var check bool = false
        		for _,e := range arg["WORKER"].(*Worker).ErrNo{
        			if e == u {
        				check=true
        				break
        			}
        		}
        		if !check{
        			fmt.Printf("errno  not match in response , server is %s ,errno is %d \n",ser,u)
        			Body.Close()
        			return ActionErrorCodeErr
        		}
        		Body.Close()
    		}else{
    			Body.Close()
    			fmt.Printf("parse errno fail ,miss erron or isnot a number,server is  %s\n",ser)
    			return ActionErrorCodeErr
    		}
    	}
    	Body.Close();
    	return nil
}
