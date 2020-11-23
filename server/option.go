package server

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"regexp"
)

var Options []Opt

type Option struct {
	ServiceName  string
	Services     []TwirpServer
	Middlewares  []mux.MiddlewareFunc
	DieHookFuncs []func()
	//
	EtcdEndpoints []string
}

type Opt func(*Option) error

func WithServiceName(val string) Opt {
	return func(opt *Option) error {
		if ok, _ := regexp.MatchString(`^[[:word:]\.]+$`, val); ok {
			opt.ServiceName = val
		} else {
			log.Fatal("allow [0-9A-Za-z_\\.]")
		}
		return nil
	}
}

func WithServices(srvs ...TwirpServer) Opt {
	return func(opt *Option) error {
		opt.Services = append(opt.Services, srvs...)
		return nil
	}
}

func WithDieHookFunc(_func func()) Opt {
	return func(opt *Option) error {
		opt.DieHookFuncs = append(opt.DieHookFuncs, _func)
		return nil
	}
}

func WithMiddlewares(middlewareFuncs ...mux.MiddlewareFunc) Opt {
	return func(opt *Option) error {
		opt.Middlewares = append(opt.Middlewares, middlewareFuncs...)
		return nil
	}
}

func WithEtcdEndpoints(endpoints []string) Opt {
	return func(opt *Option) error {
		opt.EtcdEndpoints = endpoints
		return nil
	}
}
