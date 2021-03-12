package stat

import "sync"

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
}

/*DigDeepStat - статистика по глубине и количеству сокровищ*/
func DigDeepStat(t int, depth int, treasures int, licType int) {
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

}

type licDigStat struct {
	sumDt int64
	num   int
}
