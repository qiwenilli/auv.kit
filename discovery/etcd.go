package discovery

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/pkg/transport"

	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type etcdClient struct {
	client     *clientv3.Client
	rootPrefix string
}

type WatchEvent func(rootPrefix string, event mvccpb.Event_EventType, key, val string)

var (
	etcdClientEntity *etcdClient
	onceEtcd         sync.Once
	error_disconnet  = errors.New("etcd client disconnect")
)

func InitEtcd(rootPrefix string, endpoints, key, cert, ca string, events ...WatchEvent) {

	var tlsConfig *tls.Config
	var err error
	if key != "" && cert != "" && ca != "" {
		tlsInfo := transport.TLSInfo{
			// CertFile:      "/tmp/test-certs/test-name-1.pem",
			// KeyFile:       "/tmp/test-certs/test-name-1-key.pem",
			// TrustedCAFile: "/tmp/test-certs/trusted-ca.pem",
			CertFile:      cert,
			KeyFile:       key,
			TrustedCAFile: ca,
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	onceEtcd.Do(func() {
		config := clientv3.Config{
			Endpoints:   strings.Split(endpoints, ","),
			DialTimeout: time.Second * 2,
		}
		if tlsConfig != nil {
			config.TLS = tlsConfig
		}

		cli, err := clientv3.New(config)
		if err != nil {
			fmt.Println(err)
			return
		}

		etcdClientEntity = &etcdClient{
			client:     cli,
			rootPrefix: rootPrefix,
		}

		go func() {
			for {
				for _, ep := range config.Endpoints {
					ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*500)
					// resp, err := cli.Status(ctx, ep)
					_, err := cli.Status(ctx, ep)
					if err != nil {
						fmt.Println(err)
						cancel()
						continue
					}
					// fmt.Println(resp.Header.MemberId)
					// fmt.Printf("endpoint: %s / Leader: %v\n", ep, resp.Header.MemberId == resp.Leader)
				}
				time.Sleep(time.Second)
			}
		}()

		go func() {
			rch := cli.Watch(context.TODO(), rootPrefix, clientv3.WithPrefix())
			for wresp := range rch {
				for _, ev := range wresp.Events {
					for _, event := range events {
						event(rootPrefix, ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
					}

					// switch ev.Type {
					// case mvccpb.PUT:
					// 	// ev.Kv.Key, ev.Kv.Value
					// case mvccpb.DELETE:
					// 	// ev.Kv.Key, ev.Kv.Value
					// }
				}
			}
			//
		}()

		// get 初始化所有servicename config
		if resp, err := cli.Get(context.TODO(), rootPrefix+"."+CONFIG_PREFIX, clientv3.WithPrefix()); err == nil {
			for _, kv := range resp.Kvs {
				sourceKvPut(string(kv.Key), string(kv.Value))
			}
		}

		// get 初始化所有servicename config
		if resp, err := cli.Get(context.TODO(), rootPrefix+"."+SERVICE_PREFIX, clientv3.WithPrefix()); err == nil {
			for _, kv := range resp.Kvs {
				sourceServicePut(string(kv.Key), string(kv.Value))
			}
		}

	})
}

func NewEtcdClient() *etcdClient {
	return etcdClientEntity
}

func (e *etcdClient) GetRootPrefix() string {
	return e.rootPrefix
}

func (e *etcdClient) BuildKey(typeName, key string) string {
	// auv.naming/Service.Name/Service.Version/hostname
	// auv.config/config_name
	// type = naming | config
	if e == nil {
		return ""
	}
	return e.rootPrefix + "." + typeName + "/" + key
}

func (e *etcdClient) Get(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	if e == nil {
		return nil, error_disconnet
	}
	return e.get(ctx, key, false)
}

func (e *etcdClient) GetPrefix(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	return e.get(ctx, key, true)
}

func (e *etcdClient) get(ctx context.Context, key string, isPrefix bool) (*clientv3.GetResponse, error) {
	return e.client.Get(ctx, key, clientv3.WithPrefix())
}

func (e *etcdClient) Put(ctx context.Context, key, val string) error {
	if e == nil {
		return error_disconnet
	}
	_, err := e.client.Put(ctx, key, val)
	return err
}

func (e *etcdClient) Del(ctx context.Context, key string) error {
	if e == nil {
		return error_disconnet
	}
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *etcdClient) Watch(ctx context.Context, key string, isPrefix bool) clientv3.WatchChan {
	if e == nil {
		return nil
	}
	//
	var rch clientv3.WatchChan
	if isPrefix {
		rch = e.client.Watch(ctx, key, clientv3.WithPrefix())
	} else {
		rch = e.client.Watch(ctx, key)
	}
	return rch
	// for wresp := range rch {
	// 	for _, ev := range wresp.Events {
	// 		switch ev.Type {
	// 		case mvccpb.PUT:
	// 			// ev.Kv.Key, ev.Kv.Value
	// 		case mvccpb.DELETE:
	// 			// ev.Kv.Key, ev.Kv.Value

	// 		}
	// 	}
	// }
	// return nil
}
