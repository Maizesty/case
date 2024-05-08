package main

import (
	"bufio"
	"embed_server/conf"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var (
	Our_partition          map[int]int
	Our_hot_cache          map[int]map[int]bool
	Het_partition          map[int]int
	Het_hot_cache          map[int]map[int]bool
	Local_emb              [][]float32
	Our_no_cache_partition map[int]int
)

func init() {
	//n_part := 4
	//Our_partition, Our_hot_cache = loadPartition("our", n_part)
	//Het_partition, Het_hot_cache = loadPartition("het", n_part)
	//Our_partition_avazu, Our_hot_cache_avazu = loadPartition("our", n_part)
	//Het_partition_avazu, Het_hot_cache_avazu = loadPartition("het", n_part)
	//Our_no_cache_partition, _ = loadPartition("our", n_part)
	generateEmbed()
	//fmt.Println("finish loading data")
}

func generateEmbed() {
	Local_emb = make([][]float32, conf.EmbMax) // 本机保存的嵌入参数
	for i := 0; i < conf.EmbMax; i++ {
		for j := 0; j < conf.EmbDim; j++ {
			Local_emb[i] = append(Local_emb[i], rand.Float32())
		}
	}
}

func loadPartition(file_name string, n_part int) (map[int]int, map[int]map[int]bool) {
	partition := make(map[int]int)
	hot_cache := make(map[int]map[int]bool)
	for i := 0; i < n_part; i++ {
		hot_cache[i] = make(map[int]bool)
	}
	path := "./partition_data/" + file_name + "/" + file_name + "_p.txt"
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		key, _ := strconv.Atoi(line[0])
		value, _ := strconv.Atoi(line[1])
		partition[key] = value
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < n_part; i++ {
		f_path := "./partition_data/" + file_name + "/" + file_name + "_h" + strconv.Itoa(i) + ".txt"
		data, err := os.ReadFile(f_path)
		if err != nil {
			log.Fatal(err)
		}
		hot_data := strings.Split(string(data), " ")
		for hot_i := range hot_data {
			key, _ := strconv.Atoi(hot_data[hot_i])
			hot_cache[i][key] = true
		}
	}
	return partition, hot_cache
}
