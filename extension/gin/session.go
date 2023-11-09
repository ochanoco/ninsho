package gin

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ochanoco/ninsho"
	// "ninsho"
)

const USER_SESSION_NAME = "user"
const LOGINING_SESSION_NAME = "logining_session"

func loadStruct[T any](key string, c *gin.Context) (*T, error) {
	var t T
	sess := sessions.Default(c)
	j := sess.Get(key)
	if j == nil {
		return nil, nil
	}

	buf := []byte(j.(string))
	err := json.Unmarshal(buf, &t)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize data: %v", err)
	}

	return &t, err
}

func saveStruct[T any](t T, key string, c *gin.Context) error {
	buf, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to serialize data: %v", err)
	}

	sess := sessions.Default(c)
	sess.Set(key, string(buf))
	sess.Save()

	return nil
}

func LoadUser[T ninsho.User](c *gin.Context) (*T, error) {
	jwt, err := loadStruct[T](USER_SESSION_NAME, c)

	if err != nil {
		return nil, fmt.Errorf("failed to load user: %v", err)
	}

	return jwt, nil
}

func SaveUser(user ninsho.User, c *gin.Context) error {
	err := saveStruct(user, USER_SESSION_NAME, c)

	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func LoadLoginingSession(c *gin.Context) (*ninsho.Session, error) {
	loginingSess, err := loadStruct[ninsho.Session](LOGINING_SESSION_NAME, c)

	if err != nil {
		fmt.Printf("failed to load logining session: %v", err)
	}

	return loginingSess, err
}

func SaveLoginingSession(loginingSess ninsho.Session, c *gin.Context) error {
	err := saveStruct(loginingSess, LOGINING_SESSION_NAME, c)

	if err != nil {
		return fmt.Errorf("failed to save logining session: %v", err)
	}

	return nil
}
