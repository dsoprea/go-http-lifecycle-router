package ghlr

import (
    "net/http"

    "github.com/gorilla/mux"
)

type LifecycleHandler interface {
    BeforeHandle(r *http.Request) error
    AfterHandle(r *http.Request) error

    BeforeApiHandle(r *http.Request) error
    AfterApiHandle(r *http.Request) error

    BeforeUiHandle(r *http.Request) error
    AfterUiHandle(r *http.Request) error
}

type httpHandler func(w http.ResponseWriter, r *http.Request)

type LifecycleRouter struct {
    Router *mux.Router

    lh LifecycleHandler
}

func NewLifecycleRouter(lh LifecycleHandler) *LifecycleRouter {
    r := mux.NewRouter()

    lr := &LifecycleRouter{
        Router: r,

        lh: lh,
    }

    return lr
}

// AddApiHandler registers a path that will produce data.
func (lr *LifecycleRouter) AddApiHandler(urlPath string, hh httpHandler, methods ...string) {
    f := func(w http.ResponseWriter, r *http.Request) {
        lr.lh.BeforeHandle(r)
        lr.lh.BeforeApiHandle(r)

        hh(w, r)

        lr.lh.AfterApiHandle(r)
        lr.lh.AfterHandle(r)
    }

    lr.Router.HandleFunc(urlPath, f).Methods(methods...)
}

// AddUiHandler registers a path that will produce browser content (e.g. HTML).
func (lr *LifecycleRouter) AddUiHandler(urlPath string, hh httpHandler, methods ...string) {
    f := func(w http.ResponseWriter, r *http.Request) {
        lr.lh.BeforeHandle(r)
        lr.lh.BeforeUiHandle(r)

        hh(w, r)

        lr.lh.AfterUiHandle(r)
        lr.lh.AfterHandle(r)
    }

    lr.Router.HandleFunc(urlPath, f).Methods(methods...)
}
