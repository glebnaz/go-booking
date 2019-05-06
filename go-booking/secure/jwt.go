package secure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type JWT struct {
	h header
	p payload
	s string
}
type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Id   string
	Time int64
}

func (h *header) GenerateHeaderKey() {
	h.Alg = "HS256"
	h.Typ = "JWT"
}

func (p *payload) GeneratePayloadKey(id string) {
	p.Id = id
	p.Time = time.Now().Add(time.Hour * 10).Unix()
}

func (j *JWT) GenerateToken(id string) string {
	j.h.GenerateHeaderKey()
	j.p.GeneratePayloadKey(id)
	secretKey := generateSecretKey()
	headerJSON, err := json.Marshal(j.h)
	if err != nil {
		fmt.Println("fatal err in GenerateSignature with make json header")
	}
	headerBase64 := base64.StdEncoding.EncodeToString(headerJSON)
	payloadJSON, err := json.Marshal(j.p)
	if err != nil {
		fmt.Println("fatal err in GenerateSignature with make json payload")
	}
	payloadBase64 := base64.StdEncoding.EncodeToString(payloadJSON)
	signatureNoCode := headerBase64 + "." + payloadBase64
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signatureNoCode))
	signatureCode := base64.StdEncoding.EncodeToString(h.Sum(nil))
	token := signatureNoCode + "." + signatureCode
	j.s = token
	return secretKey
}

func (j *JWT) CheckToken(secret string) bool {
	tokenArr := strings.Split(j.s, ".")
	if len(tokenArr) < 3 {
		return false
	} else {
		signatureNoCode := tokenArr[0] + "." + tokenArr[1]
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(signatureNoCode))
		tv := base64.StdEncoding.EncodeToString(h.Sum(nil))
		if tv == tokenArr[2] {
			return true
		} else {
			return false
		}
	}
}

func (t JWT) GetToken() string {
	return t.s
}
func (t *JWT) SetSignature(token string) {
	t.s = token
}
func DecodeInfo(token string) payload {
	tokenArr := strings.Split(token, ".")
	jsonString, err := base64.StdEncoding.DecodeString(tokenArr[1])
	if err != nil {
		fmt.Println("err in decode info token")
	}
	var tokenInfo payload
	err = json.Unmarshal(jsonString, &tokenInfo)
	if err != nil {
		fmt.Println("err in decode info token (decode json)")
	}
	return tokenInfo
}

func generateSecretKey() string {
	rand.Seed(time.Now().UTC().UnixNano())
	randomInt := 1 + rand.Intn(51-0)
	const letterBytes = "abOdefghijklmnopqrstuvwxyz789DEFGHIJKLMNcPQRSTUVWXYZ0123456ABC"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(randomInt)]
	}
	return string(b)
}
