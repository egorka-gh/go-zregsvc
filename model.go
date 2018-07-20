package main

import (
	"database/sql"
	"strings"
)

const (
	// err states
	errDATABASE      int = -1 // more precisely app level error
	errWrongCAPTCHA  int = -5
	errWrongCardCODE int = -10
	errWrongSTATE    int = -11
	errCardNotISSUED int = -12

	// web info states
	webSTART        int = -1000
	webCONFIRMATION int = -1001

	// client states
	clNotREGISTERED   int = 1
	clREGISTRATION    int = 5
	clWaiteREFINEMENT int = 10
	clREGISTERED      int = 100
)

//ClientState clent states & frontend messages
type ClientState struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	WebComment string `json:"web_comment" db:"web_comment"`
}

//ValidateResult response 4 frontend
type ValidateResult struct {
	Program         int    `json:"program" db:"program"`
	Card            string `json:"card" db:"card"`
	State           int    `json:"state" db:"state"`
	ErrCode         int    `json:"err" db:"err"`
	Message         string `json:"message" db:"message"`
	Captcha         string `json:"captcha"`
	CaptchaSolution string `json:"captchaSolution"`
	CaptchaState    int    `json:"captchaState"`
}

//CardProgram discount programm
type CardProgram struct {
	ID          int    `json:"id" db:"id"`
	Program     int    `json:"program" db:"program"`
	CardStart   string `json:"cardStart" db:"card_start"`
	CardEnd     string `json:"cardEnd" db:"card_end"`
	CardLen     int    `json:"cardLen" db:"card_len"`
	Active      bool   `json:"active" db:"active"`
	CheckIssued bool   `json:"checkIssued" db:"check_issued"`
}

