package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const defaultNTPServer = "0.beevik-ntp.pool.ntp.org"

func main() {
	serverFlag := flag.String(
		"server",
		defaultNTPServer,
		"NTP server address (default: 0.beevik-ntp.pool.ntp.org)",
	)

	flag.Parse()

	if *serverFlag == "" {
		exitWithError(fmt.Errorf("NTP server address cannot be empty"))
	}

	getNetworkTime(*serverFlag)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func getNetworkTime(server string) {
	response, err := ntp.Query(server)
	if err != nil {
		exitWithError(err)
	}

	if err = response.Validate(); err != nil {
		exitWithError(fmt.Errorf("invalid NTP response: %w", err))
	}

	correctedTime := response.Time

	fmt.Printf("NTP Server: %s\n", server)
	fmt.Printf("Current Time: %v\n", correctedTime)
	fmt.Printf("Clock Offset: %v\n", response.ClockOffset)
	fmt.Printf("Round Trip Time: %v\n", response.RTT)
	fmt.Printf("Precision: %v\n", response.Precision)
	fmt.Printf("Stratum: %d\n", response.Stratum)
}
