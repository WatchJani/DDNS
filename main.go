package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type DnsRecord struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"content"`
	TTL  int    `json:"ttl"`
}

func (d *DNS) updateDNSRecord(record DnsRecord) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", d.Zone_ID, d.Record_ID)

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+d.API_Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update DNS record: %s", resp.Status)
	}

	fmt.Println("DNS record updated successfully")
	return nil
}

func GetIp() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

type Configuration struct {
	Zone_ID   string `json:"ZONE_ID"`
	Record_ID string `json:"RECORD_ID"`
	API_Token string `json:"API_TOKEN"`
}

type DNS struct {
	*Configuration
}

func NewDNS() (*DNS, error) {
	conf, err := LoadConfiguration()
	if err != nil {
		return nil, err
	}

	return &DNS{
		Configuration: conf,
	}, nil
}

func newConfiguration() *Configuration {
	return &Configuration{}
}

func LoadConfiguration() (*Configuration, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := newConfiguration()
	if err := json.NewDecoder(file).Decode(config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return nil, err
	}

	return config, nil
}

func main() {
	dns_api, err := NewDNS()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		ip, err := GetIp()
		if err != nil {
			log.Println(err)
		}

		fmt.Println("My IP address is:", ip)

		if ip != "" {
			if err := dns_api.updateDNSRecord(DnsRecord{
				Type: "A",
				Name: "@",
				Data: ip,
				TTL:  3600,
			}); err != nil {
				fmt.Println("Error updating DNS record:", err)
			}
		}
	}
}
