# HITS

![Hits](https://storage.googleapis.com/hit-counter/main.png)
一种查看有多少人访问过您的网站或 GitHub 存储库的简单方法。
<p align="center">
<a href="https://circleci.com/gh/gjbae1212/hit-counter"><img src="https://circleci.com/gh/gjbae1212/hit-counter.svg?style=svg"></a>
<a href="https://hits.seeyoufarm.com"><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fcn1095%2Fhit-counter%2FREADME&count_bg=%2379C83D&title_bg=%23555555&icon=go.svg&icon_color=%2300ADD8&title=hits&edge_flat=false"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-GPL-blue.svg" alt="license" /></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/hit-counter"><img src="https://goreportcard.com/badge/github.com/gjbae1212/hit-counter" alt="Go Report Card" /></a> 
</p>

## 概述

[HITS](https://hits.seeyoufarm.com) 提供展示 **标题** 和 **今日日/总量** 访问数的 SVG 徽章。

如果您将徽章嵌入到网站、GitHub 或 Notion 上，则每个页面请求（页面点击）都会被计数。

徽章包括每天（从 GMT 开始）和总（所有）访问数。

[HITS](https://hits.seeyoufarm.com) 还显示访问量最高的 GitHub 项目。 （前10名）

[HITS](https://hits.seeyoufarm.com) 显示使用此服务的每个项目或站点的实时页面点击量（使用 Websocket）。

[HITS](https://hits.seeyoufarm.com) 是由 gjbae1212@gmail.com 使用 Golang、WebAssembly (Wasm)、HTML 制作的，目前由 Google Cloud 平台提供服务。
 
## 如何使用
### 如何生成徽章
您可以通过 [HITS](https://hits.seeyoufarm.com/#badge) 生成徽章。

![Hits](https://storage.googleapis.com/hit-counter/gen.png)

## 特征
- 显示页面上的每日页面浏览量和总页面浏览量。  
- 支持自定义样式的徽章。
- 支持徽章免费图标（https://simpleicons.org）。
- 显示您网站的图表，了解最近 6 个月内每日历史记录的数量。
- 显示 github 项目的排名。
- 显示实时流。
      
## ETC
[HITS](https://hits.seeyoufarm.com) 计算每个页面的点击量，而不存储敏感信息（IP、标头等）。 
为了防止海量请求滥用，部分请求信息会转换为本地缓存中的哈希数据，并在经过一段时间后删除。

此外，HITS 不使用 GitHub Traffic 或 Google Analytics 数据，它只是计算您网站或存储库的每个页面点击量。
  
## 执照
项目已获得 V3.0 许可。
