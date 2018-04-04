package main

import (
	"time"
	"util"
)

func longRun() {
	ticker := time.NewTicker(time.Hour * 8)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := getLatestRates(); err != nil {
					util.FailToGetCurrencyExchangeAlert()
				}
				var err error
				activeCurrencyRate, err = getAcitveCurrencyRate()
				if err != nil {
					util.Println("fail to get latest active currency rate")
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

	ticker1 := time.NewTicker(time.Hour * 1)
	stop1 := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
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
}
