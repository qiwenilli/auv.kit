package Config

import (
	"flag"
)

var (
	FlagHelp             bool
	FlagServerAddr       string
	FlagPprofEnable      bool
	FlagDebugLevel       string
	FlagAllowCrossDomain bool

	FlagLogType         string
	FlagAccessLogEnable bool
)

func init() {

	flag.BoolVar(&FlagHelp, "h", false, "")
	flag.StringVar(&FlagServerAddr, "server.addr", ":8080", "http listen")
	flag.BoolVar(&FlagPprofEnable, "pprof.enable", false, "")
	flag.BoolVar(&FlagAllowCrossDomain, "allow.cross.domain.enable", false, "")
	flag.StringVar(&FlagLogType, "log.type", "file", "nginx log")
	flag.StringVar(&FlagDebugLevel, "log.level", "debug", "trace | debug | info | warn | error | fatal | panic for work_log")
	flag.BoolVar(&FlagAccessLogEnable, "access.log.enable", false, "nginx access_log")
}

type Config struct {
}