//Client client model
type Client struct {
	Program    int    `json:"program" db:"program"`
	Card       string `json:"card" db:"card"`
	State      int    `json:"state" db:"state"`
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	PhoneCode  string `json:"phoneCode" db:"phone_code"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Gender     int    `json:"gender" db:"gender"`
	Birthday   string `json:"birthday" db:"birthday"`
	Pet        string `json:"pet" db:"pet"`
	SendPromo  bool   `json:"sendPromo" db:"send_promo"`
}

//RegisterDTO dto model
type RegisterDTO struct {
	Client Client
	Result ValidateResult
}

func loadStates() ([]ClientState, error) {
	var list []ClientState
	var ssql = "SELECT cs.id, ifnull(cs.name,'') name, ifnull(csm.web_comment,'') web_comment FROM client_state cs LEFT OUTER JOIN client_state_msg csm ON cs.id = csm.id WHERE cs.id!=0 ORDER BY cs.id"
	err := db.Select(&list, ssql)
	return list, err
}

func validateCard(result ValidateResult) ValidateResult {
	result.State = 0
	result.ErrCode = 0

	// check if card empty
	if len(strings.TrimSpace(result.Card)) == 0 {
		result.ErrCode = errWrongCardCODE
		result.Message = "Указана не верная карта"
		return result
	}

	// look for programm by card && card len
	var ssql = "SELECT pc.*" +
		" FROM programs p" +
		" INNER JOIN program_cards pc ON p.id = pc.program AND pc.active!=0" +
		" WHERE LENGTH(?) = pc.card_len AND ? BETWEEN pc.card_start AND pc.card_end AND p.external=0 AND p.active !=0"
	var prg CardProgram
	err := db.Get(&prg, ssql, result.Card, result.Card)
	if err != nil {
		if err == sql.ErrNoRows {
			//range not found
			result.ErrCode = errWrongCardCODE
			result.Message = "Указана не верная карта"
			return result
		}
		//db error
		result.ErrCode = errDATABASE
		result.Message = err.Error()
		return result
	}
	result.Program = prg.Program

	//check if card exists
	var client Client
	ssql = "SELECT program, card, state, IFNULL(surname, '') surname, IFNULL(name, '') name, IFNULL(patronymic, '') patronymic," +
		" IFNULL(phone_code, '') phone_code, IFNULL(phone, '') phone, IFNULL(email, '') email, gender," +
		" IFNULL(birthday, '') birthday, IFNULL(pet, '') pet, send_promo" +
		" FROM clients c WHERE c.program = ? AND c.card = ?"
	err = db.Get(&client, ssql, result.Program, result.Card)
	if err == nil {
		//card exists
		result.State = client.State
		if result.State >= clREGISTRATION {
			//card registered
			result.ErrCode = errWrongSTATE
			result.Message = "Карта уже зарегистрирована"
		}
	} else {
		if err == sql.ErrNoRows {
			//no card found
			if prg.CheckIssued {
				//card has to be issued
				result.ErrCode = errCardNotISSUED
				result.Message = "Карта не выдана"
			}
		} else {
			//db error
			result.ErrCode = errDATABASE
			result.Message = err.Error()
		}
	}

	if result.ErrCode != 0 {
		result.Program = 0
	}
	return result
}

func registerCard(dto *RegisterDTO) ValidateResult {
	res := dto.Result
	cli := dto.Client

	//check if card exists vs state < CL_REGISTRATION
	var client Client
	ssql := "SELECT program, card, state, IFNULL(surname, '') surname, IFNULL(name, '') name, IFNULL(patronymic, '') patronymic," +
		" IFNULL(phone_code, '') phone_code, IFNULL(phone, '') phone, IFNULL(email, '') email, gender," +
		" IFNULL(birthday, '') birthday, IFNULL(pet, '') pet, send_promo" +
		" FROM clients c WHERE c.program = ? AND c.card = ? AND c.state<5"
	err := db.Get(&client, ssql, res.Program, res.Card)

	if cli.Birthday == "" {
		cli.Birthday = "shouldSetNull"
	}
	promo := 0
	if cli.SendPromo {
		promo = 1
	}
	if err == nil {
		//card exists
		//update
		ssql =
			"UPDATE clients" +
				" SET" +
				" state = ?" +
				" ,surname = ?" +
				" ,name = ?" +
				" ,patronymic = ?" +
				" ,phone_code = ?" +
				" ,phone = ?" +
				" ,email = ?" +
				" ,gender = ?" +
				" ,birthday = STR_TO_DATE( ? , '%d.%m.%Y')" +
				" ,pet = ?" +
				" ,send_promo = ?" +
				" WHERE program = ? AND card = ?"
		_, err = db.Exec(ssql,
			clREGISTRATION,
			cli.Surname,
			cli.Name,
			cli.Patronymic,
			cli.PhoneCode,
			cli.Phone,
			cli.Email,
			cli.Gender,
			cli.Birthday,
			cli.Pet,
			promo,
			res.Program,
			res.Card)
	} else {
		if err == sql.ErrNoRows {
			//no card found
			//insert
			ssql =
				"INSERT IGNORE INTO clients" +
					" (program, card, state, surname, name, patronymic, phone_code, phone, email, gender, birthday, pet, send_promo)" +
					" VALUES" +
					" (?,?,?,?,?,?,?,?,?,?, STR_TO_DATE( ? , '%d.%m.%Y'),?,?)"
			_, err = db.Exec(ssql,
				res.Program,
				res.Card,
				clREGISTRATION,
				cli.Surname,
				cli.Name,
				cli.Patronymic,
				cli.PhoneCode,
				cli.Phone,
				cli.Email,
				cli.Gender,
				cli.Birthday,
				cli.Pet,
				promo)
		}
	}

	if err != nil {
		//db error
		res.ErrCode = errDATABASE
		res.Program = 0
		res.Message = err.Error()
	} else {
		//complited, use err state??
		//res.Program = 0
		//res.ErrCode = errWrongSTATE
		res.ErrCode = 0
		res.State = clREGISTRATION
	}

	return res
}
