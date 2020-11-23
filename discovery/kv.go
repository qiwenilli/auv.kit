package discovery

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/emirpasic/gods/maps/hashmap"
)

const (
	CONFIG_PREFIX = "config"
)

var (
	kvs        = hashmap.New()
	key_reg, _ = regexp.Compile("[0-9A-Fa-f\\._-]")
)

func KvPut(key string, val string) error {
	if key_reg.MatchString(key) {
		cli := NewEtcdClient()
		return cli.Put(context.TODO(), cli.BuildKey(CONFIG_PREFIX, key), val)
	} else {
		return errors.New("key not allow")
	}
}

func KvGet(key string) string {
	if val, ok := kvs.Get(key); ok {
		return val.(string)
	}
	return ""
}

func sourceKvPut(key, val string) {
	keySlice := strings.Split(key, "/")
	kvs.Put(keySlice[1], val)
}

func ConfigEvent() WatchEvent {
	return func(rootPrefix string, event mvccpb.Event_EventType, key, val string) {
		//auv.config/key
		cli := NewEtcdClient()
		if strings.HasPrefix(key, cli.BuildKey(CONFIG_PREFIX, "")) {
			keySlice := strings.Split(key, "/")
			if mvccpb.DELETE == event {
				kvs.Remove(keySlice[1])
			} else if mvccpb.PUT == event {
				kvs.Put(keySlice[1], val)
			}
		}
	}
}
