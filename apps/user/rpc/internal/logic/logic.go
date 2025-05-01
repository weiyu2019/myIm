package logic

import (
	"context"
	"errors"
	"myIm/apps/user/models"
	"myIm/pkg/ctxdata"
	"myIm/pkg/encrypt"
	"time"

	"myIm/apps/user/rpc/internal/svc"
	"myIm/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegistered = errors.New("phone not registered")
	ErrUserPassWrong      = errors.New("user password wrong")
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
			return nil, ErrPhoneNotRegistered
		}
		return nil, err
	}

	//密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEn.Password.String) {
		return nil, ErrUserPassWrong
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEn.Id)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
