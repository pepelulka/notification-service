package worker

import (
	"bytes"
	"context"
	"fmt"
	"notification_service/internal/config"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdClient struct {
	Mu     sync.Mutex
	Client *clientv3.Client
}

func createEtcdClient(cfg *config.EtcdConfig) (etcdClient, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return etcdClient{}, err
	}
	return etcdClient{Mu: sync.Mutex{}, Client: cli}, nil
}

func (cli *etcdClient) Close() error {
	return cli.Client.Close()
}

func (cli *etcdClient) Get(ctx context.Context, key string) (string, bool, error) {
	var (
		resp *clientv3.GetResponse
		err  error
	)
	func() {
		cli.Mu.Lock()
		defer cli.Mu.Unlock()
		resp, err = cli.Client.Get(ctx, key)
	}()
	if err != nil {
		return "", false, err
	}
	for _, kv := range resp.Kvs {
		if bytes.Equal(kv.Key, []byte(key)) {
			return string(kv.Value), true, nil
		}
	}
	return "", false, nil
}

func (cli *etcdClient) Set(ctx context.Context, key string, value string) error {
	var err error
	func() {
		cli.Mu.Lock()
		defer cli.Mu.Unlock()
		_, err = cli.Client.Put(ctx, key, value)
	}()
	return err
}
