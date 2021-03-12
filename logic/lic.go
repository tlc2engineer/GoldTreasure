package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
	"sync"
	"time"
)

var mu = new(sync.Mutex)
var numLic int = 0

/*LicGor - лицензии*/
func LicGor(chCoin chan uint32, chLic chan *models.License) {
	var coin uint32
	var wallet models.Wallet
	var free bool = false
	var randNumb = 1
	for {
		//Изменеие количества лицензий
		mu.Lock()
		if numLic >= 10 {
			mu.Unlock()
			time.Sleep(time.Millisecond)
			continue
		} else {
			numLic++
			mu.Unlock()
		}
		//----------------------------
		wallet = models.Wallet{} // новый кошелек
		free = false
		if len(chCoin) > 11 && randNumb == 3 {
			var count int = 0
			for count < 11 {
				coin = <-chCoin
				wallet = append(wallet, coin)
				count++
			}
		} else {
			if randNumb == 2 {
				select {
				case coin = <-chCoin:
					wallet = append(wallet, coin)
					free = false
				default:
					free = true
				}
			} else {
				free = true
			}
		}
		//----------------------------
		lic, err := api.PostLicense(wallet)
		if err != nil {
			stat.NewStatErr(stat.Lic)
			mu.Lock()
			numLic--
			mu.Unlock()
		} else {
			stat.NewLcStat(free)
			chLic <- lic
		}
		randNumb++
		if randNumb > 3 {
			randNumb = 1
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
