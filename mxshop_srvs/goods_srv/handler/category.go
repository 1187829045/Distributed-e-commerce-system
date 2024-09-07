package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/goods_srv/model"

	"google.golang.org/protobuf/types/known/emptypb"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/proto"
)

// 商品目录列表

func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// 获取子类目
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 创建一个 SubCategoryListResponse 对象，用来存储响应数据
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	// 尝试从数据库中获取指定 ID 的分类信息，如果没有找到，则返回错误
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	// 将获取到的分类信息填充到 CategoryInfoResponse 对象中，并赋值给响应的 Info 字段
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategorys []model.Category
	// 查询当前分类下的所有子分类，条件是 ParentCategoryID 等于请求的分类 ID
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)

	var subCategoryResponse []*proto.CategoryInfoResponse
	// 将每个子分类的信息转换为 CategoryInfoResponse 并添加到 subCategoryResponse 列表中
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	// 将子分类信息列表赋值给响应的 SubCategorys 字段
	categoryListResponse.SubCategorys = subCategoryResponse
	// 返回完整的响应对象
	return &categoryListResponse, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab
	category.ParentCategoryID = req.ParentCategory
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}
	global.DB.Save(&category)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
