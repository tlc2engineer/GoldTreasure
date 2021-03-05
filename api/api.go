package api

import (
	"Golden/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/valyala/fasthttp"
)

/*BasicPath - базовый путь*/
var BasicPath string

/*GetBalance - получение баланса*/
func GetBalance() (*models.Balance, error) {
	req := fmt.Sprintf("%s/balance", BasicPath)
	resp, err := http.Get(req)
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			bytes, err := ioutil.ReadAll(resp.Body)
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
		_, err = getError(resp.Body)
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
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release
	url := BasicPath + "/explore"
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	area := models.Area{
		PosX:  &x,
		PosY:  &y,
		SizeX: sizeX,
		SizeY: sizeY,
	}

	bts, err := json.Marshal(area)
	if err != nil {
		return nil, err
	}
	req.SetBody(bts)
	err = fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		report := models.Report{}
		err = json.Unmarshal(resp.Body(), &report)
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
	req := BasicPath + "/licenses"
	bts, err := json.Marshal(wallet)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(bts)
	resp, err := http.Post(req, "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		license := models.License{}
		json.Unmarshal(bts, &license)
		return &license, err
	}
	if resp.StatusCode != 502 {
		_, err = getError(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
}

/*DigPost -  копать*/
func DigPost(depth int64, licID int64, posX int64, posY int64) (models.TreasureList, error) {
	req := BasicPath + "/dig"
	dig := models.Dig{Depth: &depth, LicenseID: &licID, PosX: &posX, PosY: &posY}
	bts, err := json.Marshal(dig)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(bts)
	resp, err := http.Post(req, "application/json", buff)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		treasures := models.TreasureList{}
		err = json.Unmarshal(bts, &treasures)
		if err != nil {
			return nil, err
		}
		return treasures, nil
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
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
	req := BasicPath + "/licenses"
	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
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
	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
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
