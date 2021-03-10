package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
)

const maxDepth = 10

/*DigG - горутина копания*/
func DigG(ch chan DigData, cht chan models.TreasureList, chLic chan *models.License, chUsedLic chan *int64) {
	var license *models.License // лицензия
	for ddata := range ch {
		trCount := ddata.amount // число ненайденных сокровиц
		depth := 1              //глубина
		for trCount > 0 && depth <= maxDepth {
			if license == nil || *license.DigAllowed <= *license.DigUsed {
				if license != nil {
					chUsedLic <- license.ID // использованная лицензия
				}
				license = <-chLic // полцчаем лицензию
			}
			tlist, err := api.DigPost(int64(depth), *license.ID, ddata.x, ddata.y)
			if err != nil {
				stat.NewStatErr(stat.Digg)
				//fmt.Println("Dig err", err)
			} else {
				*license.DigUsed++
				depth++
				if tlist != nil {
					trCount--
					stat.NewSendTreas(len(tlist))
					stat.NewDigTlist()
					cht <- tlist
					// if len(cht) > 99 {
					// 	fmt.Println("cht chain is full")
					// }
				}
			}
		}
		stat.NewDsStat(depth, (int(ddata.amount) - int(trCount)))

	}
}

/*DigData - вспомагательная структура*/
type DigData struct {
	x, y, amount int64
}
