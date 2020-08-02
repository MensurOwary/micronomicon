package tag

import (
	"github.com/stretchr/testify/assert"
	"micron/scraper"
	"testing"
)

func TestRepository_GetAvailableTags(t *testing.T) {
	tt := []struct {
		testName     string
		rows         map[string][]scraper.Row
		size         int
		expectedTags []Tag
	}{
		{
			testName: "GetAvailableTags when there are many tags",
			rows:     getCaseMany(),
			size:     2,
			expectedTags: []Tag{
				{Name: "react"},
				{Name: "go"},
			},
		},
		{
			testName:     "GetAvailableTags when there are no tags",
			rows:         getCaseNone(),
			size:         0,
			expectedTags: []Tag{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			repository := NewRepository(&mockDatabase{
				rows: tc.rows,
			})

			result := repository.GetAvailableTags()

			assert.Equal(t, tc.size, result.Size)
			assert.ElementsMatch(t, tc.expectedTags, result.Tags)
		})
	}

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
