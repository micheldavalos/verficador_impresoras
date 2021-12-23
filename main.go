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
		colly.AllowedDomains("10.0.3.75", "10.0.3.69", "10.0.3.65"),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Checking printer in", r.URL.String())
	})

	var serverName string
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		// b, _ := r.Request.Marshal()
		serverInfo31 := string(r.Body)
		lines := strings.Split(serverInfo31, "\n")
		// fmt.Println(serverInfo31)
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
	})


	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("http://10.0.3.75/ServerInfo31.js")


}