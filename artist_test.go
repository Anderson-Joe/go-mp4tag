package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeArtistValues(t *testing.T) {
	payload := []byte("sample audio")
	fixture, _ := makeTaggableMP4(mdatHeader32(uint32(8+len(payload))), payload)
	path := writeFixtureFile(t, "artists.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	multiple := MP4Tags{
		Artists:          []string{"Artist A", "艺人 B"},
		ArtistSorts:      []string{"A, Artist", "B, Artist"},
		AlbumArtists:     []string{"Artist C", "艺人 D"},
		AlbumArtistSorts: []string{"C, Artist", "D, Artist"},
	}
	single := MP4Tags{
		Artist: "Solo", ArtistSort: "Solo sort",
		AlbumArtist: "Album solo", AlbumArtistSort: "Album solo sort",
	}
	singleValues := MP4Tags{
		Artists: []string{single.Artist}, ArtistSorts: []string{single.ArtistSort},
		AlbumArtists: []string{single.AlbumArtist}, AlbumArtistSorts: []string{single.AlbumArtistSort},
	}
	clear := MP4Tags{Artists: []string{}, ArtistSorts: []string{}, AlbumArtists: []string{}, AlbumArtistSorts: []string{}}
	for _, step := range []struct {
		name   string
		tags   MP4Tags
		delete []string
		want   MP4Tags
	}{
		{name: "single", tags: single, want: singleValues},
		{name: "multiple", tags: multiple, want: multiple},
		{name: "preserve on unrelated edit", tags: MP4Tags{Title: "Keep credits"}, want: multiple},
		{name: "replace with legacy API", tags: single, want: singleValues},
		{name: "restore multiple", tags: multiple, want: multiple},
		{name: "delete by name", delete: []string{"artist", "artistsort", "albumartist", "albumartistsort"}},
		{name: "restore after deletion", tags: multiple, want: multiple},
		{name: "clear with empty slices", tags: clear},
	} {
		t.Run(step.name, func(t *testing.T) {
			if err := file.Write(&step.tags, step.delete); err != nil {
				t.Fatal(err)
			}
			read, boxes, err := file.actualRead()
			if err != nil {
				t.Fatal(err)
			}
			for _, field := range []struct {
				atom      string
				got, want []string
				legacy    string
			}{
				{"(c)art", read.Artists, step.want.Artists, read.Artist},
				{"soar", read.ArtistSorts, step.want.ArtistSorts, read.ArtistSort},
				{"aART", read.AlbumArtists, step.want.AlbumArtists, read.AlbumArtist},
				{"soaa", read.AlbumArtistSorts, step.want.AlbumArtistSorts, read.AlbumArtistSort},
			} {
				if !slices.Equal(field.got, field.want) || field.legacy != firstTextValue(field.want) {
					t.Fatalf("%s: got %q (legacy %q), want %q", field.atom, field.got, field.legacy, field.want)
				}
				wantItems := 0
				if len(field.want) != 0 {
					wantItems = 1
				}
				if got := len(boxes.getBoxesByPath("moov.udta.meta.ilst." + field.atom)); got != wantItems {
					t.Fatalf("%s: native item count = %d, want %d", field.atom, got, wantItems)
				}
			}
			if offset := payloadOffsetIn(t, path, payload); readStcoAt(t, path) != offset {
				t.Fatal("chunk offset no longer points to the audio payload")
			}
		})
	}
}
