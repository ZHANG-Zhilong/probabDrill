package constant

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var stdLayerOnce sync.Once
var seq map[string]int
var seq2 map[int]string

//var layers *[]int

//GetSeqByName return
func GetSeqByName(name string) int {
	initStdLayer()
	if s, ok := seq[name]; ok {
		return s
	} else {
		//标准层序必须一一对应，不能存在没有覆盖的地方
		log.Fatalf("not found, name:%s\n", name)
	}
	return seq[name]
}
func GetNameBySeq(seq int) string {
	initStdLayer()
	return seq2[seq]
}

//initStdLayer init
func initStdLayer() {
	stdLayerOnce.Do(func() {
		file, err := os.Open(viper.GetString("stdLayer"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		contents := strings.Split(string(content), "\n")
		seq = make(map[string]int)
		seq2 = make(map[int]string)
		for _, item := range contents {
			items := strings.Split(item, "\t")
			id, _ := strconv.ParseInt(items[2], 10, 64)
			seq[items[0]] = int(id)
			seq2[int(id)] = items[0]
		}
	})
}
