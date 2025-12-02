package api

import (
	"context"
	"devture-matrix-corporal/corporal/httphelp"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	listenAddress       string
	handlerRegistrators []httphelp.HandlerRegistrator

	server *http.Server
}

func NewServer(listenAddress string, handlerRegistrators []httphelp.HandlerRegistrator) *Server {
	return &Server{
		listenAddress:       listenAddress,
		handlerRegistrators: handlerRegistrators,
	}
}

func (me *Server) Start() error {
	me.server = &http.Server{
		Handler: me.newHandler(),
		Addr:    me.listenAddress,
	}

	go func() {
		logrus.Infof(`Stating Api Server on %s`, me.server.Addr)
		err := me.server.ListenAndServe()
		if err != http.ErrServerClosed {
			logrus.Warn(err)
		}
	}()

	return nil
}

func (me *Server) Stop() error {
	if me.server == nil {
		return nil
	}

	return me.server.Shutdown(context.Background())
}

func (me *Server) newHandler() http.Handler {
	router := mux.NewRouter()

	for _, registrator := range me.handlerRegistrators {
		registrator.RegisterRoutesWithRouter(router)
	}

	return router
}
