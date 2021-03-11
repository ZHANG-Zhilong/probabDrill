package constant

import (
	"io/ioutil"
	"os"
	"probabDrill"
	"strconv"
	"strings"
	"sync"
)

var stdLayerOnce sync.Once
var seq *map[string]int
//var layers *[]int

//GetSeqByName return
func GetSeqByName(name string) int {
	initStdLayer()
	return (*seq)[name]
}

//initStdLayer init
func initStdLayer() {
	stdLayerOnce.Do(func() {
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
