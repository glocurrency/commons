package logger_test

import (
	"testing"

	"github.com/glocurrency/commons/logger"
	"github.com/stretchr/testify/assert"
)

func TestIsErrLog(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			"err json log",
			`{"error":"i am an error!","level":"error","msg":"error!","time":"2024-04-07T11:07:42Z"}`,
			true,
		},
		{
			"info json log",
			`{"age":100,"level":"info","msg":"hi!","time":"2024-04-07T11:17:48Z"}`,
			false,
		},
		{
			"err text log",
			`time="2024-04-07T11:10:00Z" level=error msg="error!" age=100 error="i am an error!"`,
			true,
		},
		{
			"info text log",
			`time="2024-04-07T11:17:23Z" level=info msg="hi!" age=100`,
			false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, logger.IsErrLog([]byte(test.data)))
		})
	}
}
