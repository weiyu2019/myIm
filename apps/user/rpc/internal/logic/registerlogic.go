package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"myIm/apps/user/models"
	"myIm/pkg/ctxdata"
	"myIm/pkg/encrypt"
	"myIm/pkg/wuid"
	"myIm/pkg/xerr"
	"time"

	"myIm/apps/user/rpc/internal/svc"
	"myIm/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "phone is registered")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 验证用户是否注册，根据手机号码验证
	userEn, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	if userEn != nil {
		return nil, errors.WithStack(ErrPhoneIsRegistered)
	}

	// 定义用户数据
	userEn = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEn.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEn)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert user err %v", err)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEn.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "get jwt token err %v", err)
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
