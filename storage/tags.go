package storage

func (s *Storage) GetTags() []string {
	out := []string{}
	for t, _ := range s.allowedTags {
		out = append(out, t)
	}
	return out
}

func (s *Storage) AddTag(tag string) {
	if s.allowedTags[tag] {
		return
	}

	enc := &encoder{s.metaFile}

	enc.write(opAddAllowedTag)
	enc.writeString(tag)
}
