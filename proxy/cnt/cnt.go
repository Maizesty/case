package cnt

import "sync"


type Smap struct{
	sync.RWMutex
	M map[string]int
}


var Cnter  = &Smap{M: make(map[string]int)}