package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size = 1", 1},
		{"size = 5", 5},
		{"size = 10", 10},
		{"size = 20", 20},
		{"size = 30", 30},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str1, err1 := RandomString(tt.size)
			str2, err2 := RandomString(tt.size)

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			for _, ch := range str1 {
				assert.Contains(t, Charset, string(ch))
			}
			for _, ch := range str2 {
				assert.Contains(t, Charset, string(ch))
			}

			assert.NotEqual(t, str1, str2)
		})
	}
}
