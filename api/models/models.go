package models

type INNInfo struct {
	Fam     string `json:"fam"`
	Nam     string `json:"nam"`
	Otch    string `json:"otch"`
	Bdate   string `json:"bdate"`
	Docno   string `json:"docno"`
}

type INNReady struct {
	Passport string `json:"passport"`
	Inn      string `json:"inn"`
}

