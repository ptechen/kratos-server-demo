package grpc

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	pb "github.com/ptechen/kratos-proto/demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"kratos-server-demo/internal/signal"
	"time"
)

func addTls() (res credentials.TransportCredentials) {
	c, err := credentials.NewServerTLSFromFile("./configs/certs/server.pem", "./configs/certs/server.key")
	if err != nil {
		log.Error("credentials.NewServerTLSFromFile err: %v", err)
		panic(err)
	}
	return c
}

// New new a grpc server.
func New(svc pb.DemoServer) (ws *warden.Server, err error) {
	var (
		cfg warden.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	register()
	ws = warden.NewServer(&cfg, grpc.Creds(addTls()))
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}

func register() {
	dis := discovery.New(nil)
	ins := &naming.Instance{
		Zone:     "sha1",
		Env:      env.DeployEnv,
		AppID:    "demo.service",
		Addrs:    []string{"grpc://0.0.0.0:9000", "http:0.0.0.0:8000"},
		LastTs:   time.Now().Unix(),
		Metadata: map[string]string{"weight": "10"},
	}

	cancel, err := dis.Register(context.Background(), ins)
	if err != nil {
		panic(err)
	}
	go func() {
	forTag:
		for {
			select {
			case _, ok := <-signal.ExitChan:
				if !ok {
					cancel()
					break forTag
				}
			default:
				time.Sleep(time.Second)
			}
		}
		log.Info("discovery cancel success")
	}()
}

//func addTls() (res credentials.TransportCredentials) {
//	c, err := tls.LoadX509KeyPair("./configs/certs/apiserver.pem", "./configs/certs/apiserver-key.pem")
//	if err != nil {
//		log.Error("credentials.NewServerTLSFromFile err: %v", err)
//		panic(err)
//	}
//	certPool := x509.NewCertPool()
//	ca, err := ioutil.ReadFile("./configs/certs/ca.pem")
//	if err != nil {
//		panic(err)
//	}
//	if ok := certPool.AppendCertsFromPEM(ca); !ok {
//		panic("fsfs")
//	}
//	res = credentials.NewTLS(&tls.Config{
//		Certificates:[]tls.Certificate{c},
//		ClientAuth: tls.RequireAndVerifyClientCert,
//		ClientCAs: certPool,
//	})
//	return res
//}
