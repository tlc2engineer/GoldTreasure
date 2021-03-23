package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
)

var expChainFull bool

const targetAmount = 2 //2
const segmentSize = 8  //8

/*Start - запуск базовой логики*/
func Start() {
	// базовый путь
	api.GetBasicPath()
	//------------stat-------------

	go stat.StatGor()
	//-----------------------------
	chTrlist := make(chan treasData, 20)
	chCoin := make(chan uint32, 100)
	go PostCashG(chTrlist, chCoin, true)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	chDig := make(chan DigData, 5)
	chLic := make(chan *models.License, 10)
	chUsedLic := make(chan *int64, 10)

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
	// go DigG(chDig, chTrlist, chLic, chUsedLic)
	// go DigG(chDig, chTrlist, chLic, chUsedLic)
	// go DigG(chDig, chTrlist, chLic, chUsedLic)
	// go DigG(chDig, chTrlist, chLic, chUsedLic)
	// go DigG(chDig, chTrlist, chLic, chUsedLic)

	//go exploreSegment(0, 0, 3498, 1748, 4, chDig)

	// go func() {
	// 	for {
	// 		select {
	// 		case <-time.Tick(time.Millisecond * 500):
	// 			if expChainFull != (len(chDig) > 100) {
	// 				expChainFull = len(chDig) > 100
	// 			}
	// 		}
	// 	}
	// }()

	//research(100, 100+8*20, 8)
	//research(100, 100+16*20, 16)
	//go searchSegments(0, 1841, 3491, 3491, segmentSize, targetAmount, chDig)
	//searchSegments(0, 0, 3491, 1741, segmentSize, targetAmount, chDig)

	//exploreSegment(0, 0, 3498, 1748, 4, chDig)

	//exploreSegment(0, 1750, 3498, 3498, 4, chDig)
	//go searchArea(0, 1750, 3400, 1500, segmentSize, chDig, targetAmount)
	searchArea(0, 0, 3400, 1700, segmentSize, chDig, targetAmount)
}
