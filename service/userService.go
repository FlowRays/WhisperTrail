package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"

	"github.com/FlowRays/WhisperTrail/dao"
	"github.com/FlowRays/WhisperTrail/model"
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var params = &argonParams{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

func UserRegister(user *model.User, db *model.Database) error {
	hashedPassword, err := hashPassword(user.Password, params)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = dao.CreateUser(user, db)

	return err
}

func UserLogin(user *model.User, db *model.Database) (string, error) {
	rawPassword := user.Password
	err := dao.GetUserByName(user, db)
	if err != nil {
		return "", err
	}

	valid, err := verifyPassword(user.Password, rawPassword)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("Invalid password")
	}

	return GenerateToken(user.ID)
}

func generateSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func hashPassword(password string, params *argonParams) (string, error) {
	salt, err := generateSalt(params.saltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	hashedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.memory, params.iterations, params.parallelism, b64Salt, b64Hash)

	return hashedPassword, nil
}

func decodeHash(encodedHash string) (params *argonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("Incorrect format of the encoded hash")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("Incompatible version of argon2")
	}

	params = &argonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLength = uint32(len(hash))

	return params, salt, hash, nil
}

func verifyPassword(hashedPassword, password string) (bool, error) {
	p, salt, hash, err := decodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func GenerateToken(userID uint) (string, error) {
	claims := &model.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 设置过期时间
			Issuer:    "WhisperTrail",                        // 设置签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("hs_secret_key"))) // 设置秘钥
}

func ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("hs_secret_key")), nil // 使用相同的秘钥验证签名
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}
