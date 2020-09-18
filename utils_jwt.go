package goutils

import (
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/twinj/uuid"
)

type PayloadSetting struct {
	Issuer    string
	Subject   string
	Audience  []string
	Expired   time.Duration
	TokenType string
}

type PayloadParam struct {
	UserId    string
	UserEmail string
	Created   string
}

type PayloadData struct {
	Setting PayloadSetting
	Param   PayloadParam
}

type TokenPayload struct {
	jwt.Payload
	Data PayloadData
}

type TokenResult struct {
	Status bool
	Err    error
	Msg    string
	Data   PayloadData
}

func GetToken(payload PayloadData, secret string) (string, error) {
	now := time.Now()
	hs := jwt.NewHS256([]byte(secret))
	u4 := uuid.NewV4()
	randomId := u4.String()

	tokenPayload := TokenPayload{
		Payload: jwt.Payload{
			Issuer:         payload.Setting.Issuer,
			Subject:        payload.Setting.Subject,
			Audience:       jwt.Audience(payload.Setting.Audience),
			ExpirationTime: jwt.NumericDate(now.Add(payload.Setting.Expired)),
			NotBefore:      jwt.NumericDate(now.Add(payload.Setting.Expired)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          randomId,
		},
		Data: payload,
	}

	token, err := jwt.Sign(tokenPayload, hs)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func VerifyToken(token string, secret string, audience []string) TokenResult {
	var result TokenResult
	var hs = jwt.NewHS256([]byte(secret))
	var (
		now = time.Now()
		aud = jwt.Audience(audience)
		// Validate claims "iat", "exp" and "aud".
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		audValidator = jwt.AudienceValidator(aud)

		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		tokenPayload    TokenPayload
		validatePayload = jwt.ValidatePayload(&tokenPayload.Payload, iatValidator, expValidator, audValidator)
	)

	_, err := jwt.Verify([]byte(token), hs, &tokenPayload, validatePayload)
	if err != nil {
		switch err {
		case jwt.ErrIatValidation:
			result = TokenResult{false, err, "Unable validate iat", tokenPayload.Data}
			return result
		case jwt.ErrExpValidation:
			result = TokenResult{false, err, "Expired token", tokenPayload.Data}
			return result
		case jwt.ErrAudValidation:
			result = TokenResult{false, err, "Unable validate jwt au", tokenPayload.Data}
			return result
		default:
			result = TokenResult{false, err, "Unknown error", tokenPayload.Data}
			return result
		}
	}

	result = TokenResult{true, err, "", tokenPayload.Data}

	return result
}

func VerifyExpiredToken(token string, secret string, audience []string) (bool, error, string) {
	var (
		res bool
		err error
		msg string
	)

	var hs = jwt.NewHS256([]byte(secret))

	var (
		now = time.Now()
		aud = jwt.Audience(audience)
		// Validate claims "iat", "exp" and "aud".
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		audValidator = jwt.AudienceValidator(aud)

		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		tokenPayload    TokenPayload
		validatePayload = jwt.ValidatePayload(&tokenPayload.Payload, iatValidator, expValidator, audValidator)
	)

	_, err = jwt.Verify([]byte(token), hs, &tokenPayload, validatePayload)
	if err != nil {
		switch err {
		case jwt.ErrIatValidation:
			res = false
			msg = "Unable validate iat"
			return res, err, msg
		case jwt.ErrExpValidation:
			res = true
			msg = "Expired token"
			return res, err, msg
		case jwt.ErrAudValidation:
			res = false
			msg = "Unable validate jwt aud"
			return res, err, msg
		default:
			res = false
			msg = "Unknown error"
			return res, err, msg
		}
	}

	res = false
	msg = "Unable regenerate token. Token is not expired!"
	return res, err, msg
}
