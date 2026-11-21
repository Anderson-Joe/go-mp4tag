package mp4tag

import (
	"slices"
	"testing"
)

func TestReplaceCustomValues(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	path := writeFixtureFile(t, "custom.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = file.Write(&MP4Tags{
		Custom:      map[string]string{"GENREID": "14", "OTHER": "a"},
		OtherCustom: map[string][]string{"GENREID": {"34"}, "OTHER": {"b"}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, values := range [][]string{{"1253", "34", "14"}, {"14"}, nil, {"34", "14"}, {}} {
		if err := file.Write(&MP4Tags{ReplaceCustomValues: map[string][]string{"genreid": values}}, nil); err != nil {
			t.Fatal(err)
		}
		tags, err := file.Read()
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		if first, ok := tags.Custom["GENREID"]; ok {
			got = append([]string{first}, tags.OtherCustom["GENREID"]...)
		}
		if !slices.Equal(got, values) {
			t.Fatalf("GENREID = %q, want %q", got, values)
		}
		if tags.Custom["OTHER"] != "a" || !slices.Equal(tags.OtherCustom["OTHER"], []string{"b"}) {
			t.Fatal("unrelated freeform values changed")
		}
	}
}
