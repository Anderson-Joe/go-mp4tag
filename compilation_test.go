package mp4tag

import "testing"

func TestCompilationTag(t *testing.T) {
	payload := []byte("sample audio")
	fixture, _ := makeTaggableMP4(mdatHeader32(uint32(8+len(payload))), payload)
	path := writeFixtureFile(t, "compilation.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	enabled, disabled := true, false
	for _, tt := range []struct {
		name   string
		tags   MP4Tags
		remove []string
		want   bool
	}{
		{name: "unset"},
		{name: "set", tags: MP4Tags{Compilation: &enabled}, want: true},
		{name: "preserve", tags: MP4Tags{Title: "Compilation"}, want: true},
		{name: "clear explicitly", tags: MP4Tags{Compilation: &disabled}},
		{name: "restore", tags: MP4Tags{Compilation: &enabled}, want: true},
		{name: "delete by name", remove: []string{"COMPILATION"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := file.Write(&tt.tags, tt.remove); err != nil {
				t.Fatal(err)
			}
			tags, boxes, err := file.actualRead()
			if err != nil {
				t.Fatal(err)
			}
			box := boxes.getBoxByPath("moov.udta.meta.ilst.cpil.data")
			if tt.want {
				if tags.Compilation == nil || !*tags.Compilation || box == nil {
					t.Fatal("expected an enabled native cpil tag")
				}
			} else if tags.Compilation != nil || box != nil {
				t.Fatal("expected cpil to be absent")
			}
			if offset := payloadOffsetIn(t, path, payload); readStcoAt(t, path) != offset {
				t.Fatal("chunk offset no longer points to the audio payload")
			}
		})
	}
}
