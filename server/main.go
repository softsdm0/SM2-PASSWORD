package main

import (
	"fmt"
	"log"
	"time"

	"gitee.com/ouhaoqiang/goserver"
	c "gitee.com/ouhaoqiang/goserver/utils/config"
	"gitee.com/ouhaoqiang/goserver/utils/logger"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login/casdoor"
	"gitee.com/ouhaoqiang/passwordserver/server/src/router"
	"gitee.com/ouhaoqiang/passwordserver/server/src/session"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/cache"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLoggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	gS := goserver.InitGoServer(&Server{})
	gS.Start()
}

// 创建一个结构体，实现Run()和Close()方法
type Server struct {
	app *fiber.App
}

// 运行方法
func (s *Server) Run() {
	// 临时调用goserver的ConfigPath
	// TODO: 以后goserver修改
	config.ConfigPath = c.ConfigPath
	err := config.InitConfig()
	if err != nil {
		fmt.Printf("init config err: %v\n", err)
		return
	}

	// 初始化casdoorsdk
	casdoor.InitCasdoorConfig()

	app := fiber.New(fiber.Config{
		// 开启多线程
		Prefork: false,
		// 限制body
		BodyLimit: 0,
		// 最大并发连接数
		Concurrency:  0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		ErrorHandler: func(*fiber.Ctx, error) error {
			return nil
		},
		// 关闭保持连接
		DisableKeepalive: false,
		AppName:          config.Config.Server.Name,
		// 如果设置为true，则会以更高的CPU使用率为代价大幅降低内存使用率。
		// 只有当服务器为大部分空闲的保活连接消耗过多内存时，才尝试启用此选项。这可能会将内存使用量减少50%以上。
		ReduceMemoryUsage: false,
		// 设置json解析器
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		EnableTrustedProxyCheck: false,
		EnableIPValidation:      false,
		// 如果设置为true，将打印所有路由及其方法、路径和处理程序。
		EnablePrintRoutes: true,
		// 设置颜色
		ColorScheme: fiber.Colors{},
	})

	// 日志
	app.Use(fiberLoggerMiddleware.New())

	// 添加跨域中间件
	corsCfg := config.Config.Server.Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsCfg.AllowOrigins,
		AllowHeaders:     corsCfg.AllowHeaders,
		AllowMethods:     corsCfg.AllowMethods,
		AllowCredentials: corsCfg.AllowCredentials,
	}))

	s.app = app

	// 连接数据库
	database.ConnectDB()

	// 连接缓存
	cache.ConnectCache()

	// 初始化sm2密钥
	session.InitSm2Info()

	// 设置路由
	router.SetupRoutes(app)

	// 启动服务
	http := config.Config.Server.HTTP
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", http.Bind, http.Port)))
}

// 关闭处理方法
func (s *Server) Close() error {
	logger.Logger.Debug("开始关闭服务")

	time.Sleep(time.Second * 1)
	logger.Logger.Debug("业务已经处理完成")
	logger.Logger.Debug("服务关闭")
	return nil
}
