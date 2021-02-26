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
	lst, err := api.GetLicenses()
	if err != nil {
		fmt.Println("L list", err)
	} else {
		for _, lic := range lst {
			fmt.Println("lic list error:", *lic.ID)
		}
	}
	chTrlist := make(chan models.TreasureList, 2)
	go PostCashG(chTrlist)
	//------------тестируем explore-------------
	for x := 1; x < 3500; x++ {
		for y := 1; y < 3500; y++ {
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
								count--
								chTrlist <- tlist
								// for _, treasure := range tlist {
								// 	_, err := api.PostCash(treasure)
								// 	if err != nil {
								// 		fmt.Println(err)
								// 	} else {
								// 		count--
								// 	}
								// }
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
			fmt.Println("license err:", err)
		} else {
			return lic
		}
	}
}

/*PostCashG - горутина отправки сообщений*/
func PostCashG(ch chan models.TreasureList) {
	for tlist := range ch {
		for _, treasure := range tlist {
			_, err := api.PostCash(treasure)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

/*DigData - вспомагательная структура*/
type DigData struct {
	x, y, amount int64
}

/*DigG - горутина копания*/
func DigG(ch chan DigData, cht chan models.TreasureList) {
	var license *models.License // лицензия
	for ddata := range ch {
		trCount := ddata.amount // число ненайденных сокровиц
		depth := 1              //глубина
		for trCount > 0 && depth <= 10 {
			if license == nil || *license.DigAllowed <= *license.DigUsed {
				license = updateLicense()
			}
			tlist, err := api.DigPost(int64(depth), *license.ID, ddata.x, ddata.y)
			if err != nil {
				fmt.Println(err)
			} else {
				*license.DigUsed++
				depth++
				if tlist != nil {
					trCount--
					cht <- tlist
				}
			}
		}

	}
}
