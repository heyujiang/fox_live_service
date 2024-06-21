package main

import (
	"context"
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "fox_live_service/config"
)

func main() {
	handler := router.Register()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", global.Config.GetString("Host"), global.Config.GetInt("Port")), Handler: handler}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 处理CTRL+C等中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//退出业务处理

	log.Println("Shutdown Server ...")

	cancelCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(cancelCtx); err != nil {
		log.Fatalf("fox live backnd service shut down ： %+v", err)
	}
}
