package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitArgs(t *testing.T) {
	for name, val := range map[string]struct {
		command string
		expectedResult []string
		errMessage string
	}{
		"command 1": {
			command: "/confluence edit \"abc\"",
			expectedResult: []string{"/confluence", "edit", "abc"},
		},
		"command 2": {
			command: "/confluence list     ",
			expectedResult: []string{"/confluence", "list"},
		},
		"command 3": {
			command: "/confluence      subscribe",
			expectedResult: []string{"/confluence", "subscribe"},
		},
		"command 4": {
			command: "/confluence edit \"  test     \"",
			expectedResult: []string{"/confluence", "edit", "test"},
		},
		"command 5": {
			command: "/confluence subscribe",
			expectedResult: []string{"/confluence", "subscribe"},
		},
		"command 6": {
			command: "/confluence unsubscribe \" test\"",
			expectedResult: []string{"/confluence", "unsubscribe", "test"},
		},
		"command 7": {
			command: "/confluence edit \"abc  ",
			errMessage: "quotes not closed",
		},
	} {
		t.Run(name, func(t *testing.T) {
			args, err := SplitArgs(val.command)
			if err != nil {
				assert.Equal(t, val.errMessage, err.Error())
				return
			}
			assert.Equal(t, val.expectedResult, args)
		})
	}
}
