package discovery

import (
	"context"
)

const (
	SERVICE_PREFIX = "naming"
)

func RegisterServerName(serviceName, HostName string, serviceAddr string) error {
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Put(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key), serviceAddr)
}

func GetServerWithName(serviceName string) []string {
	var serviceAddrs []string

	cli := NewEtcdClient()
	resp, err := cli.GetPrefix(context.TODO(), cli.BuildKey(SERVICE_PREFIX, serviceName))
	if err == nil {
		for _, v := range resp.Kvs {
			serviceAddrs = append(serviceAddrs, v.String())
		}
	}
	return serviceAddrs
}

func DestroyService(serviceName, HostName string) error {
	cli := NewEtcdClient()
	key := serviceName + "/" + HostName
	return cli.Del(context.TODO(), cli.BuildKey(SERVICE_PREFIX, key))
}
