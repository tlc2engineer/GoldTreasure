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

var areaPool = sync.Pool{
	New: func() interface{} {
		return new(models.Area)
	},
}
var report = models.Report{}
var repMut = new(sync.Mutex)
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
		repMut.Lock()
		defer repMut.Unlock()
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
	buff := bbufPool.Get().(*bytes.Buffer)
	defer bbufPool.Put(buff)
	buff.Reset()
	buff.WriteString("[")
	for i, coin := range wallet {
		buff.WriteString(fmt.Sprintf("%d", coin))
		if i != len(wallet)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString("]")
	req.SetBody(buff.Bytes())
	err := fasthttp.Do(req, resp)

	if err != nil {
		return nil, err
	}
	stat.NewReq(stat.Lic)
	if resp.StatusCode() == http.StatusOK {
		bts := resp.Body()
		if err != nil {
			return nil, err
		}
		license := models.License{}
		license.UnmarshalBinary(bts)
		return &license, err
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
func DigPost(depth int64, licID int64, posX int64, posY int64) (*models.TreasureList, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := BasicPath + "/dig"
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	dig := models.Dig{Depth: &depth, LicenseID: &licID, PosX: &posX, PosY: &posY}
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
		bts := resp.Body()
		if err != nil {
			return nil, err
		}
		treasures := models.TreasureList{}
		err = json.Unmarshal(bts, &treasures)
		if err != nil {
			return nil, err
		}
		return &treasures, nil
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
		wallet := models.Wallet{}
		err = json.Unmarshal(resp.Body(), &wallet)
		if err != nil {
			return nil, err
		}
		return &wallet, nil
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
	fmt.Printf("basic path: %s\n", BasicPath)

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

func getTreasuresList(bts []byte) (*models.TreasureList, error) {
	var p fastjson.Parser
	val, err := p.ParseBytes(bts)
	if err != nil {
		return nil, err
	}
	tlist, err := val.Array()
	if err != nil {
		return nil, err
	}
	list := models.TreasureList{}
	for _, treas := range tlist {
		list = append(list, models.Treasure(treas.String()))
	}
	return &list, nil

}

func getWallet(bts []byte) (*models.Wallet, error) {
	var p fastjson.Parser
	val, err := p.ParseBytes(bts)
	if err != nil {
		return nil, err
	}
	wallList, err := val.Array()
	if err != nil {
		return nil, err
	}
	wallet := models.Wallet{}
	for _, v := range wallList {
		wallet = append(wallet, uint32(v.GetInt()))
	}
	return &wallet, nil
}
