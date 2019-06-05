package main

import (
	"fmt"
	"time"
	"util"
)

func longRun(exit chan bool) {
	fmt.Println("long run service start")
	ticker := time.NewTicker(time.Hour * 8)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := getLatestRates(); err != nil {
					util.FailToGetCurrencyExchangeAlert()
				}
			case <-stop:
				return
			}
		}
	}()
	defer func() {
		ticker.Stop()
		stop <- true
	}()

	ticker1 := time.NewTicker(time.Hour * 24)
	stop1 := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker1.C:
				if err := longRunTransactionCheck(); err != nil {
					util.Printf("long run transaction check error:%#v", err)
				}
			case <-stop1:
				return
			}
		}
	}()
	defer func() {
		ticker1.Stop()
		stop1 <- true
	}()

	ticker2 := time.NewTicker(time.Hour * 24)
	stop2 := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker2.C:
				if err := agentDailyCheck(); err != nil {
					util.Printf("long run agent level daily check error:%#v", err)
				}
				if err := customerDailyCheck(); err != nil {
					util.Printf("long run agent level daily check error:%#v", err)
				}
			case <-stop2:
				return
			}
		}
	}()
	defer func() {
		ticker2.Stop()
		stop2 <- true
	}()
	<-exit
	fmt.Println("long service exit")
}
