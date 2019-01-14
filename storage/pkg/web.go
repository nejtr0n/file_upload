package main

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"log"
	"test/storage/proto"
	"test/storage/ui/http"
	"time"
)

func main()  {
	svc := web.NewService(
		web.Name("go.micro.api.storage"),
		web.Handler(
			http.NewRouter(
				proto.NewStorageService("go.micro.srv.storage", client.DefaultClient),
			),
		),
		web.RegisterTTL(time.Second * 30),
		web.RegisterInterval(time.Second * 15),
	)

	err := svc.Init()
	if err != nil {
		log.Fatal(err)
	}

	if err := svc.Run(); err != nil {
		panic(err)
	}
}