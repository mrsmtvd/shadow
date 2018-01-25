package database

type Storage interface {
	Executor() Executor
	Executors() []Executor
	Master() Executor
	Slave() Executor
	Slaves() []Executor
	AllowUseMasterAsSlave()
	DisallowUseMasterAsSlave()
	SetBalancer(Balancer)
}
