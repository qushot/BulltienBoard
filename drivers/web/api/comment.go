package api

import (
	"net/http"

	"github.com/osechiman/BulltienBoard/adapters/controllers"
	"github.com/osechiman/BulltienBoard/adapters/gateways"
	"github.com/osechiman/BulltienBoard/adapters/presenters"

	"github.com/gin-gonic/gin"
)

// postComment はPostされてきたComment(json)を保存します。
func postComment(c *gin.Context) {
	cr := gateways.GetInMemoryRepositoryInstance()
	cp := presenters.NewCommentPresenter()
	cc := controllers.NewCommentController(cr)
	cm, err := cc.AddComment(c)
	if err != nil {
		responseByError(c, err)
		return
	}

	res := cp.ConvertToHttpCommentResponse(cm)
	c.JSON(http.StatusCreated, res)
	return
}

// listComment はCommentの一覧をjsonで出力します。
func listComment(c *gin.Context) {
	r := gateways.GetInMemoryRepositoryInstance()
	cp := presenters.NewCommentPresenter()
	cc := controllers.NewCommentController(r)
	cl, err := cc.ListComment()
	if err != nil {
		responseByError(c, err)
		return
	}

	res := cp.ConvertToHttpCommentListResponse(cl)
	c.JSON(http.StatusOK, res)
	return
}
