package lang

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zhuravlovO/KPI-APZ-lab-03/painter"
)

func HttpHandler(loop *painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		cmds, err := p.Parse(in)
		if err != nil {
			http.Error(rw, "Bad script: "+err.Error(), http.StatusBadRequest)
			return
		}
		finalOps := []painter.Operation{painter.ResetOp}
		finalOps = append(finalOps, cmds...)

		fmt.Printf("Parser created %d operations, total with reset: %d\n", len(cmds), len(finalOps))
		loop.Post(painter.OperationList(finalOps))
		rw.WriteHeader(http.StatusOK)
	})
}
