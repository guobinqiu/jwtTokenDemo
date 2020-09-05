package gin

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginHandler(t *testing.T) {
	server := httptest.NewServer(mainHandler())
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	token := e.GET("/v1/user/login").
		Expect().
		Status(http.StatusOK).
		Text().
		Raw()
	t.Logf("token=%s", token)
}

func TestInfoHandler(t *testing.T) {
	server := httptest.NewServer(mainHandler())
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	token := e.GET("/v1/user/login").
		Expect().
		Status(http.StatusOK).
		Text().
		Raw()

	user := e.GET("/v1/user/info").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		Text()
	t.Logf("user=%s", user.Raw())
	user.Equal("123456")
}

func TestInfoHandlerExpire(t *testing.T) {
	server := httptest.NewServer(mainHandler())
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	token := e.GET("/v1/user/login").
		Expect().
		Status(http.StatusOK).
		Text().
		Raw()
	t.Logf("token=%s", token)

	time.Sleep(4 * time.Second) //等待token过期

	e.GET("/v1/user/info").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		Text().Contains("Invalid auth token")
}
