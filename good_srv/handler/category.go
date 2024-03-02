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
	"mxshop_srvs/good_srv/utils"
)

func NewPbCategory(c *model.Category) *pb.Category {
	if c == nil {
		return nil
	}

	return &pb.Category{
		Id:               c.ID,
		Name:             c.Name,
		ParentCategoryId: utils.Ptr2Val(c.ParentCategoryID),
		Level:            c.Level,
		IsTab:            c.IsTab,
	}
}

func NewModelCategory(c *pb.Category) *model.Category {
	if c == nil {
		return nil
	}

	category := &model.Category{
		Name:  c.Name,
		Level: c.Level,
		IsTab: c.IsTab,
	}

	category.ID = c.Id
	if c.ParentCategoryId != 0 {
		category.ParentCategoryID = &c.ParentCategoryId
	}

	return category
}

func (g *GoodsService) GetAllCategorysList(ctx context.Context, request *pb.GetAllCategorysListRequest) (*pb.GetAllCategorysListResponse, error) {
	var (
		err        error
		total      int64
		categories []*model.Category
		resp       pb.GetAllCategorysListResponse
	)

	err = global.Mysql.Model(&model.Category{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = global.Mysql.Model(&model.Category{}).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)
	for _, c := range categories {
		resp.Categories = append(resp.Categories, NewPbCategory(c))
	}

	return &resp, nil
}

func (g *GoodsService) GetSubCategory(ctx context.Context, request *pb.GetSubCategoryRequest) (*pb.GetSubCategoryResponse, error) {
	var (
		err      error
		category model.Category
		subCates []*model.Category
		resp     pb.GetSubCategoryResponse
	)

	res := global.Mysql.First(&category, request.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", request.Id)
	}

	err = global.Mysql.Where(&model.Category{ParentCategoryID: &category.ID}).Find(&subCates).Error
	if err != nil {
		return nil, err
	}

	resp.Category = NewPbCategory(&category)
	for _, c := range subCates {
		resp.SubCategories = append(resp.SubCategories, NewPbCategory(c))
	}

	return &resp, nil
}

func (g *GoodsService) CreateCategory(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.Category, error) {
	var (
		err      error
		category *model.Category
	)

	category = NewModelCategory(request.Category)

	err = global.Mysql.Create(category).Error
	if err != nil {
		return nil, err
	}

	pbCat := NewPbCategory(category)
	return pbCat, nil
}

func (g *GoodsService) DeleteCategory(ctx context.Context, request *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	var (
		err error
	)

	res := global.Mysql.Delete(&model.Category{}, request.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", request.Id)
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateCategory(ctx context.Context, request *pb.UpdateCategoryRequest) (*emptypb.Empty, error) {
	var (
		err      error
		category model.Category
	)

	res := global.Mysql.First(&category, request.Category.Id)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类(id=%d)不存在", request.Category.Id)
	}

	if request.Category.ParentCategoryId != 0 {
		category.ParentCategoryID = &request.Category.ParentCategoryId
	}
	category.Name = request.Category.Name
	category.IsTab = request.Category.IsTab
	category.Level = request.Category.Level

	if err = global.Mysql.Save(&category).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
