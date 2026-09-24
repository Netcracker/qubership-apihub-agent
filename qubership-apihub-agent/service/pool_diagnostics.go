package service

import (
	"time"

	"github.com/Netcracker/qubership-apihub-agent/utils"
)

const poolDiagnosticsInterval = time.Minute

// RunConnectionPoolDiagnostics periodically logs HTTP connection pool counters and
// OS-level open connection stats, to aid diagnosing connection exhaustion/leaks.
func RunConnectionPoolDiagnostics() {
	utils.SafeAsync(func() {
		for range time.Tick(poolDiagnosticsInterval) {
			utils.SafeAsync(func() {
				utils.LogPoolStats()
				utils.LogOSConnStats()
			})
		}
	})
}
