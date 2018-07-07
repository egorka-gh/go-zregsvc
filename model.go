package main

//ClientState clent states & frontend messages
type ClientState struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	WebComment string `json:"web_comment" db:"web_comment"`
}

//ValidateResult response 4 frontend
type ValidateResult struct {
	Program int    `json:"program" db:"program"`
	Card    string `json:"card" db:"card"`
	State   int    `json:"state" db:"state"`
	ErrCode int    `json:"err" db:"err"`
	Message string `json:"message" db:"message"`
}

//CardProgram discount programm
type CardProgram struct {
	ID          int    `json:"id" db:"id"`
	Program     int    `json:"program" db:"program"`
	CardStart   string `json:"cardStart" db:"cardStart"`
	CardEnd     string `json:"cardEnd" db:"cardEnd"`
	CardLen     int    `json:"cardLen" db:"cardLen"`
	Active      bool   `json:"active" db:"active"`
	CheckIssued bool   `json:"checkIssued" db:"checkIssued"`
}

//Client client model
type Client struct {
	Card       string `json:"card" db:"card"`
	State      int    `json:"state" db:"state"`
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	PhoneCode  string `json:"phoneCode" db:"phoneCode"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Gender     int    `json:"gender" db:"gender"`
	Birthday   string `json:"birthday" db:"birthday"`
	Pet        string `json:"pet" db:"pet"`
	SendPromo  bool   `json:"sendPromo" db:"sendPromo"`
}

func loadStates() ([]ClientState, error) {
	var list []ClientState
	var ssql = "SELECT cs.id, ifnull(cs.name,'') name, ifnull(csm.web_comment,'') web_comment FROM client_state cs LEFT OUTER JOIN client_state_msg csm ON cs.id = csm.id WHERE cs.id!=0 ORDER BY cs.id"
	err := db.Select(&list, ssql)
	return list, err
}
