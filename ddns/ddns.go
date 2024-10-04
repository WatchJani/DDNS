package ddns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const END_POINT_FOR_IP_ADDRESS = "https://api.ipify.org"

type Updater func()

type DDNS struct {
	*Configuration
	timeUpdater time.Duration
}

func New(time time.Duration) (*DDNS, error) {
	conf, err := LoadConfiguration()
	if err != nil {
		return nil, err
	}

	return &DDNS{
		Configuration: conf,
		timeUpdater:   time,
	}, nil
}

type Configuration struct {
	Zone_ID   string `json:"ZONE_ID"`
	Record_ID string `json:"RECORD_ID"`
	API_Token string `json:"API_TOKEN"`

	IP_ADDRESS string `json:"CURRENT_IP_ADDRESS"`
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

type DnsRecord struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"content"`
	TTL  int    `json:"ttl"`
}

func (d *DDNS) UpdateDNSRecord(record DnsRecord) error {
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
	resp, err := http.Get(END_POINT_FOR_IP_ADDRESS)
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

func (d *DDNS) Update(fn Updater) {
	ticker := time.NewTicker(d.timeUpdater)
	defer ticker.Stop()

	for {
		fn()
		<-ticker.C
	}
}

func UpdateConfig(filename string, dns *DDNS) error {
	data, err := json.MarshalIndent(dns, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
