package mp4tag

import "testing"

func TestPlaybackMetadataRoundTrip(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	file, err := Open(writeFixtureFile(t, "playback.m4a", fixture))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gapless := false
	if err := file.Write(&MP4Tags{
		XID: "Vendor:isrc:TEST", Comment: "Comment",
		ItunesGapless: &gapless, ItunesGenreID: 1253,
	}, nil); err != nil {
		t.Fatal(err)
	}
	if err := file.Write(&MP4Tags{Title: "Retag"}, nil); err != nil {
		t.Fatal(err)
	}
	tags, boxes, err := file.actualRead()
	if err != nil {
		t.Fatal(err)
	}
	if tags.XID != "Vendor:isrc:TEST" || tags.Comment != "Comment" ||
		tags.ItunesGapless == nil || *tags.ItunesGapless || tags.ItunesGenreID != 1253 {
		t.Fatalf("unexpected playback tags: XID %q, comment %q, gapless %v, genre ID %d",
			tags.XID, tags.Comment, tags.ItunesGapless, tags.ItunesGenreID)
	}
	for _, name := range []string{"xid ", "(c)cmt", "pgap", "geID"} {
		if count := len(boxes.getBoxesByPath("moov.udta.meta.ilst." + name)); count != 1 {
			t.Errorf("native %q items = %d, want 1", name, count)
		}
	}
}
