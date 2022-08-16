package store_handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// requestStore payload for storing new data in raft cluster
type requestStore struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	IsAppend bool        `json:"isAppend"`
}

// Store handling save to raft cluster. Store will invoke raft.Apply to make this stored in all cluster
// with acknowledge from n quorum. Store must be done in raft leader, otherwise return error.
func (h handler) Store(eCtx echo.Context) error {
	var form = requestStore{}
	if err := eCtx.Bind(&form); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	form.Key = strings.TrimSpace(form.Key)
	if form.Key == "" {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "key is empty",
		})
	}
	var err error
	if form.IsAppend {
		fmt.Println("append data into DB")
		err = h.DbHandler.SetArr(form.Key, form.Value)
	} else {
		err = h.DbHandler.Set(form.Key, form.Value)
	}

	if err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error preparing saving data payload: %s", err.Error()),
		})
	}

	//if err := applyFuture.Error(); err != nil {
	//	return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
	//		"error": fmt.Sprintf("error persisting data in raft cluster: %s", err.Error()),
	//	})
	//}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success persisting data",
		"data":    form,
	})
}
