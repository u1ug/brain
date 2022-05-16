package main

import (
	"fmt"
	lp "github.com/LdDl/fiber-long-poll/v2"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
)

func main() {
	//server setup
	server := fiber.New()
	manager, err := lp.StartLongpoll(lp.Options{
		LoggingEnabled:                 false,
		MaxLongpollTimeoutSeconds:      120,
		MaxEventBufferSize:             100,
		EventTimeToLiveSeconds:         120,
		DeleteEventAfterFirstRetrieval: false,
	})
	if err != nil {
		panic(err)
	}
	defer manager.Shutdown()
	server.Get("/check_messages", HandleTest(manager))
	//message generator
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(10-1)+1) * time.Second)
			manager.Publish("new_messages", "Test message text")
			fmt.Println("published")
		}
	}()
	//starting server
	err = server.Listen(":8080")
	if err != nil {
		fmt.Printf("Can't start server due the error: %s\n", err.Error())
	}
}

func HandleTest(manager *lp.LongpollManager) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Context().PostArgs().Set("timeout", "10")
		ctx.Context().PostArgs().Set("category", "new_messages")
		return manager.SubscriptionHandler(ctx)
	}
}
