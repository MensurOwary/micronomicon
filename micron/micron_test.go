package micron

import (
	"github.com/stretchr/testify/assert"
	"micron/scraper"
	"micron/tag"
	"testing"
)

func TestMicronsService_GetARandomMicronForTag(t *testing.T) {
	t.Run("When there are 1+ microns available", func(t *testing.T) {
		service := NewService(&mockDatabase{
			rows: getCaseMany(),
		})

		micron := service.GetARandomMicronForTag(tag.Tag{
			Name: "react",
		})

		assert.Contains(t, []Micron{
			{
				URL:   "www.react.com",
				Title: "How to get started with react",
				tag:   tag.Tag{Name: "react"},
			},
			{
				URL:   "www.hooks.com/react",
				Title: "React Hooks",
				tag:   tag.Tag{Name: "react"},
			},
		}, micron)
	})

	t.Run("When there is no micron available", func(t *testing.T) {
		service := NewService(&mockDatabase{
			rows: getCaseNone(),
		})

		micron := service.GetARandomMicronForTag(tag.Tag{
			Name: "react",
		})

		assert.Equal(t, EmptyMicron, micron)
	})

	t.Run("When there is only 1 micron available", func(t *testing.T) {
		service := NewService(&mockDatabase{
			rows: getCaseMany(),
		})

		micron := service.GetARandomMicronForTag(tag.Tag{
			Name: "go",
		})

		assert.Equal(t, Micron{
			URL:   "www.golang.org",
			Title: "Why is Golang so fast?",
			tag:   tag.Tag{Name: "go"},
		}, micron)
	})
}

type mockDatabase struct {
	rows map[string][]scraper.Row
}

func (m *mockDatabase) Database() map[string][]scraper.Row {
	return m.rows
}

func getCaseMany() map[string][]scraper.Row {
	return map[string][]scraper.Row{
		"react": {
			scraper.Row{
				Title: "How to get started with react",
				Link:  "www.react.com",
				Tags:  "react,native,javascript",
			}, scraper.Row{
				Title: "React Hooks",
				Link:  "www.hooks.com/react",
				Tags:  "react,hooks,async",
			},
		},
		"go": {
			scraper.Row{
				Title: "Why is Golang so fast?",
				Link:  "www.golang.org",
				Tags:  "go,gin,performance",
			},
		},
	}
}

func getCaseNone() map[string][]scraper.Row {
	return map[string][]scraper.Row{}
}
