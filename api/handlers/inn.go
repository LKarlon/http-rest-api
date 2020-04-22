package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/LKarlon/http-rest-api.git/api/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func INN (w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var res []models.INNInfo
	var inn models.INNReady
	var resp []models.INNReady

	err = json.Unmarshal(reqBody, &res)
	if err != nil{
		log.Fatalln(err)
	}
	for _, value := range res {
		inn, err = NewInfo(value.Fam, value.Nam, value.Otch, value.Bdate,
			value.Doctype, value.Docno)
		if err == nil {
			resp = append(resp, inn)
		}
		time.Sleep(time.Second)
	}
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil{
		log.Fatalln(err)
	}
}

func NewInfo(surname, name, patronymic, birthdate, doctype, docnumber string) (models.INNReady, error){
	urls := "https://service.nalog.ru/inn-proc.do"
	data := url.Values{
		"fam": {surname},
		"nam": {name},
		"otch": {patronymic},
		"bdate": {birthdate},
		"bplace": {""},
		"doctype": {doctype},
		"docno": {docnumber},
		"c": {"innMy"},
		"captcha": {""},
		"captchaToken": {""},
	}

	resp, err := http.PostForm(urls, data)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	var inn models.INNReady
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil{
		log.Fatalln(err)
	}
	log.Println(result)
	if _, ok := result["inn"]; ok == false{
		return inn, fmt.Errorf("need capch")
	}
	if result["code"].(float64) == 0{
		log.Println("not valid data with passport %s", docnumber)
		return inn, fmt.Errorf("not valid data")
	}
	inn.Passport = docnumber
	inn.Inn = result["inn"].(string)
	return inn, nil
}
