package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/mng-g/proto_example/coffeeshop_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// TODO: use TLS instead!
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to gPRC server ")
	}

	defer conn.Close()

	c := pb.NewCoffeeShopClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatal("error calling function getMenu")
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				log.Fatalf("can not receive %v", err)
			}

			items = resp.Items
			log.Printf("Resp received: %v", resp.Items)
		}
	}()

	<-done

	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	if err != nil {
		log.Printf("can not place order %v", err)
	}

	log.Printf("%v", receipt)

	status, err := c.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Printf("can not get order status %v", err)
	}

	log.Printf("%v", status)
}
