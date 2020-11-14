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
	FlagSwaggerUiEnable  bool

	FlagLogType         string
	FlagLogPath         string
	FlagAccessLogEnable bool
)

func init() {

	flag.BoolVar(&FlagHelp, "h", false, "")
	flag.StringVar(&FlagServerAddr, "server.addr", ":8080", "http listen")
	flag.BoolVar(&FlagPprofEnable, "pprof.enable", false, "pprof.enable (default: false)")
	flag.BoolVar(&FlagAllowCrossDomain, "allow.cross.domain.enable", false, "(default: false)")
	flag.BoolVar(&FlagSwaggerUiEnable, "swaggerui.enable", false, "(default: false)")

	flag.StringVar(&FlagLogType, "log.type", "file", "nginx log")
	flag.StringVar(&FlagLogPath, "log.path", ".", "save path for log")
	flag.StringVar(&FlagDebugLevel, "log.level", "debug", "trace | debug | info | warn | error | fatal | panic for work_log")
	flag.BoolVar(&FlagAccessLogEnable, "access.log.enable", false, "nginx access_log")
}

type Config struct {
}
