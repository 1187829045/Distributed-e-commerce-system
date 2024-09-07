package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"shop_srvs/inventory_srv/model"

	"shop_srvs/inventory_srv/global"
	"shop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存，支持本地事务。此处示例包含多个商品 [1:10,  2:5, 3:20]

	// 并发情况下，可能会出现超卖现象，因此使用分布式锁来防止这种情况。

	// 创建 Redis 客户端
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.128.128:6379", // Redis 服务器地址
	})

	// 创建 Redis 连接池
	pool := goredis.NewPool(client) // 也可以使用 redigo.NewPool(client)

	// 创建一个 Redsync 对象（rs），用于实现分布式锁
	rs := redsync.New(pool)

	// 开启数据库事务
	tx := global.DB.Begin()

	// 在扣减库存前，应该先查询订单是否已经扣减过库存，以防止重复扣减
	// 在并发环境下，这种防止重复请求的机制尤其重要，这里使用分布式锁来防止
	var existingSellDetail model.StockSellDetail
	if result := tx.Where("order_sn = ?", req.OrderSn).First(&existingSellDetail); result.RowsAffected > 0 {
		// 如果订单号已存在，说明已经处理过，直接返回
		return &emptypb.Empty{}, nil
	}
	// 创建库存扣减记录对象
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn, // 订单号
		Status:  1,           // 扣减状态，1表示扣减中
	}

	var details []model.GoodsDetail // 存储订单中的商品信息
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId, // 商品ID
			Num:   goodInfo.Num,     // 商品数量
		})

		var inv model.Inventory

		// 创建分布式锁，锁定当前商品的库存
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取Redis分布式锁异常") // 获取锁失败
		}

		// 查询商品库存信息
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 如果查询不到库存信息，回滚事务
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		// 检查库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 库存不足，回滚事务
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减库存
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv) // 保存库存变化

		// 释放分布式锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放Redis分布式锁异常") // 释放锁失败
		}
	}

	// 将订单详情写入库存扣减历史记录表
	sellDetail.Detail = details
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback() // 如果写入失败，回滚事务
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}

	// 提交事务，保存所有更改
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) TrySell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.128.128:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		//inv.Stocks -= goodInfo.Num
		inv.Freeze += goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) ConfirmSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.128.128:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) CancelSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.128.128:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?",
		//goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

type OrderInfo struct {
	OrderSn string
}

// 自动归还
// 既然是归还库存，那么我应该具体的知道每件商品应该归还多少， 但是有一个问题是什么？重复归还的问题
// 所以说这个接口应该确保幂等性， 你不能因为消息的重复发送导致一个订单的库存归还多次， 没有扣减的库存你别归还
// 如果确保这些都没有问题， 新建一张表， 这张表记录了详细的订单扣减细节，以及归还细节

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// 遍历接收到的消息
	for i := range msgs {
		var orderInfo OrderInfo
		// 将消息体中的 JSON 数据反序列化为 OrderInfo 结构体
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body) // 记录 JSON 解析错误日志
			return consumer.ConsumeSuccess, nil            // 返回消费成功，终止当前消息的处理
		}

		// 开启数据库事务
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		// 查询数据库中状态为 1 的订单销售明细
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).
			First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil // 如果未找到匹配记录，返回消费成功
		}

		// 如果查询到订单销售明细，则逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			// 更新库存，将库存数量增加对应的商品数量
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).
				Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()                          // 如果更新库存失败，回滚事务
				return consumer.ConsumeRetryLater, nil // 返回重试消费的信号
			}
		}

		// 更新订单销售明细的状态为 2（已归还库存）
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).
			Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()                          // 如果更新状态失败，回滚事务
			return consumer.ConsumeRetryLater, nil // 返回重试消费的信号
		}

		// 提交事务
		tx.Commit()
		return consumer.ConsumeSuccess, nil // 返回消费成功
	}
	return consumer.ConsumeSuccess, nil // 如果没有消息需要处理，返回消费成功
}
