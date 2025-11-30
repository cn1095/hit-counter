package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"syscall/js"
	"time"

	"github.com/goware/urlx"
)

var (
	markdownFormat = "![](%s)"
	showFormat     = "<img src=\"%s\"/>"
	linkFormat     = "&lt;a&gt;&lt;img src=\"%s\"/&gt;"
	incrPath       = "api/count/incr/badge.svg"
	keepPath       = "api/count/keep/badge.svg"
	defaultDomain  = ""
	defaultURL     = ""
	defaultWS      = ""
)

var (
	phase string
)

func parseURL(s string) (schema, host, port, path, query, fragment string, err error) {
	if s == "" {
		err = fmt.Errorf("[err] ParseURI empty uri")
	}

	url, suberr := urlx.Parse(s)
	if suberr != nil {
		err = suberr
		return
	}

	schema = url.Scheme

	host, port, err = urlx.SplitHostPort(url)
	if err != nil {
		return
	}
	if schema == "http" && port == "" {
		port = "80"
	} else if schema == "https" && port == "" {
		port = "443"
	}

	path = url.Path
	query = url.RawQuery
	fragment = url.Fragment
	return
}

func onClick() {
	value := js.Global().Get("document").Call("getElementById", "history_url").Get("value").String()
	value = strings.TrimSpace(value)
	showGraph(value)
}

// DEPRECATED
func onKeyUp() {
	value := js.Global().Get("document").Call("getElementById", "badge_url").Get("value").String()
	value = strings.TrimSpace(value)
	generateBadge(value)
}

// DEPRECATED
func generateBadge(value string) {
	schema, host, _, path, _, _, err := parseURL(value)
	markdown := ""
	link := ""
	show := ""
	incrURL := ""
	keepURL := ""
	if err != nil || (schema != "http" && schema != "https") {
		markdown = "INVALID URL"
		link = "INVALID URL"
	} else {
		normalizeURL := ""
		if path == "" || path == "/" {
			normalizeURL = fmt.Sprintf("%s://%s", schema, host)
		} else {
			normalizeURL = fmt.Sprintf("%s://%s%s", schema, host, path)
		}
		incrURL = fmt.Sprintf("%s/%s?url=%s", defaultURL, incrPath, url.QueryEscape(normalizeURL))
		keepURL = fmt.Sprintf("%s/%s?url=%s", defaultURL, keepPath, url.QueryEscape(normalizeURL))
		markdown = fmt.Sprintf(markdownFormat, incrURL, defaultURL)
		link = fmt.Sprintf(linkFormat, defaultURL, incrURL)
		show = keepURL
	}
	js.Global().Get("document").Call("getElementById", "badge_markdown").Set("innerHTML", markdown)
	js.Global().Get("document").Call("getElementById", "badge_link").Set("innerHTML", link)
	js.Global().Get("document").Call("getElementById", "embed_link").Set("innerHTML", incrURL)
	js.Global().Get("document").Call("getElementById", "badge_show").Set("src", show)
}

func showGraph(value string) {
	schema, _, _, _, _, _, err := parseURL(value)
	if err != nil || (schema != "http" && schema != "https") {
		js.Global().Get("document").Call("getElementById", "history_view").Set("innerHTML", "Not Found")
	} else {
		go func(v string) {
			res, err := http.Get(fmt.Sprintf("%s/api/count/graph/dailyhits.svg?url=%s", defaultURL, v))
			if err != nil {
				js.Global().Get("document").Call("getElementById", "history_view").Set("innerHTML", "Error")
				return
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			encodedBody := base64.StdEncoding.EncodeToString(body)
			if err != nil {
				js.Global().Get("document").Call("getElementById", "history_view").Set("innerHTML", "Error")
				return
			}
			js.Global().Get("document").Call("getElementById", "history_view").Set("innerHTML", "<img class=\"graph_img\" src=\"data:image/svg+xml;base64,"+encodedBody+"\"></div>")
		}(value)
	}
}

func registerCallbacks() {
	// It will be processing when a url input field will be received a event of keyboard up.
	// DEPRECATED
	js.Global().Set("generateBadge", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		onKeyUp()
		return nil
	}))

	js.Global().Set("showGraph", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Get("document").Call("getElementById", "history_button").Set("disabled", true)
		js.Global().Get("document").Call("getElementById", "history_view").Set("innerHTML", `<div class="spinner-border" role="status">
		<span class="sr-only">Loading...</span>
		</div>`)
		onClick()
		js.Global().Get("document").Call("getElementById", "history_button").Set("disabled", false)
		return nil
	}))

	// connect websocket
	connectWebsocket()
}

func connectWebsocket() {
	ws := js.Global().Get("WebSocket").New(defaultWS)
	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("websocket 连接")
		return nil
	}))
	ws.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := args[0].Get("code").Int()
		println(fmt.Sprintf("websocket 关闭 %d\n", code))
		if code == 1000 {
			println("websocket bye!")
		} else {
			go func() {
				select {
				case <-time.After(time.Second * 10):
					connectWebsocket()
				}
			}()
		}
		return nil
	}))
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		p := js.Global().Get("document").Call("createElement", "p")
		p.Set("innerHTML", args[0].Get("data"))
		js.Global().Get("document").Call("getElementById", "stream_view").Call("prepend", p)
		return nil
	}))
	ws.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := args[0].Get("code").String()
		println(fmt.Sprintf("websocket 错误 %s\n", code))
		if "ECONNREFUSED" == code {
			go func() {
				select {
				case <-time.After(time.Second * 10):
					connectWebsocket()
				}
			}()
		} else {
			println("websocket 再见!")
		}
		return nil
	}))
}

func main() {
	// 打印当前阶段
	println("START GO WASM ", phase)

	// 从浏览器动态获取当前页面域名、协议等信息
	location := js.Global().Get("window").Get("location")

	// 获取协议，例如 "http:" 或 "https:"
	protocol := location.Get("protocol").String()

	// 获取 host（包含端口），例如 "example.com:8080"
	host := location.Get("host").String()

	// 获取 hostname（不含端口），例如 "example.com"
	hostname := location.Get("hostname").String()

	// 获取 port，例如 "8080"
	port := location.Get("port").String()

	// 根据协议和 host 生成 API 地址
	if protocol == "http:" {
		defaultURL = "http://" + host
		defaultWS = "ws://" + host + "/ws"
	} else {
		defaultURL = "https://" + host
		defaultWS = "wss://" + host + "/ws"
	}

	// 也可以仅用 hostname，按需选择
	_ = hostname
	_ = port

	// 设置 defaultDomain 为动态获取的 host
	defaultDomain = host

	// 注册回调函数
	registerCallbacks()

	// 阻塞 main，让 WebAssembly 一直运行
	c := make(chan struct{}, 0)
	<-c
}
