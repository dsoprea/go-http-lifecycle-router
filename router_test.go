package ghlr

import (
    "testing"
    "bytes"
    "net/http"
    "net/http/httptest"
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
    return &testLifecycleHandler{}
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
    f := func(w http.ResponseWriter, r *http.Request) {
        isHandled = true
    }

    lr.AddApiHandler("/", f, "GET")


    s := httptest.NewServer(lr.Router)
    testUrl := s.URL + "/"

    r, err := http.NewRequest("GET", testUrl, nil)
    if err != nil {
        t.Fatal(err)
    }


    response, err := http.DefaultClient.Do(r)
    if err != nil {
        t.Fatal(err)
    }


    actualCode := response.StatusCode
    expectedCode := http.StatusOK

    if actualCode != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 actualCode, expectedCode)
    }


    b := new(bytes.Buffer)
    b.ReadFrom(response.Body)
    actualBody := b.String()

    expectedBody := ""

    if actualBody != expectedBody {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 actualBody, expectedBody)
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

func Test_UiHandler(t *testing.T) {
    tlh := newTestLifecycleHandler()
    lr := NewLifecycleRouter(tlh)


    isHandled := false
    f := func(w http.ResponseWriter, r *http.Request) {
        isHandled = true
    }

    lr.AddUiHandler("/", f, "GET")


    s := httptest.NewServer(lr.Router)
    testUrl := s.URL + "/"

    r, err := http.NewRequest("GET", testUrl, nil)
    if err != nil {
        t.Fatal(err)
    }


    response, err := http.DefaultClient.Do(r)
    if err != nil {
        t.Fatal(err)
    }


    actualCode := response.StatusCode
    expectedCode := http.StatusOK

    if actualCode != expectedCode {
        t.Fatalf("handler returned unexpected code: (%d) != (%d)",
                 actualCode, expectedCode)
    }


    b := new(bytes.Buffer)
    b.ReadFrom(response.Body)
    actualBody := b.String()

    expectedBody := ""

    if actualBody != expectedBody {
        t.Fatalf("handler returned unexpected body: [%v] != [%v]",
                 actualBody, expectedBody)
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
