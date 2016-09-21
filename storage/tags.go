package storage

import "errors"

func (s *Storage) GetTags() []string {
	content := s.readEntries(nil)

	out := []string{}
	for t, _ := range content.tags.allowedTags {
		out = append(out, t)
	}
	return out
}

func (s *Storage) AddTags(tags []string) error {
	content, f := s.openAndReadEntries()
	defer f.Close()

	for _, newTag := range tags {
		if _, exists := content.tags.allowedTags[newTag]; !exists {
			content.tags.maxTagId++
			content.tags.allowedTags[newTag] = content.tags.maxTagId
		}
	}

	if len(content.tags.allowedTags) == len(content.tags.tagNames) {
		return errors.New("No new tags")
	}

	enc := &encoder{f}
	enc.write(opSetAllowedTags)
	enc.writeAllowedTagsMap(content.tags.allowedTags)
	return nil
}

func (s *Storage) DeleteTag(tagToDelete string) error {
	content, f := s.openAndReadEntries()
	defer f.Close()

	delete(content.tags.allowedTags, tagToDelete)

	if len(content.tags.allowedTags) == len(content.tags.tagNames) {
		return errors.New("No tags have been deleted")
	}

	enc := &encoder{f}
	enc.write(opSetAllowedTags)
	enc.writeAllowedTagsMap(content.tags.allowedTags)
	return nil
}
