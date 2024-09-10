package main

import (
	server "diplom/infrastructure/Server"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

func main() {
	//ctx, _ := context.WithCancel(context.Background())
	//rpc := "https://polygon-testnet.public.blastapi.io"

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := server.NewServer(logger, ":8080", nil)
	go func() {

		server.Start()
	}()

	time.Sleep(5 * time.Second)
	server.Restart()
}
