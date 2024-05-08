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
	EmbMax32			int32
	EmbServerIps [4]string
	TestTimes    int
	TestEmbSize  int
)

func init() {
	Cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'cube.ini': %v", err)
	}
	Port, _ = Cfg.Section("").Key("HTTP_PORT").Int()
	Rank, _ = Cfg.Section("").Key("RANK").Int()
	EmbDim, _ = Cfg.Section("").Key("EMB_DIM").Int()
	EmbMax, _ = (Cfg.Section("").Key("EMB_MAX").Int())
	EmbMax32 = int32(EmbMax)
	TestTimes, _ = Cfg.Section("").Key("TEST_TIMES").Int()
	TestEmbSize, _ = Cfg.Section("").Key("TEST_EMB_SIZE").Int()
	loadConf(Cfg)
	//fmt.Println("finish conf init")
}

func loadConf(Cfg *ini.File) {
	sec, err := Cfg.GetSection("cube")
	if err != nil {
		log.Fatalf("Fail to get section 'cube': %v", err)
	}
	for i := 0; i < 4; i++ {
		EmbServerIps[i] = sec.Key("CUBE" + strconv.Itoa(i)).String()
	}
}
