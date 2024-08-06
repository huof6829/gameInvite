package sys

import (
	"net/http"

	"github.com/Savvy-Gameing/backend/internal/logic/sys"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TelegramHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys.NewTelegramLogic(r.Context(), svcCtx)
		err := l.Telegram()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
