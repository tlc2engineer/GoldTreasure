package stat

import (
	"fmt"
	"time"
)

/*StatGor - вывод статистики*/
func StatGor(licStatChan <-chan bool, areaStatChan chan Area, digChan chan int, digAmountChan chan int, coinChan chan int) {
	var freeLicNum int = 0
	var payLicNum int = 0
	var areas int = 0
	var amounts int = 0
	var digged int = 0
	var diggedAmounts int = 0
	var coinSum int = 0

	go func() {
		for dig := range digChan {
			digged += dig
		}
	}()
	go func() {
		for digAmount := range digAmountChan {
			diggedAmounts += digAmount
		}
	}()
	go func() {
		for coins := range coinChan {
			coinSum += coins
		}
	}()
	go func() {
		for lic := range licStatChan {
			if lic {
				freeLicNum++
			} else {
				payLicNum++
			}
		}
	}()
	go func() {
		for area := range areaStatChan {
			areas++
			amounts += area.Amount
		}
	}()
	for {
		select {
		case <-time.Tick(time.Minute):
			fmt.Printf("LicFree: %d, LicPay: %d Areas: %d Amounts: %d Digged: %d DiggedAmounts: %d Coins: %d\n", freeLicNum, payLicNum, areas, amounts, digged, diggedAmounts, coinSum)
		}
	}

}

/*Area - область с сокровищами*/
type Area struct {
	Amount int
}
