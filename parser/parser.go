package icsParser

import (
	"encoding/json"
	"regexp"
	"strings"
)

type VCalendar struct {
	prodid  string
	version string
	method  string
	event   Vevent
}

type Vuser struct {
	name   string
	status bool
	mail   string
	cutype string
	role   string
}

type Vevent struct {
	created      string
	dtstart      string
	dtend        string
	dtstamp      string
	creator      Vuser
	participants []Vuser
	title        string
	description  string
	location     string
	lastedit     string
}

func GetJson(body *string) *string {
	prepared := prepareString(body)
	(*prepared) = (*prepared)[1 : len(*prepared)-1]
	m := make(map[string]interface{})
	parseBlocks(&m, *prepared)
	json, _ := json.Marshal(m)
	jsonString := string(json)
	return &jsonString
}

func prepareString(body *string) *[]string {
	firstSplit := (strings.Split(*body, "\n"))
	optimize(&firstSplit)
	return &firstSplit
}

func split(data *[]string) {

}

func parseBlocks(m *map[string]interface{}, block []string) int {
	i := 0
	for ; i < len(block); i++ {
		value := block[i]
		key, val := getKeyValueOf(&value)
		if "begin" == strings.ToLower(key) {
			_m := make(map[string]interface{})
			i++
			i += parseBlocks(&_m, block[i:])
			(*m)[val] = (&_m)
			continue
		}
		if "end" == strings.ToLower(key) {
			return i
		}
		(*m)[key] = val
	}
	return i
}

func optimize(array *[]string) {
	newArray := []string{}
	currentKey := 0
	for _, item := range *array {
		if len(item) == 0 {
			continue
		}
		if item[0] == ' ' {
			newArray[currentKey-1] = newArray[currentKey-1] + item[1:]
			continue
		}
		newArray = append(newArray, item)
		currentKey++
	}
	*array = newArray
}

func getKeyValueOf(item *string) (string, string) {
	reg, _ := regexp.Compile("[:;]")
	splittedString := reg.Split(*item, 2)
	return splittedString[0], splittedString[1]
}
