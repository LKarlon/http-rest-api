package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/LKarlon/http-rest-api.git/api/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type APIServer struct {
	config            *Config
	logger            *logrus.Logger
	router            *mux.Router
	readyData         map[string]string
	DataForProcessing []models.INNInfo
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:            config,
		logger:            logrus.New(),
		router:            mux.NewRouter(),
		readyData:         map[string]string{},
		DataForProcessing: []models.INNInfo{},
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configRouter()
	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configRouter() {
	s.router.HandleFunc("/inn", s.INN)
}

// Worker работает с ReadyData и DataForProcessing
func (s *APIServer) Worker() {
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
			if _, ok := result["ERROR"]; ok == true {			//Если требуется капча,
				fmt.Println("need captcha. waite 20 sec") 		//ждем 20 секунд и пробуем снова
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
			s.readyData[value.Docno] = result["inn"].(string)   //Добавляем данные в базу готовых паспортов
			time.Sleep(time.Second)
		}
		s.DataForProcessing = s.DataForProcessing[:0]			//Очищаем массив с данными для обработки
	}
}

// INN берет данные из DataForProcessing и отправляет ответ клиенту
func (s *APIServer) INN(w http.ResponseWriter, r *http.Request) {
	fmt.Print(s.DataForProcessing)
	fmt.Println(" - INN START")
	fmt.Println(s.readyData)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var req []models.INNInfo
	var inn models.INNReady
	var resp []models.INNReady

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, value := range req {
		if val, ok := s.readyData[value.Docno]; ok == true {
			inn.Passport = value.Docno
			inn.Inn = val
			resp = append(resp, inn)
		}
		if _, ok := s.readyData[value.Docno]; ok == false {				//если готовых данных нет,
			s.DataForProcessing = append(s.DataForProcessing, value) 	//добавляем в данные для обработки
		}
	}
	fmt.Print(s.DataForProcessing)
	fmt.Println(" - INN END")

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Fatalln(err)
	}
}
