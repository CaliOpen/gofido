// A test that uses a mock.
package main_test

import (
    "encoding/json"
    "errors"
    "fmt"
    . "github.com/CaliOpen/gofido"
    "github.com/CaliOpen/gofido/store"
    "github.com/gin-gonic/gin"
    "github.com/golang/mock/gomock"
    "github.com/tstranex/u2f"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRegisterRequest(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mock := store.NewMockStoreInterface(ctrl)

    user := gomock.Any()
    appId := "https://localhost:123456"
    trustedFacets := []string{appId}
    challenge, err := u2f.NewChallenge(appId, trustedFacets)
    registrations := &[]u2f.Registration{}
    response := u2f.NewWebRegisterRequest(challenge, *registrations)
    mock.EXPECT().NewChallenge(user).Return(*challenge, nil)
    mock.EXPECT().GetRegistrations(user).Return(*registrations, nil)
    s := FidoServer{}
    s.Store = mock

    gin.SetMode(gin.TestMode)
    r := gin.Default()
    register_url := fmt.Sprintf("/api/%s/register", user)
    r.GET(register_url, s.RegisterRequest)

    req, err := http.NewRequest(http.MethodGet, register_url, nil)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
    }
    body, _ := ioutil.ReadAll(w.Result().Body)
    expected := &u2f.WebRegisterRequest{}
    _ = json.Unmarshal(body, expected)
    if expected.AppID != response.AppID {
        t.Errorf("Expected %s for appId %s", expected.AppID, response.AppID)
    }
    request := expected.RegisterRequests[0]
    if request.Version != "U2F_V2" {
        t.Errorf("Invalid request version %s", request.Version)
    }
    if request.Challenge != store.EncodeBase64(challenge.Challenge) {
        t.Errorf("Invalid challenge response")
    }
}

func TestRegisterRequestKO(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mock := store.NewMockStoreInterface(ctrl)

    user := gomock.Any()
    appId := "https://localhost:123456"
    trustedFacets := []string{appId}
    challenge, err := u2f.NewChallenge(appId, trustedFacets)
    registrations := &[]u2f.Registration{}
    mock.EXPECT().NewChallenge(user).Return(*challenge, nil)
    mock.EXPECT().GetRegistrations(user).Return(*registrations, errors.New("test"))
    s := FidoServer{}
    s.Store = mock

    gin.SetMode(gin.TestMode)
    r := gin.Default()
    register_url := fmt.Sprintf("/api/%s/register", user)
    r.GET(register_url, s.RegisterRequest)

    req, err := http.NewRequest(http.MethodGet, register_url, nil)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusInternalServerError, w.Code)
    }
}
