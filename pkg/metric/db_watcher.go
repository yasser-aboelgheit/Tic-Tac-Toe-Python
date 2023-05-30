package metric

import (
	"context"
	"database/sql"
	"time"

)

var (
	sqlDbStatsWaitCountCounter           = Counter{Name: "sql.db.stats.wait_count"}
	sqlDbStatsWaitDurationSecondsCounter = Counter{Name: "sql.db.stats.wait_duration_seconds"}
	sqlDbStatsMaxIdleClosedCounter       = Counter{Name: "sql.db.stats.max_idle_closed"}
	sqlDbStatsMaxIdletimeClosedCounter   = Counter{Name: "sql.db.stats.max_idletime_closed"}
	sqlDbStatsMaxLifetimeClosedCounter   = Counter{Name: "sql.db.stats.max_lifetime_closed"}
	sqlDbStatsMaxOpenConnectionsGauge    = NamedGauge{Name: "sql.db.stats.max_open_connections"}
	sqlDbStatsConnectionsOpenGauge       = NamedGauge{Name: "sql.db.stats.connections.open"}
	sqlDbStatsConnectionsInUseGauge      = NamedGauge{Name: "sql.db.stats.connections.in_use"}
	sqlDbStatsConnectionsIdleGauge       = NamedGauge{Name: "sql.db.stats.connections.idle"}
)

type dbLogger interface {
	Warnw(ctx context.Context, msg string, attrs map[string]interface{})
}

// WatchStats watches and reports db stats.
func WatchStats(ctx context.Context, db *sql.DB, name string, lgr dbLogger) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats := db.Stats()
			reportDBStats(name, stats)
		}
	}
}

func reportDBStats(dbname string, stats sql.DBStats) {
	const rate = 1.0

	tags := []string{"db.name:" + dbname}

	sqlDbStatsMaxOpenConnectionsGauge.Set(float64(stats.MaxOpenConnections), tags, rate)

	// Pool Status
	sqlDbStatsConnectionsOpenGauge.Set(float64(stats.OpenConnections), tags, rate)
	sqlDbStatsConnectionsInUseGauge.Set(float64(stats.InUse), tags, rate)
	sqlDbStatsConnectionsIdleGauge.Set(float64(stats.Idle), tags, rate)

	// Counters
	sqlDbStatsWaitCountCounter.Count(stats.WaitCount, tags, rate)
	sqlDbStatsWaitDurationSecondsCounter.Count(int64(stats.WaitDuration.Seconds()), tags, rate)
	sqlDbStatsMaxIdleClosedCounter.Count(stats.MaxIdleClosed, tags, rate)
	sqlDbStatsMaxIdletimeClosedCounter.Count(stats.MaxIdleTimeClosed, tags, rate)
	sqlDbStatsMaxLifetimeClosedCounter.Count(stats.MaxLifetimeClosed, tags, rate)
}
