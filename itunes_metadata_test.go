package mp4tag

import (
	"bytes"
	"os"
	"testing"
)

func TestNativeItunesMetadata(t *testing.T) {
	audio := []byte("sample audio")
	fixture, _ := makeTaggableMP4(mdatHeader32(uint32(8+len(audio))), audio)
	path := writeFixtureFile(t, "itunes.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	music, oldMovie := ItunesMediaTypeMusic, ItunesMediaType(0)
	for _, step := range []struct {
		name           string
		tags           MP4Tags
		remove         []string
		wantID         uint32
		wantStorefront uint32
		wantType       *ItunesMediaType
	}{
		{name: "absent"},
		{name: "set", tags: MP4Tags{ItunesCatalogID: 4294967295, ItunesStorefrontID: 143441, ItunesMediaType: &music}, wantID: 4294967295, wantStorefront: 143441, wantType: &music},
		{name: "preserve", tags: MP4Tags{Title: "Keep iTunes metadata"}, wantID: 4294967295, wantStorefront: 143441, wantType: &music},
		{name: "replace", tags: MP4Tags{ItunesCatalogID: 1717034230, ItunesStorefrontID: 143465, ItunesMediaType: &oldMovie}, wantID: 1717034230, wantStorefront: 143465, wantType: &oldMovie},
		{name: "remove", remove: []string{"ITUNESCATALOGID", "ITUNESCOUNTRYID", "ITUNESMEDIATYPE"}},
	} {
		t.Run(step.name, func(t *testing.T) {
			if err := file.Write(&step.tags, step.remove); err != nil {
				t.Fatal(err)
			}
			tags, boxes, err := file.actualRead()
			if err != nil {
				t.Fatal(err)
			}
			if tags.ItunesCatalogID != step.wantID {
				t.Errorf("catalog ID = %d, want %d", tags.ItunesCatalogID, step.wantID)
			}
			if tags.ItunesStorefrontID != step.wantStorefront {
				t.Errorf("storefront ID = %d, want %d", tags.ItunesStorefrontID, step.wantStorefront)
			}
			if (tags.ItunesMediaType == nil) != (step.wantType == nil) ||
				(step.wantType != nil && tags.ItunesMediaType != nil && *tags.ItunesMediaType != *step.wantType) {
				t.Error("read-back media type differs from expected")
			}
			if step.name == "set" {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				for _, native := range []struct {
					atom    string
					payload []byte
				}{
					{"cnID", []byte{0xff, 0xff, 0xff, 0xff}},
					{"sfID", []byte{0, 2, 0x30, 0x51}},
					{"stik", []byte{1}},
				} {
					box := boxes.getBoxByPath("moov.udta.meta.ilst." + native.atom + ".data")
					want := append([]byte("data\x00\x00\x00\x15\x00\x00\x00\x00"), native.payload...)
					if box == nil || !bytes.Equal(data[box.StartOffset+4:box.EndOffset], want) {
						t.Fatalf("%s does not contain the expected native integer data", native.atom)
					}
				}
			}
			if offset := payloadOffsetIn(t, path, audio); readStcoAt(t, path) != offset {
				t.Fatal("chunk offset no longer points to the audio payload")
			}
		})
	}
}
