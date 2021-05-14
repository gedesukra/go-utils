package goutils

import (
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/twinj/uuid"
)

type JwtPayload struct {
	Payload jwt.Payload
	Data    map[string]interface{}
}

type JwtResult struct {
	Status bool
	Err    error
	Msg    string
	Data   JwtPayload
}

func JwtTokenGet(issuer string, subject string, audience []string, expired time.Duration, payload map[string]interface{}, secret string) (string, error) {
	now := time.Now()
	hs := jwt.NewHS256([]byte(secret))
	u4 := uuid.NewV4()
	randomId := u4.String()

	jwtPayload := JwtPayload{
		Payload: jwt.Payload{
			Issuer:         issuer,
			Subject:        subject,
			Audience:       jwt.Audience(audience),
			ExpirationTime: jwt.NumericDate(now.Add(expired)),
			NotBefore:      jwt.NumericDate(now.Add(expired)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          randomId,
		},
		Data: payload,
	}

	token, err := jwt.Sign(jwtPayload, hs)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func JwtTokenVerify(token string, audience []string, secret string) JwtResult {
	var result JwtResult
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
		jwtPayload      JwtPayload
		validatePayload = jwt.ValidatePayload(&jwtPayload.Payload, iatValidator, expValidator, audValidator)
	)

	_, err := jwt.Verify([]byte(token), hs, &jwtPayload, validatePayload)
	if err != nil {
		switch err {
		case jwt.ErrIatValidation:
			result = JwtResult{false, err, "Unable validate iat", jwtPayload}
			return result
		case jwt.ErrExpValidation:
			result = JwtResult{false, err, "Expired token", jwtPayload}
			return result
		case jwt.ErrAudValidation:
			result = JwtResult{false, err, "Unable validate jwt au", jwtPayload}
			return result
		default:
			result = JwtResult{false, err, "Unknown error", jwtPayload}
			return result
		}
	}

	result = JwtResult{true, err, "", jwtPayload}

	return result
}
