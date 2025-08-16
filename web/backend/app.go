package backend

import (
	"database/sql"
	mydb "go-demoservice/db"
	"go-demoservice/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func App(db *sql.DB) {
	var msg utils.Message

	e := echo.New()

	e.GET("/order/:order_uid", func(c echo.Context) error {
		order_uid := c.Param("order_uid")
		utils.BackendLogger.Debug(order_uid)
		errCreateMessage := mydb.GetInfoAndCreateMessage(db, order_uid, &msg)
		if errCreateMessage != nil {
			utils.BackendLogger.Error(errCreateMessage)
			return c.String(http.StatusNotFound, "Info wasn't found")
		}
		return c.JSON(http.StatusOK, msg)
	})

	e.Start(":8080")

}
