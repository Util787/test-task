package usecase

import "context"

func (s *SortUsecase) SaveAndSort(ctx context.Context, num int) ([]int, error) {
	err := s.arrStorage.SaveNum(ctx, num)
	if err != nil {
		return nil, err
	}

	arr, err := s.arrStorage.GetArr(ctx)
	if err != nil {
		return nil, err
	}

	arr = quickSort(arr)

	return arr, nil
}

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivot := arr[len(arr)/2]
	var left, equal, right []int
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v == pivot {
			equal = append(equal, v)
		} else {
			right = append(right, v)
		}
	}

	res := quickSort(left)
	res = append(res, equal...)
	res = append(res, quickSort(right)...)
	return res
}
