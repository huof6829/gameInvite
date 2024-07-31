package invite

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/Savvy-Gameing/backend/common/response"
	"github.com/Savvy-Gameing/backend/internal/logic/invite"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/internal/types"
)

func InviteVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InviteVerifyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		validate := validator.New()
		if err := validate.Struct(&req); err != nil {
			response.HttpResult(r, w, nil, err)
			return
		}

		l := invite.NewInviteVerifyLogic(r.Context(), svcCtx)
		resp, err := l.InviteVerify(&req)
		response.HttpResult(r, w, resp, err)
	}
}
