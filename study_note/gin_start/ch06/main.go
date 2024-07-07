package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sale_master/study_note/gin_start/ch06/proto"
)

func main() {
	router := gin.Default()

	router.GET("/moreJSON", moreJSON)
	router.GET("/someProtoBuf", returnProto)

	router.Run(":8083")
}

func returnProto(c *gin.Context) {
	course := []string{"python", "go", "微服务"}
	user := &proto.Teacher{
		Name:   "bobby",
		Course: course,
	}
	c.ProtoBuf(http.StatusOK, user)
}

func moreJSON(c *gin.Context) {
	var msg struct {
		//定义了一个名为 Name 的结构体字段，并使用了结构体标签来指定在进行JSON序列化和反序列化时，这个字段应该映射到JSON对象中的 user 字段。
		Name    string `json:"user"`
		Message string
		Number  int
	}
	msg.Name = "bobby"
	msg.Message = "这是一个测试json"
	msg.Number = 20

	c.JSON(http.StatusOK, msg)
}
