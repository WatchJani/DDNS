package main

import (
	"fmt"
	"log"
	"os"
	"root/cmd"
	"root/ddns"
	"time"
)

func main() {
	path := cmd.Flags()

	conf, err := ddns.LoadConfiguration(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ddns_api := ddns.New(5*time.Minute, conf)

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
