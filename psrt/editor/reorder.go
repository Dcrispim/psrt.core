package editor

import "fmt"

func moveSliceItem[T any](items []T, from, to int) ([]T, error) {
	n := len(items)
	if from < 0 || from >= n {
		return items, fmt.Errorf("source position %d out of range", from)
	}
	if to < 0 || to >= n {
		return items, fmt.Errorf("target position %d out of range", to)
	}
	if from == to {
		return items, nil
	}
	item := items[from]
	rest := append([]T{}, items[:from]...)
	rest = append(rest, items[from+1:]...)
	rest = append(rest[:to], append([]T{item}, rest[to:]...)...)
	return rest, nil
}

func moveBeforeIndex[T any](items []T, from, ref int) ([]T, error) {
	if from == ref {
		return items, nil
	}
	to := ref
	if from < ref {
		to = ref - 1
	}
	return moveSliceItem(items, from, to)
}

func moveAfterIndex[T any](items []T, from, ref int) ([]T, error) {
	if from == ref {
		return items, nil
	}
	to := ref + 1
	if from < ref {
		to = ref
	}
	if to >= len(items) {
		to = len(items) - 1
	}
	return moveSliceItem(items, from, to)
}

func moveToIndex[T any](items []T, from, to int) ([]T, error) {
	n := len(items)
	if to < 0 {
		to = 0
	}
	if to >= n {
		to = n - 1
	}
	return moveSliceItem(items, from, to)
}

func moveByDelta[T any](items []T, from, delta int) ([]T, error) {
	return moveToIndex(items, from, from+delta)
}
