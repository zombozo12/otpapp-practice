package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"otpapp-native/config"
	"otpapp-native/models"
	"otpapp-native/redis"
	"regexp"
	"strings"
)

type (
	reqData struct {
		Number string `json:"number"`
	}

	valData struct {
		Number string `json:"number"`
		Code   string `json:"code"`
	}
)

var sb strings.Builder

func PhoneRequest(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("OTP Request Error: %+v", err)
		return
	}

	var rd reqData
	if err := json.Unmarshal(req, &rd); err != nil {
		badRequest("Unmarshal Failed", w)
		log.Printf("OTP Data Error : %+v", err)
		return
	}

	if len(rd.Number) == 0 {
		badRequest("Phone number cannot be empty", w)
		log.Print("Empty request")
		return
	}

	// validate phone format
	val := IsPhoneValid(rd.Number)
	if val != true {
		badRequest("Phone number is invalid", w)
		log.Print("Invalid phone number")
		return
	}

	// string builder. somehow ("phone:" + rd.Number) didn't work
	sb.WriteString("phone:" + rd.Number)

	// get redis phone data
	red, err := redis.String("GET", sb.String())
	if err != nil {
		badRequest("Redis Error", w)
		log.Printf("Redis GET Error : %+v", err)
		return
	}

	var rdsPhone redis.RedisPhone

	log.Println(red != "")

	if red == "" {
		log.Printf("Redis GET key : %s | Value: %+v", "number", red)
		if err = json.Unmarshal([]byte(red), &rdsPhone); err != nil {
			badRequest("JSON Unmarshal Error", w)
			log.Printf("JSON Unmarshal Error : %+v", err)
			return
		}

		okResponse("You're already requested OTP", &rdsPhone, w)
		return
	}

	// get database instance
	tx, err := config.GetDBInstance().Beginx()
	if err != nil {
		badRequest("DB Error", w)
		log.Printf("DB Error : %+v", err)
		return
	}
	defer tx.Rollback()

	code := OTPGenerator(6)

	// insert phone db
	resPhone, err := models.InsertPhoneDB(rd.Number, code, tx)
	if err != nil {
		badRequest("Insert Phone DB Error", w)
		log.Printf("Insert Phone DB Error : %+v", err)
		return
	}

	if err = tx.Commit(); err != nil {
		badRequest("Commit Phone DB Error", w)
		log.Printf("Commit Phone DB Error : %+v", err)
		return
	}

	// convert struct
	rdsPhone.Number = resPhone.Number
	rdsPhone.Code = resPhone.Code
	rdsPhone.ExpiredAt = resPhone.ExpiredAt

	jsonPhone, err := json.Marshal(&rdsPhone)
	if err != nil {
		badRequest("JSON Phone Marshal Error", w)
		log.Printf("JSON Phone Marshal Error : %+v", err)
		return
	}

	res, err := redis.String("SET", sb.String(), string(jsonPhone), "EX", 60)
	if err != nil {
		badRequest("Redis SET Error", w)
		log.Printf("Redis SET Error : %+v", err)
		return
	}

	if res != "OK" {
		badRequest("Redis SET Failed", w)
		log.Printf("Redis SET Failed : %+v", res)
		return
	}

	okResponse("Here's your OTP code", &rdsPhone, w)
}

func PhoneValidate(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("OTP Request Error: %+v", err)
		return
	}

	var vd valData
	if err := json.Unmarshal(req, &vd); err != nil {
		badRequest("Unmarshal Failed", w)
		log.Printf("OTP Data Error : %+v", err)
		return
	}

	if len(vd.Number) == 0 {
		badRequest("Phone number cannot be empty", w)
		log.Print("Empty phone number request")
		return
	}

	if len(vd.Code) == 0 {
		badRequest("Invalid code length", w)
		log.Print("Invalid code length")
		return
	}

	// validate phone format
	val := IsPhoneValid(vd.Number)
	if val != true {
		badRequest("Phone number is invalid", w)
		log.Print("Invalid phone number")
		return
	}

	// string builder. somehow ("phone:" + rd.Number) didn't work
	sb.WriteString("phone:" + vd.Number)

	// get redis phone data
	red, err := redis.String("GET", sb.String())
	if err != nil {
		badRequest("Redis Error", w)
		log.Printf("Redis GET Error : %+v", err)
		return
	}

	var rdsPhone redis.RedisPhone

	if red == "" {
		badRequest("Ouch it's seem your code has been invalidated or invalid, please try again.", w)
		log.Printf("Redis not found for phone number : %s", vd.Number)
		return
	}

	log.Printf("Redis GET key : %s | Value: %+v", "number", red)

	// unmarshal redis data
	if err = json.Unmarshal([]byte(red), &rdsPhone); err != nil {
		badRequest("JSON Unmarshal Error", w)
		log.Printf("JSON Unmarshal Error : %+v", err)
		return
	}

	if rdsPhone.Number != vd.Number {
		badRequest("OTP Code is invalid", w)
		log.Println("OTP Code is invalid")
		return
	}

	res, err := redis.Int("DEL", sb.String())
	if err != nil {
		badRequest("Redis DEL Error", w)
		log.Printf("Redis DEL Error: %+v", err)
		return
	}

	if res == 0 {
		badRequest("Redis DEL Error", w)
		log.Printf("Redis DEL Error: %+v", err)
		return
	}

	// get database instance
	tx, err := config.GetDBInstance().Beginx()
	if err != nil {
		badRequest("DB Error", w)
		log.Printf("DB Error : %+v", err)
		return
	}
	defer tx.Rollback()

	_, err = models.UpdatePhoneValidateDB(vd.Number, vd.Code, tx)
	if err != nil {
		badRequest("Update Phone DB Error", w)
		log.Printf("Update Phone DB Error : %+v", err)
		return
	}

	okResponse("Success validating OTP code", nil, w)
	return
}

func IsPhoneValid(number string) bool {
	r, _ := regexp.Compile("^(\\+62|62|0)8[1-9][0-9]{6,9}$")
	return r.MatchString(number)
}

func OTPGenerator(n int) string {
	var letters = []rune("1234567890")
	l := make([]rune, n)
	for i := range l {
		l[i] = letters[rand.Intn(len(letters))]
	}

	return string(l)
}
