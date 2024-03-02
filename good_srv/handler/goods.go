package handler

import (
	pb "albelt.top/mxshop_protos/albelt/good_srv/go"
	"context"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mxshop_srvs/good_srv/global"
	"mxshop_srvs/good_srv/model"
	"mxshop_srvs/good_srv/utils"
	"strconv"
)

const (
	MaxPageSize = 1000
)

type GoodsService struct {
	pb.UnimplementedGoodsServiceServer
}

func NewPbGoods(g *model.Goods) *pb.Good {
	if g == nil {
		return nil
	}

	return &pb.Good{
		Id:              g.ID,
		Name:            g.Name,
		GoodsSn:         g.GoodsSn,
		ClickNum:        g.ClickNum,
		SoldNum:         g.SoldNum,
		FavNum:          g.FavNum,
		MarketPrice:     g.MarketPrice,
		ShopPrice:       g.ShopPrice,
		GoodsBrief:      g.GoodsBrief,
		ShipFree:        g.ShipFree,
		GoodsFrontImage: g.GoodsFrontImage,
		IsNew:           g.IsNew,
		IsHot:           g.IsHot,
		OnSale:          g.OnSale,
		CreateTime:      timestamppb.New(g.CreatedAt),
		CategoryId:      g.CategoryID,
		Category:        NewPbCategory(&g.Category),
		BrandId:         g.BrandID,
		Brand:           NewPbBrand(&g.Brand),
	}
}

func NewModelGoods(g *pb.Good) *model.Goods {
	if g == nil {
		return nil
	}

	good := &model.Goods{
		CategoryID:      g.CategoryId,
		Category:        utils.Ptr2Val(NewModelCategory(g.Category)),
		BrandID:         g.BrandId,
		Brand:           utils.Ptr2Val(NewModelBrand(g.Brand)),
		OnSale:          g.OnSale,
		ShipFree:        g.ShipFree,
		IsNew:           g.IsNew,
		IsHot:           g.IsHot,
		Name:            g.Name,
		GoodsSn:         g.GoodsSn,
		ClickNum:        g.ClickNum,
		SoldNum:         g.SoldNum,
		FavNum:          g.FavNum,
		MarketPrice:     g.MarketPrice,
		ShopPrice:       g.ShopPrice,
		GoodsBrief:      g.GoodsBrief,
		GoodsFrontImage: g.GoodsFrontImage,
	}

	good.ID = g.Id

	return good
}

/* 商品查询:
   1. 使用es查询商品的id列表
   2. 使用id列表在mysql中查询具体的业务数据
*/
func (g *GoodsService) GoodsList(ctx context.Context, req *pb.GoodsListRequest) (*pb.GoodsListResponse, error) {

	var (
		err      error
		total    int32
		goodsIds []int32
		goods    []*model.Goods
		resp     pb.GoodsListResponse
	)

	// 组装es查询条件
	q := elastic.NewBoolQuery()
	if req.IsNew {
		q.Filter(elastic.NewTermQuery("is_new", true))
	}
	if req.IsHot {
		q.Filter(elastic.NewTermQuery("is_host", true))
	}
	if req.Keywords != "" {
		q.Must(elastic.NewMultiMatchQuery(req.Keywords, "name", "goods_brief"))
	}
	if req.PriceMin > 0 {
		q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		q.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}

	// 使用商品分类条件
	if req.CategoryId > 0 {
		cateIds, err := getLevel3CategoryIds(ctx, req.CategoryId)
		if err != nil {
			return nil, err
		}

		values := make([]interface{}, 0, len(cateIds))
		for _, id := range cateIds {
			values = append(values, id)
		}

		if len(cateIds) > 0 {
			q.Filter(elastic.NewTermsQuery("category_id", values...))
		}
	}

	// 查询es
	res, err := global.EsCli.Search(model.EsGoods{}.IndexName()).
		Query(q).From(int((req.Page - 1) * req.Count)).Size(int(req.Count)).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 获取total和goodsIds
	total = int32(res.Hits.TotalHits.Value)
	for _, hit := range res.Hits.Hits {
		if id, err := strconv.ParseInt(hit.Id, 10, 64); err == nil {
			goodsIds = append(goodsIds, int32(id))
		} else {
			zap.S().Warnf("EsGoods id = %s, parse failed:%s", hit.Id, err.Error())
		}
	}
	zap.S().Infof("total:%d, goodsIds:%v", total, goodsIds)

	// 根据goodsIds查询mysql获取数据
	err = global.Mysql.Where("id IN ?", goodsIds).Find(&goods).Error

	resp.Total = total
	for _, g := range goods {
		resp.Goods = append(resp.Goods, NewPbGoods(g))
	}
	return &resp, nil
}

