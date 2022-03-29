package main

import (
	"client/api"
	"client/socket"
	"fmt"
	"os/exec"
	"sync"
)

func main()  {
	sessionToken, err := api.CreateSession()
	if err != nil {
		fmt.Println("Session token is nil")
		return
	}

	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fmt.Sprintf("http://127.0.0.1:5500/index.html?state=%v", sessionToken)).Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := socket.Socket{
		SessionToken: sessionToken,
		MessageHandler: func(message socket.Message) {
			if message.Type == "AUTHENTICATED" {
				fmt.Println(message.Payload)
			}
		},
	}

	err = s.Connect()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
