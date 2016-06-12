package action

import (
	"fmt"
	"strings"
)

type ActionFunc func(arg map[string]interface{}) error
type ActionGenFunc func(string) func(map[string]interface{}) error

type ActionCenter struct {
	funcs    map[string]func(arg map[string]interface{}) error
	genFuncs map[string]func(string) func(map[string]interface{}) error
}

func NewActionCenter() *ActionCenter {
	this := &ActionCenter{
		make(map[string]func(arg map[string]interface{}) error),
		make(map[string]func(string) func(map[string]interface{}) error),
	}
	//	this.Register("DriverExtraInfo", DriverExtraInfo, nil)
	this.Register("Transit", Transit, nil)
	this.Register("HttpAccess", HttpAccess, nil)
	this.Register("TimeFormat", TimeFormat, nil)
	this.Register("FieldAdd", nil, GenFieldAdd)
	return this
}

func (this *ActionCenter) Register(funcName string,
	fun func(arg map[string]interface{}) error,
	genFun func(string) func(map[string]interface{}) error) {
	if fun != nil {
		if _, ok := this.funcs[funcName]; ok == false {
			this.funcs[funcName] = fun
		} else {
			panic(fmt.Sprintf("%v is exist", funcName))
		}
	} else if genFun != nil {
		if _, ok := this.genFuncs[funcName]; ok == false {
			this.genFuncs[funcName] = genFun
		} else {
			panic(fmt.Sprintf("%v is exist", funcName))
		}
	} else {
		panic("fun and genfun are all nil")
	}
}

func (this *ActionCenter) GetAction(funcName string) func(arg map[string]interface{}) error {
	if v, ok := this.funcs[funcName]; ok == true {
		return v
	} else {
		for k, v := range this.genFuncs {
			if strings.Contains(strings.ToLower(funcName), strings.ToLower(k)) {
				return v(funcName)
			}
		}
	}
	panic(fmt.Sprintf("there is no func for %v", funcName))
	return nil
}
