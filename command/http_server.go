package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"payton/route"
	"time"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http服务",
	Long:  "http服务相关命令",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	httpCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动",
	Long:  "启动http服务",
	Run: func(cmd *cobra.Command, args []string) {
		HttpServerStart()
	},
}

func HttpServerStart() {
	var err error
	e := echo.New()
	e.HideBanner = true
	// 中间件
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// 注册路由
	route.RegisterRoute(e)

	httpListen := viper.GetString("http.listen")
	go func() {
		if err := e.Start(httpListen); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
