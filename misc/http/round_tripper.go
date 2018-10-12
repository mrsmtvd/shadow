package http

import (
	h "net/http"
)

type roundTripperFunc = func(r *h.Request, next h.RoundTripper) (w *h.Response, err error)

type RoundTripperWrapper interface {
	h.RoundTripper
	Next(next interface{}) roundTripper
}

type roundTripper struct {
	f    roundTripperFunc
	next RoundTripperWrapper
}

func NewRoundTripper(f roundTripperFunc) RoundTripperWrapper {
	return roundTripper{
		f: f,
	}
}

func (t roundTripper) Next(next interface{}) roundTripper {
	if t.next != nil {
		t.next = t.next.Next(next)
	} else {
		if f, ok := next.(roundTripperFunc); ok {
			t.next = NewRoundTripper(f)
		} else if o, ok := next.(h.RoundTripper); ok {
			t.next = NewRoundTripper(func(r *h.Request, _ h.RoundTripper) (w *h.Response, err error) {
				return o.RoundTrip(r)
			})
		} else {
			panic("Unknown type of next round tripper")
		}
	}

	return t
}

func (t roundTripper) RoundTrip(r *h.Request) (w *h.Response, err error) {
	return t.f(r, t.next)
}
