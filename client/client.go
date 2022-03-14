package client

import (
	"sync"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

type Conf struct {
	Name         string        `yaml:"name"`
	ServerName   string        `yaml:"server_name"`
	RegistryAddr string        `yaml:"registry_addr"`
	PoolTTL      time.Duration `yaml:"pool_ttl"`
	Retries      int           `yaml:"retries"`
	PoolSize     int           `yaml:"pool_size"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	Client       client.Client
}

var ServerConf = make(map[string]Conf)

var RpcPoolLock = new(sync.RWMutex)

func NewClients(conf []Conf) error {
	for k := range conf {
		conf[k].NewRpcClient()
		ServerConf[conf[k].Name] = conf[k]
	}
	return nil
}

func GetClient(name string) (client.Client,bool) {
	c, ok := ServerConf[name]
	return c,ok
}

func GetConf(name string) Conf {
	return ServerConf[name]
}

func (c *Conf) NewRpcClient() {
	client.NewClient = grpc.NewClient
	// etcd 注册
	registry := etcd.NewRegistry(registry.Addrs(c.RegistryAddr))
	c.buildDefault()
	// 初始化客户端
	conn := client.NewClient(
		client.Registry(registry),         // 配置注册
		client.Retries(c.Retries),         // 重试次数
		client.DialTimeout(c.DialTimeout), //超时时间
		client.PoolTTL(c.PoolTTL),
		client.PoolSize(c.PoolSize),
	)
	c.Client = conn
}

func (c *Conf) buildDefault() {
	if c.DialTimeout == 0 {
		c.DialTimeout = 10 * time.Second
	}
	if c.PoolSize == 0 {
		c.PoolSize = 10
	}
	if c.PoolTTL == 0 {
		c.PoolTTL = 10 * time.Second
	}
}
