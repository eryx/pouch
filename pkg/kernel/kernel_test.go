package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKernelVersion(t *testing.T) {
	version, err := GetKernelVersion()
	assert.Equal(t, nil, err)

	println(version.String())
}

func Benchmark_GetKernelVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKernelVersion()
	}
}

func Benchmark_GetKernelVersionByVarCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKernelVersionByVarCache()
	}
}

func Benchmark_GetKernelVersionByVarCacheWithTTL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKernelVersionByVarCacheWithTTL()
	}
}

func Benchmark_GetKernelVersionByUnix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKernelVersionByUnix()
	}
}

func Benchmark_GetKernelVersionBySyscall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKernelVersionBySyscall()
	}
}
