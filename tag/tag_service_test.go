package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagsService_GetTagById(t *testing.T) {

	t.Run("Get an existing tag", func(t *testing.T) {
		service := NewService(nil, NewRepository(&mockDatabase{
			rows: getCaseMany(),
		}))

		tag := service.GetTagById("react")

		assert.Equal(t, Tag{
			Name: "react",
		}, tag)
	})

	t.Run("Tag does not exist", func(t *testing.T) {
		service := NewService(nil, NewRepository(&mockDatabase{
			rows: getCaseNone(),
		}))

		tag := service.GetTagById("react")

		assert.Equal(t, EmptyTag, tag)
	})
}
