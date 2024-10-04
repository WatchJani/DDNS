package main

import (
	"fmt"
	"log"
	"os"
	"rot/ddns"
	"time"
)

func main() {
	ddns_api, err := ddns.New(5 * time.Minute)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("My current Ip addres", ddns_api.IP_ADDRESS)

	ddns_api.Update(func() {
		ip, err := ddns.GetIp()
		if err != nil {
			log.Println(err)
			return
		}

		if ip != ddns_api.IP_ADDRESS {
			ddns_api.IP_ADDRESS = ip
			fmt.Println("My IP new address is:", ip)

			if err := ddns_api.UpdateDNSRecord(ddns.DnsRecord{
				Type: "A",
				Name: "@",
				Data: ddns_api.IP_ADDRESS,
				TTL:  3600,
			}); err != nil {
				fmt.Println("Error updating DNS record:", err)
				return
			}

			ddns.UpdateConfig("./config.json", ddns_api)
		}
	})
}
