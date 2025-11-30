package main // import "github.com/cn1095/hit-counter"  
  
import (  
	"embed"  
	"flag"  
	"log"  
	"os"  
	"runtime"  
  
	"github.com/cn1095/hit-counter/internal"  
  
	"path/filepath"  
  
	"github.com/cn1095/hit-counter/env"  
	"github.com/labstack/echo/v4"  
)  
  
//go:embed public/* view/*  
var embeddedFiles embed.FS // 嵌入 public 和 view 目录中的所有文件  
  
var (  
	address = flag.String("port", ":8080", "服务器监听地址和端口 (例如: :8080, 0.0.0.0:8080)")  
	tls     = flag.Bool("tls", false, "启用 Let's Encrypt 自动 TLS")  
	redis   = flag.String("redis", "", "Redis 服务器地址 (逗号分隔多个地址, 覆盖 REDIS_ADDRS)")  
	redisPassword = flag.String("redis-password", "", "Redis 服务器密码 (覆盖 REDIS_PASSWORD)")
	phase   = flag.String("phase", "", "部署阶段: local, development, production (覆盖 PHASE)")  
	debug   = flag.Bool("debug", false, "启用调试模式 (覆盖 DEBUG)")  
	logPath = flag.String("log", "", "日志文件路径 (覆盖 LOG_PATH)")  
	sentry  = flag.String("sentry", "", "Sentry DSN 错误追踪 (覆盖 SENTRY_DSN)")  
)  
  
func main() {  
	flag.Parse()  
  
	// 使用命令行参数覆盖环境变量  
	if *redis != "" {  
		os.Setenv("REDIS_ADDRS", *redis)  
	}  
	if *redisPassword != "" {  
		os.Setenv("REDIS_PASSWORD", *redisPassword)  
	} 
	if *phase != "" {  
		os.Setenv("PHASE", *phase)  
	}  
	if *debug {  
		os.Setenv("DEBUG", "true")  
	}  
	if *logPath != "" {  
		os.Setenv("LOG_PATH", *logPath)  
	}  
	if *sentry != "" {  
		os.Setenv("SENTRY_DSN", *sentry)  
	}  
  
	runtime.GOMAXPROCS(runtime.NumCPU())  
  
	// 初始化 sentry  
	name, _ := os.Hostname()  
	if err := internal.InitSentry(env.GetSentryDSN(), env.GetPhase(), env.GetPhase(),  
		name, true, env.GetDebug()); err != nil {  
		log.Println(err)  
	}  
  
	e := echo.New()  
  
	// 配置 echo 服务器选项  
	var opts []Option  
  
	// 调试选项  
	opts = append(opts, WithDebugOption(env.GetDebug()))  
  
	var dir string  
	var file string  
	if env.GetLogPath() != "" {  
		dir, file = filepath.Split(env.GetLogPath())  
	}  
  
	// 日志选项  
	logger, err := internal.NewLogger(dir, file)  
	if err != nil {  
		log.Panic(err)  
	}  
	opts = append(opts, WithLogger(logger))  
  
	// 添加中间件  
	if err := AddMiddleware(e, opts...); err != nil {  
		log.Panic(err)  
	}  
	  
	// 设置静态文件目录为嵌入文件系统  
	staticFS := echo.MustSubFS(embeddedFiles, "public")  
	e.StaticFS("/", staticFS)  
	  
	// 添加路由  
	if err := AddRoute(e, env.GetRedisAddrs()[0]); err != nil {  
		log.Panic(err)  
	}  
  
	if *tls {  
		// 启动 TLS 服务器，使用 Let's Encrypt 证书  
		e.StartAutoTLS(*address)  
	} else {  
		e.Start(*address)  
	}  
}
