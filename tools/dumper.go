package tools

import (
	"encoding/json"
	"log"
)

// Verbose flag to switch verbose
var Verbose = true

// Dumper verbose logs
func Dumper(infos ...interface{}) {
	if Verbose {
		j, _ := json.MarshalIndent(infos, "", "\t")
		log.Println(string(j))
	}
}
