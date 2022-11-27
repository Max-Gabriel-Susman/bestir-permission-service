package database

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/Max-Gabriel-Susman/bestir-go-kit/bestirlog"
	"go.uber.org/zap"
)

type StatsReporter struct {
	db  *sql.DB
	ddc *statsd.Client
	lg  *bestirlog.ZapLogger
}

func NewStatsReporter(db *sql.DB, ddc *statsd.Client, lg *bestirlog.ZapLogger) *StatsReporter {
	return &StatsReporter{db: db, ddc: ddc, lg: lg}
}

func (sr *StatsReporter) EmitStats(_ context.Context, tags []string, rate float64) error {
	s := sr.db.Stats()

	// Maximum number of open connections to the database.
	if err := sr.ddc.Set("db.max_open_connections", strconv.Itoa(s.MaxOpenConnections), tags, rate); err != nil {
		return err
	}

	// Pool Status
	// The number of established connections both in use and idle.
	if err := sr.ddc.Set("db.pool.open_connections", strconv.Itoa(s.OpenConnections), tags, rate); err != nil {
		return err
	}
	// The number of connections currently in use.
	if err := sr.ddc.Set("db.pool.in_use", strconv.Itoa(s.InUse), tags, rate); err != nil {
		return err
	}
	// The number of idle connections.
	if err := sr.ddc.Set("db.pool.idle", strconv.Itoa(s.Idle), tags, rate); err != nil {
		return err
	}

	// Counters
	// The total number of connections waited for.
	if err := sr.ddc.Set("db.counters.wait_count", strconv.FormatInt(s.WaitCount, 10), tags, rate); err != nil {
		return err
	}
	// The total time blocked waiting for a new connection.
	if err := sr.ddc.Timing("db.counters.wait_duration", s.WaitDuration, tags, rate); err != nil {
		return err
	}
	// The total number of connections closed due to SetMaxIdleConns.
	if err := sr.ddc.Set("db.counters.max_idle_closed", strconv.FormatInt(s.MaxIdleClosed, 10), tags, rate); err != nil {
		return err
	}
	// The total number of connections closed due to SetConnMaxIdleTime.
	if err := sr.ddc.Set("db.counters.max_idle_time_closed", strconv.FormatInt(s.MaxIdleTimeClosed, 10), tags, rate); err != nil {
		return err
	}
	// The total number of connections closed due to SetConnMaxLifetime.
	if err := sr.ddc.Set("db.counters.max_lifetime_closed", strconv.FormatInt(s.MaxLifetimeClosed, 10), tags, rate); err != nil {
		return err
	}
	return nil
}

func (sr *StatsReporter) ReportDBStats(ctx context.Context, tags []string, rate float64) {
	ticker := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sr.EmitStats(ctx, tags, rate); err != nil {
				sr.lg.Error(ctx, "Error Emitting DB Stats", zap.Error(err))
			}
		}
	}
}
