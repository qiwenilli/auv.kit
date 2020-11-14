package auv

import (
	log "github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"

	"regexp"
)

var Options []Opt

type Option struct {
	ServiceName  string
	Services     []TwirpServer
	Interceptors []twirp.Interceptor
	DieHookFunc  func()
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

// func WithHostName(val string) Opt {
// 	return func(opt *Option) error {
// 		return nil
// 	}
// }

func WithServices(srv TwirpServer) Opt {
	return func(opt *Option) error {
		opt.Services = append(opt.Services, srv)
		return nil
	}
}

func WithInterceptors(interceptor twirp.Interceptor) Opt {
	return func(opt *Option) error {
		opt.Interceptors = append(opt.Interceptors, interceptor)
		return nil
	}
}

func WithDieHookFunc(_func func()) Opt {
	return func(opt *Option) error {
		opt.DieHookFunc = _func
		return nil
	}
}
