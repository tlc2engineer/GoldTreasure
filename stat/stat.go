package stat

import (
	"fmt"
	"runtime"
	"time"
)

/*ErrType - тип ошибки*/
type ErrType int

const (
	Exp ErrType = iota //exp
	Digg
	Cash
	Lic
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
				switch es.Type() {
				case Exp:
					exErrors++
				case Digg:
					digErrors++
				case Cash:
					cashErrors++
				case Lic:
					licErrors++
				}
			}
		}
	}()

	for {
		select {
		case <-time.Tick(time.Minute):
			fmt.Printf("LicFree: %d, LicPay: %d Areas: %d Amounts: %d Digged: %d DiggedAmounts: %d Coins: %d\n", freeLicNum, payLicNum, areas, amounts, digged, diggedAmounts, coinSum)
			fmt.Printf("Errors: %d, ExpErr: %d DigErr: %d CashErr: %d LicErr: %d \n", errors, exErrors, digErrors, cashErrors, licErrors)
			fmt.Printf("DigTreas: %d,DigTlist: %d SendTlist: %d\n", digTreasures, digTlist, sendTlist)
		case <-time.After(time.Minute * 3):
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			// For info on each, see: https://golang.org/pkg/runtime/#MemStats
			fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
			fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
			fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
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
	Type() ErrType
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
