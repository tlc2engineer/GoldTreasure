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
	license = updateLicense()
	//------------------------------

	//------------тестируем explore-------------
	for x := 1; x <= 3500; x++ {
		for y := 1; y <= 3500; y++ {
			amount, err := api.Explore(int64(x), int64(y))
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					count := *amount
					depth := 1
					for count > 0 && depth <= 10 {
						if *license.DigAllowed <= *license.DigUsed {
							license = updateLicense()
						}
						tlist, err := api.DigPost(int64(depth), *license.ID, int64(x), int64(y))
						if err != nil {
							fmt.Println(err)
						} else {
							*license.DigUsed++
							depth++
							if tlist != nil {
								for _, treasure := range tlist {
									_, err := api.PostCash(treasure)
									if err != nil {
										fmt.Println(err)
									} else {
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

func updateLicense() *models.License {
	for {
		lic, err := api.PostLicense()
		if err != nil {
			fmt.Println("lic err:", err)
		} else {

			fmt.Println("license:", *lic.ID, *lic.DigUsed, *lic.DigAllowed)
			return lic
		}
	}
}
