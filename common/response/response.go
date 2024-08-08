package response

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		httpx.OkJsonCtx(r.Context(), w, Response{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	} else {
		logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)

		switch e := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if e.Number == uint16(ResourceExisted.GetErrCode()) {
				httpx.OkJsonCtx(r.Context(), w, Response{
					Code:    int(e.Number),
					Message: ResourceExisted.GetErrMsg(),
					Data:    nil,
				})
				return
			} else {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, Response{
					Code:    int(e.Number),
					Message: e.Error(),
					Data:    nil,
				})
				return
			}

		case validator.ValidationErrors:
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, Response{
				Code:    ParameterErr.GetErrCode(),
				Message: ParameterErr.GetErrMsg(),
				Data:    nil,
			})
			return

		case *CodeError:
			// statusHttp := 200
			// switch e.GetErrCode() {
			// case ParameterErr.GetErrCode(), JWTErr.GetErrCode(), DBErr.GetErrCode(), RedisErr.GetErrCode(), ServerErr.GetErrCode():
			// 	statusHttp = http.StatusBadRequest
			// }

			httpx.WriteJsonCtx(r.Context(), w, 200, Response{
				Code:    e.GetErrCode(),
				Message: e.GetErrMsg(),
				Data:    nil,
			})

		default:
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, Response{
				Code:    ServerErr.GetErrCode(),
				Message: e.Error(),
				Data:    nil,
			})
		}
	}
}
