package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mbSrvs/user_srv/global"
	"mbSrvs/user_srv/model"

	"mbSrvs/user_srv/proto"
)

type UserServer struct{}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) (userInfo proto.UserInfoResponse) {
	userInfo = proto.UserInfoResponse{
		Id:       uint64(user.Id),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   uint32(user.Gender),
		Role:     uint32(user.Role),
	}
	if user.Birthday != nil {
		userInfo.Birthday = uint64(user.Birthday.Unix())
	}
	return
}

func (s *UserServer) GetUserList(c context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	zap.S().Info(global.DB)
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = uint32(result.RowsAffected)

	//分页数据
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

func (s *UserServer) GetUserByMobile(c context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.Where("mobile = ?", req.Mobile).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo := ModelToResponse(user)
	return &userInfo, nil
}
func (s *UserServer) GetUserById(c context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.First(&user, req.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo := ModelToResponse(user)
	return &userInfo, nil
}
func (s *UserServer) CreateUser(c context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.Where(map[string]interface{}{"mobile": req.Mobile}).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile

	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, genPwd := password.Encode(req.Password, options)

	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, genPwd)
	user.NickName = req.NickName

	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	userInfo := ModelToResponse(user)

	return &userInfo, nil
}
func (s *UserServer) UpdateUser(c context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.Birthday), 0)

	user.Birthday = &birthDay
	user.NickName = req.NickName
	user.Gender = uint8(req.Gender)
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}
func (s *UserServer) CheckPassword(c context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}

	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Status: check}, nil
}
