package stat

type level struct {
	depth          int
	totalTreasures int
	totalTime      int
	total          int
}

var levels map[int]*level

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
}

/*DepthStat - статистика по глубине и количеству сокровищ*/
func DepthStat(t int, depth int, treasures int) {
	level := levels[depth]
	level.total++
	level.totalTime += t
	level.totalTreasures += treasures

}
