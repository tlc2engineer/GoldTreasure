package main

import (
	"Golden/api"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "localhost"
	}
	req := fmt.Sprintf("http://%s:8000/balance", address)
	fmt.Println(req)
	for {
		resp, err := http.Get(req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		fmt.Println("again")
		time.Sleep(time.Second)
	}
	// базовый путь
	api.GetBasicPath()
	// тестируем balance
	//
	//------------тестируем explore-------------
	// for x := 1; x <= 20; x++ {
	// 	for y := 1; y <= 20; y++ {
	// 		amount, err := api.Explore(int64(x), int64(y))
	// 		if err != nil {
	// 			fmt.Println("Exp err:", err)
	// 		} else {
	// 			fmt.Println(x, y, *amount)
	// 		}
	// 	}
	// }
	//---------post license-----------
	lic, err := api.PostLicense()
	if err != nil {
		fmt.Println("lic err:", err)
	} else {
		fmt.Println("license:", *lic.ID, *lic.DigUsed, *lic.DigAllowed)
	}
	//------------------------------
	log.Println("GO!")
	address = os.Getenv("ADDRESS")
	log.Printf("Address:%s\n", address)
	if address == "" {
		address = "localhost"
	}

	//------------тестируем explore-------------
	for x := 1; x <= 20; x++ {
		for y := 1; y <= 20; y++ {
			amount, err := api.Explore(int64(x), int64(y))
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				fmt.Println(x, y, *amount)
			}
		}
	}

}
