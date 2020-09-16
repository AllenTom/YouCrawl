package youcrawl

import (
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func TestDefaultCookieStore_GetCookie(t *testing.T) {
	store := DefaultCookieStore{}
	testCookieJar := cookiejar.Jar{}
	store.Store("test1", &testCookieJar)
	getTestCookieResult := store.GetCookie("test1")
	if getTestCookieResult == nil {
		t.Error("no match key,but store before")
	}
}

func TestDefaultCookieStore_GetCookieKey(t *testing.T) {
	store := DefaultCookieStore{}
	testCookieJar := cookiejar.Jar{}
	store.Store("default", &testCookieJar)
	middleware := CookieMiddleware{
		Store: &store,
		GetKey: func(c *http.Client, r *http.Request, ctx *Context) string {
			return "default"
		},
	}
	ctx := Context{Cookie: &cookiejar.Jar{}}
	client := &http.Client{}
	middleware.Process(client, nil, &ctx)
	middleware.RequestCallback(client, nil, &ctx)

}
