package usecases

import (
	"reflect"
	"sort"
	"testing"

	"github.com/osechiman/BulltienBoard/adapters/gateways"
	"github.com/osechiman/BulltienBoard/entities"
	"github.com/osechiman/BulltienBoard/entities/valueobjects"
)

func TestBulletinBoardUsecase_AddBulletinBoard(t *testing.T) {
	repository := gateways.GetInMemoryRepositoryInstance()
	repository.DeleteAll()

	bid, _ := valueobjects.NewBulletinBoardID("")
	b, _ := entities.NewBulletinBoard(bid, "bulletin board title")

	type fields struct {
		Repository BulletinBoardRepositorer
	}
	type args struct {
		bb entities.BulletinBoard
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "エンティティの登録が正常に出来る",
			fields: fields{
				Repository: repository,
			},
			args: args{
				bb: b,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbu := &BulletinBoardUsecase{
				Repository: tt.fields.Repository,
			}
			if err := bbu.AddBulletinBoard(tt.args.bb); (err != nil) != tt.wantErr {
				t.Errorf("AddBulletinBoard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("BulletinBoardの登録数がBulletinBoardLimitを超えて登録された場合、エラーが返却される", func(t *testing.T) {
		repository.DeleteAll()
		bbu := &BulletinBoardUsecase{
			Repository: repository,
		}
		for i := 0; i < BulletinBoardLimit; i++ {
			bid, _ = valueobjects.NewBulletinBoardID("")
			b, _ := entities.NewBulletinBoard(bid, "bulletin board title")
			repository.AddBulletinBoard(b)
		}

		bid, _ = valueobjects.NewBulletinBoardID("")
		b, _ := entities.NewBulletinBoard(bid, "bulletin board last")

		wantErr := true
		if err := bbu.AddBulletinBoard(b); (err != nil) != wantErr {
			t.Errorf("AddBulletinBoard() error = %v, wantErr %v", err, wantErr)
		}
	})

}

func TestBulletinBoardUsecase_GetBulletinBoardByID(t *testing.T) {
	repository := gateways.GetInMemoryRepositoryInstance()
	repository.DeleteAll()

	bid, _ := valueobjects.NewBulletinBoardID("")
	b, _ := entities.NewBulletinBoard(bid, "bulletin board title")
	repository.AddBulletinBoard(b)
	type fields struct {
		Repository BulletinBoardRepositorer
	}
	type args struct {
		ID               valueobjects.BulletinBoardID
		threadRepository ThreadRepositorer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.BulletinBoard
		wantErr bool
	}{
		{
			name: "BulletinBoardIDからentities.BulletinBoardが取得出来る",
			fields: fields{
				Repository: repository,
			},
			args: args{
				ID:               bid,
				threadRepository: repository,
			},
			want: entities.BulletinBoard{
				ID:      bid,
				Title:   "bulletin board title",
				Threads: []entities.Thread{},
			},
			wantErr: false,
		},
		{
			name: "BulletinBoardIDが存在しない値だった場合、エラーが返却される",
			fields: fields{
				Repository: repository,
			},
			args: args{
				ID:               valueobjects.BulletinBoardID{},
				threadRepository: repository,
			},
			want:    entities.BulletinBoard{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbu := &BulletinBoardUsecase{
				Repository: tt.fields.Repository,
			}
			got, err := bbu.GetBulletinBoardByID(tt.args.ID, tt.args.threadRepository)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBulletinBoardByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBulletinBoardByID() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("BulletinBoardにThreadが登録されていた場合はThreadの内容も返却される", func(t *testing.T) {
		repository.DeleteAll()
		bid, _ = valueobjects.NewBulletinBoardID("")
		b, _ := entities.NewBulletinBoard(bid, "bulletin board title")
		bbu := &BulletinBoardUsecase{
			Repository: repository,
		}
		repository.AddBulletinBoard(b)

		tid, _ := valueobjects.NewThreadID("")
		th, _ := entities.NewThread(tid, bid, "thread title")
		repository.AddThread(th)

		wantErr := false
		b.Threads = append(b.Threads, th)
		want := b

		got, err := bbu.GetBulletinBoardByID(bid, repository)
		if (err != nil) != wantErr {
			t.Errorf("GetBulletinBoardByID() error = %v, wantErr %v", err, wantErr)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetBulletinBoardByID() got = %v, want %v", got, want)
		}
	})
}

func TestBulletinBoardUsecase_ListBulletinBoard(t *testing.T) {
	repository := gateways.GetInMemoryRepositoryInstance()
	repository.DeleteAll()

	bid, _ := valueobjects.NewBulletinBoardID("")
	b, _ := entities.NewBulletinBoard(bid, "bulletin board title")
	repository.AddBulletinBoard(b)

	bid1, _ := valueobjects.NewBulletinBoardID("")
	b1, _ := entities.NewBulletinBoard(bid1, "bulletin board1 title")
	repository.AddBulletinBoard(b1)

	want := append([]entities.BulletinBoard{}, b, b1)

	type fields struct {
		Repository BulletinBoardRepositorer
	}
	tests := []struct {
		name    string
		fields  fields
		want    []entities.BulletinBoard
		wantErr bool
	}{
		{
			name: "[]entities.BulletinBoardが取得出来る",
			fields: fields{
				Repository: repository,
			},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbu := &BulletinBoardUsecase{
				Repository: tt.fields.Repository,
			}
			got, err := bbu.ListBulletinBoard()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBulletinBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Sliceの順序はソートせずに返却する仕様なので、テスト時には一度ソートをして値が等価であるかを検証します。
			sort.Slice(got, func(i, j int) bool {
				return got[i].ID.String() < got[j].ID.String()
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID.String() < tt.want[j].ID.String()
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListBulletinBoard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBulletinBoardUsecase(t *testing.T) {
	repository := gateways.GetInMemoryRepositoryInstance()
	repository.DeleteAll()

	type args struct {
		r BulletinBoardRepositorer
	}
	tests := []struct {
		name string
		args args
		want *BulletinBoardUsecase
	}{
		{
			name: "オブジェクトが正常に生成される",
			args: args{
				r: repository,
			},
			want: &BulletinBoardUsecase{
				Repository: repository,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBulletinBoardUsecase(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBulletinBoardUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}
