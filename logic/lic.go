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

	var free bool

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
			free = false
		default:
			wallet = models.Wallet{} // бесплатная
			free = true

		}
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
