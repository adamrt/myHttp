package handler

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/hconn7/vxp/internal/request"
)

func (h *HandlerCfg) handlerDefault(w io.Writer, r *request.Request) {
	fmt.Println("handler triggered")
	data, err := os.ReadFile("../../internal/static/index.html")
	if err != nil {
		h.WriteError(w, 500, "Error Reading HTML file")
		fmt.Print(err)
		return
	}
	conentLength := len(data)
	datl := strconv.Itoa(conentLength)

	h.RespondWithJson(w, 200, "Conent-Length:"+datl, "")
	w.Write([]byte(data))
}
func (h *HandlerCfg) handlerPage(w io.Writer, r *request.Request) {
	fmt.Println("handler triggered")
	data, err := os.ReadFile("../../internal/static/adv.html")
	if err != nil {
		h.WriteError(w, 500, "Error Reading HTML file")
		fmt.Print(err)
		return
	}
	conentLength := len(data)
	datl := strconv.Itoa(conentLength)

	h.RespondWithJson(w, 200, "Conent-Length:"+datl, "")
	w.Write([]byte(data))
}
