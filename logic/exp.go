package logic

import (
	"Golden/api"
	"Golden/stat"
	"fmt"
	"time"
)

func divideSegment(x, y, size int64, ch chan DigData) int {
	sum := 0
	amount, err := api.Explore(x, y, size, size)
	if err != nil {
		//fmt.Println("Exp err:", err)
		stat.NewStatErr(stat.Exp)
		return 0
	}

	if *amount != 0 {
		if size >= 4 {
			money := int(*amount)
			tsum := 0
		m1:
			for x1 := x; x1 < x+size; x1 += size / 2 {
				for y1 := y; y1 < y+size; y1 += size / 2 {
					am := divideSegment(x1, y1, size/2, ch)
					sum += am
					tsum += am
					if money == tsum {
						break m1
					}
				}
			}
			if money != tsum {
				fmt.Printf("------t: %d fact: %d ", money, tsum)
			}
		} else {
			money := int(*amount)
			sum += exploreArea(int(x), int(y), int(x+size), int(y+size), ch, money)
		}

	}
	return sum
}

func searchSegments(x0, y0, xe, ye, size, limit int, ch chan DigData) {
	for x := x0; x < xe; x += size {
		for y := y0; y < ye; y += size {
			amount, err := api.Explore(int64(x), int64(y), int64(size), int64(size))
			if err != nil {
				//fmt.Println("Exp err:", err)
				stat.NewStatErr(stat.Exp)
				return
			}
			if int(*amount) >= limit {
				if !expChainFull {
					go divideSegment(int64(x), int64(y), int64(size), ch)
				} else {
					divideSegment(int64(x), int64(y), int64(size), ch)
				}

			}
		}
	}
}

/*exploreArea - исследование области*/
func exploreArea(xbg, ybg, xend, yend int, ch chan DigData, targetMoney int) int {
	sum := 0
m1:
	for x := xbg; x < xend; x++ {
		for y := ybg; y < yend; y++ {
			amount, err := api.Explore(int64(x), int64(y), 1, 1)
			if err != nil {
				stat.NewStatErr(stat.Exp)
				//fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					digData := DigData{x: int64(x), y: int64(y), amount: int64(*amount)}

					stat.NewArStat(int(*amount))
					ch <- digData
					sum += int(*amount)
					if targetMoney == sum {
						break m1
					}

				}
			}
		}
	}
	if sum != targetMoney {
		stat.NewExpAreaErr()
		//fmt.Printf("Exp error t:%d s:%d\n", targetMoney, sum)
	}
	return targetMoney
}

/*exploreSegment - исследование сегмента*/
func exploreSegment(xbg, ybg, xend, yend, size int, ch chan DigData) int {
	sum := 0
	for x := xbg; x < xend; x += size {
		for y := ybg; y < yend; y += size {
			amount, err := api.Explore(int64(x), int64(y), int64(size), int64(size))
			if err != nil {
				stat.NewStatErr(stat.Exp)
				//fmt.Println("Exp err:", err)
			} else {
				if *amount != 0 {
					if size >= 4 {
						money := int(*amount)
						tsum := 0
					m1:
						for x1 := x; x1 < x+size; x1 += size / 2 {
							for y1 := y; y1 < y+size; y1 += size / 2 {
								am := exploreSegment(x1, y1, x1+size/2, y1+size/2, size/2, ch)
								sum += am
								tsum += am
								if money == tsum {
									break m1
								}
							}
						}
						if money != tsum {
							fmt.Printf("t: %d fact: %d ", money, tsum)
						}
					} else {
						money := int(*amount)
						sum += exploreArea(x, y, x+size, y+size, ch, money)
					}

				}
			}
		}
	}
	return sum

}

func research(start, end, step int) {
	for x := start; x < end; x += step {
		tbg := time.Now()
		for y := start; y < end; y += step {
			amount, err := api.Explore(int64(x), int64(y), int64(step), int64(step))
			if err == nil {
				fmt.Printf("%2d ", *amount)
			}
		}
		fmt.Printf(" ms%d\n", int(time.Since(tbg).Milliseconds())/((end-start)/step))
	}
}

/*segment - сегмент*/
type segment struct {
	x, y, xSize, ySize int
	amount             int
}

/*explore - исследование сегмента*/
func (seg *segment) explore() error {
	amount, err := api.Explore(int64(seg.x), int64(seg.y), int64(seg.xSize), int64(seg.ySize))
	if err != nil {
		stat.NewStatErr(stat.Exp)
		return err
	}
	seg.amount = int(*amount)
	return nil
}

/*divide2 - разделение на два сегмента*/
func (seg *segment) divide2() (seg1 *segment, seg2 *segment) {

	if seg.xSize > seg.ySize {
		seg1 = newSegment(seg.x, seg.y, seg.xSize/2, seg.ySize)
		err := seg1.explore()
		for err != nil {
			err = seg1.explore()
		}
		seg2 = newSegment(seg.x+seg.xSize/2, seg.y, seg.xSize/2, seg.ySize)
		seg2.amount = seg.amount - seg1.amount
		return
	}
	seg1 = newSegment(seg.x, seg.y, seg.xSize, seg.ySize/2)
	err := seg1.explore()
	for err != nil {
		err = seg1.explore()
	}
	seg2 = newSegment(seg.x, seg.y+seg.ySize/2, seg.xSize, seg.ySize/2)
	seg2.amount = seg.amount - seg1.amount
	return
}

/*newSegment - новый сегмент*/
func newSegment(x, y, xSize, ySize int) *segment {
	seg := new(segment)
	seg.x = x
	seg.y = y
	seg.xSize = xSize
	seg.ySize = ySize
	return seg
}

func (seg *segment) setSegment(x, y, xSize, ySize int) {
	seg.x = x
	seg.y = y
	seg.xSize = xSize
	seg.ySize = ySize
}

/*searchArea - исследование заданной области*/
func searchArea(xStart, yStart, sizeX, sizeY, step int, ch chan DigData, limit int) {
	for x := xStart; x+step < xStart+sizeX; x += step {
		for y := yStart; y+step < yStart+sizeY; y += step {
			seg := newSegment(x, y, step, step)
			err := seg.explore()
			if err != nil {
				fmt.Println(err)
			}
			if seg.amount >= limit {
				//if !expChainFull {
				go explore(seg, ch)
				//	} else {
				//		explore(seg, ch)
				//	}
			}
		}
	}

}

/*explore - рекурсивная функция исследования*/
func explore(seg *segment, ch chan DigData) {
	if seg.xSize == 1 && seg.ySize == 1 {
		stat.NewArStat(seg.amount)
		ch <- DigData{x: int64(seg.x), y: int64(seg.y), amount: int64(seg.amount)}
		return
	}
	seg1, seg2 := seg.divide2()
	if seg1.amount > 0 {

		explore(seg1, ch)

	}
	if seg2.amount > 0 {
		explore(seg2, ch)
	}
}
