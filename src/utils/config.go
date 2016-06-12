package utils

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
)

func LoadConfig(path string, obj interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var configStr string = string(bytes)
	re1 := regexp.MustCompile("##.*##")
	all := re1.FindAllString(configStr, -1)
	for j, i := range all {
		subfile, err := ioutil.ReadFile(strings.Trim(i, "##"))
		if err != nil {
			return err
		}
		addstr := strings.Trim(string(subfile), " ")
		if j == (len(all) - 1) {
			addstr = strings.Trim(addstr, ",\n\t")
		} else {
			if addstr[len(addstr)-1] != ',' {
				addstr += ","
			}
		}
		configStr = strings.Replace(configStr, i, addstr, 1)
	}
	if err := json.Unmarshal([]byte(configStr), obj); err != nil {
		return err
	}
	return nil
}
