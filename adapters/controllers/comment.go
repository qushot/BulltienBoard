package controllers

import (
	"github.com/osechiman/BulltienBoard/adapters/gateways"
	"github.com/osechiman/BulltienBoard/entities"
	"github.com/osechiman/BulltienBoard/entities/errorobjects"
	"github.com/osechiman/BulltienBoard/entities/valueobjects"
	"github.com/osechiman/BulltienBoard/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// CommentController はCommentRepositorerのコントローラーです。
// 初期化時に渡すリポジトリ以外を利用したい場合はそれぞれメソッドの引数で受け取ってください。
type CommentController struct {
	// Repository はこのコントローラーで利用するメインのリポジトリです。
	// このコントローラーで利用するメインのリポジトリです。
	Repository usecases.CommentRepositorer
}

// Comment はリクエストされてきたPost値を受け取る為のStructです。
type Comment struct {
	ID       string // ID はCommentのIDです。
	ThreadID string `validate:"required"` // ThreadIDer はThreadのIDです。
	Text     string `validate:"required"` // Text はユーザーが入力した文字列です。
}

// NewCommentController はCommentControllerを初期化します。
func NewCommentController(r usecases.CommentRepositorer) *CommentController {
	return &CommentController{Repository: r}
}

// AddComment はPostされてきたデータを元にCommentを追加します。
// コマンド・クエリの原則からは外れますがAPIのレスポンスに登録したデータを返却するためにエンティティを返します。
func (cc *CommentController) AddComment(c *gin.Context) (entities.Comment, error) {
	pc := Comment{}
	err := c.BindJSON(&pc)
	if err != nil {
		return entities.Comment{}, errorobjects.NewParameterBindingError(err)
	}

	validate := validator.New()
	err = validate.Struct(pc)
	if err != nil {
		return entities.Comment{}, errorobjects.NewMissingRequiredFieldsError(err)
	}

	cid, err := valueobjects.NewCommentID("")
	if err != nil {
		return entities.Comment{}, err
	}

	tid, err := convertIDToThreadID(pc.ThreadID)
	if err != nil {
		return entities.Comment{}, err
	}

	ct, err := valueobjects.NewCommentTime(-1)
	if err != nil {
		return entities.Comment{}, err
	}

	cm, err := entities.NewComment(cid, tid, pc.Text, ct)
	if err != nil {
		return entities.Comment{}, err
	}

	cu := usecases.NewCommentUsecase(cc.Repository)
	return cm, cu.AddComment(cm, gateways.GetInMemoryRepositoryInstance())
}

// ListComment はCommentの一覧を取得します。
func (cc *CommentController) ListComment() ([]entities.Comment, error) {
	tu := usecases.NewCommentUsecase(cc.Repository)
	return tu.ListComment()
}

// ListCommentByThreadID は指定されたThreadIDを持つComment一覧を取得します。
func (cc *CommentController) ListCommentByThreadID(tID string) ([]entities.Comment, error) {
	tid, err := convertIDToThreadID(tID)
	if err != nil {
		return nil, err
	}
	tu := usecases.NewCommentUsecase(cc.Repository)
	return tu.ListCommentByThreadID(tid)
}

// convertIDToCommentID は文字列のCommentIDをentities.CommentIDに変換します。
func convertIDToCommentID(ID string) (valueobjects.CommentID, error) {
	return valueobjects.NewCommentID(ID)
}

// convertCreatAtToCommentTime は渡された数値をentities.CommentTimeに変換します。
func convertCreatAtToCommentTime(unixTime int64) (valueobjects.CommentTime, error) {
	return valueobjects.NewCommentTime(unixTime)
}
