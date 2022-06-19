package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func New() *Router {
	return &Router{gin.Default()}
}

type Context struct {
	*gin.Context
}

// type HandlerFunc func(order.Context)

func (r *Router) GET(relativePath string, handler func(*Context)) {
	r.Engine.GET(relativePath, func(c *gin.Context) {
		handler(&Context{c})
	})
}

func (r *Router) POST(relativePath string, handler func(*Context)) {
	r.Engine.POST(relativePath, func(c *gin.Context) {
		handler(&Context{c})
	})
}

func (r *Router) PUT(relativePath string, handler func(*Context)) {
	r.Engine.PUT(relativePath, func(c *gin.Context) {
		handler(&Context{c})
	})
}

func (r *Router) DEL(relativePath string, handler func(*Context)) {
	r.Engine.DELETE(relativePath, func(c *gin.Context) {
		handler(&Context{c})
	})
}

func (r *Router) ListenAndServe(ctx context.Context, cancel context.CancelFunc) func() {
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return func() {
		defer cancel()

		<-ctx.Done()
		cancel()
		fmt.Println("shutting down gracefully, press Ctrl+C again to force")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(timeoutCtx); err != nil {
			fmt.Println(err)
		}
	}
}
