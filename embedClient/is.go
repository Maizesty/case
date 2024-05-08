package main

import (
	"context"
	"embed_client/conf"
	embed "embed_client/embed"
	is "embed_client/infer"

	// "time"
	tf "github.com/galeone/tensorflow/tensorflow/go"
)


func floatsConcat(A [] float32,B [] float32,index int){
	cnt := 0
	for i := index * conf.EmbDim;i< (index +1 ) * conf.EmbDim;i++{
		A[i] = B[cnt]
		cnt += 1
	}
}


func cacheQuery(keys []int) ([][]float32,[]bool,bool){
	result := make([][]float32, 1)
	result[0] = make([]float32, conf.EmbDim * conf.EmbNum)
	need_key := make([]bool, len(keys))
	perfect := true
	for i, key := range keys {
		if lv, ok := cache.Get(key); ok {
			if vals ,ok := lv.([]float32);ok{
				floatsConcat(result[0],vals,i)
				// result[i] = vals
			}
		}else{
			perfect = false
			need_key[i] = true
		}
		
	}
	return result,need_key,perfect
}


func queryPS(keys []int,needKeys [] bool)map[int][]float32{
	result := make(map[int][]float32)
	var queryKeys []int32
	var queryIndex []int
	for k,v := range(needKeys){
		if v {
			queryKeys = append(queryKeys, int32(keys[k]))
			queryIndex = append(queryIndex, k)
		}
	}
	// ctx, _ := context.WithTimeout(context.Background(), time.Second)
	req := &embed.EmbReqData{Keys: queryKeys}
	r, err := EmbClient.Lookup(Ctx,req)
	if err != nil {
		return nil
	} 
	
	eles :=r.EmbVectors
	for k,v := range(eles){
		key := queryKeys[k]
		val := v.Element
		cache.Set(key,val,1)
		result[queryIndex[k]] = val
	}
	return result
}


func serveModel(input [][]float32) interface{} {
	if Model!=nil {
		inputTensor, _ := tf.NewTensor(input)
		results := Model.Exec([]tf.Output{
				Model.Op("StatefulPartitionedCall", 0),
		}, map[tf.Output]*tf.Tensor{
				Model.Op("serving_default_input_2", 0): inputTensor,
		})
		predictions := results[0]
		// fmt.Println(predictions.Value())
		return predictions.Value()
	}
	return nil
}


type server struct{ 
	is.MyServiceServer
}

func (s *server) Infer(ctx context.Context, in *is.ReqData) (*is.Resp, error) {
	keys := make([]int,len(in.Keys))
	for k,v := range(in.Keys){
		keys[k] = int(v)
	}
	model_input, need_keys, perfect := cacheQuery(keys)
	psResult := queryPS(keys,need_keys)
	for k, v :=range(psResult){
		floatsConcat(model_input[0],v,k)
		// model_input[k] = v
	}
	// fmt.Println(len(model_input[0]))
	result := serveModel(model_input).([][]float32)
	CntTotal.Add(1)
	if perfect{
		CntHit.Add(1)
	}

	return &is.Resp{Output:  result[0][0] ,Perfect: perfect}, nil
}
