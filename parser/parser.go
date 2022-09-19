package icsParser

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
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

func ParallelGetJson(body *string) *string {
	rq := regexp.MustCompile("END:VEVENT")
	list := rq.Split(*body, -1)
	chunk := make(chan string, len(list))
	results := make(chan *string, len(list))
	function := func() {
		data := <-chunk
		data += "\nEND:VEVENT"
		results <- GetJson(&(data))
	}
	for _, v := range list {
		chunk <- v
		go function()
	}
	resArray := make([]string, len(list))
	for i := range resArray {
		str := <-results
		resArray[i] = *str
	}
	resStr := "[" + strings.Join(resArray, ",") + "]"
	return &resStr
}

func GetJson(body *string) *string {
	prepared := prepareString(body)
	//(*prepared) = (*prepared)[1 : len(*prepared)-1]
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
			if val == "VEVENT" {
				val = val + strconv.Itoa(i)
			}
			(*m)[val] = (&_m)
			continue
		}
		if "end" == strings.ToLower(key) {
			return i
		}
		if slices.Contains([]string{
			"DESCRIPTION",
			"LOCATION",
		}, key) {
			(*m)[key] = val
			continue
		}
		decomposeVal, err := parseRow(&val)
		if err != nil {
			(*m)[key] = val
		} else {
			(*m)[key] = &decomposeVal
		}
	}
	return i
}

func parseRow(row *string) (*map[string]string, error) {
	if !strings.Contains(*row, "=") && !strings.Contains(*row, ":") {
		return nil, errors.New("cannot split string")
	}
	regex := regexp.MustCompile("[;]")
	m := make(map[string]string)
	items := regex.Split((*row), -1)
	for _, v := range items {
		splittedVal := strings.Split(v, "=")

		if len(splittedVal) == 1 {
			splitBy := strings.Split(splittedVal[0], ":")
			if len(splitBy) >= 2 {
				for len(splitBy) > 0 {
					m[splitBy[0]] = splitBy[1]
					splitBy = splitBy[2:]
				}
			} else {
				m[splittedVal[0]] = ""
			}
		} else {
			splitBy := strings.Split(splittedVal[1], ":")
			if len(splitBy) > 2 {
				m[splittedVal[0]] = splitBy[0]
				splitBy = splitBy[1:]
				for len(splitBy) > 0 {
					m[splitBy[0]] = splitBy[1]
					splitBy = splitBy[2:]
				}
			} else {
				m[splittedVal[0]] = splittedVal[1]
			}
		}
	}
	return &m, nil
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
