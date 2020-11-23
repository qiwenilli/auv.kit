package discovery

import (
	"context"
	"regexp"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/emirpasic/gods/sets/hashset"
)

const (
	SERVICE_PREFIX = "naming"
)

var (
	srvs = make(map[string]*hashset.Set, 1)
)

func RegisterServerName(serviceName, HostName string, serviceAddr string) error {
	if ok, err := regexp.MatchString("[0-9A-Fa-f\\._-]", serviceName); !ok {
		return err
	}
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Put(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key), serviceAddr)
}

func GetServerWithName(serviceName string) string {
	if list, ok := srvs[serviceName]; ok {
		srvSlice := list.Values()
		if srvsLen := len(srvSlice); srvsLen > 0 {
			idx := GetModVal(uint64(srvsLen))
			return srvSlice[idx].(string)
		}
	}
	return ""
}

func DestroyService(serviceName, HostName string) error {
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Del(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key))
}

func sourceServicePut(key, val string) {
	keySlice := strings.Split(key, "/")
	serviceName := keySlice[1]
	list, ok := srvs[serviceName]
	if !ok {
		list = hashset.New()
		srvs[serviceName] = list
	}
	list.Add(val)
}

func ServiceEvent() WatchEvent {
	return func(rootPrefix string, event mvccpb.Event_EventType, key, val string) {
		//auv.naming/serviceName/hostname/addr
		if strings.Contains(key, rootPrefix+"."+SERVICE_PREFIX) {
			keySlice := strings.Split(key, "/")
			serviceName := keySlice[1]
			list, ok := srvs[serviceName]
			if !ok {
				list = hashset.New()
				srvs[serviceName] = list
			}
			if mvccpb.DELETE == event {
				list.Remove(val)
			} else if mvccpb.PUT == event {
				list.Add(val)
			}
		}
	}
}
