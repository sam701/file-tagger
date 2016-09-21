package storage

type StorageFile struct {
	Id                uint64
	Name              string
	Period            string
	Tags              []tagIdType
	CreationTimestamp uint64
}

func (f *StorageFile) removeTag(tag tagIdType) {
	for i, t := range f.Tags {
		if t == tag {
			f.Tags[i] = f.Tags[len(f.Tags)-1]
			f.Tags = f.Tags[:len(f.Tags)-1]
		}
	}
}

func (f *StorageFile) Match(tags []string, tagIdMap map[string]tagIdType) bool {
	if len(f.Tags) < len(tags) {
		return false
	}
	for _, ta := range tags {
		match := false
		tagId := tagIdMap[ta]
		for _, t := range f.Tags {
			if t == tagId {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}
