package grpc

import (
	"github.com/kihamo/shadow"
	"google.golang.org/grpc"
)

type Component interface {
	shadow.Component

	GetServiceInfo() map[string]grpc.ServiceInfo
}
