package model

// EsGoods 是用于 Elasticsearch 中存储商品信息的结构体
type EsGoods struct {
	ID          int32   `json:"id"`          // 商品ID
	CategoryID  int32   `json:"category_id"` // 分类ID
	BrandsID    int32   `json:"brands_id"`
	OnSale      bool    `json:"on_sale"`
	ShipFree    bool    `json:"ship_free"`
	IsNew       bool    `json:"is_new"`
	IsHot       bool    `json:"is_hot"`
	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

// GetIndexName 返回 Elasticsearch 中索引的名称
func (EsGoods) GetIndexName() string {
	return "goods"
}

// 返回 Elasticsearch 中的映射（Mapping）配置，定义了字段的数据类型和分析器等信息

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
