package handler

import (
	"errors"
	"fmt"
	"io"

	"github.com/hconn7/vxp/internal/request"
)

// type handlerFunc func(method, target, handler)

type handler func(w io.Writer, r *request.Request)

type StatusCode map[int]string

type HandlerCfg struct {
	StatusCode map[int]string
	Router     *Router
}

func NewHandler() *HandlerCfg {
	r := NewRouter()
	h := &HandlerCfg{
		StatusCode: make(map[int]string),
		Router:     r,
	}
	h.LoadCodes()
	h.InitRouters()
	return h
}
func (h *HandlerCfg) InitRouters() {
	h.Router.SetRouters("GET", "/", h.handlerDefault)
	h.Router.SetRouters("GET", "/favicon.ico", func(w io.Writer, r *request.Request) {
		h.RespondWithJson(w, 204, "", "")
	})
	h.Router.SetRouters("GET", "/hconn", h.handlerPage)

}

func (h *HandlerCfg) LoadCodes() {
	h.InitCode(200, "OK")
	h.InitCode(201, "Created")
	h.InitCode(204, "No Content")
	h.InitCode(301, "Moved Permanently")
	h.InitCode(302, "Found")
	h.InitCode(400, "Bad Request")
	h.InitCode(401, "Unauthorized")
	h.InitCode(403, "Forbidden")
	h.InitCode(404, "Not Found")
	h.InitCode(500, "Internal Server Error")
	return

}

func (h *HandlerCfg) InitCode(code int, msg string) {
	h.StatusCode[code] = msg
	return
}

func (h *HandlerCfg) RespondWithJson(w io.Writer, code int, headers string, payload string) error {
	_, ok := h.CheckValidity(code)
	if ok {
		h.WriteStatus(w, code)
		h.WriteHeaders(w, headers)
		w.Write([]byte(payload))
		return nil
	}
	return errors.New("Error, code likely doesn't exist?")
}

func (h *HandlerCfg) WriteError(w io.Writer, code int, msg string) {
	h.RespondWithJson(w, 500, "", msg)
}
func (h *HandlerCfg) WriteStatus(w io.Writer, code int) {
	if msg, ok := h.CheckValidity(code); ok {
		m := fmt.Sprintf("HTTP/1.1 %v %s\r\n", code, msg)
		w.Write([]byte(m))
		return
	} else {
		w.Write([]byte("HTTP/1,1 500 Internal Server Error \r\n"))
		return
	}
}
func (h *HandlerCfg) WriteHeaders(w io.Writer, headers string) {
	msg := fmt.Sprintf("%s\r\n\r\n", headers)
	w.Write([]byte(msg))
	return
}

func (h *HandlerCfg) CheckValidity(code int) (string, bool) {
	if msg, ok := h.StatusCode[code]; ok {
		return msg, true
	}
	return "", false
}
