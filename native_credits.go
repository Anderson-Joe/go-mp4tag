package mp4tag

type nativeCreditList struct {
	atom   string
	values *[]string
}

func (tags *MP4Tags) nativeCreditLists() [5]nativeCreditList {
	return [5]nativeCreditList{
		{"\xa9arg", &tags.Arrangers},
		{"\xa9con", &tags.Conductors},
		{"\xa9sne", &tags.Engineers},
		{"perf", &tags.Performers},
		{"\xa9prd", &tags.Producers},
	}
}
