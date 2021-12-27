package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type IPList struct {
	IpList []IP `json:"ips"`
}

type IP struct {
	IP string
}

type Printer struct {
	IP         string `json:"ip"`
	ServerName string
	Model      string
	Manuf      string
	DStatus    string
}

func GetIPs(fileName string) []string {
	ips := []string{}

	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Can't open file", fileName)
		return ips
	}
	log.Println("Successful open file", fileName)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// fmt.Println("Impresoras a buscar: ")
	var ipList IPList
	err = json.Unmarshal(byteValue, &ipList)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(ipList.IpList); i++ {
		ips = append(ips, ipList.IpList[i].IP)
	}
	return ips
}

func main() {
	ips := GetIPs("ip.json")
	log.Println("IP list", ips)

	var printerList []Printer

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(ips...),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Checking printer in", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		printer := Printer{}
		printer.IP = r.Request.URL.Host
		url := r.Request.URL.String()
		//log.Println("Visited", url)
		body := string(r.Body)
		lines := strings.Split(body, "\n")

		if url == "http://"+r.Request.URL.Host+"/ServerInfo31.js" {
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "myServer.name") {
					code := strings.Split(lines[i], "=")
					if len(code) == 2 {
						serverName := strings.TrimSpace(code[1])
						printer.ServerName = strings.Trim(serverName, "'';")
						//log.Println("Server Name: ", serverName)
					}
					break
				}
			}
		} else if url == "http://"+r.Request.URL.Host+"/DeviceInfo32.js" {
			// fmt.Println(len(lines))
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "d1.model =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						model := strings.TrimSpace(code[1])
						printer.Model = strings.Trim(model, "'';")
						//log.Println(model)
					}
				} else if strings.Contains(lines[i], "d1.manuf =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						manuf := strings.TrimSpace(code[1])
						printer.Manuf = strings.Trim(manuf, "'';")
						//log.Println(manuf)
					}
				} else if strings.Contains(lines[i], "d1.Dstatus =") {
					code := strings.Split(lines[i], "=")
					// fmt.Println(len(code))
					if len(code) >= 2 {
						dstatus := strings.TrimSpace(code[1])
						dstatus = strings.Trim(dstatus, "'")
						printer.DStatus = strings.TrimSuffix(dstatus, "';//paper empty, error, ready")
						//log.Println(dstatus)
					}
				}
			}
		}
		log.Println(printer)
		printerList = append(printerList, printer)
		wg.Done()
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		wg.Done()
	})

	for _, ip := range ips {
		link := "http://" + ip + "/ServerInfo31.js"
		wg.Add(1)
		go c.Visit(link)
	}

	wg.Wait()

	//c.Visit("http://10.0.3.75/DeviceInfo32.js")

}
