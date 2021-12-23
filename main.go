package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Printers struct {
	Printers []Printer `json:"printers"`
}

type Printer struct {
	IP string `json:"ip"`
}

func GetIPs(fileName string) []string  {
	ips := []string{}

	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("No se puede abrir el archivo", fileName)
		return ips
	}
	fmt.Println("Se abrió el archivo ip.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	fmt.Println("Impresoras a buscar: ")
	var printers Printers
	json.Unmarshal(byteValue, &printers)
	for i := 0; i < len(printers.Printers); i++ {
		ips = append(ips, printers.Printers[i].IP)
	}

	return ips
}

func main() {
	
	ips := GetIPs("ip.json")
	log.Println(ips)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("10.0.3.75", "10.0.3.69", "10.0.3.73"),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Checking printer in", r.URL.String())
	})

	var (
		serverName string
		model string
		manuf string
		dstatus string
	)
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		log.Println("Visited", url)
		body := string(r.Body)
		lines := strings.Split(body, "\n")
		if url == "http://" +  r.Request.URL.Host + "/ServerInfo31.js" {
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "myServer.name") {
					code := strings.Split(lines[i], "=")
					if len(code) == 2 {
						serverName = strings.TrimSpace(code[1])
						serverName = strings.Trim(serverName, "'';")
						log.Println(serverName)
					}
					break
				}
			}
		} else if url == "http://" +  r.Request.URL.Host + "/DeviceInfo32.js" {
			// fmt.Println(len(lines))
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "d1.model =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						model = strings.TrimSpace(code[1])
						model = strings.Trim(model, "'';")
						log.Println(model)
					}
				} else if strings.Contains(lines[i], "d1.manuf =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						manuf = strings.TrimSpace(code[1])
						manuf = strings.Trim(manuf, "'';")
						log.Println(manuf)
					}
				} else if strings.Contains(lines[i], "d1.Dstatus =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						dstatus = strings.TrimSpace(code[1])
						dstatus = strings.Trim(dstatus, "'")
						dstatus = strings.TrimSuffix(dstatus, "';//paper empty, error, ready")
						log.Println(dstatus)
					}
				}
			}
		}
	})




	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("http://10.0.3.73/ServerInfo31.js")
	c.Visit("http://10.0.3.73/DeviceInfo32.js")


}