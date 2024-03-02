package initial

import (
	"context"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"mxshop_srvs/good_srv/global"
	"mxshop_srvs/good_srv/model"
	"os"
	"strconv"
)

func InitEs() {
	// 初始化es客户端
	esServerUrl := "http://ubuntu-learn:9200"
	logger := log.New(os.Stdout, "[es]: ", log.LstdFlags)
	ctx := context.Background()

	cli, err := elastic.NewClient(
		elastic.SetURL(esServerUrl),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetInfoLog(logger),
	)
	if err != nil {
		panic(err)
	}

	res, _, err := cli.Ping(esServerUrl).Do(ctx)
	if err != nil {
		panic(err)
	}

	global.EsCli = cli
	zap.S().Infof("InitEs ok, es version:%s", res.Version.Number)

	// 创建索引
	goodsIndex := &model.EsGoods{}
	err = createEsIndex(cli, ctx, goodsIndex.IndexName(), goodsIndex.Mapping())
	if err != nil {
		panic(err)
	}

	// 同步数据
	//err = syncMysqlToEs(global.Mysql, global.EsCli, context.Background())
	//if err != nil {
	//	panic(err)
	//}
}

func createEsIndex(cli *elastic.Client, ctx context.Context, indexName string, mapping string) error {
	exists, err := cli.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := cli.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
		zap.S().Infof("createEsIndex %s.", indexName)
	}

	return nil
}

// 同步mysql的goods表到es中
func syncMysqlToEs(mysqlCli *gorm.DB, esCli *elastic.Client, ctx context.Context) error {
	var (
		goods []model.Goods
		err   error
	)

	if err = mysqlCli.Find(&goods).Error; err != nil {
		return err
	}

	for _, g := range goods {
		esG := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandID:     g.BrandID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			ShopPrice:   g.ShopPrice,
			GoodsBrief:  g.GoodsBrief,
		}

		// 使用DB的主键作为es的id，避免数据重复
		_, err = esCli.Index().Index(esG.IndexName()).
			BodyJson(esG).Id(strconv.FormatInt(int64(g.ID), 10)).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
