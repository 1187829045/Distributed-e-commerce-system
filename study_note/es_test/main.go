package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

/*
	对于商品的操作包括

搜索
添加
更新
删除
需要同步到es
service还是web
写在service层，如果我在web层调用商品保存，mysql写成功但是es保存失败，就会导致无法回滚。但是写在mysql事中，就可以避免这个情况。
*/
const goodsmapping = `
{
	"mappings":{
	   "properties":{
			"name":{
				"type":"text",
                "analyzer":"ik_max_word"
            },
            "id":{
               "type":"integer"
            }
       }
	}
}`

type Account struct {
	AccountNumber int32  `json:"account_number"`
	FirstName     string `json:"first_name"`
}

func main() {
	logger := log.New(os.Stdout, "Info: ", log.LstdFlags)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.128.128:9200"), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}
	q := elastic.NewMatchQuery("address", "street")
	//src, err := q.Source()
	//if err != nil {
	//	panic(err)
	//}
	//data, err := json.Marshal(src)
	//if err != nil {
	//	panic(err)
	//}
	//got := string(data)
	//fmt.Println(got)
	//获取数据
	result, err := client.Search().Index("user").Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}
	total := result.Hits.TotalHits.Value
	fmt.Printf("total hits: %d\n", total)
	for _, value := range result.Hits.Hits {
		account := Account{}
		_ = json.Unmarshal(value.Source, &account)
		//if jsonData, err := value.Source.MarshalJSON(); err == nil {
		//	fmt.Println(string(jsonData))
		//} else {
		//	panic(err)
		//}
		fmt.Println(account)
	}

	//保存数据到es
	//account := Account{AccountNumber: 123456, FirstName: "firstname"}
	//put1, err := client.Index().
	//	Index("myuser").
	//	BodyJson(account).
	//	Do(context.Background())
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//fmt.Printf("Indexed account %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	//mapping 的构建
	createIndex, err := client.CreateIndex("mygoods").BodyString(goodsmapping).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}
