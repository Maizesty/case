package main

import(

	"log"
	"case_proxy/load"
	"case_proxy/conf"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"case_proxy/rpc_proxy"
)
var config *conf.Config
var err error
var rpcProxy *rpc_proxy.PRCProxy

func main(){
	config, err = conf.ReadConfig("conf/config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}
	load.LoadPartition(config.PartitionFile)
	router := mux.NewRouter()
	for _, l := range config.Location {
		rpcProxy, err = rpc_proxy.NewPRCProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		// start health check

		router.Handle("/predict", rpcProxy)
	}
	// name := "World"
	svr := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}
	config.Print()

	err := svr.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve error: %s", err)
	}
}