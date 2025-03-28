package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseUintParam extracts and converts a URL parameter to uint
func ParseUintParam(c *gin.Context, param string) (uint, error) {
	strID := c.Param(param)
	id64, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}
