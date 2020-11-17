package auv

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
	// "unsafe"

	"github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
	// "github.com/twitchtv/twirp"
	"github.com/arl/statsviz"
	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mkevac/debugcharts"

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
	s.withWorkLog(serverOpt.ServiceName)
	if auvconfig.FlagPprofEnable {
		s.withPprof()
	}
	if auvconfig.FlagSwaggerUiEnable {
		s.withSwaggerUi()
	}
	s.withService(serverOpt.Services...)
	s.withSignal(serverOpt.DieHookFunc)

	//
	var middlewares []mux.MiddlewareFunc
	middlewares = append(middlewares, auvhttp.MiddlewareTraceId)
	if auvconfig.FlagAllowCrossDomain {
		middlewares = append(middlewares, auvhttp.MiddlewareForCrossDomain)
	}
	middlewares = append(middlewares, serverOpt.Middlewares...)
	s.mux.Use(middlewares...)

	// s.withAccessLog(serverOpt.ServiceName)

	// use gzip
	// handlers.CompressHandler(http.DefaultServeMux)

	s.handler = s.mux
	s.HttpServer.Handler = s.handler
	for _, path := range s.pathRules {
		log.Info(path)
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
			s.pathRules = append(s.pathRules, fmt.Sprintf("%s%s", srv.PathPrefix(), typ.Method(i).Name))
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

func (s *server) withWorkLog(serviceName string) {
	work_log := utils.CreateRotatelogs("work_log", serviceName, auvconfig.FlagLogPath)
	stdout := io.MultiWriter(os.Stdout, work_log)
	log.SetOutput(stdout)
	log.SetFormatter(new(internal.LogFormatter))
	// add filename to log
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	// add serviceName hook
	serviceNameHook := &internal.ServiceNameHook{ServiceName: serviceName}
	log.AddHook(serviceNameHook)
}

func (s *server) withPprof() {
	s.WithHandle("/debug/statsviz/", statsviz.Index)
	s.WithHandleFunc("/debug/statsviz/ws", statsviz.Ws)

	debugRouter := s.mux.PathPrefix("/debug/")
	debugRouter.Handler(http.DefaultServeMux)

	vv := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m")
	if vv.Kind() == reflect.Map {
		for _, element := range vv.MapKeys() {
			log.Info(element)
		}
	}
}

func (s *server) withSwaggerUi() {
	uiHandle := func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(200)

		html, _ := utils.FlateDecode(internal.SwggerHtml)
		resp.Write(html)
	}
	s.WithHandleFunc("/swagger", uiHandle)
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
