package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	srlndk "github.com/srl-wim/protos"
	"github.com/srl-wim/srl-ndk-git/agent"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "agent_name", "ndk-git")
	a, err := agent.NewAgent(ctx, "ndk-git")
	if err != nil {
		log.Printf("failed to create agent: %v", err)
		os.Exit(1)
	}
	a.StartGRPCServer()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go a.KeepAlive(ctx, time.Minute)

	// Config notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartConfigNotificationStream(ctx)
		for {
			select {
			case event := <-appIDChan:
				log.Printf("Config notification: %+v", event)
				for _, item := range event.Notification {
					switch x := item.SubscriptionTypes.(type) {
					case *srlndk.Notification_Config:
						resp := item.GetConfig()
						if resp.Data != nil {
							a.HandleConfigEvent(resp.Op, resp.Key, &resp.Data.Json)
						} else {
							a.HandleConfigEvent(resp.Op, resp.Key, nil)
						}
					default:
						log.Printf("\nGot unhandled message %s ", x)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// AppId notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartAppIDNotificationStream(ctx, 0)
		for {
			select {
			case event := <-appIDChan:
				log.Printf("appID notification: %+v", event)
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// BFDSession notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartBFDSessionNotificationStream(ctx, nil, nil, 0)
		for {
			select {
			case event := <-appIDChan:
				log.Printf("BFDSession notification: %+v", event)
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// NwInst notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartNwInstNotificationStream(ctx)
		for {
			select {
			case event := <-appIDChan:
				log.Printf("NwInst notification: %+v", event)
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// Interface notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartInterfaceNotificationStream(ctx, "")
		for {
			select {
			case event := <-appIDChan:
				log.Printf("Interface notification: %+v", event)
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// LLDPNeighbor notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		appIDChan := a.StartLLDPNeighNotificationStream(ctx, "", "", "")
		for {
			select {
			case event := <-appIDChan:
				log.Printf("LLDPNeighbor notification: %+v", event)
			case <-ctx.Done():
				return
			}
		}
	}()
	//time.Sleep(time.Second)
	// Route notifications
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	appIdChan := app.StartRouteNotificationStream(ctx, "", nil, 0)
	// 	for {
	// 		select {
	// 		case event := <-appIdChan:
	// 			log.Printf("Route notification: %+v", event)
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	wg.Wait()
}
