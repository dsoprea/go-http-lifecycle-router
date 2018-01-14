package ghlr

import (
    "bytes"

    "net/http"
    "net/http/httptest"

    "github.com/gorilla/mux"
    "github.com/dsoprea/go-logging"
)

type bytesCloser struct {
    *bytes.Buffer
}

func (bc bytesCloser) Close() error {
    return nil
}

func DoRequest(router *mux.Router, method string, relPath string, requestBody string) (code int, body string, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    s := httptest.NewServer(router)
    testUrl := s.URL + "/" + relPath

    r, err := http.NewRequest(method, testUrl, nil)
    log.PanicIf(err)

    if requestBody != "" {
        b := bytes.NewBufferString(requestBody)
        r.Body = bytesCloser{ Buffer: b }
    }

    response, err := http.DefaultClient.Do(r)
    log.PanicIf(err)

    b := new(bytes.Buffer)
    b.ReadFrom(response.Body)
    actualBody := b.String()

    return response.StatusCode, actualBody, nil
}
