package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/anacrolix/log"

	"github.com/segmentio/ksuid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTStorage struct {
	ExpiresAt    int
	PublicString string
	PublicKey    *rsa.PublicKey
	privateKey   *rsa.PrivateKey
}

func ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) < 2 {
		return ""
	}
	return headerSplit[1]
}

func (j *JWTStorage) VerifyUser(c *gin.Context) {
	token, err := jwt.Parse(ExtractToken(c), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.PublicKey, nil
	})

	if err != nil {
		log.Printf("Error parsing JWT: %s\n", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user", claims["sub"])
	} else {
		log.Printf("Invalid JWT Format\n")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
		return
	}

	c.Next()
}

// NewAnonToken creates an anonymous jwt
func (j *JWTStorage) NewAnonToken() string {
	claims := jwt.StandardClaims{
		Subject:   ksuid.New().String(),
		ExpiresAt: time.Now().Add(time.Duration(j.ExpiresAt) * time.Hour).Unix(),
		Issuer:    "Quirk",
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	ss, _ := token.SignedString(j.privateKey)
	return ss
}

// InitJWT returns a new JWTStorage struct based on parameters
// Errors if unable to read RSA keys
func InitJWT(expiry int, keyType, privateKey, publicKey string) (*JWTStorage, error) {
	if expiry <= 0 {
		return nil, fmt.Errorf("JWT expiry [%d] must be a positive integer", expiry)
	}

	jwtInfo := &JWTStorage{ExpiresAt: expiry}
	var err error

	switch keyType {
	case "file":
		err = jwtInfo.readKeysFile(privateKey, publicKey)
	default:
		err = fmt.Errorf("unsupported keytype [%s]", keyType)
	}

	return jwtInfo, err
}

func (j *JWTStorage) readKeysFile(privateKey, publicKey string) error {
	publicBytes, err := ioutil.ReadFile(publicKey)
	if err != nil {
		return fmt.Errorf("error reading public key [%s]: %s", publicKey, err.Error())
	}
	j.PublicString = string(publicBytes)

	j.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return fmt.Errorf("error parsing public key [%s]: %s", publicKey, err.Error())
	}

	privateBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return fmt.Errorf("error reading private key [%s]: %s", publicKey, err.Error())
	}

	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return fmt.Errorf("error parsing private key [%s]: %s", publicKey, err.Error())
	}
	return nil
}
