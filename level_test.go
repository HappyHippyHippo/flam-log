package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelFrom(t *testing.T) {
	scenarios := []struct {
		name string
		val  any
		def  []Level
		want Level
	}{
		{name: "from level", val: Info, want: Info},
		{name: "from int", val: int(Warning), want: Warning},
		{name: "from invalid int", val: 99, want: None},
		{name: "from invalid int with default", val: 99, def: []Level{Error}, want: Error},
		{name: "from string", val: "debug", want: Debug},
		{name: "from invalid string", val: "invalid", want: None},
		{name: "from invalid string with default", val: "invalid", def: []Level{Fatal}, want: Fatal},
		{name: "from other type", val: 1.23, want: None},
		{name: "from other type with default", val: true, def: []Level{Notice}, want: Notice},
		{name: "from nil", val: nil, want: None},
		{name: "from nil with default", val: nil, def: []Level{Info}, want: Info},
	}

	for _, scenario := range scenarios {
		test := scenario
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, LevelFrom(test.val, test.def...))
		})
	}
}

func TestLevelNameAndLevelMap(t *testing.T) {
	t.Parallel()

	t.Run("check if all levels in LevelName are also in LevelMap", func(t *testing.T) {
		t.Parallel()

		for level, name := range LevelName {
			levelFromName, ok := LevelMap[name]
			assert.True(t, ok, fmt.Sprintf("level name '%s' not in LevelMap", name))
			assert.Equal(t, level, levelFromName, fmt.Sprintf("level name '%s' maps to a different level in LevelMap", name))
		}
	})

	t.Run("check if all levels in LevelMap are also in LevelName", func(t *testing.T) {
		t.Parallel()

		for name, level := range LevelMap {
			nameFromLevel, ok := LevelName[level]
			assert.True(t, ok, fmt.Sprintf("level %d not in LevelName", level))
			assert.Equal(t, name, nameFromLevel, fmt.Sprintf("level %d maps to a different name in LevelName", level))
		}
	})
}
