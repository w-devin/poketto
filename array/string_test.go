package array

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIn(t *testing.T) {
	arr := []string{"hello", "world"}

	require.Equal(t, true, In("hello", arr))
	require.Equal(t, false, In("Hello", arr))
	require.Equal(t, true, InFold("Hello", arr))
}

func TestDistinct(t *testing.T) {
	arr := []string{"hello", "world", "Hello"}

	require.Equal(t, true, In("hello", arr))
	require.Equal(t, true, In("Hello", arr))

	newArr := Distinct(arr)
	require.Equal(t, true, In("hello", newArr))
	require.Equal(t, true, In("Hello", newArr))

	newArr = DistinctFold(arr)
	require.Equal(t, true, In("hello", newArr))
	require.Equal(t, false, In("Hello", newArr))
}
