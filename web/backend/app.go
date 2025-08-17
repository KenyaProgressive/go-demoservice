package backend

import (
	"context"
	"database/sql"
	mydb "go-demoservice/db"
	"go-demoservice/utils"
	"net/http"
	"sync"
	"time"

	_ "go-demoservice/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title go-demoservice (prod. Nikita Parfenov)
// @version 1.0
// @description Документация к первому проекту "Wb-Техношколы"
// @host localhost:8080
// @BasePath /
func App(db *sql.DB, cacheMap map[string]utils.Message, wg *sync.WaitGroup, ctx context.Context) {

	defer wg.Done()

	e := echo.New()
	e.Static("/", "web/frontend")

	e.GET("/order/:order_uid", func(c echo.Context) error {
		// Getting order_data by order_uid
		return orderUidHandler(c, cacheMap, db)
	})

	e.GET("/docs/*", echoSwagger.WrapHandler)

	go func() {
		utils.BackendLogger.Info("Server was started")
		if errLaunchServer := e.Start(":8080"); errLaunchServer != http.ErrServerClosed {
			utils.BackendLogger.Error(errLaunchServer)
		}
	}()

	<-ctx.Done()

	utils.BackendLogger.Info("Server shutting down")

	shuttingDownContext, cancelShuttingDown := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelShuttingDown()

	if errShuttingDown := e.Shutdown(shuttingDownContext); errShuttingDown != nil {
		utils.BackendLogger.Errorf("Problem with shutting down server: %s", errShuttingDown)
	}

	utils.BackendLogger.Info("Server successfully stopped")

}

// orderUidHandler godoc
// @Summary Получить заказ по UUID (order_uid)
// @Description Возвращает информацию о заказе
// @Param order_uid path string true "Order ID"
// @Success 200 {object} utils.Message
// @Failure 404 {string} string "Info wasn't found"
// @Router /order/{order_uid} [get]
func orderUidHandler(c echo.Context, cacheMap map[string]utils.Message, db *sql.DB) error {
	var msg utils.Message

	startOperation := time.Now()

	order_uid := c.Param("order_uid")
	if _, ok := cacheMap[order_uid]; ok {
		utils.BackendLogger.Infof("Message got from map by %f sec-s", time.Since(startOperation).Seconds())
		return c.JSON(http.StatusOK, cacheMap[order_uid])
	}

	errCreateMessage := mydb.GetInfoAndCreateMessage(db, order_uid, &msg)

	if errCreateMessage != nil {
		utils.BackendLogger.Error(errCreateMessage)
		return c.String(http.StatusNotFound, "Info wasn't found")
	}

	cacheMap[order_uid] = msg

	utils.BackendLogger.Infof("Message got from DB by %f sec-s", time.Since(startOperation).Seconds())

	return c.JSON(http.StatusOK, msg)
}
