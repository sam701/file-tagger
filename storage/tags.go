package storage

import "errors"

func (s *Storage) GetTags() []string {
	out := []string{}
	for t, _ := range s.allowedTags {
		out = append(out, t)
	}
	return out
}

func (s *Storage) AddTag(tag string) {
	if _, exists := s.allowedTags[tag]; exists {
		return
	}

	enc := &encoder{s.metaFile}

	enc.write(opAddAllowedTag)
	s.maxTagId++
	enc.write(s.maxTagId)
	enc.writeString(tag)
}

func (s *Storage) DeleteTag(tag string) error {
	tagId, exists := s.allowedTags[tag]
	if !exists {
		return errors.New("Tag " + tag + " does not exist")
	}

	enc := &encoder{s.metaFile}

	enc.write(opRenameTag)
	enc.write(tagId)
	return nil
}
