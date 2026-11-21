package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeContributorIDs(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	file, err := Open(writeFixtureFile(t, "ids.m4a", fixture))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	artists := []uint64{42, 4294967297}
	composers := []uint64{4294967298, 19}
	for _, step := range []struct {
		input              MP4Tags
		artists, composers []uint64
		legacyArtist       int32
	}{
		{MP4Tags{ItunesArtistIDs: artists, ItunesComposerIDs: composers}, artists, composers, 42},
		{MP4Tags{Title: "Retag"}, artists, composers, 42},
		{MP4Tags{ItunesArtistID: 7}, []uint64{7}, composers, 7},
		{MP4Tags{ItunesArtistIDs: []uint64{}, ItunesComposerIDs: []uint64{}}, nil, nil, -1},
	} {
		if err := file.Write(&step.input, nil); err != nil {
			t.Fatal(err)
		}
		got, err := file.Read()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got.ItunesArtistIDs, step.artists) || !slices.Equal(got.ItunesComposerIDs, step.composers) || got.ItunesArtistID != step.legacyArtist {
			t.Fatalf("atID=%v cmID=%v legacy=%v", got.ItunesArtistIDs, got.ItunesComposerIDs, got.ItunesArtistID)
		}
	}
}
