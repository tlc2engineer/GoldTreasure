package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
	"fmt"
)

/*PostCashG - горутина отправки сообщений о деньгах*/
func PostCashG(ch chan treasData, chCoins chan uint32, toLic bool) {
	for tData := range ch {
		sum := 0
		for _, treasure := range tData.tlist {
			w, err := api.PostCash(treasure)
			for err != nil && err.Error() == "Status not ok:503" {
				stat.NewStatErr(stat.Cash)
				w, err = api.PostCash(treasure)
			}
			if err != nil {
				stat.NewStatErr(stat.Cash)
				fmt.Println("Post cash err", err)
			} else {
				stat.NewSendTlist()
				stat.NewCoinStat(len(*w))
				sum += len(*w)
				if w != nil && toLic {
					for _, coin := range *w {
						chCoins <- coin
					}
				}
			}

		}
		stat.DepthStat(int(tData.dt), tData.depth, sum)
	}
}

type treasData struct {
	tlist       models.TreasureList
	x, y, depth int
	dt          int64
}
