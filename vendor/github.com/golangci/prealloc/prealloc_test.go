package prealloc

import "testing"

func BenchmarkNoPreallocate(b *testing.B) {
	existing := make([]int64, 10, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Don't preallocate our initial slice
		var init []int64
		for _, element := range existing {
			init = append(init, element)
		}
	}
}

func BenchmarkPreallocate(b *testing.B) {
	existing := make([]int64, 10, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Preallocate our initial slice
		init := make([]int64, 0, len(existing))
		for _, element := range existing {
			init = append(init, element)
		}
	}
}
