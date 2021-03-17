package stat

import (
	"encoding/base64"
	"fmt"
	"sync"
)

type level struct {
	depth          int
	totalTreasures int
	totalTime      int
	total          int
	min, max       int
}

var levels map[int]*level
var licDep map[int]licDigStat
var mut *sync.Mutex = new(sync.Mutex)
var treasureMap [][]byte

func init() {
	levels = make(map[int]*level)
	for i := 1; i <= 10; i++ {
		levels[i] = new(level)
		levels[i].depth = i
	}
	licStatMap = make(map[int]int)
	for i := 1; i <= 21; i++ {
		licStatMap[i] = 0
	}
	licDep = make(map[int]licDigStat)
	treasureMap = make([][]byte, 600)
	for i := 0; i < 600; i++ {
		treasureMap[i] = make([]byte, 600)
	}
}

/*DigDeepStat - статистика по глубине и количеству сокровищ*/
func DigDeepStat(t int, depth int, treasures int, licType int, x int, y int) {
	level := levels[depth]
	level.total++
	level.totalTime += t
	level.totalTreasures += treasures
	if level.min == 0 && level.max == 0 {
		level.min = treasures
		level.max = treasures
	}
	if level.min > treasures {
		level.min = treasures
	}
	if level.max < treasures {
		level.max = treasures
	}
	mut.Lock()
	stat, ok := licDep[licType]
	if !ok {
		licDep[licType] = licDigStat{sumDt: int64(t), num: 1}
	} else {
		stat.num++
		stat.sumDt += int64(t)
		licDep[licType] = stat
	}
	mut.Unlock()
	if x < 300 && y < 300 {
		treasureMap[x][y] = byte(depth)
	}

}

type licDigStat struct {
	sumDt int64
	num   int
}

type mapUnit struct {
	deep byte
}

func createByteArr(trMapData [][]byte) []byte {
	retArray := make([]byte, 0)
	var coord int = 0
	for x := 0; x < 300; x++ {
		for y := 0; y < 300; y++ {
			if treasureMap[x][y] > 0 {
				pCoord := x*300 + y
				dL := pCoord - coord
				var b0, b1 uint8 = uint8(dL >> 8), uint8(dL & 0xff)
				retArray = append(retArray, b0, b1, trMapData[x][y])
				coord = pCoord
			}
		}
	}
	return retArray
}

func PrintMap() {
	fmt.Print("Map:")
	trBytes := createByteArr(treasureMap)
	str := base64.StdEncoding.EncodeToString(trBytes)
	fmt.Print(str)
}
