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

func NewPbBrand(b *model.Brand) *pb.Brand {
	if b == nil {
		return nil
	}

	return &pb.Brand{
		Id:   b.ID,
		Name: b.Name,
		Logo: b.Logo,
	}
}

func NewModelBrand(b *pb.Brand) *model.Brand {
	if b == nil {
		return nil
	}

	brand := model.Brand{
		Name: b.Name,
		Logo: b.Logo,
	}
	brand.ID = b.Id

	return &brand
}

func (g *GoodsService) BrandList(ctx context.Context, request *pb.BrandListRequest) (*pb.BrandListResponse, error) {
	var (
		err    error
		total  int64
		brands []*model.Brand
		resp   pb.BrandListResponse
	)

	err = global.Mysql.Model(&model.Brand{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = global.Mysql.
		Scopes(Paginate(uint64(request.Page), uint64(request.Count))).
		Find(&brands).
		Error
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)
	for _, b := range brands {
		resp.Brands = append(resp.Brands, NewPbBrand(b))
	}

	return &resp, nil
}

func (g *GoodsService) CreateBrand(ctx context.Context, request *pb.CreateBrandRequest) (*pb.Brand, error) {
	var (
		err error
	)

	// 根据名称查找品牌是否存在
	res := global.Mysql.
		Where(&model.Brand{Name: request.Brand.Name}).
		First(&model.Brand{})
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已经存在")
	}

	// 创建
	brand := NewModelBrand(request.Brand)
	if err = global.Mysql.Create(brand).Error; err != nil {
		return nil, err
	}

	pbBrand := NewPbBrand(brand)
	return pbBrand, nil
}

func (g *GoodsService) DeleteBrand(ctx context.Context, request *pb.DeleteBrandRequest) (*emptypb.Empty, error) {
	var (
		err error
	)

	res := global.Mysql.First(&model.Brand{}, request.BrandId)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌(id=%d)不存在", request.BrandId)
	}

	if err = global.Mysql.Delete(&model.Brand{}, request.BrandId).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateBrand(ctx context.Context, request *pb.UpdateBrandRequest) (*emptypb.Empty, error) {
	var (
		err   error
		brand model.Brand
	)

	res := global.Mysql.First(&brand, request.Brand.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌(id=%d)不存在", request.Brand.Id)
	}

	brand.Name = request.Brand.Name
	brand.Logo = request.Brand.Logo

	if err = global.Mysql.Save(&brand).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
