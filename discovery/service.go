package discovery

import (
	"context"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/emirpasic/gods/maps/hashmap"
)

const (
	SERVICE_PREFIX = "naming"
)

var (
	Srvs = hashmap.New()
)

func RegisterServerName(serviceName, HostName string, serviceAddr string) error {
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Put(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key), serviceAddr)
}

func GetServerWithName(serviceName string) string {

	srvs := Srvs.Values()
	_ = srvs
	if srvsLen := len(srvs); srvsLen > 0 {
		idx := GetModVal(uint64(srvsLen))
		return srvs[idx].(string)
	}
	return ""
}

func DestroyService(serviceName, HostName string) error {
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Del(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key))
}

func ServiceEvent() WatchEvent {
	return func(rootPrefix string, event mvccpb.Event_EventType, key, val string) {
		//auv.naming/serviceName/hostname/addr
		if strings.Contains(key, rootPrefix+"."+SERVICE_PREFIX) {
			keySlice := strings.Split(key, "/")
			if mvccpb.DELETE == event {
				Srvs.Remove(keySlice[1])
			} else if mvccpb.PUT == event {
				Srvs.Put(keySlice[1], val)
			}
		}
	}
}
