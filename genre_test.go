package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeGenreValues(t *testing.T) {
	payload := []byte("sample audio")
	fixture, _ := makeTaggableMP4(mdatHeader32(uint32(8+len(payload))), payload)
	path := writeFixtureFile(t, "genres.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for _, tt := range []struct {
		tags MP4Tags
		want []string
	}{
		{MP4Tags{CustomGenre: "Pop"}, []string{"Pop"}},
		{MP4Tags{CustomGenres: []string{"Rock", "R&B/Soul", "国语流行"}}, []string{"Rock", "R&B/Soul", "国语流行"}},
		{MP4Tags{Title: "Preserve genres"}, []string{"Rock", "R&B/Soul", "国语流行"}},
		{MP4Tags{CustomGenres: []string{"Jazz", "Blues"}}, []string{"Jazz", "Blues"}},
		{MP4Tags{CustomGenre: "Pop"}, []string{"Pop"}},
		{MP4Tags{CustomGenres: []string{}}, nil},
	} {
		if err := file.Write(&tt.tags, nil); err != nil {
			t.Fatal(err)
		}
		tags, boxes, err := file.actualRead()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(tags.CustomGenres, tt.want) {
			t.Fatalf("genres = %q, want %q", tags.CustomGenres, tt.want)
		}
		wantFirst := ""
		if len(tt.want) > 0 {
			wantFirst = tt.want[0]
		}
		if tags.CustomGenre != wantFirst {
			t.Fatalf("legacy genre = %q, want %q", tags.CustomGenre, wantFirst)
		}
		count := 0
		if len(tt.want) > 0 {
			count = 1
		}
		if got := len(boxes.getBoxesByPath("moov.udta.meta.ilst.(c)gen")); got != count {
			t.Fatalf("native genre box count = %d, want %d", got, count)
		}
		if bytesAt := payloadOffsetIn(t, path, payload); readStcoAt(t, path) != bytesAt {
			t.Fatalf("chunk offset no longer points to audio payload")
		}
		if len(tags.Custom) != 0 {
			t.Fatalf("unexpected freeform tags: %v", tags.Custom)
		}
	}
}
