package jsonrpc

import (
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/metrics"
)

func withMetrics[T any](
	metrics *metrics.ChainWebAPIMetrics,
	method string,
	f func() (T, error),
) (T, error) {
	started := time.Now()
	ret, err := f()
	metrics.EVMRPCCall(method, err == nil, time.Since(started))
	return ret, err
}
