package backend

import (
	"database/sql"
	mydb "go-demoservice/db"
	"go-demoservice/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func App(db *sql.DB, cacheMap map[string]utils.Message) {
	var msg utils.Message

	e := echo.New()

	e.GET("/order/:order_uid", func(c echo.Context) error {

		startOperation := time.Now()

		order_uid := c.Param("order_uid")
		if _, ok := cacheMap[order_uid]; ok {
			utils.BackendLogger.Infof("Message got from map by %s sec-s", time.Since(startOperation).String())
			return c.JSON(http.StatusOK, cacheMap[order_uid])
		}

		errCreateMessage := mydb.GetInfoAndCreateMessage(db, order_uid, &msg)

		if errCreateMessage != nil {
			utils.BackendLogger.Error(errCreateMessage)
			return c.String(http.StatusNotFound, "Info wasn't found")
		}

		cacheMap[order_uid] = msg

		utils.BackendLogger.Infof("Message got from DB by %s sec-s", time.Since(startOperation).String())
		return c.JSON(http.StatusOK, msg)
	})

	// e.GET("/", func(c echo.Context) error {
	// 	return nil
	// })

	e.Start(":8080")

}
