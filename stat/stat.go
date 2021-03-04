package stat

import (
	"fmt"
	"time"
)

/*StatChan - канал статистики*/
var statChan chan Stat = make(chan Stat)

/*StatGor - вывод статистики*/
func StatGor() {
	var freeLicNum int = 0
	var payLicNum int = 0
	var areas int = 0
	var amounts int = 0
	var digged int = 0
	var diggedAmounts int = 0
	var coinSum int = 0
	var errors int = 0
	var exErrors int = 0
	var digErrors int = 0
	var cashErrors int = 0
	var licErrors int = 0
	go func() {
		for statEnt := range statChan {
			switch statEnt.(type) {
			case AreaStat:
				area, _ := statEnt.(AreaStat)
				areas++
				amounts += area.GetAmount()
			case LicStat:
				lic, _ := statEnt.(LicStat)
				if lic.IsFree() {
					freeLicNum++
				} else {
					payLicNum++
				}
			case DigStat:
				dig, _ := statEnt.(DigStat)
				digged += dig.GetDigg()
				diggedAmounts += dig.GetAmounts()
			case CoinStat:
				cs := statEnt.(CoinStat)
				coinSum += cs.GetCoins()
			case ErrStat:
				es := statEnt.(ErrStat)
				errors++
				switch es.ErrName() {
				case "Exp":
					exErrors++
				case "Dig":
					digErrors++
				case "Cash":
					cashErrors++
				case "Lic":
					licErrors++
				}
			}
		}
	}()
	// go func() {
	// 	for dig := range digChan {
	// 		digged += dig
	// 	}
	// }()
	// go func() {
	// 	for digAmount := range digAmountChan {
	// 		diggedAmounts += digAmount
	// 	}
	// }()
	// go func() {
	// 	for coins := range coinChan {
	// 		coinSum += coins
	// 	}
	// }()
	// go func() {
	// 	for lic := range licStatChan {
	// 		if lic {
	// 			freeLicNum++
	// 		} else {
	// 			payLicNum++
	// 		}
	// 	}
	// }()
	// go func() {
	// 	for area := range areaStatChan {
	// 		areas++
	// 		amounts += area.Amount
	// 	}
	// }()
	for {
		select {
		case <-time.Tick(time.Minute):
			fmt.Printf("LicFree: %d, LicPay: %d Areas: %d Amounts: %d Digged: %d DiggedAmounts: %d Coins: %d\n", freeLicNum, payLicNum, areas, amounts, digged, diggedAmounts, coinSum)
			fmt.Printf("Errors: %d, ExpErr: %d DigErr: %d CashErr: %d LicErr: %d \n", errors, exErrors, digErrors, cashErrors, licErrors)
		}
	}

}

/*Area - область с сокровищами*/
type Area struct {
	Amount int
}

/*Stat - статистика*/
type Stat interface {
	StatName() string
}

/*AreaStat - статистика исследованной обрасти*/
type AreaStat interface {
	Stat
	GetAmount() int
}

/*LicStat - статистика лицензий*/
type LicStat interface {
	Stat
	IsFree() bool
}

/*DigStat - статистика копаний*/
type DigStat interface {
	Stat
	GetDigg() int
	GetAmounts() int
}

/*CoinStat - количество денег*/
type CoinStat interface {
	Stat
	GetCoins() int
}

/*ErrStat - статистика ошибок*/
type ErrStat interface {
	Stat
	ErrName() string
}
