package mp4tag

import (
	"slices"
	"testing"
)

func TestNativeRoleAndGenreValues(t *testing.T) {
	fixture, _ := makeTaggableMP4(mdatHeader32(9), []byte{0})
	file, err := Open(writeFixtureFile(t, "native-roles.m4a", fixture))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	fields := &MP4Tags{
		CustomGenres:   []string{"Mandopop", "Music", "Pop"},
		ItunesGenreIDs: []uint64{1253, 34, 14},
	}
	for _, role := range fields.nativeCreditLists() {
		*role.values = []string{"One", "Two"}
	}
	if err := file.Write(fields, nil); err != nil {
		t.Fatal(err)
	}
	if err := file.Write(&MP4Tags{Title: "Retag"}, nil); err != nil {
		t.Fatal(err)
	}
	got, err := file.Read()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got.ItunesGenreIDs, fields.ItunesGenreIDs) || got.ItunesGenreID != 1253 ||
		!slices.Equal(got.CustomGenres, fields.CustomGenres) || got.Conductor != "One" {
		t.Fatalf("geID=%v, genre=%v, conductor=%q", got.ItunesGenreIDs, got.CustomGenres, got.Conductor)
	}
	for _, role := range got.nativeCreditLists() {
		if !slices.Equal(*role.values, []string{"One", "Two"}) {
			t.Errorf("native %q = %q", role.atom, *role.values)
		}
	}
}
