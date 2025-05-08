package logic

import (
	"context"
	"github.com/pkg/errors"
	"myIm/apps/user/models"
	"myIm/pkg/ctxdata"
	"myIm/pkg/encrypt"
	"myIm/pkg/xerr"
	"time"

	"myIm/apps/user/rpc/internal/svc"
	"myIm/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "phone not registered")
	ErrUserPassWrong      = xerr.New(xerr.SERVER_COMMON_ERROR, "user password wrong")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证用户是否注册，根据手机号码验证
	userEn, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegistered)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	//密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEn.Password.String) {
		return nil, errors.WithStack(ErrUserPassWrong)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEn.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get jwt token err %v", err)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
