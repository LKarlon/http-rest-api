package service

import "github.com/LKarlon/http-rest-api.git/api/models"

type innWorker interface{
	GetInn (passports []string)(models.INNReady, error)
}


type service struct{
	worker innWorker
}


func (s *service) GetInn(passports []string)(data models.INNReady, err error){
	s.worker.GetInn(passports)
}

func NewService() service{
	return
}