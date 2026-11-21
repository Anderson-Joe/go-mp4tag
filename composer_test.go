package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeComposerValues(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	file, err := Open(writeFixtureFile(t, "composers.m4a", fixture))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	names := []string{"Composer A", "Composer B"}
	sorts := []string{"A, Composer", "B, Composer"}
	for _, step := range []struct {
		tags         MP4Tags
		names, sorts []string
	}{
		{MP4Tags{Composers: names, ComposerSorts: sorts}, names, sorts},
		{MP4Tags{Title: "Keep composer metadata"}, names, sorts},
		{MP4Tags{Composer: "Solo", ComposerSort: "Solo sort"}, []string{"Solo"}, []string{"Solo sort"}},
		{MP4Tags{Composers: []string{}, ComposerSorts: []string{}}, nil, nil},
	} {
		if err := file.Write(&step.tags, nil); err != nil {
			t.Fatal(err)
		}
		tags, boxes, err := file.actualRead()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(tags.Composers, step.names) || !slices.Equal(tags.ComposerSorts, step.sorts) {
			t.Fatalf("composers = %q, sorts = %q, want %q and %q", tags.Composers, tags.ComposerSorts, step.names, step.sorts)
		}
		if len(step.names) > 0 {
			if tags.Composer != step.names[0] || tags.ComposerSort != step.sorts[0] {
				t.Fatal("legacy fields did not expose the first value")
			}
			for _, name := range []string{"(c)wrt", "soco"} {
				if boxCount := len(boxes.getBoxesByPath("moov.udta.meta.ilst." + name)); boxCount != 1 {
					t.Fatalf("native %s item count = %d, want 1", name, boxCount)
				}
			}
		}
	}
}
