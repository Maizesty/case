package conf

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
)

var (
	Port         int
	Rank         int
	EmbDim       int
	EmbMax       int
	EmbNum			 int
	EmbServerIps []string
	TestTimes    int
	TestEmbSize  int
	ModelPath    string
	EmbServerNum int
	CacheSize int64
)

func init() {
	Cfg, err := ini.Load("conf/conf_dev.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'cube.ini': %v", err)
	}
	Port, _ = Cfg.Section("").Key("HTTP_PORT").Int()
	Rank, _ = Cfg.Section("").Key("RANK").Int()
	EmbDim, _ = Cfg.Section("").Key("EMB_DIM").Int()
	EmbMax, _ = Cfg.Section("").Key("EMB_MAX").Int()
	TestTimes, _ = Cfg.Section("").Key("TEST_TIMES").Int()
	TestEmbSize, _ = Cfg.Section("").Key("TEST_EMB_SIZE").Int()
	EmbServerNum,_ = Cfg.Section("").Key("EMB_PS_NUM").Int()
	ModelPath = Cfg.Section("").Key("MODEL_PATH").String()
	CacheSize,_ = Cfg.Section("").Key("CACHE_SIZE").Int64()
	EmbNum,_ = Cfg.Section("").Key("EMB_NUM").Int()
	loadConf(Cfg)
	//fmt.Println("finish conf init")
}

func loadConf(Cfg *ini.File) {
	sec, err := Cfg.GetSection("cube")
	if err != nil {
		log.Fatalf("Fail to get section 'cube': %v", err)
	}
	EmbServerIps = make([]string, EmbServerNum)
	for i := 0; i < EmbServerNum; i++ {
		EmbServerIps[i] = sec.Key("CUBE" + strconv.Itoa(i)).String()
	}
}
