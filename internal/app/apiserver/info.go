package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/LKarlon/http-rest-api.git/api/models"
	"log"
	"net/http"
	"net/url"
)

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
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(result)
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

