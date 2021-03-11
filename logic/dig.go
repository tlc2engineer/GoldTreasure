package logic

import (
	"Golden/api"
	"Golden/models"
	"Golden/stat"
	"time"
)

const maxDepth = 10

/*DigG - горутина копания*/
func DigG(ch chan DigData, cht chan treasData, chLic chan *models.License, chUsedLic chan *int64) {
	var license *models.License // лицензия
	for ddata := range ch {
		trCount := ddata.amount // число ненайденных сокровиц
		depth := 1              //глубина
		for trCount > 0 && depth <= maxDepth {
			if license == nil || *license.DigAllowed <= *license.DigUsed {
				if license != nil {
					chUsedLic <- license.ID // использованная лицензия
				}
				license = <-chLic // получаем лицензию
			}
			tbg := time.Now()
			tlist, err := api.DigPost(int64(depth), *license.ID, ddata.x, ddata.y)
			dt := time.Since(tbg).Milliseconds()
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
					cht <- treasData{tlist: tlist, x: int(ddata.x), y: int(ddata.y), depth: int(depth - 1), dt: dt}
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
