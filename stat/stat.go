package stat

import (
	"fmt"
	"time"
)

/*ReqType - тип ошибки*/
type ReqType int

const (
	Exp ReqType = iota //exp
	Digg
	Cash
	Lic
)

/*StatChan - канал статистики*/
var statChan chan Stat = make(chan Stat)
var numReq = 0
var numExpReq, numDigReq, numLicReq, numCashReq int

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
	n := 0
	for {
		select {
		case <-time.Tick(time.Minute):
			n++
			fmt.Println("N:", n)
			fmt.Printf("LicFree: %d, LicPay: %d Areas: %d Amounts: %d Digged: %d DiggedAmounts: %d Coins: %d\n", freeLicNum, payLicNum, areas, amounts, digged, diggedAmounts, coinSum)
			fmt.Printf("Errors: %d, ExpErr: %d DigErr: %d CashErr: %d LicErr: %d \n", errors, exErrors, digErrors, cashErrors, licErrors)
			fmt.Printf("DigTreas: %d,DigTlist: %d SendTlist: %d \n", digTreasures, digTlist, sendTlist)
			fmt.Printf("Req: %d,Exp: %d,Dig: %d,Lic: %d, Cash: %d\n", numReq, numExpReq, numDigReq, numLicReq, numCashReq)
			if n == 10 {
				for _, level := range levels {
					fmt.Printf("Depth: %d,time: %5.2f,treas: %5.2f \n", level.depth, (float64(level.totalTime))/float64(level.total), float64(level.totalTreasures)/float64(level.total))
				}
			}
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
	Type() ReqType
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
