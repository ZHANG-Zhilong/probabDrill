package constant

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once
var seq *map[string]int64

//GetSeqByName return
func GetSeqByName(name string) int64 {
	init1()
	return (*seq)[name]
}


//init1 init
func init1() {
	once.Do(func() {
		file, err := os.Open(StdLayer)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		contents := strings.Split(string(content), "\n")
		nameSeq := make(map[string]int64)
		for _, item := range contents {
			items := strings.Split(item, "\t")
			id, _ := strconv.ParseInt(items[2], 10, 64)
			nameSeq[items[0]] = id
		}
		seq = &nameSeq
	})
}
