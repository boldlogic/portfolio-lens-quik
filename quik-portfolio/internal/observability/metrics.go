package observability

import (
	"database/sql"
	"errors"
	"time"

	"github.com/boldlogic/packages/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	statusOK    = "ok"
	statusError = "error"
)

var durationBuckets = []float64{.0005, .001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5}

type Recorder interface {
	ObserveRepository(operation string, duration time.Duration, err error)
	ObserveDBQuery(operation string, duration time.Duration, err error)
}

type noopRecorder struct{}

func Noop() Recorder {
	return noopRecorder{}
}

func (noopRecorder) ObserveRepository(string, time.Duration, error) {}
func (noopRecorder) ObserveDBQuery(string, time.Duration, error)    {}

type Metrics struct {
	repositoryDuration *prometheus.HistogramVec
	dbQueryDuration    *prometheus.HistogramVec
}

func New(reg metrics.Registry) *Metrics {
	if reg == nil {
		return &Metrics{}
	}

	return &Metrics{
		repositoryDuration: registerOrGetHistogramVec(reg, newRepositoryDuration()),
		dbQueryDuration:    registerOrGetHistogramVec(reg, newDBQueryDuration()),
	}
}

func (m *Metrics) ObserveRepository(operation string, duration time.Duration, err error) {
	if m == nil || m.repositoryDuration == nil {
		return
	}
	m.repositoryDuration.WithLabelValues(operation, status(err)).Observe(duration.Seconds())
}

func (m *Metrics) ObserveDBQuery(operation string, duration time.Duration, err error) {
	if m == nil || m.dbQueryDuration == nil {
		return
	}
	m.dbQueryDuration.WithLabelValues(operation, status(err)).Observe(duration.Seconds())
}

func EnsureRecorder(rec Recorder) Recorder {
	if rec == nil {
		return Noop()
	}
	return rec
}

func RegisterDBStats(reg metrics.Registry, db *sql.DB) error {
	if reg == nil || db == nil {
		return nil
	}

	collectors := []prometheus.Collector{
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "quik_portfolio_db_pool_open_connections",
			Help: "Current number of established database connections.",
		}, func() float64 { return float64(db.Stats().OpenConnections) }),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "quik_portfolio_db_pool_in_use",
			Help: "Current number of database connections in use.",
		}, func() float64 { return float64(db.Stats().InUse) }),
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Name: "quik_portfolio_db_pool_wait_count_total",
			Help: "Total number of waits for a database connection.",
		}, func() float64 { return float64(db.Stats().WaitCount) }),
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Name: "quik_portfolio_db_pool_wait_duration_seconds_total",
			Help: "Total time blocked waiting for a database connection.",
		}, func() float64 { return db.Stats().WaitDuration.Seconds() }),
	}

	for _, collector := range collectors {
		if err := registerOrIgnoreAlreadyRegistered(reg, collector); err != nil {
			return err
		}
	}
	return nil
}

func newRepositoryDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "quik_portfolio_repository_duration_seconds",
		Help:    "Repository operation duration in seconds.",
		Buckets: durationBuckets,
	}, []string{"operation", "status"})
}

func newDBQueryDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "quik_portfolio_db_query_duration_seconds",
		Help:    "Database query duration in seconds.",
		Buckets: durationBuckets,
	}, []string{"operation", "status"})
}

func registerOrGetHistogramVec(reg prometheus.Registerer, h *prometheus.HistogramVec) *prometheus.HistogramVec {
	if err := reg.Register(h); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.HistogramVec)
		}
		panic(err)
	}
	return h
}

func registerOrGetCounterVec(reg prometheus.Registerer, c *prometheus.CounterVec) *prometheus.CounterVec {
	if err := reg.Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.CounterVec)
		}
		panic(err)
	}
	return c
}

func registerOrIgnoreAlreadyRegistered(reg prometheus.Registerer, collector prometheus.Collector) error {
	if err := reg.Register(collector); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return nil
		}
		return err
	}
	return nil
}

func status(err error) string {
	if err != nil {
		return statusError
	}
	return statusOK
}
