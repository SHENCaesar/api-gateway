package main

import (
	"log"
	"net"

	"github.com/SHENCaesar/api-gateway/kitex/dal"
	demo "github.com/SHENCaesar/api-gateway/kitex/kitex_gen/demo/studentservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	dal.Init()
	addr, _ := net.ResolveTCPAddr("tcp", ":9999")
	svr := demo.NewServer(new(StudentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "student",
			},
		),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
