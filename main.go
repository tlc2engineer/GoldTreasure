package main

import (
	"Golden/api"
	"Golden/models"
	"sync"

	"fmt"
	"net/http"
	"os"
	"time"
)

var numLic int = 0
var mu = new(sync.Mutex)

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
	chTrlist := make(chan models.TreasureList, 5)
	chCoin := make(chan uint32, 100)
	go PostCashG(chTrlist, chCoin, true)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	chDig := make(chan DigData, 5)
	chLic := make(chan *models.License, 10)
	chUsedLic := make(chan *int64)

	go func() {
		for l := range chUsedLic {
			if *l > 0 {
				mu.Lock()
				numLic--
				mu.Unlock()
			}

		}
	}()
	go LicGor(chCoin, chLic)
	go LicGor(chCoin, chLic)
	go LicGor(chCoin, chLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)
	go DigG(chDig, chTrlist, chLic, chUsedLic)

	//------------тестируем explore-------------
	for x := 1; x < 3500; x++ {
		for y := 1; y < 3500; y++ {
			amount, err := api.Explore(int64(x), int64(y))
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					digData := DigData{x: int64(x), y: int64(y), amount: int64(*amount)}
					chDig <- digData

				}
			}
		}
	}

}

func updateLicense() *models.License {

	for {
		wallet := models.Wallet{}
		lic, err := api.PostLicense(wallet)
		if err != nil {
			//fmt.Println("license err:", err)
		} else {
			return lic
		}
	}
}

/*PostCashG - горутина отправки сообщений*/
func PostCashG(ch chan models.TreasureList, chCoins chan uint32, toLic bool) {
	for tlist := range ch {
		for _, treasure := range tlist {
			w, err := api.PostCash(treasure)
			for err != nil && err.Error() == "Status not ok:503" {
				w, err = api.PostCash(treasure)
			}
			if err != nil {
				fmt.Println("Post cash err", err)
			} else {
				if w != nil && toLic {
					for _, coin := range *w {
						chCoins <- coin
					}
				}
			}

		}
	}
}

/*DigData - вспомагательная структура*/
type DigData struct {
	x, y, amount int64
}

/*DigG - горутина копания*/
func DigG(ch chan DigData, cht chan models.TreasureList, chLic chan *models.License, chUsedLic chan *int64) {
	var license *models.License // лицензия
	for ddata := range ch {
		trCount := ddata.amount // число ненайденных сокровиц
		depth := 1              //глубина
		for trCount > 0 && depth <= 10 {
			if license == nil || *license.DigAllowed <= *license.DigUsed {
				if license != nil {
					chUsedLic <- license.ID // использованная лицензия
				}
				license = <-chLic // полцчаем лицензию
			}
			tlist, err := api.DigPost(int64(depth), *license.ID, ddata.x, ddata.y)
			if err != nil {
				fmt.Println("Dig err", err)
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

/*LicGor - лицензии*/
func LicGor(chCoin chan uint32, chLic chan *models.License) {
	var coin uint32
	var wallet models.Wallet

	var numCoin, numFree int
	for {
		mu.Lock()
		if numLic >= 10 {
			mu.Unlock()
			time.Sleep(time.Millisecond)
			continue
		} else {
			numLic++
			mu.Unlock()
		}
		select {

		case coin = <-chCoin:
			wallet = models.Wallet{} // платная
			wallet = append(wallet, coin)
			numCoin++
		default:
			wallet = models.Wallet{} // бесплатная
			numFree++
		}
		lic, err := api.PostLicense(wallet)
		if err != nil {
			mu.Lock()
			numLic--
			mu.Unlock()
		} else {

			chLic <- lic
		}
	}
}
