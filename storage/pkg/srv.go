package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"log"
	"test/storage/app/service"
	"test/storage/proto"
	"time"
)

func main()  {
	svc := micro.NewService(
		micro.Name("go.micro.srv.storage"),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
	)

	svc.Init()

	err := svc.Server().Init(
		server.Wait(true), // Graceful shutdown
	)

	if err != nil {
		log.Fatal(err)
	}

	// Register Handlers
	err = proto.RegisterStorageHandler(svc.Server(), service.NewService())
	if err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}