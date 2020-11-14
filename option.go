package auv

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"regexp"
)

var Options []Opt

type Option struct {
	ServiceName string
	Services    []TwirpServer
	Middlewares []mux.MiddlewareFunc
	DieHookFunc func()
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

func WithServices(srv TwirpServer) Opt {
	return func(opt *Option) error {
		opt.Services = append(opt.Services, srv)
		return nil
	}
}

func WithDieHookFunc(_func func()) Opt {
	return func(opt *Option) error {
		opt.DieHookFunc = _func
		return nil
	}
}

func WithMiddlewares(middlewareFunc mux.MiddlewareFunc) Opt {
	return func(opt *Option) error {
		opt.Middlewares = append(opt.Middlewares, middlewareFunc)
		return nil
	}
}
