package api

import (
	"Golden/models"
	"bytes"
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf(fmt.Sprintf("STATUS NOT 200:%s", err))
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
	return nil, err

}

/*PostLicense - запрос лицензии*/
func PostLicense() (*models.License, error) {
	req := BasicPath + "/licenses"
	fmt.Println("license post", req)
	wallet := models.Wallet{}
	issue := IssueLicense{Args: wallet}
	bts, err := json.Marshal(issue)
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
	return nil, fmt.Errorf("Status not ok:%d", resp.StatusCode)
}

type IssueLicense struct {
	Args models.Wallet `json:"args"`
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
