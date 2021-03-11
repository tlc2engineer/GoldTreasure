package stat

var licStatMap map[int]int
var diffPrice int

/*NewLicStat - статистика лицензий*/
func NewLicStat(money, numDigg int) {
	if licStatMap[money] != 0 && licStatMap[money] != numDigg {
		diffPrice++
	}
	licStatMap[money] = numDigg
}
