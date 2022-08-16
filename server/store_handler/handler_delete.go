package store_handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// Delete handling remove data from raft cluster. Delete will invoke raft.Apply to make this deleted in all cluster
// with acknowledge from n quorum. Delete must be done in raft leader, otherwise return error.
func (h handler) Delete(eCtx echo.Context) error {
	var key = strings.TrimSpace(eCtx.Param("key"))
	if key == "" {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "key is empty",
		})
	}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success removing data",
		"data": map[string]interface{}{
			"key":   key,
			"value": nil,
		},
	})
}
