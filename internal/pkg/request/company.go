package request

type Company struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
	Alias   string `json:"alias"`
	Giro    string `json:"giro"`
}
