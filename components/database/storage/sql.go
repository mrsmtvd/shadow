package storage

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/kihamo/shadow/components/database"
)

type SQL struct {
	mutex            sync.RWMutex
	useMasterAsSlave int64
	masterExecutor   *SQLExecutor
	slaveExecutors   []*SQLExecutor
	balancer         database.Balancer
}

func NewSQL(driver string, masterDSN string, slavesDSN []string, options map[string]string, allowUseMasterAsSlave bool) (s *SQL, err error) {
	if masterDSN == "" {
		return nil, errors.New("DSN of master is empty")
	}

	s = &SQL{
		slaveExecutors: make([]*SQLExecutor, 0, len(slavesDSN)),
	}

	if s.masterExecutor, err = NewSQLExecutor(driver, masterDSN, options); err != nil {
		return nil, err
	}

	if len(slavesDSN) > 0 {
		for _, dsn := range slavesDSN {
			executor, err := NewSQLExecutor(driver, dsn, options)
			if err != nil {
				return nil, err
			}

			s.slaveExecutors = append(s.slaveExecutors, executor)
		}
	} else {
		allowUseMasterAsSlave = true
	}

	if allowUseMasterAsSlave {
		s.AllowUseMasterAsSlave()
	} else {
		s.DisallowUseMasterAsSlave()
	}

	return s, nil
}

func (s *SQL) Executor() database.Executor {
	return s.Slave()
}

func (s *SQL) Executors() []database.Executor {
	executors := make([]database.Executor, 0, len(s.slaveExecutors)+1)
	executors = append(executors, s.Master())
	executors = append(executors, s.Slaves()...)

	return executors
}

func (s *SQL) Master() database.Executor {
	return s.masterExecutor
}

func (s *SQL) Slave() database.Executor {
	s.mutex.RLock()
	balancer := s.balancer
	s.mutex.RUnlock()

	if balancer == nil {
		return s.Master()
	}

	return balancer.Get()
}

func (s *SQL) Slaves() []database.Executor {
	executors := make([]database.Executor, 0, len(s.slaveExecutors))

	for _, executor := range s.slaveExecutors {
		executors = append(executors, executor)
	}

	return executors
}

func (s *SQL) AllowUseMasterAsSlave() {
	current := atomic.LoadInt64(&s.useMasterAsSlave)
	if current == 1 {
		return
	}

	atomic.StoreInt64(&s.useMasterAsSlave, 1)

	s.mutex.RLock()
	if s.balancer != nil {
		s.balancer.Set(s.Executors())
	}
	s.mutex.RUnlock()
}

func (s *SQL) DisallowUseMasterAsSlave() {
	if len(s.slaveExecutors) == 0 {
		return
	}

	current := atomic.LoadInt64(&s.useMasterAsSlave)
	if current == 0 {
		return
	}

	atomic.StoreInt64(&s.useMasterAsSlave, 0)

	s.mutex.RLock()
	if s.balancer != nil {
		s.balancer.Set(s.Slaves())
	}
	s.mutex.RUnlock()
}

func (s *SQL) SetBalancer(balancer database.Balancer) {
	executors := make([]database.Executor, 0, len(s.slaveExecutors)+1)
	executors = append(executors, s.Slaves()...)

	if atomic.LoadInt64(&s.useMasterAsSlave) == 1 {
		executors = append(executors, s.Master())
	}

	balancer.Set(executors)

	s.mutex.Lock()
	s.balancer = balancer
	s.mutex.Unlock()
}

func (s *SQL) SetMaxIdleConns(n int) {
	s.masterExecutor.DB().SetMaxIdleConns(n)

	for _, executor := range s.slaveExecutors {
		executor.DB().SetMaxIdleConns(n)
	}
}

func (s *SQL) SetMaxOpenConns(n int) {
	s.masterExecutor.DB().SetMaxOpenConns(n)

	for _, executor := range s.slaveExecutors {
		executor.DB().SetMaxOpenConns(n)
	}
}

func (s *SQL) SetConnMaxLifetime(d time.Duration) {
	s.masterExecutor.DB().SetConnMaxLifetime(d)

	for _, executor := range s.slaveExecutors {
		executor.DB().SetConnMaxLifetime(d)
	}
}

func (s *SQL) TraceOn(logger gorp.GorpLogger) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.masterExecutor.executor.(*gorp.DbMap).TraceOn("", logger)

	for _, executor := range s.slaveExecutors {
		executor.executor.(*gorp.DbMap).TraceOn("", logger)
	}
}

func (s *SQL) TraceOff() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.masterExecutor.executor.(*gorp.DbMap).TraceOff()

	for _, executor := range s.slaveExecutors {
		executor.executor.(*gorp.DbMap).TraceOff()
	}
}

func (s *SQL) SetTypeConverter(converter gorp.TypeConverter) {
	s.masterExecutor.executor.(*gorp.DbMap).TypeConverter = converter

	for _, executor := range s.slaveExecutors {
		executor.executor.(*gorp.DbMap).TypeConverter = converter
	}
}

func (s *SQL) AddTableWithName(i interface{}, name string) {
	s.masterExecutor.executor.(*gorp.DbMap).AddTableWithName(i, name)

	for _, executor := range s.slaveExecutors {
		executor.executor.(*gorp.DbMap).AddTableWithName(i, name)
	}
}

func (s *SQL) CreateTablesIfNotExists() error {
	return s.masterExecutor.executor.(*gorp.DbMap).CreateTablesIfNotExists()
}

func (s *SQL) Dialect() string {
	return s.masterExecutor.dialect
}
