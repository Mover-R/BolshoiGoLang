package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type x2Request struct {
	Val int
}

type x2Response struct {
	Val int
}

func main() {
	engine := gin.New()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	engine.POST("x2", func(ctx *gin.Context) {
		var x2req x2Request
		if err := ctx.BindJSON(&x2req); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		resp := x2Response{
			Val: x2req.Val * 2,
		}

		ctx.JSON(http.StatusOK, resp)
	})

	engine.Run(":8090")
}
