package utils

import (
	"net/http"
	"net/http/httptrace"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// PoolStats holds live connection counters for a named HTTP transport.
type PoolStats struct {
	InFlight      int64
	TotalRequests int64
	NewConns      int64
	ReusedConns   int64
}

var (
	poolStatsMu sync.Mutex
	poolStats   = map[string]*PoolStats{}
)

type instrumentedTransport struct {
	underlying http.RoundTripper
	stats      *PoolStats
}

func (t *instrumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt64(&t.stats.InFlight, 1)
	atomic.AddInt64(&t.stats.TotalRequests, 1)
	defer atomic.AddInt64(&t.stats.InFlight, -1)

	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			if info.Reused {
				atomic.AddInt64(&t.stats.ReusedConns, 1)
			} else {
				atomic.AddInt64(&t.stats.NewConns, 1)
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	return t.underlying.RoundTrip(req)
}

// RegisterPoolStats wraps rt with connection-pool counters tracked under name and returns
// the wrapper to use as the transport in its place.
func RegisterPoolStats(name string, rt http.RoundTripper) http.RoundTripper {
	stats := &PoolStats{}
	poolStatsMu.Lock()
	poolStats[name] = stats
	poolStatsMu.Unlock()
	return &instrumentedTransport{underlying: rt, stats: stats}
}

// LogPoolStats logs the current counters for every registered HTTP connection pool.
func LogPoolStats() {
	poolStatsMu.Lock()
	defer poolStatsMu.Unlock()
	for name, stats := range poolStats {
		log.Infof("HTTP pool %q: in_flight=%d total_requests=%d new_conns=%d reused_conns=%d",
			name,
			atomic.LoadInt64(&stats.InFlight),
			atomic.LoadInt64(&stats.TotalRequests),
			atomic.LoadInt64(&stats.NewConns),
			atomic.LoadInt64(&stats.ReusedConns),
		)
	}
}
