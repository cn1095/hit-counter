package main  
  
import (  
	"fmt"  
	"net/http"  
	"strconv"  
  
	"github.com/cn1095/hit-counter/handler"  
	api_handler "github.com/cn1095/hit-counter/handler/api"  
	"github.com/cn1095/hit-counter/internal"  
	"github.com/labstack/echo/v4"  
)  
  
func AddRoute(e *echo.Echo, redisAddr string) error {  
	if e == nil {  
		return fmt.Errorf("[错误] AddRoute参数为空")  
	}  
  
	// 传递 embeddedFiles 给 NewHandler  
	h, err := handler.NewHandler(redisAddr, embeddedFiles)  
	if err != nil {  
		return fmt.Errorf("[错误] AddRoute %w", err)  
	}  
  
	api, err := api_handler.NewHandler(h)  
	if err != nil {  
		return fmt.Errorf("[错误] AddRoute %w", err)  
	}  
  
	// error handler  
	e.HTTPErrorHandler = h.Error  
  
	// WebAssembly 文件  
	e.GET("/hits.wasm", h.Wasm)  
  
	// WebSocket  
	e.GET("/ws", h.WebSocket)  
  
	// 主页  
	e.GET("/", h.Index)  
  
	// 图标  
	e.GET("/icon/all.json", h.IconAll)  
	e.GET("/icon/:icon", h.Icon)  
  
	// 健康检查  
	e.GET("/healthcheck", h.HealthCheck)  
  
	// API 组 /api/count  
	g1, err := groupApiCount()  
	if err != nil {  
		return fmt.Errorf("[错误] AddRoute %w", err)  
	}  
	count := e.Group("/api/count", g1...)  
	// 徽章  
	count.GET("/keep/badge.svg", api.KeepCount)  
	count.GET("/incr/badge.svg", api.IncrCount)  
  
	// 图表  
	count.GET("/graph/dailyhits.svg", api.DailyHitsInRecently)  
  
	// API 组 /api/rank  
	g2, err := groupApiRank()  
	if err != nil {  
		return fmt.Errorf("[错误] AddRoute %w", err)  
	}  
	rank := e.Group("/api/rank", g2...)  
	_ = rank  
  
	return nil  
}  
  
func groupApiCount() ([]echo.MiddlewareFunc, error) {  
	var chain []echo.MiddlewareFunc  
	// Add param  
	paramFunc := func(h echo.HandlerFunc) echo.HandlerFunc {  
		return func(c echo.Context) error {  
			hitctx := c.(*handler.HitCounterContext)  
  
			// check a url is invalid or not.  
			url := hitctx.QueryParam("url")  
			if url == "" {  
				return echo.NewHTTPError(http.StatusBadRequest, "未找到 URL 查询字符串")  
			}  
  
			schema, host, _, path, _, _, err := internal.ParseURL(url)  
			if err != nil {  
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("无效的 URL 查询字符串 %s", url))  
			}  
  
			if !internal.StringInSlice(schema, []string{"http", "https"}) {  
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("不支持的协议/模式 %s", schema))  
			}  
  
			// extract required parameters  
			title := hitctx.QueryParam("title")  
			titleBg := hitctx.QueryParam("title_bg")  
			countBg := hitctx.QueryParam("count_bg")  
			edgeFlat, _ := strconv.ParseBool(hitctx.QueryParam("edge_flat"))  
			icon := hitctx.QueryParam("icon")  
			iconColor := hitctx.QueryParam("icon_color")  
  
			// insert params to context.  
			hitctx.Set("host", host)  
			hitctx.Set("path", path)  
			hitctx.Set("title", title)  
			hitctx.Set("title_bg", titleBg)  
			hitctx.Set("count_bg", countBg)  
			hitctx.Set("edge_flat", edgeFlat)  
			hitctx.Set("icon", icon)  
			hitctx.Set("icon_color", iconColor)  
  
			return h(hitctx)  
		}  
	}  
	chain = append(chain, paramFunc)  
	return chain, nil  
}  
  
func groupApiRank() ([]echo.MiddlewareFunc, error) {  
	var chain []echo.MiddlewareFunc  
	return chain, nil  
}
