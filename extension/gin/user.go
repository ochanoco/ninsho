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
	if user == nil {
		return nil, nil
	}

	buf := []byte(user.(string))
	err := json.Unmarshal(buf, &jwt)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize user: %v", err)
	}

	return &jwt, nil
}
