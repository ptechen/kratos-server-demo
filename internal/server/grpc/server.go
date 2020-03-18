package grpc

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/bilibili/kratos/pkg/naming/etcd"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	pb "github.com/ptechen/kratos-proto/demo/api"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	ws = warden.NewServer(&cfg)
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}

//func register() {
//	var (
//		cfg struct{
//			Endpoints []string
//		}
//		ct  paladin.TOML
//	)
//	if err := paladin.Get("etcd.toml").Unmarshal(&ct); err != nil {
//		return
//	}
//	if err := ct.Get("etcd").UnmarshalTOML(&cfg); err != nil {
//		return
//	}
//	cli := &clientv3.Config{
//		Endpoints:            cfg.Endpoints,
//		AutoSyncInterval:     0,
//		DialTimeout: time.Second * 3,
//		DialOptions: []grpc.DialOption{grpc.WithBlock()},
//		DialKeepAliveTime:    0,
//		DialKeepAliveTimeout: 0,
//		MaxCallSendMsgSize:   0,
//		MaxCallRecvMsgSize:   0,
//		TLS:                  nil,
//		Username:             "",
//		Password:             "",
//		RejectOldCluster:     false,
//		Context:              nil,
//		LogConfig:            nil,
//		PermitWithoutStream:  false,
//	}
//	et, _ := etcd.New(cli)
//	ins := &naming.Instance{
//		Zone:     "sha1",
//		Env:      env.DeployEnv,
//		AppID:    "demo.service",
//		Addrs:    []string{"grpc://0.0.0.0:9000", "http:0.0.0.0:8000"},
//		LastTs:   time.Now().Unix(),
//		Metadata: map[string]string{"weight": "10"},
//	}
//	//et := &etcd.EtcdBuilder{}
//	cancel, _ := et.Register(context.Background(), ins)
//
//	defer cancel()
//}

//func register() {
//	conf := &discovery.Config{
//		Nodes: []string{"192.168.3.241:7171"},
//		Zone:  "sh1",
//		Env:   env.DeployEnv,
//	}
//	dis := discovery.New(conf)
//	ins := &naming.Instance{
//		Zone:     "sha1",
//		Env:      env.DeployEnv,
//		AppID:    "demo.service",
//		Addrs:    []string{"grpc://0.0.0.0:9000", "http:0.0.0.0:8000"},
//		LastTs:   time.Now().Unix(),
//		Metadata: map[string]string{"weight": "10"},
//	}
//
//	cancel, _ := dis.Register(context.Background(), ins)
//
//	defer cancel()
//}

func register() {
	var (
		cfg struct {
			Endpoints []string
		}
		ct paladin.TOML
	)
	if err := paladin.Get("etcd.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err := ct.Get("etcd").UnmarshalTOML(&cfg); err != nil {
		return
	}
	config := &clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	builder, err := etcd.New(config)

	if err != nil {
		fmt.Println("etcd 连接失败")
		return
	}

	ins := &naming.Instance{
		Zone:     "z1",
		Env:      env.DeployEnv,
		AppID:    "demo.service",
		Addrs:    []string{"grpc://0.0.0.0:9000", "http:0.0.0.0:8000"},
		LastTs:   time.Now().Unix(),
		Metadata: map[string]string{"weight": "10"},
	}
	//et := &etcd.EtcdBuilder{}
	cancel, err := builder.Register(context.Background(), ins)
	if err != nil {
		return
	}
	fmt.Println(cancel)

}
