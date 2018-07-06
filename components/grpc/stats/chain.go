package stats

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type ChainHandler struct {
	stats.Handler

	mutex    sync.RWMutex
	handlers []stats.Handler
}

func WithStatsHandlerClientChain(handlers ...stats.Handler) grpc.DialOption {
	return grpc.WithStatsHandler(NewChainHandler(handlers...))
}

func WithStatsHandlerServerChain(handlers ...stats.Handler) grpc.ServerOption {
	return grpc.StatsHandler(NewChainHandler(handlers...))
}

func NewChainHandler(handlers ...stats.Handler) *ChainHandler {
	return &ChainHandler{
		handlers: handlers,
	}
}

func (h *ChainHandler) Handlers() []stats.Handler {
	h.mutex.RLock()
	tmp := make([]stats.Handler, len(h.handlers))
	copy(tmp, h.handlers)
	h.mutex.RUnlock()

	return tmp
}

func (h *ChainHandler) AddHandler(handler stats.Handler) {
	h.mutex.Lock()
	h.handlers = append(h.handlers, handler)
	h.mutex.Unlock()
}

func (h *ChainHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	for _, handler := range h.Handlers() {
		ctx = handler.TagConn(ctx, info)
	}

	return ctx
}

func (h *ChainHandler) HandleConn(ctx context.Context, stat stats.ConnStats) {
	for _, handler := range h.Handlers() {
		handler.HandleConn(ctx, stat)
	}
}

func (h *ChainHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	for _, handler := range h.Handlers() {
		ctx = handler.TagRPC(ctx, info)
	}

	return ctx
}

func (h *ChainHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	for _, handler := range h.Handlers() {
		handler.HandleRPC(ctx, stat)
	}
}
