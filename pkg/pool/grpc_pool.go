package pool

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pool struct {
	conns []*grpc.ClientConn
	mu    sync.Mutex
	index int
	size  int
}

func NewPool(targetAddress string, poolSize int) (*Pool, error) {
	pool := &Pool{
		conns: make([]*grpc.ClientConn, poolSize),
		size:  poolSize,
	}

	for i := 0; i < poolSize; i++ {
		conn, err := grpc.NewClient(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy": "round_robin",
			"methodConfig": [{
			  "retryPolicy": {
				"maxAttempts": 4,
				"initialBackoff": "0.1s",
				"maxBackoff": "1s",
				"backoffMultiplier": 2.0,
				"retryableStatusCodes": ["UNAVAILABLE"]
			  }
			}]
		  }`))
		if err != nil {
			pool.CloseAll()
			return nil, err
		}
		pool.conns[i] = conn
	}

	return pool, nil
}

func (p *Pool) Get() *grpc.ClientConn {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn := p.conns[p.index]
	p.index = (p.index + 1) % p.size
	return conn
}

func (p *Pool) CloseAll() {
	for _, conn := range p.conns {
		if conn != nil {
			conn.Close()
		}
	}
}
