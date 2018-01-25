package database

type Balancer interface {
	Get() Executor
	Set([]Executor)
}
