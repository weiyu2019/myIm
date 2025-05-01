package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"myIm/apps/user/models"

	"myIm/apps/user/rpc/internal/svc"
	"myIm/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		userEntitys []*models.Users
		err         error
	)
	if len(in.Phone) > 0 {
		userEn, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEn)
		}
	} else if len(in.Name) > 0 {
		userEntitys, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}
	if err != nil {
		return nil, err
	}
	var resp []*user.UserEntity
	copier.Copy(&resp, userEntitys)
	return &user.FindUserResp{
		User: resp,
	}, nil
}
