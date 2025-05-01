package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"myIm/apps/user/rpc/internal/config"
	"myIm/apps/user/rpc/internal/svc"
	"myIm/apps/user/rpc/user"
	"path/filepath"
	"reflect"
	"testing"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}

func TestLoginLogic_Login(t *testing.T) {
	type fields struct {
		ctx    context.Context
		svcCtx *svc.ServiceContext
		Logger logx.Logger
	}
	type args struct {
		in *user.LoginReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.LoginResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoginLogic{
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
				Logger: tt.fields.Logger,
			}
			got, err := l.Login(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}
