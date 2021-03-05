package main

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
	"sync"

	"fmt"
	"net/http"
	"os"
	"time"
)

var numLic int = 0
var mu = new(sync.Mutex)

var areaStatChan = make(chan stat.Area)
var digStatChan = make(chan int)
var digAmountStatChan = make(chan int)
var coinStatChan = make(chan int)

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
	//------------stat-------------
	licStatChan := make(chan bool)
	go stat.StatGor(licStatChan, areaStatChan, digStatChan, digAmountStatChan, coinStatChan)
	//-----------------------------
	chTrlist := make(chan models.TreasureList, 5)
	chCoin := make(chan uint32, 100)
	go PostCashG(chTrlist, chCoin, true)
	go PostCashG(chTrlist, chCoin, false)
	go PostCashG(chTrlist, chCoin, false)
	chDig := make(chan DigData, 100)
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
	go LicGor(chCoin, chLic, licStatChan)
	go LicGor(chCoin, chLic, licStatChan)
	go LicGor(chCoin, chLic, licStatChan)
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

	//go exploreSegment(0, 0, 3498, 1748, 4, chDig)

	//go searchSegments(0, 1741, 3491, 3491, 8, 4, chDig)
	searchSegments(0, 0, 3491, 1741, 8, 4, chDig)
	//exploreSegment(0, 0, 3498, 1748, 4, chDig)

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
				coinStatChan <- len(*w)
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
		digStatChan <- depth
		digAmountStatChan <- (int(ddata.amount) - int(trCount))

	}
}

/*LicGor - лицензии*/
func LicGor(chCoin chan uint32, chLic chan *models.License, statChan chan bool) {
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
			mu.Lock()
			numLic--
			mu.Unlock()
		} else {
			statChan <- free
			chLic <- lic
		}
	}
}

/*exploreArea - исследование области*/
func exploreArea(xbg, ybg, xend, yend int, ch chan DigData, targetMoney int) int {
	sum := 0
m1:
	for x := xbg; x < xend; x++ {
		for y := ybg; y < yend; y++ {
			amount, err := api.Explore(int64(x), int64(y), 1, 1)
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					digData := DigData{x: int64(x), y: int64(y), amount: int64(*amount)}
					// if len(ch) > 98 {
					// 	fmt.Println("Chain full!")
					// }
					areaStatChan <- stat.Area{Amount: int(*amount)}
					ch <- digData
					sum += int(*amount)
					if targetMoney == sum {
						break m1
					}

				}
			}
		}
	}
	if sum != targetMoney {
		fmt.Printf("Exp error t:%d s:%d\n", targetMoney, sum)
	}
	return targetMoney
}

/*exploreSegment - исследование сегмента*/
func exploreSegment(xbg, ybg, xend, yend, size int, ch chan DigData) int {
	sum := 0
	for x := xbg; x < xend; x += size {
		for y := ybg; y < yend; y += size {
			amount, err := api.Explore(int64(x), int64(y), int64(size), int64(size))
			if err != nil {
				fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					if size >= 4 {
						money := int(*amount)
						tsum := 0
					m1:
						for x1 := x; x1 < x+size; x1 += size / 2 {
							for y1 := y; y1 < y+size; y1 += size / 2 {
								am := exploreSegment(x1, y1, x1+size/2, y1+size/2, size/2, ch)
								sum += am
								tsum += am
								if money == tsum {
									break m1
								}
							}
						}
						if money != tsum {
							fmt.Printf("t: %d fact: %d ", money, tsum)
						}
					} else {
						money := int(*amount)
						sum += exploreArea(x, y, x+size, y+size, ch, money)
					}

				}
			}
		}
	}
	return sum

}

func divideSegment(x, y, size int64, ch chan DigData) int {
	sum := 0
	amount, err := api.Explore(x, y, size, size)
	if err != nil {
		fmt.Println("Exp err:", err)
		return 0
	}

	if *amount != 0 {
		if size >= 4 {
			money := int(*amount)
			tsum := 0
		m1:
			for x1 := x; x1 < x+size; x1 += size / 2 {
				for y1 := y; y1 < y+size; y1 += size / 2 {
					am := divideSegment(x1, y1, size/2, ch)
					sum += am
					tsum += am
					if money == tsum {
						break m1
					}
				}
			}
			if money != tsum {
				fmt.Printf("------t: %d fact: %d ", money, tsum)
			}
		} else {
			money := int(*amount)
			sum += exploreArea(int(x), int(y), int(x+size), int(y+size), ch, money)
		}

	}
	return sum
}

func searchSegments(x0, y0, xe, ye, size, limit int, ch chan DigData) {
	for x := x0; x < xe; x += size {
		for y := y0; y < ye; y += size {
			amount, err := api.Explore(int64(x), int64(y), int64(size), int64(size))
			if err != nil {
				fmt.Println("Exp err:", err)
				return
			}
			if int(*amount) >= limit {
				go func(x1, y1, size1, amount1 int, ch chan DigData) {
					fact := divideSegment(int64(x1), int64(y1), int64(size1), ch)
					if fact != amount1 {
						fmt.Printf("search t: %d fact: %d ", amount, fact)
					}
				}(x, y, size, int(*amount), ch)
			}
		}
	}
}
