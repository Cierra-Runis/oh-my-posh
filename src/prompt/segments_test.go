package prompt

import (
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/config"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime"

	"github.com/stretchr/testify/assert"
)

func TestRenderBlock(t *testing.T) {
	engine := New(&runtime.Flags{
		IsPrimary: true,
	})
	block := &config.Block{
		Segments: []*config.Segment{
			{
				Type:       "text",
				Template:   "Hello",
				Foreground: "red",
				Background: "blue",
			},
			{
				Type:       "text",
				Template:   "World",
				Foreground: "red",
				Background: "blue",
			},
		},
	}

	prompt, length := engine.writeBlockSegments(block)
	assert.Equal(t, "\x1b[44m\x1b[31mHello\x1b[0m\x1b[44m\x1b[31mWorld\x1b[0m", prompt)
	assert.Equal(t, 10, length)
}

func TestCanRenderSegment(t *testing.T) {
	cases := []struct {
		Case     string
		Executed map[string]bool
		Needs    []string
		Expected bool
	}{
		{
			Case:     "No cross segment dependencies",
			Expected: true,
		},
		{
			Case:     "Cross segment dependencies, nothing executed",
			Expected: false,
			Needs:    []string{"Foo"},
		},
		{
			Case:     "Cross segment dependencies, available",
			Expected: true,
			Executed: map[string]bool{
				"Foo": true,
			},
			Needs: []string{"Foo"},
		},
	}
	for _, c := range cases {
		segment := &config.Segment{
			Type:  "text",
			Needs: c.Needs,
		}

		engine := &Engine{}
		got := engine.canRenderSegment(segment, c.Executed)

		assert.Equal(t, c.Expected, got, c.Case)
	}
}
