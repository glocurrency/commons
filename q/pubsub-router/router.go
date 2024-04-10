package pubsubrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		if err := RewriteForPubSub(req); err != nil {
			instrumentation.NoticeError(context.Background(), err, "cannot rewrite",
				instrumentation.WithField("prefix", r.prefix))
		}
	}
	r.Engine.ServeHTTP(w, req)
}

func RewriteForPubSub(req *http.Request) error {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("cannot read body: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(data))

	var msg q.PubSubMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("cannot unmarshal '%s': %w", string(data), err)
	}

	name := msg.GetName()
	if name == "" {
		return nil
	}

	parsedURI, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return fmt.Errorf("cannot parse uri '%s': %w", req.RequestURI, err)
	}

	parsedURI.Path = strings.TrimRight(parsedURI.Path, "/") + "/" + name
	req.RequestURI = parsedURI.String()

	req.URL.Path = strings.TrimRight(req.URL.Path, "/") + "/" + name
	return nil
}
