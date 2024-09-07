package initialize

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
)

func InitEs() {
	// 初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "shop", log.LstdFlags)
	var err error
	// 使用 `elastic.NewClient` 初始化一个 Elasticsearch 客户端。
	global.EsClient, err = elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}
	// 检查 Elasticsearch 中是否存在指定的索引（index），此处的索引名称通过 `model.EsGoods{}.GetIndexName()` 获取。
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		// 如果索引不存在，则创建一个新的索引。
		_, err = global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).
			Do(context.Background())

		if err != nil {
			panic(err)
		}
	}
}
