package handler

import (
	"html/template"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"

	servercontext "github.com/gjbae1212/hit-counter/internal/app/server/context"
	"github.com/gjbae1212/hit-counter/internal/counter"
	"github.com/gjbae1212/hit-counter/web"
	"github.com/labstack/echo/v4"
	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Index(t *testing.T) {
	scores := []*counter.Score{}
	for i := range 100 {
		scores = append(scores, &counter.Score{
			Name:  "github.com/a/" + strconv.Itoa(i),
			Value: int64(i),
		})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Value > scores[j].Value
	})

	tests := []struct {
		name   string
		expect func(t *testing.T) *handler
		ctx    *servercontext.HitCounterContext
		output error
	}{
		{
			name: "success",
			expect: func(t *testing.T) *handler {
				// get index html file.
				indexHtml, err := web.GetIndexHtml("local")
				if err != nil {
					panic(perrors.WithStack(err))
				}
				// parse index template.
				indexTemplate, err := template.New("index").Parse(string(indexHtml))
				if err != nil {
					panic(perrors.WithStack(err))
				}

				c := counter.NewMockCounter(t)
				c.EXPECT().GetRankTotalByLimit(mock.Anything, "github.com", mock.Anything).Return(scores, nil)
				h := &handler{
					indexTemplate: indexTemplate,
					counter:       c,
				}
				return h
			},
			ctx: servercontext.NewHitCounterContext(echo.New().NewContext(
				httptest.NewRequest("GET", "http://localhost", nil), httptest.NewRecorder())),
			output: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := tc.expect(t)
			output := h.Index(tc.ctx)
			assert.Equal(t, tc.output, output)
		})
	}
}
