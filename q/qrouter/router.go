package qrouter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/instrumentation"
	"github.com/glocurrency/commons/q"
)

type router struct {
	*gin.Engine
	prefix string
}

func NewRouter(engine *gin.Engine, prefix string) *router {
	return &router{Engine: engine, prefix: prefix}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, r.prefix) {
		if err := RewriteForQ(req); err != nil {
			instrumentation.NoticeError(context.Background(), err, "cannot rewrite",
				instrumentation.WithField("prefix", r.prefix))
		}
	}
	r.Engine.ServeHTTP(w, req)
}

func RewriteForQ(req *http.Request) error {
	msg, err := q.NewQMessage(req)
	if err != nil {
		return fmt.Errorf("cannot parse q message: %w", err)
	}

	if msg.Name == "" {
		return nil
	}

	parsedURI, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return fmt.Errorf("cannot parse uri '%s': %w", req.RequestURI, err)
	}

	parsedURI.Path = strings.TrimRight(parsedURI.Path, "/") + "/" + msg.Name
	req.RequestURI = parsedURI.String()

	req.URL.Path = strings.TrimRight(req.URL.Path, "/") + "/" + msg.Name
	return nil
}
