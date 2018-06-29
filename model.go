package main

//ClientState clent states & frontend messages
type ClientState struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	WebComment string `json:"сomment" db:"web_comment"`
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

func loadStates() ([]ClientState, error) {
	var list []ClientState
	var ssql = "SELECT cs.id, ifnull(cs.name,'') name, ifnull(csm.web_comment,'') web_comment FROM client_state cs LEFT OUTER JOIN client_state_msg csm ON cs.id = csm.id WHERE cs.id!=0 ORDER BY cs.id"
	err := db.Select(&list, ssql)
	return list, err
}
