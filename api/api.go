package api

import (
	"Golden/models"
	"Golden/stat"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

/*BasicPath - базовый путь*/
var BasicPath string
var dig = models.Dig{}

var areaPool = sync.Pool{
	New: func() interface{} {
		return new(models.Area)
	},
}

var reportPool = sync.Pool{
	New: func() interface{} { return new(models.Report) },
}

var bbufPool = sync.Pool{
	New: func() interface{} {
		bts := make([]byte, 100)
		return bytes.NewBuffer(bts)
	},
}

/*GetBalance - получение баланса*/
func GetBalance() (*models.Balance, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := fmt.Sprintf("%s/balance", BasicPath)
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.Set("Content-Type", "application/json")
	err := fasthttp.Do(req, resp)
	if err == nil {
		if resp.StatusCode() == http.StatusOK {
			bytes := resp.Body()
			if err != nil {
				return nil, err
			}
			balance := models.Balance{}
			err = json.Unmarshal(bytes, &balance)
			if err != nil {
				return nil, err
			}
			return &balance, nil
		}
		_, err = getBtsError(resp.Body())
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}

/*Explore - разведка точки x,y*/
func Explore(x, y, sizeX, sizeY int64) (*models.Amount, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	defer fasthttp.ReleaseResponse(resp)
	url := BasicPath + "/explore"
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	// area := areaPool.Get().(*models.Area)
	// defer areaPool.Put(area)
	// area.PosX = &x
	// area.PosY = &y
	// area.SizeX = sizeX
	// area.SizeY = sizeY

	buff := bbufPool.Get().(*bytes.Buffer)
	defer bbufPool.Put(buff)
	buff.Reset()
	n, err := fmt.Fprintf(buff, `{"posX":%d,"posY":%d,"sizeX":%d,"sizeY":%d}`, x, y, sizeX, sizeY)
	if err != nil {
		return nil, err
	}

	req.SetBody(buff.Bytes()[:n])
	err = fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	stat.NewReq(stat.Exp)
	if resp.StatusCode() == http.StatusOK {
		report := reportPool.Get().(*models.Report)
		defer reportPool.Put(report)
		err = report.UnmarshalJSON(resp.Body())

		if err != nil {
			return nil, err
		}
		return report.Amount, nil

	}
	me, err := getBtsError(resp.Body())
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Not 200:%d %s", *me.Code, *me.Message)

}

/*PostLicense - запрос лицензии*/
func PostLicense(wallet models.Wallet) (*models.License, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)

	defer fasthttp.ReleaseResponse(resp)
	url := BasicPath + "/licenses"
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	bbuf := bbufPool.Get().(*bytes.Buffer)
	defer bbufPool.Put(bbuf)
	bbuf.Reset()
	bbuf.WriteString("[")
	for i := 0; i < len(wallet); i++ {
		coin := wallet[i]
		bbuf.WriteString(fmt.Sprintf("%d", coin))
		if i < len(wallet)-1 {
			bbuf.WriteString(",")
		}
	}
	bbuf.WriteString("]")
	req.SetBody(bbuf.Bytes())
	err := fasthttp.Do(req, resp)

	if err != nil {
		return nil, err
	}
	stat.NewReq(stat.Lic)
	if resp.StatusCode() == http.StatusOK {
		license := new(models.License)
		err := license.UnmarshalBinary(resp.Body())
		if err != nil {
			return nil, err
		}
		return license, err
	}
	if resp.StatusCode() != 502 {
		_, err = getBtsError(resp.Body())
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode())
}

/*DigPost -  копать*/
func DigPost(depth int64, licID int64, posX int64, posY int64) (models.TreasureList, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := BasicPath + "/dig"
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	//dig := models.Dig{Depth: &depth, LicenseID: &licID, PosX: &posX, PosY: &posY}
	dig.Depth = &depth
	dig.LicenseID = &licID
	dig.PosX = &posX
	dig.PosY = &posY
	bts, err := dig.MarshalJSON()
	if err != nil {
		return nil, err
	}
	req.SetBody(bts)
	err = fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	stat.NewReq(stat.Digg)
	if resp.StatusCode() == http.StatusOK {
		return getTreasureList(resp.Body())
	}
	if resp.StatusCode() == 404 {
		return nil, nil
	}
	_, err = getBtsError(resp.Body())
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode())
}

/*PostCash - посылка cash*/
func PostCash(treasure models.Treasure) (*models.Wallet, error) {
	url := BasicPath + "/cash"
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	bts, err := json.Marshal(treasure)
	if err != nil {
		return nil, err
	}
	req.SetBody(bts)
	err = fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	stat.NewReq(stat.Cash)
	if resp.StatusCode() == http.StatusOK {
		return getWallet(resp.Body())
	}
	_, err = getBtsError(resp.Body())
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode())
}

/*GetLicenses - получение списка лицензий*/
func GetLicenses() (models.LicenseList, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := BasicPath + "/licenses"
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.Set("Content-Type", "application/json")
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		bts := resp.Body()
		if err != nil {
			return nil, err
		}
		licList := models.LicenseList{}
		err = json.Unmarshal(bts, &licList)
		if err != nil {
			return nil, err
		}
		return licList, nil
	}
	_, err = getBtsError(resp.Body())
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode())
}

/*GetBasicPath - получение базового пути*/
func GetBasicPath() {
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "localhost"
	}
	BasicPath = fmt.Sprintf("http://%s:8000", address)
	//fmt.Printf("basic path: %s\n", BasicPath)

}

func getError(rc io.ReadCloser) (*models.Error, error) {
	bts, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	error := models.Error{}
	json.Unmarshal(bts, &error)
	//fmt.Println(*error.Message, *error.Code)
	return &error, nil

}

func getBtsError(bts []byte) (*models.Error, error) {

	error := models.Error{}
	json.Unmarshal(bts, &error)
	//fmt.Println(*error.Message, *error.Code)
	return &error, nil

}

func getTreasureList(bts []byte) (models.TreasureList, error) {
	trList := models.TreasureList{}
	var p fastjson.Parser
	val, err := p.ParseBytes(bts)
	if err != nil {
		return nil, err
	}
	arr, err := val.Array()
	if err != nil {
		return nil, err
	}
	for _, treasure := range arr {
		trList = append(trList, models.Treasure(treasure.GetStringBytes()))
	}
	return trList, nil
}

func getWallet(bts []byte) (*models.Wallet, error) {
	wallet := models.Wallet{}
	var p fastjson.Parser
	val, err := p.ParseBytes(bts)
	if err != nil {
		return nil, err
	}
	arr, err := val.Array()
	if err != nil {
		return nil, err
	}
	for _, coin := range arr {
		wallet = append(wallet, uint32(coin.GetInt()))
	}
	return &wallet, nil
}
