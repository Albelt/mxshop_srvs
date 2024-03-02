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

func NewPbCategoryBrand(cb *model.CategoryBrandRelation) *pb.CategoryBrand {
	if cb == nil {
		return nil
	}

	return &pb.CategoryBrand{
		Id:         cb.ID,
		BrandId:    cb.BrandsID,
		Brand:      NewPbBrand(&cb.Brands),
		CategoryId: cb.CategoryID,
		Category:   NewPbCategory(&cb.Category),
	}
}

func NewModelCategoryBrand(cb *pb.CategoryBrand) *model.CategoryBrandRelation {
	if cb == nil {
		return nil
	}

	categoryBrand := &model.CategoryBrandRelation{
		CategoryID: cb.CategoryId,
		BrandsID:   cb.BrandId,
	}
	categoryBrand.ID = cb.Id
	if cb.Brand != nil {
		categoryBrand.Brands = *NewModelBrand(cb.Brand)
	}
	if cb.Category != nil {
		categoryBrand.Category = *NewModelCategory(cb.Category)
	}

	return categoryBrand
}

func (g *GoodsService) CategoryBrandList(ctx context.Context, request *pb.CategoryBrandListRequest) (*pb.CategoryBrandListResponse, error) {
	var (
		err            error
		total          int64
		categoryBrands []*model.CategoryBrandRelation
		resp           pb.CategoryBrandListResponse
	)

	err = global.Mysql.Model(&model.CategoryBrandRelation{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = global.Mysql.
		Scopes(Paginate(uint64(request.Page), uint64(request.Count))).
		Preload("Category").Preload("Brands").
		Find(&categoryBrands).
		Error
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)
	for _, cb := range categoryBrands {
		resp.CategoryBrands = append(resp.CategoryBrands, NewPbCategoryBrand(cb))
	}

	return &resp, nil
}

func (g *GoodsService) GetBrandsOfCategory(ctx context.Context, request *pb.GetBrandsOfCategoryRequest) (*pb.GetBrandsOfCategoryResponse, error) {
	var (
		err            error
		category       model.Category
		categoryBrands []*model.CategoryBrandRelation
		resp           pb.GetBrandsOfCategoryResponse
	)

	res := global.Mysql.First(&category, request.CategoryId)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", request.CategoryId)
	}

	err = global.Mysql.
		Where(&model.CategoryBrandRelation{CategoryID: request.CategoryId}).
		Preload("Brands").
		Find(&categoryBrands).Error
	if err != nil {
		return nil, err
	}

	resp.Category = NewPbCategory(&category)
	for _, cb := range categoryBrands {
		resp.Brands = append(resp.Brands, NewPbBrand(&cb.Brands))
	}

	return &resp, nil
}

func (g *GoodsService) CreateCategoryBrand(ctx context.Context, request *pb.CreateCategoryBrandRequest) (*pb.CategoryBrand, error) {
	var (
		err error
	)

	categoryId := request.CategoryBrand.CategoryId
	res := global.Mysql.First(&model.Category{}, categoryId)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", categoryId)
	}

	brandId := request.CategoryBrand.BrandId
	res = global.Mysql.First(&model.Brand{}, brandId)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌(id=%d)不存在", brandId)
	}

	categoryBrand := NewModelCategoryBrand(request.CategoryBrand)
	err = global.Mysql.Create(categoryBrand).Error
	if err != nil {
		return nil, err
	}

	err = global.Mysql.
		Preload("Category").Preload("Brands").
		First(&categoryBrand, categoryBrand.ID).Error
	if err != nil {
		return nil, err
	}

	pbCategoryBrand := NewPbCategoryBrand(categoryBrand)
	return pbCategoryBrand, nil
}

func (g *GoodsService) DeleteCategoryBrand(ctx context.Context, request *pb.DeleteCategoryBrandRequest) (*emptypb.Empty, error) {
	res := global.Mysql.Delete(&model.CategoryBrandRelation{}, request.Id)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类(id=%d)不存在", request.Id)
	}

	return &emptypb.Empty{}, nil
}
