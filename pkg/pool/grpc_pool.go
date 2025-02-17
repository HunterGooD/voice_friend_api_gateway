package pool

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// Pool TODO: pool connection with generic for Rest and others clients
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

func (p *Pool) GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startIndex := p.index

	for {
		conn := p.conns[p.index]
		p.index = (p.index + 1) % p.size

		// if  Ready to return
		if conn != nil && conn.GetState() == connectivity.Ready {
			return conn, nil
		}

		// if any connection not ready return error
		if p.index == startIndex {
			return nil, errors.Wrap(grpc.ErrClientConnClosing, "Error all connection is closing")
		}

		if ctx.Err() != nil {
			return nil, errors.Wrap(ctx.Err(), "Error getting connection context close")
		}
	}
}

func (p *Pool) CloseAll() {
	for _, conn := range p.conns {
		if conn != nil {
			conn.Close()
		}
	}
}
