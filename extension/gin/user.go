package gin

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUser[T any](c *gin.Context) (*T, error) {
	var jwt T

	sess := sessions.Default(c)
	user := sess.Get("user")
	if user == "" {
		return nil, fmt.Errorf("user not logged in")
	}

	err := json.Unmarshal([]byte(user.(string)), &jwt)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to deserialize user")
	}

	return &jwt, nil
}
