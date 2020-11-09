package auv

import (
	// "context"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"sync"

	"github.com/keepeye/logrus-filename"
	"github.com/qiwenilli/auv.kit/utils"
	log "github.com/sirupsen/logrus"
	"reflect"
	// "github.com/twitchtv/twirp"
)

type TwirpServer interface {
	http.Handler
	ServiceDescriptor() ([]byte, int)
	ProtocGenTwirpVersion() string
	PathPrefix() string
}

var (
	FlagHelp        bool
	FlagServerAddr  string
	FlagPprofEnable bool
	FlagDebugLevel  string
	//
	onceServer   sync.Once
	serverEntity *server
)

const ()

func init() {

	log.SetOutput(os.Stdout)
	filenameHook := filename.NewHook()
	filenameHook.Field = "f"
	log.AddHook(filenameHook)
	//
	flag.BoolVar(&FlagHelp, "h", false, "")
	flag.StringVar(&FlagServerAddr, "server.addr", "", "http listen")
	flag.BoolVar(&FlagPprofEnable, "pprof.enable", false, "")
	flag.StringVar(&FlagDebugLevel, "log.level", "debug", "trace | debug | info | warn | error | fatal | panic")
}

type server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer() *server {
	onceServer.Do(func() {
		if flag.Parsed() {
			log.Fatal("flag parse() to doing ")
		}
		flag.Parse()
		if FlagHelp {
			flag.Usage()
			os.Exit(0)
		}
		//
		if logLevel, err := log.ParseLevel(FlagDebugLevel); err != nil {
			log.Fatal(err)
		} else {
			log.SetLevel(logLevel)
		}

		//
		serverEntity = &server{}
	})
	return serverEntity
}

func (s *server) WithService(srvs ...TwirpServer) *server {
	mux := http.NewServeMux()
	for _, srv := range srvs {
		if srv == nil {
			continue
		}
		mux.Handle(srv.PathPrefix(), srv)
		////
		typ := reflect.TypeOf(srv)
		for i := 0; i < typ.NumMethod(); i++ {
			log.Println(srv.PathPrefix(), typ.Method(i).Name, typ.Method(i).Type)
		}
	}
	s.mux = mux
	//
	return s
}

func (s *server) Run() {
	pprofService(s.mux)
	if FlagServerAddr != "" {
		if addrSlice := strings.Split(FlagServerAddr, ":"); len(addrSlice) == 2 && addrSlice[0] == "" {
			s.addr = fmt.Sprintf("%s%s", utils.LocalIP(), FlagServerAddr)
		} else {
			s.addr = fmt.Sprintf("%s", FlagServerAddr)
		}
	} else {
		s.addr = fmt.Sprintf("%s:8080", utils.LocalIP())
	}
	log.Info("listen: http://", s.addr)
	//
	http.ListenAndServe(s.addr, s.mux)
}

func pprofService(mux *http.ServeMux) {
	if !FlagPprofEnable {
		return
	}
	log.Info("pprof enable")
	log.Info("router: ", "/debug/pprof/")
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}
