package scc

import (
	"fmt"
	"strings"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/smallnest/rpcx/codec"
)

func init() {
	zookeeper.Register()
}

type zookeeperStore struct {
	zookeeperServers []string
	options          *store.Config
	kv               store.Store
	codec            Codec
	appName          string
	watchList        map[string]chan struct{}
}

func NewZookeeperStore(appName string, zookeeperServers []string) (Store, error) {
	var (
		c = &zookeeperStore{
			zookeeperServers: zookeeperServers,
			codec:            codec.JSONCodec{},
			watchList:        make(map[string]chan struct{}, 0),
		}
		kv  store.Store
		err error
	)
	appName = strings.Trim(appName, "/")
	if appName == "" {
		return nil, ErrUnspecifiedAppName
	}

	c.appName = appName

	if kv, err = libkv.NewStore(store.ZK, c.zookeeperServers, c.options); err != nil {
		return nil, err
	}
	c.kv = kv

	return c, nil
}

func (c *zookeeperStore) SetOption(options *store.Config) error {
	c.options = options
	kv, err := libkv.NewStore(store.ETCD, c.zookeeperServers, c.options)
	if err != nil {
		return err
	}
	c.kv = kv
	return nil
}

func (c *zookeeperStore) SetCodec(codec Codec) {
	c.codec = codec
}

func (c *zookeeperStore) GetCodec() Codec {
	return c.codec
}

func (c *zookeeperStore) combinKey(key string) string {
	return strings.Trim(fmt.Sprintf("%s/%s", c.appName, strings.Trim(key, "/")), "/")
}

func (c *zookeeperStore) Get(key string, value interface{}) error {
	var (
		kvPaire *store.KVPair
		err     error
	)

	kvPaire, err = c.kv.Get(c.combinKey(key))
	if err != nil {
		return err
	}
	if c.codec == nil {
		return ErrUnspecifiedCodec
	}
	return c.codec.Decode(kvPaire.Value, value)
}
func (c *zookeeperStore) Set(key string, value interface{}) error {
	var (
		data []byte
		err  error
	)
	if data, err = c.codec.Encode(value); err != nil {
		return err
	}
	return c.kv.Put(c.combinKey(key), data, &store.WriteOptions{IsDir: false})
}

func (c *zookeeperStore) Watch(key string, cb func(data []byte)) (err error) {
	var (
		stopCh chan struct{}
		events <-chan *store.KVPair
	)
	if c.codec == nil {
		return ErrUnspecifiedCodec
	}

	stopCh = make(chan struct{})
	c.watchList[key] = stopCh
	if events, err = c.kv.Watch(c.combinKey(key), stopCh); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case pair := <-events:
				if pair == nil {
					continue
				}
				cb(pair.Value)
			}
		}
	}()
	return nil
}

func (c *zookeeperStore) StopWatch(key string) error {
	var (
		stopCh chan struct{}
		ok     bool
	)
	if stopCh, ok = c.watchList[key]; !ok {
		return ErrStopWatchFailed
	}
	stopCh <- struct{}{}
	return nil
}
