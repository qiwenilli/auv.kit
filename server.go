package auv

import (
	// "context"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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
	flag.StringVar(&FlagServerAddr, "server.addr", "", "http listen (default: 8080)")
	flag.BoolVar(&FlagPprofEnable, "pprof.enable", false, "")
	flag.StringVar(&FlagDebugLevel, "log.level", "debug", "trace | debug | info | warn | error | fatal | panic")
}

type server struct {
	name        string
	addr        string
	mux         *http.ServeMux
	HttpServer  *http.Server
	dieHookFunc func()
}

func NewServer() *server {
	onceServer.Do(func() {
		if flag.Parsed() {
			log.Fatal("dont need exec : flag parse()")
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
		serverEntity = &server{
			HttpServer: &http.Server{},
		}
	})
	return serverEntity
}

func (s *server) WithName(name string) *server {
	s.name = name
	return s
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
			methodName := typ.Method(i).Name
			switch methodName {
			case "PathPrefix", "ProtocGenTwirpVersion", "ServeHTTP", "ServiceDescriptor":
				continue
			}
			log.Debugf("%s%s", srv.PathPrefix(), typ.Method(i).Name)
		}
	}
	s.mux = mux
	s.HttpServer.Handler = mux
	s.buildServer()
	//
	return s
}

func (s *server) WithDieHookFun(hook func()) *server {
	s.dieHookFunc = hook
	return s
}

func (s *server) buildServer() {
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
	s.HttpServer.Addr = s.addr
}

func (s *server) signal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPWR)
	go func() {
		fmt.Println("exit signal", <-c)
		signal.Stop(c)
		if s.dieHookFunc != nil {
			s.dieHookFunc()
		}
		os.Exit(0)
	}()
}

func (s *server) Run() {

	err := s.HttpServer.ListenAndServe()
	if err != nil {
		log.Error(err)
	} else {
		s.signal()
		log.Info("listen: http://", s.addr)
	}
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

func swagger() {

}
