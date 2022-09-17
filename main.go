package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	icsParser "depik.com/ics/parser"
)

func parse(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	bodyptr := string(body)

	fmt.Fprintf(res, "%s", *icsParser.GetJson(&bodyptr))
}

func main() {
	http.HandleFunc("/", parse)
	http.ListenAndServe(":8080", nil)
}
