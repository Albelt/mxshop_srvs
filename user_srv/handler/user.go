package handler

import (
	pb "albelt.top/mxshop_protos/albelt/user_srv/go"
	"context"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
)

const (
	MaxPageSize = 1000
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func Paginate(page, pageSize uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

func NewPbUser(u *model.User) *pb.User {
	if u == nil {
		return nil
	}

	ret := &pb.User{
		Id:       u.ID,
		Mobile:   u.Mobile,
		NickName: u.NickName,
		Password: u.Password,
		Gender:   pb.Gender(u.Gender),
		Role:     int32(u.Role),
	}

	if u.Birthday != nil {
		ret.Birthday = timestamppb.New(*u.Birthday)
	}

	return ret
}

func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	var (
		users []*model.User
		total int64
		resp  pb.GetUserListResponse
		err   error
	)

	// total
	err = global.Mysql.Find(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// list
	err = global.Mysql.
		Scopes(Paginate(uint64(req.Page), uint64(req.Size))).
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}

	// 返回数据
	resp.Total = int32(total)
	for _, u := range users {
		resp.Users = append(resp.Users, NewPbUser(u))
	}

	return &resp, nil
}

func (s *UserService) GetUserByMobile(ctx context.Context, req *pb.GetUserByMobileRequest) (*pb.User, error) {
	var (
		user   *model.User
		pbUser *pb.User
		err    error
	)

	err = global.Mysql.Where("mobile = ?", req.Mobile).First(&user).Error
	if err != nil {
		return nil, err
	}

	pbUser = NewPbUser(user)
	return pbUser, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	var (
		user   *model.User
		pbUser *pb.User
		err    error
	)

	err = global.Mysql.First(&user, req.Id).Error
	if err != nil {
		return nil, err
	}

	pbUser = NewPbUser(user)
	return pbUser, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	var (
		user *model.User
		err  error
	)

	spew.Dump(user)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		err = tx.Where(&model.User{Mobile: req.Mobile}).First(&user).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.AlreadyExists, "user(mobile=%s) already exists", req.Mobile)
		}

		// 不存在则创建
		user = &model.User{
			Mobile:   req.Mobile,
			NickName: req.NickName,
		}
		user.GenPassword(req.Password)

		err = tx.Create(&user).Error
		if err != nil {
			return err
		}

		// 返回nil提交事务
		return nil
	})

	if err != nil {
		return nil, err
	}

	return NewPbUser(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserResponse) (*pb.User, error) {
	var (
		user *model.User
		err  error
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		_ = tx.Where("id = ?", req.Id).First(&user).Error
		if user == nil {
			return status.Errorf(codes.NotFound, "user(id=%d) not found", req.Id)
		}

		birthday := req.BirthDay.AsTime()
		user.Birthday = &birthday
		user.Gender = uint8(req.Gender)
		user.NickName = req.NickName

		updateMap := map[string]interface{}{
			"birthday":  user.Birthday,
			"gender":    user.Gender,
			"nick_name": user.NickName,
		}

		err = tx.Model(user).Updates(updateMap).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return NewPbUser(user), nil
}

func (s *UserService) CheckPassWord(ctx context.Context, req *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {
	var (
		user *model.User
		ok   bool
	)

	user = &model.User{Password: req.EncryptedPassword}
	ok = user.VerifyPassword(req.Password)

	return &pb.CheckPassWordResponse{Success: ok}, nil
}
