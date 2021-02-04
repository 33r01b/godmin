package server

const (
	CtxKeyUser      ctxKey = iota
	CtxKeyRequestID ctxKey = iota
)

type ctxKey int8
