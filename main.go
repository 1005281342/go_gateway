package main

import (
	"flag"
	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/e421083458/go_gateway/grpc_proxy_router"
	"github.com/e421083458/go_gateway/http_proxy_router"
	"github.com/e421083458/go_gateway/router"
	"github.com/e421083458/go_gateway/tcp_proxy_router"
	"os"
	"os/signal"
	"syscall"
)

//endpoint dashboard后台管理  server代理服务器
//config ./conf/prod/ 对应配置文件夹

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	// 解析参数
	flag.Parse()

	// 校验endpoint
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 校验config
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	lib.InitModule(*config)
	defer lib.Destroy()

	if *endpoint == "dashboard" {

		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {

		dao.ServiceManagerHandler.LoadOnce()
		dao.AppManagerHandler.LoadOnce()

		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()
		go func() {
			tcp_proxy_router.TcpServerRun()
		}()
		go func() {
			grpc_proxy_router.GrpcServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpServerStop()
		grpc_proxy_router.GrpcServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}
}
