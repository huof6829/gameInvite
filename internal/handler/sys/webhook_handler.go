package sys

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/Savvy-Gameing/backend/internal/logic/sys"
	"github.com/Savvy-Gameing/backend/internal/svc"
)

func WebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys.NewWebhookLogic(r.Context(), svcCtx)

		err := l.Webhook(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
