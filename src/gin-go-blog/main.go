package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	setting "gin-blog/pkg/settting"
	"gin-blog/routers"
)

func main() {

	// endless.DefaultReadTimeOut = setting.ReadTimeout
	// endless.DefaultWriteTimeOut = setting.WriteTimeout
	// endless.DefaultMaxHeaderBytes = 1 << 20
	// endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	router := routers.InitRouter()

	// server := endless.NewServer(endPoint, router)
	// server.BeforeBegin = func(add string) {
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }

	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Printf("Server err: %v", err)
	// }

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")

	// s.ListenAndServe()
}
