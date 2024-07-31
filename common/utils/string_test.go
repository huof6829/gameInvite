package utils

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConvertToFloat(t *testing.T) {
	require := require.New(t)

	actual := ConvertToFloat(9999)
	require.Equal("9999", actual)
	actual = ConvertToFloat(10001)
	require.Equal("10.0k", actual)
	actual = ConvertToFloat(10151)
	require.Equal("10.2k", actual)
	actual = ConvertToFloat(10141)
	require.Equal("10.1k", actual)
}

func TestTimeParse(t *testing.T) {
	require := require.New(t)

	t.Log(time.Now().UnixMilli())
	sec, err := strconv.ParseInt("1721708926", 10, 64)
	require.NoError(err)
	tu := time.Unix(sec, 0)
	t.Log(tu)
}
