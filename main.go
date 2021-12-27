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

type PrinterList struct {
	Printers []Printer
}

func (p *PrinterList) FindPrinter(ip string) *Printer {
	for i := range p.Printers {
		if p.Printers[i].IP == ip {
			//fmt.Printf("Address of i=%d:\t%p\n", printer, &printer)
			return &p.Printers[i]
		}
	}
	return nil
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

	var printerList PrinterList

	for _, ip := range ips {
		printer := Printer{IP: ip}
		printerList.Printers = append(printerList.Printers, printer)
	}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(ips...),
	)

	c01 := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(ips...),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Checking printer in", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		//log.Println("Visited", url)
		body := string(r.Body)
		lines := strings.Split(body, "\n")

		printer := printerList.FindPrinter(r.Request.URL.Host)
		//fmt.Printf("Address of p=&i=%d:\t%p\n", *printer, printer)

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
		}
		//log.Println(printer)
		//printerList = append(printerList, printer)
		wg.Done()
	})

	c01.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		//ip := r.Request.URL.Host
		//log.Println("Visited", url)
		body := string(r.Body)
		//fmt.Println(body)
		lines := strings.Split(body, "\n")

		printer := printerList.FindPrinter(r.Request.URL.Host)

		if url == "http://"+r.Request.URL.Host+"/DeviceInfo32.js" {
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
						dstatus = strings.Trim(dstatus, "''")
						printer.DStatus = strings.TrimSuffix(dstatus, ";//paper empty, error, ready")
						//log.Println(dstatus)
					}
				}
			}
		}
		wg.Done()
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		wg.Done()
	})

	c01.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		wg.Done()
	})

	for _, ip := range ips {
		link := "http://" + ip + "/ServerInfo31.js"
		wg.Add(1)
		go c.Visit(link)
	}
	wg.Wait()

	for _, ip := range ips {
		link := "http://" + ip + "/DeviceInfo32.js"
		wg.Add(1)
		go c01.Visit(link)
	}

	wg.Wait()

	fmt.Println(printerList)
}
