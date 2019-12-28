package main

import (
	"fmt"
	"github.com/ccchieh/ginHelper"
	"github.com/ccchieh/ginHelper_example/handler"
	"github.com/ccchieh/ginHelper_example/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	ginHelper.Build(new(handler.Helper), r)

	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8080)

	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Service run in http://", addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
