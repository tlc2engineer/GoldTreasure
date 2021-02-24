package main

import (
	"Golden/api"
	"Golden/client"
	"Golden/client/operations"
	"Golden/models"
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
	conf := client.TransportConfig{Host: address + ":" + "8000", BasePath: "/", Schemes: []string{"http"}}
	client := client.NewHTTPClientWithConfig(nil, &conf)
	clientService := client.Operations

	game(clientService)

}

func game(cs operations.ClientService) {
	var x, y int64
	license := models.License{}
	log.Println("Start game")
	for x = 1; x <= 3500; x++ {
		for y = 1; y <= 3500; y++ {
			area := models.Area{
				PosX:  &x,
				PosY:  &y,
				SizeX: 1,
				SizeY: 1,
			}
			ep := operations.ExploreAreaParams{Args: &area}
			eaok, err := cs.ExploreArea(&ep)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Millisecond * 100)
				continue
			}
			amount := *eaok.GetPayload().Amount
			if int64(amount) == 0 {
				continue
			}
			var depth int64 = 1
			left := amount
			for depth <= 10 && left > 0 {
				for license.ID == nil || *license.DigAllowed >= *license.DigAllowed {
					lp := operations.IssueLicenseParams{}
					lok, err := cs.IssueLicense(&lp)
					if err != nil {
						fmt.Println(err)
						time.Sleep(time.Millisecond * 100)
						continue
					}
					license = *lok.Payload
				}
				dig := models.Dig{LicenseID: license.ID, PosX: &x, PosY: &y, Depth: &depth}
				dp := operations.DigParams{Args: &dig}
				dok, err := cs.Dig(&dp)
				if err != nil {
					fmt.Println(err)
					time.Sleep(time.Millisecond * 100)
					continue
				}

				tlist := dok.Payload
				*license.DigUsed++
				depth++
				if tlist != nil {
					for _, tr := range tlist {
						cp := operations.CashParams{Args: tr}
						cok, err := cs.Cash(&cp)
						if err != nil {
							fmt.Println(err)
							time.Sleep(time.Millisecond * 100)
							continue
						}
						cash := cok.Payload
						if cash != nil && len(cash) > 0 {
							left--
						}
					}
				}

			}
		}
	}
}
