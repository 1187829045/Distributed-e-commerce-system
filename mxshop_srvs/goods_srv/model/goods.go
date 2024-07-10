package model

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"shop_srvs/goods_srv/global"
)

// 类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
// 实际开发过程中 尽量设置为不为null
// 这些类型我们使用int32还是int

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"` //子目录
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`                      //几级类目
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`                          //是否展现在Tab栏
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

// 商品分类和品牌的关系

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goods_category_brand"
}

//轮播图

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"` //订单编号
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"` //商品简介
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"` //封面图
}

// 同步商品到ES中，降低耦合性
// AfterCreate 是 Goods 结构体的方法，用于在 Goods 对象创建后执行一些操作

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	// 构建 EsGoods 对象，将 Goods 对象的字段映射到 EsGoods 对象
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,      // 设置商品的收藏数量
		MarketPrice: g.MarketPrice, // 设置商品的市场价格
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	// 使用 Elasticsearch 客户端将 esModel 索引到 Elasticsearch 中
	_, err = global.EsClient.Index(). // 创建索引请求
						Index(esModel.GetIndexName()). // 指定索引名称
						BodyJson(esModel).             // 指定索引内容为 esModel 的 JSON 表示
						Id(strconv.Itoa(int(g.ID))).
						Do(context.Background())
	if err != nil { // 如果索引操作出错
		return err // 返回错误
	}

	// 返回 nil 表示操作成功
	return nil
}

//商品的更新
// AfterUpdate 是一个 GORM 钩子函数，在使用 GORM（一个用于 Go 语言的 ORM 库）更新 Goods 对象后会被自动调用。

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

//在使用 GORM（一个用于 Go 语言的 ORM 库）删除Goods 对象后会被自动调用。

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
