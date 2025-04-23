package handler

import (
	"io"

	"github.com/hconn7/vxp/internal/state"
)

type method string
type target string

type handler func(w io.Writer, r *state.Request)

type handlerFunc func(method, target, handler)
