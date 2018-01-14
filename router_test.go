package ghlr

import (
    "testing"
    "reflect"

    "encoding/json"
    "net/http"

    "github.com/dsoprea/go-logging"
)


type testLifecycleHandler struct {
    beforeHandle bool
    afterHandle bool

    beforeApiHandle bool
    afterApiHandle bool

    beforeUiHandle bool
    afterUiHandle bool
}

func newTestLifecycleHandler() *testLifecycleHandler {
    return new(testLifecycleHandler)
}

func (tlh *testLifecycleHandler) BeforeHandle(r *http.Request) error {
    tlh.beforeHandle = true

    return nil
}

func (tlh *testLifecycleHandler) AfterHandle(r *http.Request) error {
    tlh.afterHandle = true

    return nil
}

func (tlh *testLifecycleHandler) BeforeApiHandle(r *http.Request) error {
    tlh.beforeApiHandle = true

    return nil
}

func (tlh *testLifecycleHandler) AfterApiHandle(r *http.Request) error {
    tlh.afterApiHandle = true

    return nil
}

func (tlh *testLifecycleHandler) BeforeUiHandle(r *http.Request) error {
    tlh.beforeUiHandle = true

    return nil
}

func (tlh *testLifecycleHandler) AfterUiHandle(r *http.Request) error {
    tlh.afterUiHandle = true

    return nil
}

func Test_ApiHandler(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)


    isHandled := false
    f := func(w http.ResponseWriter, r *http.Request, d map[string]interface{}) map[string]interface{} {
        isHandled = true

        return map[string]interface{} {
            "aa": 123,
        }
    }

    lr.AddApiHandler("/", f, []string { "GET" }, false)


    code, body, err := DoRequest(lr.Router, "GET", "", "")
    log.PanicIf(err)

    expectedCode := http.StatusOK

    if code != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 code, expectedCode)
    }


    decoded := map[string]interface{} {}
    err = json.Unmarshal([]byte(body), &decoded)
    log.PanicIf(err)

    // expected has a float because apparently all numbers decode to float
    // (the safest option for the parser).
    expected := map[string]interface{} {
        "aa": 123.0,
    }

    if reflect.DeepEqual(decoded, expected) != true {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 decoded, expected)
    }


    if isHandled != true {
        t.Fatalf("handler was not hit")
    }


    if tlh.beforeHandle != true {
        t.Errorf("handler did not hit beforeHandle trigger")
    }

    if tlh.beforeApiHandle != true {
        t.Errorf("handler did not hit beforeApiHandle trigger")
    }

    if tlh.beforeUiHandle != false {
        t.Errorf("handler hit but should've missed beforeUiHandle")
    }

    if tlh.afterUiHandle != false {
        t.Errorf("handler hit but should've missed afterUiHandle")
    }

    if tlh.afterApiHandle != true {
        t.Errorf("handler did not hit afterApiHandle trigger")
    }

    if tlh.afterHandle != true {
        t.Errorf("handler did not hit afterHandle trigger")
    }
}


type testHttpError struct {

}

func (the testHttpError) Error() string {
    return "unmanaged error message"
}

func (the testHttpError) HttpErrorMessage() string {
    return "managed error message"
}

func (the testHttpError) HttpErrorCode() int {
    return 555
}

func Test_ApiHandler_Error(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)


    f := func(w http.ResponseWriter, r *http.Request, d map[string]interface{}) map[string]interface{} {
        panic(testHttpError{})
    }

    lr.AddApiHandler("/", f, []string { "GET" }, false)

    code, body, err := DoRequest(lr.Router, "GET", "", "")
    log.PanicIf(err)


    expectedCode := 555

    if code != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 code, expectedCode)
    }


    expectedBody := "managed error message\n"

    if body != expectedBody {
        t.Fatalf("failed handler returned unexpected body: [%v] != [%v]",
                 body, expectedBody)
    }
}

func Test_UiHandler(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)


    isHandled := false
    f := func(w http.ResponseWriter, r *http.Request) {
        isHandled = true
    }

    lr.AddUiHandler("/", f, "GET")

    code, body, err := DoRequest(lr.Router, "GET", "", "")
    log.PanicIf(err)

    expectedCode := http.StatusOK

    if code != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 code, expectedCode)
    }


    expectedBody := ""

    if body != expectedBody {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 body, expectedBody)
    }


    if isHandled != true {
        t.Fatalf("handler was not hit")
    }


    if tlh.beforeHandle != true {
        t.Errorf("handler did not hit beforeHandle trigger")
    }

    if tlh.beforeApiHandle != false {
        t.Errorf("handler hit but should've missed beforeApiHandle trigger")
    }

    if tlh.beforeUiHandle != true {
        t.Errorf("handler did not hit beforeUiHandle")
    }

    if tlh.afterUiHandle != true {
        t.Errorf("handler did not hit afterUiHandle")
    }

    if tlh.afterApiHandle != false {
        t.Errorf("handler hit but should've missed afterApiHandle trigger")
    }

    if tlh.afterHandle != true {
        t.Errorf("handler did not hit afterHandle trigger")
    }

}

func Test_ApiHandler_DecodeAndPass_WantDecode(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)

    isHandled := false
    f := func(w http.ResponseWriter, r *http.Request, incoming map[string]interface{}) map[string]interface{} {
        isHandled = true

        expected := map[string]interface{}{
            "incoming-data": 123.0,
        }

        if reflect.DeepEqual(incoming, expected) != true {
            t.Fatalf("handler received unexpected request body: [%v] != [%v]",
                     incoming, expected)
        }

        return map[string]interface{} {
            "aa": 123,
        }
    }

    lr.AddApiHandler("/", f, []string { "GET" }, true)


    code, body, err := DoRequest(lr.Router, "GET", "", `{ "incoming-data": 123 }`)
    log.PanicIf(err)

    expectedCode := http.StatusOK

    if code != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 code, expectedCode)
    }


    decoded := map[string]interface{} {}
    err = json.Unmarshal([]byte(body), &decoded)
    log.PanicIf(err)

    // expected has a float because apparently all numbers decode to float
    // (the safest option for the parser).
    expected := map[string]interface{} {
        "aa": 123.0,
    }

    if reflect.DeepEqual(decoded, expected) != true {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 decoded, expected)
    }


    if isHandled != true {
        t.Fatalf("handler was not hit")
    }
}

func Test_ApiHandler_DecodeAndPass_NoDecode(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)

    isHandled := false
    f := func(w http.ResponseWriter, r *http.Request, incoming map[string]interface{}) map[string]interface{} {
        isHandled = true

        if incoming != nil {
            t.Fatalf("Did not expect incoming request-body to be decoded.")
        }

        return map[string]interface{} {
            "aa": 123,
        }
    }

    lr.AddApiHandler("/", f, []string { "GET" }, false)


    code, body, err := DoRequest(lr.Router, "GET", "", `{ "incoming-data": 123 }`)
    log.PanicIf(err)

    expectedCode := http.StatusOK

    if code != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 code, expectedCode)
    }


    decoded := map[string]interface{} {}
    err = json.Unmarshal([]byte(body), &decoded)
    log.PanicIf(err)

    // expected has a float because apparently all numbers decode to float
    // (the safest option for the parser).
    expected := map[string]interface{} {
        "aa": 123.0,
    }

    if reflect.DeepEqual(decoded, expected) != true {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 decoded, expected)
    }


    if isHandled != true {
        t.Fatalf("handler was not hit")
    }
}
