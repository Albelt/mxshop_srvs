package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_srvs/good_srv/global"
	pb "albelt.top/mxshop_protos/albelt/good_srv/go"
	"mxshop_srvs/good_srv/model"
)

func NewPbBanner(b *model.Banner) *pb.Banner {
	if b == nil {
		return nil
	}

	return &pb.Banner{
		Id:    b.ID,
		Index: b.Index,
		Image: b.Image,
		Url:   b.Url,
	}
}

func NewModelBanner(b *pb.Banner) *model.Banner {
	if b == nil {
		return nil
	}

	banner := &model.Banner{
		Image: b.Image,
		Url:   b.Url,
		Index: b.Index,
	}
	banner.ID = b.Id

	return banner
}

func (g *GoodsService) BannerList(ctx context.Context, request *pb.BannerListRequest) (*pb.BannerListResponse, error) {
	var (
		err     error
		total   int64
		banners []*model.Banner
		resp    pb.BannerListResponse
	)

	err = global.Mysql.Model(&model.Banner{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = global.Mysql.Find(&banners).Error
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)
	for _, b := range banners {
		resp.Banners = append(resp.Banners, NewPbBanner(b))
	}

	return &resp, nil
}

func (g *GoodsService) CreateBanner(ctx context.Context, request *pb.CreateBannerRequest) (*pb.Banner, error) {
	var (
		err error
	)

	banner := NewModelBanner(request.Banner)

	if err = global.Mysql.Save(&banner).Error; err != nil {
		return nil, err
	}

	pbBanner := NewPbBanner(banner)
	return pbBanner, nil
}

func (g *GoodsService) DeleteBanner(ctx context.Context, request *pb.DeleteBannerRequest) (*emptypb.Empty, error) {
	var (
		err error
	)

	res := global.Mysql.First(&model.Banner{}, request.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "轮播图(id=%d)不存在", request.Id)
	}

	err = global.Mysql.Delete(&model.Banner{}, request.Id).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateBanner(ctx context.Context, request *pb.UpdateBannerRequest) (*emptypb.Empty, error) {
	var (
		err    error
		banner model.Banner
	)

	res := global.Mysql.First(&banner, request.Banner.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图(id=%d)不存在", request.Banner.Id)
	}

	banner.Url = request.Banner.Url
	banner.Index = request.Banner.Index
	banner.Image = request.Banner.Image

	if err = global.Mysql.Save(&banner).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
