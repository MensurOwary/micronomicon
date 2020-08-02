package tag

func (t *Repository) GetAvailableTags() Tags {
	database := t.database.Database()
	keys := make([]Tag, 0, len(database))
	for tag := range database {
		keys = append(keys, Tag{
			Name: tag,
		})
	}
	return Tags{
		Tags: keys,
		Size: len(keys),
	}
}

var EmptyTag = Tag{}

func (t *tagsService) GetTagById(name string) Tag {
	for _, tag := range t.tagDb.GetAvailableTags().Tags {
		if tag.Name == name {
			return tag
		}
	}
	return EmptyTag
}
