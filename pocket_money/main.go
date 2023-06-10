package main

import (
	"fmt"
	"log"
	"time"
)

const LastDate = "2022-09-26"
const Week = 7 * 24 * time.Hour
const Payout = 0.5

func main() {
	ld, err := time.Parse("2006-01-02", LastDate)
	if err != nil {
		log.Panic("invalid last date")
	}

	pd := ld.Add(Week)
	amount := 0.5

	for pd.Before(time.Now()) {
		fmt.Printf("pay out at %v: %v, sum: €%.2f\n", pd, Payout, amount)
		pd = pd.Add(Week)
		amount += Payout
	}
}
