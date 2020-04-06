package usecases

import (
	"vspro/entities"
	"vspro/entities/errorobjects"
	"vspro/entities/valueobjects"
)

// ThreadUsecase はThreadに対するUsecaseを定義するものです。
type ThreadUsecase struct {
	Repository ThreadRepositorer // Repositorer は外部データソースに存在するentities.Threadを操作する際に利用するインターフェースです。
}

// NewThreadUsecase はThreadUsecaseを初期化します。
func NewThreadUsecase(r ThreadRepositorer) *ThreadUsecase {
	return &ThreadUsecase{Repository: r}
}

// GetThreadByID は指定されたvalueobjects.ThreadIDを持つentities.Threadを取得します。
func (tu *ThreadUsecase) GetThreadByID(ID valueobjects.ThreadID, commentRepository CommentRepositorer) (entities.Thread, error) {
	cl, err := commentRepository.ListCommentByThreadID(ID)
	if err != nil {
		switch err.(type) {
		case *errorobjects.NotFoundError:
			cl = make([]entities.Comment, 0)
		default:
			return entities.Thread{}, err
		}
	}

	t, err := tu.Repository.GetThreadByID(ID)
	if err != nil {
		return entities.Thread{}, err
	}

	t.Comments = cl
	return t, nil
}

// AddThread はentities.Threadを追加します。
func (tu *ThreadUsecase) AddThread(t entities.Thread, bulletinBoardRepository BulletinBoardRepositorer) error {
	_, err := bulletinBoardRepository.GetBulletinBoardByID(t.BulletinBoardID.Get())
	if err != nil {
		return err
	}
	return tu.Repository.AddThread(t)
}

// ListThread はentities.Threadの一覧を取得します。
func (tu *ThreadUsecase) ListThread() ([]entities.Thread, error) {
	return tu.Repository.ListThread()
}

// ListThreadByBulletinBoardID は指定されたvalueobjects.BulletinBoardIDを持つentities.Threadの一覧を取得します。
func (tu *ThreadUsecase) ListThreadByBulletinBoardID(bID valueobjects.BulletinBoardID) ([]entities.Thread, error) {
	return tu.Repository.ListThreadByBulletinBoardID(bID)
}
