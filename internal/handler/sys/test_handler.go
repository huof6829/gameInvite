package sys

import (
	"net/http"

	"github.com/Savvy-Gameing/backend/internal/logic/sys"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys.NewTestLogic(r.Context(), svcCtx)
		resp, err := l.Test()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
