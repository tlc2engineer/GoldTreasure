package main

import (
	"Golden/api"
	"Golden/models"

	"fmt"
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
	var license *models.License
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
	for license == nil {
		lic, err := api.PostLicense()
		if err != nil {
			fmt.Println("lic err:", err)
		} else {
			license = lic
			fmt.Println("license:", *lic.ID, *lic.DigUsed, *lic.DigAllowed)
		}
	}
	//------------------------------

	//------------тестируем explore-------------
	for x := 1; x <= 20; x++ {
		for y := 1; y <= 20; y++ {
			amount, err := api.Explore(int64(x), int64(y))
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				fmt.Println(x, y, *amount)
				if *amount != 0 {
					count := *amount
					depth := 1
					for count > 0 && depth < 10 && *license.DigAllowed > 0 {
						tlist, err := api.DigPost(int64(depth), *license.ID, int64(x), int64(y))
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("DIG!")
							depth++
							if tlist != nil {
								for _, treasure := range tlist {
									fmt.Println("FIND!")
									wallet, err := api.PostCash(treasure)
									if err != nil {
										fmt.Println(err)
									} else {
										fmt.Println("Post good", wallet)
										count--
									}
								}
							}
						}
					}
				}
			}
		}
	}

}
