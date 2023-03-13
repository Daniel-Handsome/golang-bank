package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/daniel/master-golang/api"
	db "github.com/daniel/master-golang/db"
	sqlc "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/gapi"
	"github.com/daniel/master-golang/pb"
	"github.com/daniel/master-golang/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var ENV_PATH = ".env"

func main() {
	config, err := utils.LoadConfig(ENV_PATH)
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.InitDatabase(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := sqlc.NewStore(db)
	// runGinServer(config, store)
	go runGrpcServer(config, store)
	runGrpcGateWayServer(config, store)
}

func runGrpcGateWayServer(config utils.Config, store sqlc.Store) {
	fmt.Println("starting HTTP GateWay server...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", "8080"))
	if err != nil {
		log.Fatalf("Failed to Listen: %v \n", err)
	}

	// register
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("server start error message : %v", err)
	}


	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatalf("register grpc server err message: %v", err)
	}

	mux  := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	
	if err := http.Serve(lis, mux); err != nil {
		log.Fatalf("started serve error mesage : %v", err)
	}

}

func runGrpcServer(config utils.Config, store sqlc.Store) {
	fmt.Println("starting gRPC server...")
	// 預設0.0.0.0 但grpc要有ip 所以不能用0.0.0.0所以 要用127.0.0.1
	// 這邊是因為在docker內 docekr跟本地不同命名空間 但Docker 启动容器时，它会为容器分配一个独立的 IP 地址
	// 如果您在容器内部运行的服务在主机上绑定了 0.0.0.0，则意味着该服务对所有网络中的主机都可以访问。在这种情况下，主机可以通过与容器分配的 IP 地址进行通信，从而访问容器内部的服务。
	// 因此，当容器内部的服务绑定到 0.0.0.0 时，本地主机可以利用容器的 IP 地址通过 Docker 访问该服务。
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Grpc_port))
	if err != nil {
		log.Fatalf("Failed to Listen: %v \n", err)
	}

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server: ")
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	// https://stackoverflow.com/questions/41424630/why-do-we-need-to-register-reflection-service-on-grpc-server
	// 某些輔助工具需要看能用甚麼function 這邊利用反射 不一定要使用
	reflection.Register(grpcServer)

	log.Printf("start listening on port %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("grpc server start error")
	}
}

func runGinServer(config utils.Config, store sqlc.Store) {
	router, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(router.Run())
}