package stat

var licStatMap map[int]int
var diffPrice int
var numLic = 0
var sumLicTime int

/*NewLicStat - статистика лицензий*/
func NewLicStat(money, numDigg int) {
	if licStatMap[money] != 0 && licStatMap[money] != numDigg {
		diffPrice++
	}
	licStatMap[money] = numDigg
}

/*NewLicTime - замер среднего времени получения лицензии*/
func NewLicTime(dt int) {
	numLic++
	sumLicTime += dt
}
