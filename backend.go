package main

import (
	"flag"
	"fmt"

	"github.com/Savvy-Gameing/backend/internal/config"
	"github.com/Savvy-Gameing/backend/internal/handler"
	"github.com/Savvy-Gameing/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/backend-api.yaml", "the config file")

func main() {
	flag.Parse()

	*configFile = "/etc/backend-api.yaml"

	var c config.Config
	conf.MustLoad(*configFile, &c)
	fmt.Printf("c: %+v\n", c)

	c.Host = "0.0.0.0"
	c.Port = 8443
	// c.CertFile = "/etc/YOURPUBLIC.pem"
	// c.KeyFile = "/etc/YOURPRIVATE.key"
	c.TgWebHook = "https://game-invite.vercel.app:8443"
	c.TgPublicPem = "/etc/YOURPUBLIC.pem"
	c.TgPrivateKey = "/etc/YOURPRIVATE.key"

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*")) // note: modify in Nginx
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	group := service.NewServiceGroup()
	handler.RegisterJob(ctx, group)
	group.Start()
	defer group.Stop()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
