package action

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func iToInt64(i interface{}) int64 {
	switch i.(type) {
	case json.Number:
		tmpPid := i.(json.Number)
		pid, _ := tmpPid.Int64()
		return pid
	case int64:
		return i.(int64)
	case string:
		val, _ := strconv.ParseInt(i.(string), 10, 64)
		return val
	}
	return 0
}

func iToString(i interface{}) string {
	switch i.(type) {
	case json.Number:
		tmpV := i.(json.Number)
		v, _ := tmpV.Int64()
		return strconv.FormatInt(v, 10)
	case int64:
		tmpV := i.(int64)
		return strconv.FormatInt(tmpV, 10)
	case string:
		return i.(string)
	}

	return ""
}

func iToFloat64(i interface{}) float64 {
	switch i.(type) {
	case json.Number:
		tmpV := i.(json.Number)
		v, _ := tmpV.Float64()
		return v
	case int64:
		tmpV := i.(int64)
		return float64(tmpV)
	case string:
		v, err := strconv.ParseFloat(i.(string), 64)
		if err != nil {
			return v
		}
	}

	return 0

}

func TimeFormat(arg map[string]interface{}) error {
	var str_time string
	if value, ok := arg["ct"]; ok {
		ct := iToInt64(value)
		str_time = time.Unix(ct, 0).Format("2006-01-02 15:04:05")
	} else {
		str_time = time.Now().Format("2006-01-02 15:04:05")
	}
	arg["timestamp"] = str_time
	return nil
}

func GenFieldAdd(str string) func(map[string]interface{}) error {
	var key string
	var val interface{}
	part := strings.Split(str, ":")
	if len(part) < 3 {
		return nil
	}
	key = part[1]
	if a, err := strconv.ParseInt(part[2], 10, 64); err != nil {
		val = strings.Trim(part[2], "\"")
	} else {
		val = a
	}
	return func(m map[string]interface{}) error {
		m[key] = val
		return nil
	}
}

func Transit(arg map[string]interface{}) error {
	worker := arg["WORKER"].(*Worker)
	transit := worker.Transit
	res := make(map[string]interface{})
	if len(transit) != 0 {
		for argk, argv := range arg {
			delete(arg, argk)
			if newkey, ok := transit[argk]; ok == true {
				if newkey == "" {
					newkey = argk
				}
				res[newkey] = argv
			}
		}
		for k, v := range res {
			arg[k] = v
		}
		arg["WORKER"] = worker
	}

	return nil
}
