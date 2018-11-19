package system

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdStore struct {
	EtcdClient *clientv3.Client
}

var etcdStore *EtcdStore

func InitEtcd(conf map[string]interface{}) {
	if etcdStore.EtcdClient != nil {
		return
	}
	etcdStore = new(EtcdStore)
	etcdConf := clientv3.Config{
		Endpoints:   conf["endpoints"].([]string),
		DialTimeout: time.Duration(5),
	}
	client, err := clientv3.New(etcdConf)

	if err != nil {
		panic("etcd init fail")
	}
	etcdStore.EtcdClient = client

}
