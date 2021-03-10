package stat

var digTreasures, digTlist, sendTlist int
var expAreaErr int

type areaStat struct {
	amounts int
}

func (as areaStat) StatName() string {
	return "AS"
}

func (as areaStat) GetAmount() int {
	return as.amounts
}

/*NewArStat - новая статистика для исследованной области*/
func NewArStat(ams int) {
	statChan <- areaStat{amounts: ams}
}

type licStat struct {
	free bool
}

func (ls licStat) StatName() string {
	return "LS"
}

func (ls licStat) IsFree() bool {
	return ls.free
}

/*NewLcStat - статистика лицензий*/
func NewLcStat(free bool) {
	statChan <- licStat{free: free}
}

type digStat struct {
	amounts int
	digged  int
}

func (ds digStat) StatName() string {
	return "DS"
}

func (ds digStat) GetAmounts() int {
	return ds.amounts
}
func (ds digStat) GetDigg() int {
	return ds.digged
}

/*NewDsStat - статистика копаний*/
func NewDsStat(digged int, amounts int) {
	statChan <- digStat{amounts: amounts, digged: digged}
}

type coinStat struct {
	coins int
}

func (cs coinStat) StatName() string {
	return "CS"
}

func (cs coinStat) GetCoins() int {
	return cs.coins
}

/*NewCoinStat - статистика монет*/
func NewCoinStat(coins int) {
	statChan <- coinStat{coins: coins}
}

type error struct {
	_type ReqType
}

func (err error) StatName() string {
	return "Err"
}

func (err error) Type() ReqType {
	return err._type
}

/*NewStatErr - новая ошибка*/
func NewStatErr(_type ReqType) {
	statChan <- error{_type: _type}
}

/*NewSendTreas - новые откопанные сокровища*/
func NewSendTreas(treasures int) {
	digTreasures += treasures
}

/*NewDigTlist - новый откопанный список сокровищ*/
func NewDigTlist() {
	digTlist++
}

/*NewSendTlist - новый отосланный список сокровищ*/
func NewSendTlist() {
	sendTlist++
}

/*NewExpAreaErr - непонятная ошибка с подсчетом в функции exploreArea*/
func NewExpAreaErr() {
	expAreaErr++
}

/*NewReq - новый запрос*/
func NewReq(tp ReqType) {
	numReq++
	switch tp {
	case Exp:
		numExpReq++
	case Digg:
		numDigReq++
	case Cash:
		numCashReq++
	case Lic:
		numLicReq++

	}
}
