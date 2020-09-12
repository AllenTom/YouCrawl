package youcrawl

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

type CookieStore interface {
	GetCookie(key string) *cookiejar.Jar
	SetCookie(key string, jar *cookiejar.Jar)
	GetOrCreate(key string) *cookiejar.Jar
}
type DefaultCookieStore struct {
	sync.Map
}

func (s *DefaultCookieStore) GetCookie(key string) *cookiejar.Jar {
	rawJar, _ := s.Load(key)
	jar := rawJar.(*cookiejar.Jar)
	return jar
}

func (s *DefaultCookieStore) GetOrCreate(key string) *cookiejar.Jar {
	newJar, _ := cookiejar.New(&cookiejar.Options{})
	rawJar, _ := s.Map.LoadOrStore(key, newJar)
	jar := rawJar.(*cookiejar.Jar)
	return jar
}

func (s *DefaultCookieStore) SetCookie(key string, jar *cookiejar.Jar) {
	s.Store(key, jar)
}

type CookieMiddleware struct {
	Store  CookieStore
	GetKey func(c *http.Client, r *http.Request, ctx *Context) string
}

func (m *CookieMiddleware) RequestCallback(c *http.Client, r *http.Request, ctx *Context) {
	key := "Default"
	if m.GetKey != nil {
		key = m.GetKey(c, r, ctx)
	}
	jar := c.Jar.(*cookiejar.Jar)
	m.Store.SetCookie(key, jar)
	ctx.Cookie = jar
}

func (m *CookieMiddleware) Process(c *http.Client, r *http.Request, ctx *Context) {
	key := "Default"
	if m.GetKey != nil {
		key = m.GetKey(c, r, ctx)
	}
	jar := m.Store.GetOrCreate(key)
	c.Jar = jar
}

type CookieMiddlewareOption struct {
	GetKey func(c *http.Client, r *http.Request, ctx *Context) string
}

func NewCookieMiddleware(option CookieMiddlewareOption) *CookieMiddleware {
	store := DefaultCookieStore{}
	return &CookieMiddleware{
		Store:  &store,
		GetKey: option.GetKey,
	}
}
