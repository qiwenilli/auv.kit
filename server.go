package auv

import (
	// "context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
	// "github.com/twitchtv/twirp"
	"github.com/gorilla/mux"

	apachelog "github.com/lestrrat-go/apache-logformat"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/qiwenilli/auv.kit/internal"
	auvconfig "github.com/qiwenilli/auv.kit/internal/config"
	auvhttp "github.com/qiwenilli/auv.kit/internal/http"
	"github.com/qiwenilli/auv.kit/utils"
)

type TwirpServer interface {
	http.Handler
	ServiceDescriptor() ([]byte, int)
	ProtocGenTwirpVersion() string
	PathPrefix() string
}

var (
	onceServer   sync.Once
	serverEntity *server
)

type server struct {
	addr        string
	pathRules   []string
	mux         *mux.Router
	handler     http.Handler
	dieHookFunc func()

	HttpServer *http.Server
}

func NewServer() *server {
	onceServer.Do(func() {
		if flag.Parsed() {
			log.Fatal("dont need exec : flag parse()")
		}
		flag.Parse()
		if auvconfig.FlagHelp {
			flag.Usage()
			os.Exit(0)
		}
		//
		if logLevel, err := log.ParseLevel(auvconfig.FlagDebugLevel); err != nil {
			log.Fatal(err)
		} else {
			log.SetLevel(logLevel)
		}

		//
		serverEntity = &server{
			HttpServer: &http.Server{
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			},
			mux: mux.NewRouter(),
		}
	})
	return serverEntity
}

func (s *server) Run(opts ...Opt) {

	serverOpt := &Option{}
	for _, optFunc := range opts {
		optFunc(serverOpt)
	}
	log.Info(serverOpt)
	//
	s.mux = mux.NewRouter()
	if auvconfig.FlagPprofEnable {
		s.withPprof()
	}
	s.withService(serverOpt.Services...)
	s.withSignal(serverOpt.DieHookFunc)
	if auvconfig.FlagAllowCrossDomain {
		s.mux.Use(auvhttp.MiddlewareForCrossDomain)
	}
	if auvconfig.FlagAccessLogEnable {
		s.withAccessLog(serverOpt.ServiceName)
		s.HttpServer.Handler = s.handler
	} else {
		s.HttpServer.Handler = s.mux
	}
	s.withAddr()
	err := s.HttpServer.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}

func (s *server) withService(srvs ...TwirpServer) *server {
	for _, srv := range srvs {
		if srv == nil {
			continue
		}
		s.WithHandle(srv.PathPrefix(), srv)
		//
		typ := reflect.TypeOf(srv)
		for i := 0; i < typ.NumMethod(); i++ {
			methodName := typ.Method(i).Name
			switch methodName {
			case "PathPrefix", "ProtocGenTwirpVersion", "ServeHTTP", "ServiceDescriptor":
				continue
			}
			log.Infof("%s%s", srv.PathPrefix(), typ.Method(i).Name)
		}
	}
	return s
}

func (s *server) withAddr() {
	if addrSlice := strings.Split(auvconfig.FlagServerAddr, ":"); len(addrSlice) == 2 && addrSlice[0] == "" {
		s.addr = fmt.Sprintf("%s%s", utils.LocalIP(), auvconfig.FlagServerAddr)
	} else {
		s.addr = fmt.Sprintf("%s", auvconfig.FlagServerAddr)
	}
	s.HttpServer.Addr = s.addr
	log.Info("listen: http://", s.addr)
}

func (s *server) withSignal(dieHookFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPWR)
	go func() {
		switch <-c {
		case os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPWR:
			log.Info("exit program...")
			signal.Stop(c)
			if dieHookFunc != nil {
				dieHookFunc()
			}
			os.Exit(0)
		}
	}()
}

func (s *server) withAccessLog(serviceName string) {

	createRotatelogs := func(log_prefix string) *rotatelogs.RotateLogs {
		log_f := fmt.Sprintf("%s_%s", log_prefix, serviceName)
		rl, err := rotatelogs.New(
			log_f+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(log_f),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			log.Fatalf("failed to create rotatelogs: %s", err)
		}
		return rl
	}

	access_log := createRotatelogs("access_log")
	s.handler = apachelog.CombinedLog.Wrap(s.mux, access_log)

	//
	work_log := createRotatelogs("work_log")
	stdout := io.MultiWriter(os.Stdout, work_log)
	log.SetOutput(stdout)
	log.SetFormatter(new(internal.LogFormatter))
	// add filename to log
	filenameHook := filename.NewHook()
	filenameHook.Field = "f"
	log.AddHook(filenameHook)
}

func (s *server) withPprof() {
	s.WithHandleFunc("/debug/pprof/", pprof.Index)
	s.WithHandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.WithHandleFunc("/debug/pprof/profile", pprof.Profile)
	s.WithHandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.WithHandleFunc("/debug/pprof/trace", pprof.Trace)

	s.WithHandle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	s.WithHandle("/debug/pprof/heap", pprof.Handler("heap"))
	s.WithHandle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	s.WithHandle("/debug/pprof/block", pprof.Handler("block"))
}

func swagger() {

}

func (s *server) WithHandle(path string, handler http.Handler) *server {
	s.pathRules = append(s.pathRules, path)
	s.mux.Handle(path, handler)
	return s
}

func (s *server) WithHandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *server {
	s.pathRules = append(s.pathRules, path)
	s.mux.HandleFunc(path, f)
	return s
}
