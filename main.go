package main

import (
	"flag"
	"fmt"

	"github.com/lihongsheng/pay-gateway/internal/config"
	"github.com/lihongsheng/pay-gateway/internal/server"
	"github.com/lihongsheng/pay-gateway/internal/svc"
	"github.com/lihongsheng/pay-gateway/payment"
	"github.com/lihongsheng/pay-gateway/refund"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		payment.RegisterPaymentServer(grpcServer, server.NewPaymentServer(ctx))
		refund.RegisterRefundServer(grpcServer, server.NewRefundServer(ctx))

		//if c.Mode == service.DevMode || c.Mode == service.TestMode {
		//	reflection.Register(grpcServer)
		//}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
