package ghlr

import (
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/dsoprea/go-logging"
)

var (
    routerLogger = log.NewLogger("ghlr.router")
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


type HttpErrorMessage interface {
    HttpErrorMessage() string
}

type HttpErrorCode interface {
    HttpErrorCode() int
}

// AddApiHandler registers a path that will produce data.
func (lr *LifecycleRouter) AddApiHandler(urlPath string, hh httpApiHandler, methods ...string) {
    f := func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err == nil {
                return
            }

            routerLogger.Errorf(nil, err, "There was a problem while handling the request.")

            hec, ok := err.(HttpErrorCode)
            if ok == true {
                code := hec.HttpErrorCode()
                w.WriteHeader(code)
            } else {
                w.WriteHeader(http.StatusInternalServerError)
            }

            hem, ok := err.(HttpErrorMessage)
            if ok == true {
                hem = hem

                message := hem.HttpErrorMessage()
                w.Write([]byte(message))
            } else {
                w.Write([]byte("There was a problem while handling the request."))
            }

            w.Write([]byte { '\n' })
        }()

        lr.lh.BeforeHandle(r)
        lr.lh.BeforeApiHandle(r)

// TODO(dustin): Later, we can extend this to respond as appropriate based on "Accept" headers.

        w.Header().Set("Content-Type", "application/json")

        output := hh(w, r)

        b, err := json.MarshalIndent(output, "", "  ")
        log.PanicIf(err)

        w.Write(b)
        w.Write([]byte { '\n' })

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

        lr.lh.AfterUiHandle(r)
        lr.lh.AfterHandle(r)
    }

    r := lr.Router.HandleFunc(urlPath, f)
    if len(methods) > 0 {
        r.Methods(methods...)
    }
}
