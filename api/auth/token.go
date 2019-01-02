package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

func ExtractUser(c *gin.Context) string {
	return c.GetHeader("User")
}

func VerifyUser(c *gin.Context) {
	user := ExtractUser(c)
	// TODO: Verify user ID is in the proper GUID format
	if user == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid User authentication",
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
		IssuedAt: time.Now().Unix(),
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
