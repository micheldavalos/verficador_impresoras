package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type Printers struct {
	Printers []Printer `json:"printers"`
}

type Printer struct {
	IP string `json:"ip"`
}

func main()  {
	fileName := "ip.json"
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("No se puede abrir el archivo", fileName)
		return
	}
	fmt.Println("Se abrió el archivo ip.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	fmt.Println("Impresoras a buscar: ")
	var printers Printers
	json.Unmarshal(byteValue, &printers)
	for i := 0; i < len(printers.Printers); i++ {
		fmt.Println("Ip: " + printers.Printers[i].IP)
	}

}