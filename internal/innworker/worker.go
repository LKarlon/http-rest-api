package innworker

import "github.com/LKarlon/http-rest-api.git/api/models"

type storage interface{
	get()
	put()
}


type worker struct{
	readyData	map[string]string
	DataForProcessing []models.INNInfo
	storage storage
}

func (s *worker) GetInn(passports []string) (data models.INNReady){

}
func (s *worker) UpdateMap(){

}