package discovery

import (
	"context"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/emirpasic/gods/maps/hashmap"
)

const (
	CONFIG_PREFIX = "config"
)

var (
	kvs = hashmap.New()
)

func Put(key string, val string) error {
	cli := NewEtcdClient()
	return cli.Put(context.TODO(), cli.BuildKey(CONFIG_PREFIX, key), val)
}

func Get(key string) string {
	if val, ok := kvs.Get(key); ok {
		return val.(string)
	}
	return ""
}

func ConfigEvent() WatchEvent {
	return func(rootPrefix string, event mvccpb.Event_EventType, key, val string) {
		//auv.config/serviceName
		if strings.Contains(key, rootPrefix+"."+CONFIG_PREFIX) {
			keySlice := strings.Split(key, "/")
			if mvccpb.DELETE == event {
				kvs.Remove(keySlice[1])
			} else if mvccpb.PUT == event {
				kvs.Put(keySlice[1], val)
			}
		}
	}
}
