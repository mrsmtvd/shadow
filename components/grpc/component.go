package grpc

import (
	"github.com/mrsmtvd/shadow"
	"google.golang.org/grpc"
)

type Component interface {
	shadow.Component

	GetServiceInfo() map[string]grpc.ServiceInfo
}
