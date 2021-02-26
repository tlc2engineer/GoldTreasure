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
)

/*BasicPath - базовый путь*/
var BasicPath string

/*GetBalance - получение баланса*/
func GetBalance() (*models.Balance, error) {
	req := fmt.Sprintf("%s/balance", BasicPath)
	fmt.Println("balance", req)
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
func Explore(x, y int64) (*models.Amount, error) {
	req := BasicPath + "/explore"
	fmt.Println("explore", req)
	area := models.Area{
		PosX:  &x,
		PosY:  &y,
		SizeX: 1,
		SizeY: 1,
	}
	bts, err := json.Marshal(area)
	if err != nil {
		return nil, err
	}
	responseBody := bytes.NewBuffer(bts)
	resp, err := http.Post(req, "application/json", responseBody)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		report := models.Report{}
		err = json.Unmarshal(bts, &report)
		if err != nil {
			return nil, err
		}
		return report.Amount, nil

	}
	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, err

}

/*PostLicense - запрос лицензии*/
func PostLicense() (*models.License, error) {
	req := BasicPath + "/licenses"
	fmt.Println("license post", req)
	wallet := models.Wallet{}
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
	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
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

	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
}

/*PostCash - посылка cash*/
func PostCash(treasure models.Treasure) (*models.Wallet, error) {
	req := BasicPath + "/cash"

	bts, err := json.Marshal(treasure)
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
		wallet := models.Wallet{}
		err = json.Unmarshal(bts, &wallet)
		if err != nil {
			return nil, err
		}
		return &wallet, nil
	}
	_, err = getError(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
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
	fmt.Println(error.Message, error.Code)
	return &error, nil

}
