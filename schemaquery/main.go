package main

import (
	"flag"
	"fmt"
)

var userName = flag.String("u", "root", "DB user name")
var password = flag.String("p", "12345", "DB password")
var host = flag.String("h", "cnpvgvb1ep052.pvgl.sap.corp:3306", "DB host name")
var eshopId = flag.Int64("q", -1, "eshop tenant id")
var eshopName = flag.String("sn", "", "eshop name")

func PrintArgs() {
	fmt.Println("user:", *userName)
	fmt.Println("password:", *password)
	fmt.Println("host:", *host)
}

func main() {
	flag.Parse()

	if *eshopId < 0 && len(*eshopName) <= 0 {
		fmt.Println("You must provide a valid eshop id or eshop name")
		flag.PrintDefaults()
		return
	}
	DSNWithoutSchema = *userName + ":" + *password + "@tcp(" + *host + ")/"

	if *eshopId > 0 {
		fmt.Println("Use eshop id to get eshop schema name. eshop id:", *eshopId)
		PrintArgs()
		names := GetSchemaById(*eshopId)
		fmt.Println("schema name : ", names)
		return
	} else {
		fmt.Println("Use eshop name to get schema name. eshop name:", *eshopName)
		PrintArgs()
		names := GetSchemaByName(*eshopName)
		fmt.Println("schema name : ", names)
	}
}
