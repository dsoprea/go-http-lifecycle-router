package ghlr

import (
    "strings"

    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"
    "github.com/dsoprea/go-logging"
    "github.com/go-errors/errors"
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
type httpApiHandler func(w http.ResponseWriter, r *http.Request, data map[string]interface{}) map[string]interface{}

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
func (lr *LifecycleRouter) AddApiHandler(urlPath string, hh httpApiHandler, methods []string, decodeAndPassBody bool) {
    f := func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err == nil {
                return
            }

            routerLogger.Errorf(nil, err, "There was a problem while handling the request.")

            // If we have a stackframe-wrapped error, get the original error
            // out of it. It will mask any managed errors that we might need to
            // display back to the user.
            errError, ok := err.(*errors.Error)
            var errHtml error
            if ok == true {
                errHtml = errError.Err
            } else {
                errHtml = err.(error)
            }

            // Determine if there's a custom code to return.

            hec, ok := errHtml.(HttpErrorCode)
            code := http.StatusInternalServerError
            if ok == true {
                code = hec.HttpErrorCode()
            }

            // Determine if there's a custom message to return.

            hem, ok := errHtml.(HttpErrorMessage)
            if ok == true {
                hem = hem

                message := hem.HttpErrorMessage()
                http.Error(w, message, code)
            } else {
                w.Write([]byte("There was a problem while handling the request."))
                w.Write([]byte { '\n' })
            }
        }()

        var d map[string]interface{}
        if decodeAndPassBody == true {
            d = map[string]interface{} {}

            ct := r.Header.Get("Content-Type")
            ct = strings.ToLower(ct)
            if ct != "" && ct != "application/json" {
                log.Panicf("content-type not supported")
            }

            j := json.NewDecoder(r.Body)
            defer r.Body.Close()

            err := j.Decode(&d)
            log.PanicIf(err)
        }

        lr.lh.BeforeHandle(r)
        lr.lh.BeforeApiHandle(r)

// TODO(dustin): Later, we can extend this to respond as appropriate based on "Accept" headers.

        w.Header().Set("Content-Type", "application/json")

        output := hh(w, r, d)

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
func (lr *LifecycleRouter) AddUiHandler(urlPath string, hh httpUiHandler) {
    f := func(w http.ResponseWriter, r *http.Request) {
        lr.lh.BeforeHandle(r)
        lr.lh.BeforeUiHandle(r)

        w.Header().Set("Content-Type", "text/html")

        hh(w, r)

        lr.lh.AfterUiHandle(r)
        lr.lh.AfterHandle(r)
    }

    r := lr.Router.HandleFunc(urlPath, f)
    r.Methods("GET")
}
