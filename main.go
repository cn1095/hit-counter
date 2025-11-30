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
var embeddedFiles embed.FS  
  
func EmbeddedFiles() embed.FS {  
	return embeddedFiles  
}  
  
var (  
	address       = flag.String("port", ":8080", "服务器监听地址和端口")  
	tls           = flag.Bool("tls", false, "启用 Let's Encrypt 自动 TLS")  
	redis         = flag.String("redis", "", "Redis 服务器地址 (覆盖 REDIS_ADDRS)")  
	redisPassword = flag.String("redis-password", "", "Redis 服务器密码 (覆盖 REDIS_PASSWORD)")  
	phase         = flag.String("phase", "", "部署阶段: local, development, production (覆盖 PHASE)")  
	debug         = flag.Bool("debug", false, "启用调试模式 (覆盖 DEBUG)")  
	logPath       = flag.String("log", "", "日志文件路径 (覆盖 LOG_PATH)")  
	sentry        = flag.String("sentry", "", "Sentry DSN 错误追踪 (覆盖 SENTRY_DSN)")  
)  
  
// 配置结构体，优先使用命令行参数  
type Config struct {  
	Address       string  
	TLS           bool  
	RedisAddr     string  
	RedisPassword string  
	Phase         string  
	Debug         bool  
	LogPath       string  
	SentryDSN     string  
}  
  
// 获取配置，命令行参数优先于环境变量  
func getConfig() *Config {  
	config := &Config{  
		Address:       *address,  
		TLS:           *tls,  
		RedisAddr:     *redis,  
		RedisPassword: *redisPassword,  
		Phase:         *phase,  
		Debug:         *debug,  
		LogPath:       *logPath,  
		SentryDSN:     *sentry,  
	}  
  
	// 只有在命令行参数为空时才回退到环境变量  
	if config.RedisAddr == "" {  
		if addrs := env.GetRedisAddrs(); len(addrs) > 0 {  
			config.RedisAddr = addrs[0]  
		}  
	}  
	  
	if config.RedisPassword == "" {  
		config.RedisPassword = env.GetRedisPassword()  
	}  
	  
	if config.Phase == "" {  
		config.Phase = env.GetPhase()  
	}  
	  
	if !config.Debug {  
		config.Debug = env.GetDebug()  
	}  
	  
	if config.LogPath == "" {  
		config.LogPath = env.GetLogPath()  
	}  
	  
	if config.SentryDSN == "" {  
		config.SentryDSN = env.GetSentryDSN()  
	}  
  
	return config  
}  
  
func main() {  
	flag.Parse()  
  
	config := getConfig()  
  
	runtime.GOMAXPROCS(runtime.NumCPU())  
  
	// 初始化 sentry，使用配置中的值  
	if config.SentryDSN != "" {  
    	name, _ := os.Hostname()  
    	if err := internal.InitSentry(config.SentryDSN, config.Phase, config.Phase,  
        	name, true, config.Debug); err != nil {  
        	log.Println(err)  
    	}  
	}
  
	e := echo.New()  
  
	// 配置 echo 服务器选项  
	var opts []Option  
	opts = append(opts, WithDebugOption(config.Debug))  
  
	var dir string  
	var file string  
	if config.LogPath != "" {  
		dir, file = filepath.Split(config.LogPath)  
	}  
  
	logger, err := internal.NewLogger(dir, file)  
	if err != nil {  
		log.Panic(err)  
	}  
	opts = append(opts, WithLogger(logger))  
  
	if err := AddMiddleware(e, opts...); err != nil {  
		log.Panic(err)  
	}  
	    
	staticFS := echo.MustSubFS(embeddedFiles, "public")  
	e.StaticFS("/", staticFS)  
	    
	// 验证 Redis 地址  
	if config.RedisAddr == "" {  
		log.Panic("Redis 地址未设置，请使用 -redis 参数或设置 REDIS_ADDRS 环境变量")  
	}  
  
	// 添加路由，直接使用配置中的 Redis 地址  
	if err := AddRoute(e, config.RedisAddr, config.RedisPassword); err != nil {  
    	log.Panic(err)  
	} 
  
	if config.TLS {  
		e.StartAutoTLS(config.Address)  
	} else {  
		e.Start(config.Address)  
	}  
}
