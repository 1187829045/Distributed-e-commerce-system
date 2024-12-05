package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/model"
	"shop_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// 生成订单SN号

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

// 获取用户的购物车列表信息

func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}

// 添加商品到购物车

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车- 合并
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		//如果记录已经存在，则合并购物车记录, 更新操作
		shopCart.Nums += req.Nums
	} else {
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	//更新购物车记录，更新数量和选中状态
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).
		Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

//对用户订单的分页查询，并将结果封装成 OrderListResponse 结构体返回给客户端。

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 声明一个订单信息的切片，用于存储查询结果
	var orders []model.OrderInfo

	// 声明一个订单列表响应结构体，用于存储返回给客户端的数据
	var rsp proto.OrderListResponse

	// 查询符合条件的订单总数
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	// 将总数赋值给响应中的 Total 字段
	rsp.Total = int32(total)

	// 分页查询
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)

	// 将查询到的订单信息转换为响应格式，并添加到响应的 Data 切片中
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 返回响应结构体
	return &rsp, nil
}
func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// 声明一个订单信息结构体，用于存储查询到的订单数据
	var order model.OrderInfo

	// 声明一个订单详情响应结构体，用于返回给客户端的数据
	var rsp proto.OrderInfoDetailResponse

	// 查询订单信息，确保该订单属于请求中的用户
	// 这里通过订单的ID和用户的ID进行查询，以确保只有当前用户的订单才能被查询到
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		// 如果查询不到对应的订单，返回 "订单不存在" 错误
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	// 将查询到的订单信息填充到 OrderInfoResponse 结构体中
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID               // 订单ID
	orderInfo.UserId = order.User         // 用户ID
	orderInfo.OrderSn = order.OrderSn     // 订单编号
	orderInfo.PayType = order.PayType     // 支付方式
	orderInfo.Status = order.Status       // 订单状态
	orderInfo.Post = order.Post           // 邮寄方式
	orderInfo.Total = order.OrderMount    // 订单总金额
	orderInfo.Address = order.Address     // 收货地址
	orderInfo.Name = order.SignerName     // 收货人姓名
	orderInfo.Mobile = order.SingerMobile // 收货人手机

	// 将订单信息添加到响应结构体中
	rsp.OrderInfo = &orderInfo

	// 查询订单对应的商品信息
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		// 如果查询订单商品信息失败，返回错误
		return nil, result.Error
	}

	// 遍历查询到的订单商品，将每个商品的信息填充到响应中
	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,      // 商品ID
			GoodsName:  orderGood.GoodsName,  // 商品名称（可以通过商品ID查询跨服务获取）
			GoodsPrice: orderGood.GoodsPrice, // 商品价格
			GoodsImage: orderGood.GoodsImage, // 商品图片
			Nums:       orderGood.Nums,       // 商品数量
		})
	}

	// 返回包含订单详情和商品信息的响应结构体
	return &rsp, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// 反序列化消息体为 OrderInfo 对象
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	// 从上下文中获取父Span
	parentSpan := opentracing.SpanFromContext(o.Ctx)

	// 初始化商品ID数组和购物车对象数组，创建商品数量映射
	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	// 开始一个新的Span用于选择购物车中的商品
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	// 从数据库中查询选中结算的购物车商品
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		// 如果没有选中结算的商品，返回回滚状态
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	// 遍历购物车商品，收集商品ID和数量
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 跨服务调用商品微服务批量查询商品信息
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		// 查询商品信息失败，返回回滚状态
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()

	// 计算订单总金额并准备订单商品信息和库存信息
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 跨服务调用库存微服务扣减库存
	queryInvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))

	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{OrderSn: orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		// 扣减库存失败，返回回滚状态
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	queryInvSpan.Finish()

	// 开始数据库事务，保存订单信息
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		// 保存订单失败，回滚事务并返回提交状态
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	// 批量插入订单商品
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_order_goods", opentracing.
		ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		// 插入订单商品失败，回滚事务并返回提交状态
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	// 删除购物车中的已结算商品
	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart",
		opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		// 删除购物车记录失败，回滚事务并返回提交状态
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	deleteShopCartSpan.Finish()

	// 发送延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.128.128:9876"}))
	if err != nil {
		// 生成producer失败，抛出异常
		panic("生成producer失败")
	}

	// 启动producer
	if err = p.Start(); err != nil {
		// 启动producer失败，抛出异常
		panic("启动producer失败")
	}

	// 创建新的延时消息
	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	// 发送同步延时消息
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		// 发送延时消息失败，回滚事务并记录错误
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	// 提交事务
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	// 1. 检查订单是否存在
	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		// 如果订单不存在，意味着本地事务没有成功提交，返回 Rollback 状态
		return primitive.CommitMessageState
	}

	// 2. 检查订单状态是否已更新为“已支付”或“已创建”等预期状态
	if orderInfo.Status != "TRADE_SUCCESS" || orderInfo.Status != "WAIT_BUYER_PAY" {
		return primitive.CommitMessageState
	}
	// 3. 检查订单的商品是否成功扣减库存
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: orderInfo.ID}).Find(&orderGoods); result.RowsAffected == 0 {
		// 如果订单商品记录不存在，可能意味着库存扣减未成功
		return primitive.CommitMessageState
	}

	// 4. 其他可能的检查，例如库存微服务确认，或其他相关表的状态
	// 在这里你可以加入跨服务调用或数据库检查，以确认库存扣减操作已成功完成

	// 如果所有检查通过，返回 Commit 状态，表示本地事务成功
	return primitive.RollbackMessageState
}

// 创建订单

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	//商品的金额不能由前端传递，万一有爬虫就可能以很低的价格购买，从数据库中查询
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/
	// 创建一个订单监听器实例，用于事务的回调
	orderListener := OrderListener{Ctx: ctx}
	// 创建一个RocketMQ事务生产者
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{"192.168.128.128:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败: %s", err.Error())
		return nil, err
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败: %s", err.Error())
		return nil, err
	}

	// 生成订单信息
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId), // 生成订单编号
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	// 将订单信息序列化为JSON字符串
	jsonString, _ := json.Marshal(order)

	// 发送事务消息到RocketMQ中，消息主题为`order_reback`
	_, err = p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}

	// 检查订单监听器中的状态码，如果不为OK，则表示事务失败
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	// 返回订单信息响应，包含订单ID、订单编号以及订单总金额
	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: order.OrderSn, Total: orderListener.OrderAmount}, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	// 更新订单状态，首先查询订单编号对应的订单，然后更新状态
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).
		Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		//查询订单的支付状态，如果已支付什么都不做，如果未支付，归还库存
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			//归还库存，我们可以模仿order中发送一个消息到 order_reback中去
			//修改订单的状态为已支付
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.128.128:9876"}))
			if err != nil {
				panic("生成producer失败")
			}

			if err = p.Start(); err != nil {
				panic("启动producer失败")
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}

			//if err = p.Shutdown(); err != nil {panic("关闭producer失败")}
			return consumer.ConsumeSuccess, nil
		}
	}
	return consumer.ConsumeSuccess, nil
}
