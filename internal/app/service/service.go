package service

import (
	"encoding/json"
	"fmt"
	"github.com/LKarlon/http-rest-api.git/api/models"
	"github.com/LKarlon/http-rest-api.git/internal/app/store"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Service interface {
	GetInn(w http.ResponseWriter, r *http.Request)
	Worker()
}

type service struct {
	readyData         map[string]string
	DataForProcessing []models.INNInfo
	store             *store.Store
}

func (s *service) configureStore(config *store.Config) error {
	st := store.New(config)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func NewService(config *store.Config) Service {
	s := &service{
		readyData:         map[string]string{},
		DataForProcessing: []models.INNInfo{},
	}
	if err := s.configureStore(config); err != nil {
		log.Println(err)
	}
	return s
}

func (s *service) Worker() {
	for {
		urls := "https://service.nalog.ru/inn-proc.do"
		if len(s.DataForProcessing) == 0 {  // Если данных для обработки нет,
			time.Sleep(time.Second) 		// ждем секунду и запускаем цикл заново
			continue
		}
		for _, value := range s.DataForProcessing {
			data := url.Values{
				"fam":          {value.Fam},
				"nam":          {value.Nam},
				"otch":         {value.Otch},
				"bdate":        {value.Bdate},
				"bplace":       {""},
				"doctype":      {"21"},
				"docno":        {value.Docno},
				"c":            {"innMy"},
				"captcha":      {""},
				"captchaToken": {""},
			}
			resp, err := http.PostForm(urls, data)
			if err != nil {
				log.Fatalln(err)
			}
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(result)
			if _, ok := result["ERROR"]; ok == true {
				fmt.Println("need captcha. waite 20 sec")
				time.Sleep(time.Second * 20)
				fmt.Println("after waiting")
				resp, err = http.PostForm(urls, data)
				if err != nil {
					log.Fatalln(err)
				}
				err = json.NewDecoder(resp.Body).Decode(&result)
				if _, ok := result["inn"]; ok == false {
					log.Fatalln("need captcha")
				}
			}
			if result["code"].(float64) == 0 {
				fmt.Println("not valid data with passport %s", value.Docno)
				continue
			}
			s.readyData[value.Docno] = result["inn"].(string)
			s.store.Inn().Add(value.Docno, result["inn"].(string))
			time.Sleep(time.Second)
		}
		s.DataForProcessing = s.DataForProcessing[:0]
	}
}

func (s *service) GetInn(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var req []models.INNInfo
	//var inn models.INNReady
	var resp []models.INNReady

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, value := range req {
		if val, err := s.store.Inn().FindInn(value.Docno); err == nil{
			resp = append(resp, val)
		}
		/*if val, ok := s.readyData[value.Docno]; ok == true {
			inn.Passport = value.Docno
			inn.Inn = val
			resp = append(resp, inn)
		}*/
		if _, ok := s.readyData[value.Docno]; ok == false { //если готовых данных нет,
			s.DataForProcessing = append(s.DataForProcessing, value) //добавляем в данные для обработки
		}
	}
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Fatalln(err)
	}
}
