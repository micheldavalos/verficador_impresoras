package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
}