package ghlr

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/dsoprea/go-logging"
)

type LifecycleHandler interface {
    BeforeHandle(r *http.Request) error
    AfterHandle(r *http.Request) error

    BeforeApiHandle(r *http.Request) error
    AfterApiHandle(r *http.Request) error

    BeforeUiHandle(r *http.Request) error
    AfterUiHandle(r *http.Request) error
}

type httpUiHandler func(w http.ResponseWriter, r *http.Request)
type httpApiHandler func(w http.ResponseWriter, r *http.Request) map[string]interface{}

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
func (lr *LifecycleRouter) AddApiHandler(urlPath string, hh httpApiHandler, methods ...string) {
    f := func(w http.ResponseWriter, r *http.Request) {
        lr.lh.BeforeHandle(r)
        lr.lh.BeforeApiHandle(r)

// TODO(dustin): Later, we can extend this to respond as appropriate based on "Accept" headers.

        w.Header().Set("Content-Type", "application/json")

        output := hh(w, r)

        b, err := json.MarshalIndent(output, "", "  ")
        log.PanicIf(err)

        w.WriteHeader(http.StatusOK)
        w.Write(b)

        lr.lh.AfterApiHandle(r)
        lr.lh.AfterHandle(r)
    }

    r := lr.Router.HandleFunc(urlPath, f)
    if len(methods) > 0 {
        r.Methods(methods...)
    }
}

// AddUiHandler registers a path that will produce browser content (e.g. HTML).
func (lr *LifecycleRouter) AddUiHandler(urlPath string, hh httpUiHandler, methods ...string) {
    f := func(w http.ResponseWriter, r *http.Request) {
        lr.lh.BeforeHandle(r)
        lr.lh.BeforeUiHandle(r)

        w.Header().Set("Content-Type", "text/html")

        hh(w, r)

        w.WriteHeader(http.StatusOK)

        lr.lh.AfterUiHandle(r)
        lr.lh.AfterHandle(r)
    }

    r := lr.Router.HandleFunc(urlPath, f)
    if len(methods) > 0 {
        r.Methods(methods...)
    }
}
