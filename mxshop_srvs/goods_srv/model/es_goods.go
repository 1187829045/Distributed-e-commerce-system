package model

// EsGoods 是用于 Elasticsearch 中存储商品信息的结构体
type EsGoods struct {
	ID          int32   `json:"id"`           // 商品ID
	CategoryID  int32   `json:"category_id"`  // 分类ID
	BrandsID    int32   `json:"brands_id"`    // 品牌ID
	OnSale      bool    `json:"on_sale"`      // 是否上架
	ShipFree    bool    `json:"ship_free"`    // 是否包邮
	IsNew       bool    `json:"is_new"`       // 是否新品
	IsHot       bool    `json:"is_hot"`       // 是否热销
	Name        string  `json:"name"`         // 商品名称
	ClickNum    int32   `json:"click_num"`    // 点击数
	SoldNum     int32   `json:"sold_num"`     // 销量
	FavNum      int32   `json:"fav_num"`      // 收藏数
	MarketPrice float32 `json:"market_price"` // 市场价
	GoodsBrief  string  `json:"goods_brief"`  // 商品简介
	ShopPrice   float32 `json:"shop_price"`   // 商城价
}

// GetIndexName 返回 Elasticsearch 中索引的名称
func (EsGoods) GetIndexName() string {
	return "goods"
}

// GetMapping 返回 Elasticsearch 中的映射（Mapping）配置，定义了字段的数据类型和分析器等信息
func (EsGoods) GetMapping() string {
	goodsMapping := `
	{
		"mappings" : {
			"properties" : {
				"brands_id" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"fav_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"ship_free" : {
					"type" : "boolean"
				},
				"shop_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
	return goodsMapping
}
