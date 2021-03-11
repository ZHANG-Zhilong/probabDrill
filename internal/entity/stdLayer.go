package entity

import (
	"io/ioutil"
	"os"
	"probabDrill"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once
var seq *map[string]int
//var layers *[]int

//GetSeqByName return
func GetSeqByName(name string) int {
	init1()
	return (*seq)[name]
}

//init1 init
func init1() {
	once.Do(func() {
		file, err := os.Open(probabDrill.StdLayer)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		contents := strings.Split(string(content), "\n")
		nameSeq := make(map[string]int)
		for _, item := range contents {
			items := strings.Split(item, "\t")
			id, _ := strconv.ParseInt(items[2], 10, 64)
			nameSeq[items[0]] = int(id)
		}
		seq = &nameSeq
	})
}
