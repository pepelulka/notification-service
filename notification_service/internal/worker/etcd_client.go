package worker

import "sync"

type EtcdClient struct {
	mu sync.Mutex
	kApi *clientv3. // Ended here...
}
