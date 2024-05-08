package rpc_scheduler

const (
	// ContentAwareCoarseScheduler = "content-aware-coarse"
	// ContentAwareScheduler       = "content-aware"
	// P2CScheduler                = "p2c"
	// RandomScheduler             = "random"
	R2Scheduler                 = "round-robin"
	// LeastLoadScheduler          = "least-load"
	CommiunityAwareScheduler		= "commiunity-aware"
	// CommiunityAwareP2CScheduler		= "commiunity-aware-p2c"
	NaiveCacheScheduler										= "NaiveCache"
)


var Modes =[] string{
	// "content-aware-coarse",
	// "content-aware",
	// "p2c",
	// "random",
	"round-robin",
	// "least-load",
	"commiunity-aware",
	// "commiunity-aware-p2c",
	"NaiveCache",

}