//go:build !solution

package genericsum

import (
	"math/cmplx"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type numTp interface {
	int | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64
}

func Min[Tp constraints.Ordered](a Tp, b Tp) Tp {
	if b < a {
		return b
	}
	return a
}

func SortSlice[Tp constraints.Ordered](a []Tp) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}

func MapsEqual[T, U comparable](m1, m2 map[T]U) bool {
	if len(m1) != len(m2) {
		return false
	}

	for key, value := range m1 {
		if v, ok := m2[key]; !ok || v != value {
			return false
		}
	}

	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	merged := make(chan T)
	var wg sync.WaitGroup

	wg.Add(len(chs))

	for _, ch := range chs {
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				merged <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

type mtxTp interface {
	numTp | complex64 | complex128
}

func IsHermitianMatrix[T mtxTp](m [][]T) bool {
	for x := range m {
		for y := range m[x] {
			if !isHermitianElement(m, x, y) {
				return false
			}
		}
	}
	return true
}

func isHermitianElement[T mtxTp](m [][]T, x, y int) bool {
	fst := any(m[x][y])
	scd := any(m[y][x])

	if fstInt, ok := fst.(int); ok {
		if scdInt, ok := scd.(int); ok {
			return fstInt == scdInt
		}
		return false
	}

	if fstComplex64, ok := fst.(complex64); ok {
		if scdComplex64, ok := scd.(complex64); ok {
			return cmplx.Conj(complex128(fstComplex64)) == complex128(scdComplex64)
		}
		return false
	}

	if fstComplex128, ok := fst.(complex128); ok {
		if scdComplex128, ok := scd.(complex128); ok {
			return cmplx.Conj(fstComplex128) == scdComplex128
		}
		return false
	}

	return false
}
