package stat

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
	_type ErrType
}

func (err error) StatName() string {
	return "Err"
}

func (err error) Type() ErrType {
	return err._type
}

/*NewStatErr - новая ошибка*/
func NewStatErr(_type ErrType) {
	statChan <- error{_type: _type}
}