// 通过categoryId获取level=3的分类的id列表
func getLevel3CategoryIds(ctx context.Context, categoryId int32) ([]int32, error) {
	const (
		level1Sql = `SELECT id FROM good_srv_categories WHERE parent_category_id IN (SELECT id FROM good_srv_categories WHERE parent_category_id = ?)`
		level2Sql = `SELECT id FROM good_srv_categories WHERE parent_category_id = ?`
	)

	var (
		err      error
		category model.Category
		ids      []int32
	)

	res := global.Mysql.First(&category, categoryId)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", categoryId)
	}

	switch category.Level {
	case 1:
		if err = global.Mysql.Raw(level1Sql, categoryId).Scan(&ids).Error; err != nil {
			return nil, err
		}
	case 2:
		if err = global.Mysql.Raw(level2Sql, categoryId).Scan(&ids).Error; err != nil {
			return nil, err
		}
	case 3:
		ids = append(ids, category.ID)
	default:
		return nil, status.Errorf(codes.Internal, "分类(id=%d)有误", categoryId)
	}

	return ids, nil
}

func (g *GoodsService) BatchGetGoods(ctx context.Context, req *pb.BatchGetGoodsRequest) (*pb.GoodsListResponse, error) {
	var (
		err   error
		goods []*model.Goods
		resp  pb.GoodsListResponse
	)

	err = global.Mysql.Where("id IN ?", req.GoodsIds).Find(&goods).Error
	if err != nil {
		return nil, err
	}

	resp.Total = int32(len(goods))
	for _, g := range goods {
		resp.Goods = append(resp.Goods, NewPbGoods(g))
	}

	return &resp, nil
}

func (g *GoodsService) CreateGoods(ctx context.Context, req *pb.CreateGoodsRequest) (*pb.Good, error) {
	var (
		err error
	)

	good := NewModelGoods(req.Good)

	// 检查分类、品牌
	res := global.Mysql.First(&model.Category{}, good.CategoryID)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "分类(id=%d)不存在", good.CategoryID)
	}

	res = global.Mysql.First(&model.Brand{}, good.BrandID)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌(id=%d)不存在", good.BrandID)
	}

	err = global.Mysql.Save(good).Error
	if err != nil {
		return nil, err
	}

	pbGood := NewPbGoods(good)
	return pbGood, nil
}

func (g *GoodsService) DeleteGoods(ctx context.Context, req *pb.DeleteGoodsRequest) (*emptypb.Empty, error) {
	res := global.Mysql.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品(id=%d)不存在", req.Id)
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsService) UpdateGoods(ctx context.Context, req *pb.UpdateGoodsRequest) (*emptypb.Empty, error) {
	var good model.Goods

	res := global.Mysql.First(&good, req.Good.Id)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品(id=%d)不存在", req.Good.Id)
	}

	good.Name = req.Good.Name
	good.GoodsSn = req.Good.GoodsSn
	good.GoodsBrief = req.Good.GoodsBrief
	good.ShopPrice = req.Good.ShopPrice
	good.MarketPrice = req.Good.MarketPrice

	err := global.Mysql.Save(&good).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *GoodsService) GetGoodsDetail(ctx context.Context, req *pb.GetGoodsDetailRequest) (*pb.Good, error) {
	var (
		good model.Goods
	)

	res := global.Mysql.First(&good, req.Id)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品(id=%d)不存在", req.Id)
	}

	pbGood := NewPbGoods(&good)
	return pbGood, nil
}
