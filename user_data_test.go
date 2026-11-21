package mp4tag

import "testing"

func TestNativeISRCAndRecordLabel(t *testing.T) {
	payload := []byte("ISRC-LABEL-AUDIO")
	fixture, _ := makeTaggableMP4(mdatHeader32(uint32(8+len(payload))), payload)
	path := writeFixtureFile(t, "native-user-data.m4a", fixture)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for _, step := range []struct {
		input       MP4Tags
		remove      []string
		isrc, label string
	}{
		{MP4Tags{ISRC: "HKA972302100", RecordLabel: "華語唱片"}, nil, "HKA972302100", "華語唱片"},
		{MP4Tags{Title: "Retag"}, nil, "HKA972302100", "華語唱片"},
		{MP4Tags{RecordLabel: "Replacement label"}, nil, "HKA972302100", "Replacement label"},
		{MP4Tags{}, []string{"isrc", "recordlabel"}, "", ""},
	} {
		if err := file.Write(&step.input, step.remove); err != nil {
			t.Fatal(err)
		}
		tags, err := file.Read()
		if err != nil {
			t.Fatal(err)
		}
		if tags.ISRC != step.isrc || tags.RecordLabel != step.label {
			t.Fatalf("ISRC = %q, label = %q; want %q and %q", tags.ISRC, tags.RecordLabel, step.isrc, step.label)
		}
		if stored, actual := readStcoAt(t, path), payloadOffsetIn(t, path, payload); stored != actual {
			t.Fatalf("stco = %d, media starts at %d", stored, actual)
		}
	}
}
