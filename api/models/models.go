package models

type INNInfo struct {
	Fam     string `json:"fam"`
	Nam     string `json:"nam"`
	Otch    string `json:"otch"`
	Bdate   string `json:"bdate"`
	Doctype string `json:"doctype"`
	Docno   string `json:"docno"`
	C       string `json:"c"`
}

type INNReady struct {
	Passport string `json:"passport"`
	Inn      string `json:"inn"`
}
