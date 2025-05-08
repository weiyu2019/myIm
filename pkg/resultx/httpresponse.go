package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"myIm/pkg/xerr"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func OKHandle(_ context.Context, data any) any {
	return Success(data)
}

func ErrHandle(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errCode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errCode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zerr.CodeMsg); ok {
			errCode = e.Code
			errmsg = e.Msg
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				errCode = int(gstatus.Code())
				errmsg = gstatus.Message()
			}
		}
		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)
		return http.StatusBadRequest, Fail(errCode, errmsg)
	}

}
