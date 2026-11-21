package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeDirectorValues(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	file, err := Open(writeFixtureFile(t, "directors.m4a", fixture))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	names := []string{"Director A", "Director B"}
	for _, step := range []struct {
		tags   MP4Tags
		remove []string
		want   []string
	}{
		{tags: MP4Tags{Directors: names}, want: names},
		{tags: MP4Tags{Title: "Keep directors"}, want: names},
		{tags: MP4Tags{Director: "Solo"}, want: []string{"Solo"}},
		{tags: MP4Tags{Directors: []string{}}},
		{tags: MP4Tags{Directors: names}, want: names},
		{remove: []string{"director"}},
	} {
		if err := file.Write(&step.tags, step.remove); err != nil {
			t.Fatal(err)
		}
		tags, boxes, err := file.actualRead()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(tags.Directors, step.want) || tags.Director != firstTextValue(step.want) {
			t.Fatalf("director = %q, directors = %q, want %q", tags.Director, tags.Directors, step.want)
		}
		wantBoxes := 0
		if len(step.want) != 0 {
			wantBoxes = 1
		}
		if count := len(boxes.getBoxesByPath("moov.udta.meta.ilst.(c)dir")); count != wantBoxes {
			t.Fatalf("native director items = %d, want %d", count, wantBoxes)
		}
	}
}
