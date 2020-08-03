package tag

// Fetches all the available tags
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

// Represents a non-existent resource
var EmptyTag = Tag{}

// Fetches the tag by its name/id
func (t *TagsService) GetTagByID(name string) Tag {
	for _, tag := range t.tagDb.GetAvailableTags().Tags {
		if tag.Name == name {
			return tag
		}
	}
	return EmptyTag
}
