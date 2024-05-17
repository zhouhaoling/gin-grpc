package common

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(r *gin.Engine, serverName string, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	//保证下面的优雅启停
	go func() {
		log.Printf("%s runing in %s \n", serverName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//sigint 用户发送INTR字符（ctrl+C）触发 kill -2
	//sigterm 结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down preject %s... \n", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %s \n", serverName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")

	}
	log.Printf("%s stop success... \n", serverName)
}
