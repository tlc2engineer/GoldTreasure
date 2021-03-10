package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
	"fmt"
)

/*PostCashG - горутина отправки сообщений о деньгах*/
func PostCashG(ch chan models.TreasureList, chCoins chan uint32, toLic bool) {
	for tlist := range ch {
		for _, treasure := range tlist {
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
				if w != nil && toLic {
					for _, coin := range *w {
						chCoins <- coin
					}
				}
			}

		}
	}
}
