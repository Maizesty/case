package load

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)


var (
	Partition map[int32]int
	HotCache  map[int32]bool
)


func LoadPartition(basePath string){
	Partition = make(map[int32]int)
	HotCache = make(map[int32]bool)
	for i := 0;i<4;i++{
		path := basePath + "p_" + strconv.Itoa(i)
		fmt.Println(path)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		partition_data := strings.Split(string(data), " ")
		for _,p_data := range partition_data{
			key, _ := strconv.Atoi(p_data)
			key32 := int32(key)
			Partition[key32] = i
		}
	}
	path := basePath + "p_h"
	data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
	hot_data := strings.Split(string(data), " ")
	for _,hot_key := range(hot_data){
		key,_ := strconv.Atoi(hot_key)
		key32 := int32(key)
		HotCache[key32] = true
	}
	
}