package main

import (
	"fmt"
	"os"

	client "github.com/micro/services/clients/go"
)

func main() {
	c := client.NewClient(os.Getenv("MICRO_TOKEN"))
	rsp, err := c.HelloworldService.Call(client.HelloworldCallRequest{
		Name: "Janos",
	})
	fmt.Println(rsp, err)
}
