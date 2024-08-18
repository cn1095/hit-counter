package handler

import (
	"bytes"
	"fmt"
	"iter"
	"net/http"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/labstack/echo/v4"
	perrors "github.com/pkg/errors"
)

// Index is API for main page.
func (h *handler) Index(c echo.Context) error {
	const (
		topK = 10
	)

	var ranks []string
	rangeTopK, err := h.getGitHubTopK(c, topK)
	if err != nil {
		return perrors.WithStack(err)
	}
	for rank := range rangeTopK {
		ranks = append(ranks, rank)
	}

	buf := bytes.NewBuffer([]byte{})
	if err := h.indexTemplate.Execute(buf, struct {
		Ranks []string
	}{Ranks: ranks}); err != nil {
		return perrors.WithStack(err)
	}

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

// getGitHubTopK returns range loop for getting top k projects of github.com.
func (h *handler) getGitHubTopK(c echo.Context, k int) (iter.Seq[string], error) {
	hitCtx := c.(*servercontext.HitCounterContext)

	const (
		group = "github.com"
	)
	limit := k * 2

	// get top k *2 projects.
	scores, err := h.counter.GetRankTotalByLimit(hitCtx.GetContext(), group, limit)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return func(yield func(string) bool) {
		rankNum := 0
		exist := mapset.NewSetWithSize[string](len(scores))
		for _, score := range scores {
			if rankNum == k {
				return
			}

			path := strings.TrimSpace(score.Name)
			if strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}

			// add projects if score-name is equal to /profile/project.
			seps := strings.Split(path, "/")
			if len(seps) == 3 && !exist.ContainsOne(path) {
				exist.Add(path)
				rankNum++
				if !yield(fmt.Sprintf("[%d] %s%s", rankNum, group, path)) {
					return
				}
			}
		}
	}, nil
}
