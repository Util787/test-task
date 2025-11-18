package usecase

import "context"

type ArrStorage interface {
	SaveNum(ctx context.Context, num int) error
	GetArr(ctx context.Context) ([]int, error)
}

type SortUsecase struct {
	arrStorage ArrStorage
}

func NewSortUsecase(arrStorage ArrStorage) *SortUsecase {
	return &SortUsecase{
		arrStorage: arrStorage,
	}
}
