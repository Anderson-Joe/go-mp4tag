package mp4tag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

type userDataText struct {
	atom, value string
}

func userDataChanges(existing, input *MP4Tags, deletes []string) []userDataText {
	var changes []userDataText
	for _, field := range []struct {
		atom, name, old, value string
	}{
		{"isr", "isrc", existing.ISRC, input.ISRC},
		{"lab", "recordlabel", existing.RecordLabel, input.RecordLabel},
	} {
		if field.value != "" && field.value != field.old ||
			field.old != "" && (containsStr(deletes, field.name) || containsStr(deletes, "alltags")) {
			changes = append(changes, userDataText{field.atom, field.value})
		}
	}
	return changes
}

func (mp4 MP4) readUserDataText(boxes MP4Boxes, name string) (string, error) {
	box := boxes.getBoxByPath("moov.udta." + name)
	if box == nil {
		return "", nil
	}
	var header [4]byte
	reader := io.NewSectionReader(mp4.f, box.StartOffset+box.HeaderSize, box.BoxSize-box.HeaderSize)
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return "", fmt.Errorf("read %s text header: %w", name, err)
	}
	text := make([]byte, binary.BigEndian.Uint16(header[:2]))
	if _, err := io.ReadFull(reader, text); err != nil {
		return "", fmt.Errorf("read %s text: %w", name, err)
	}
	if bytes.HasPrefix(text, []byte{0xfe, 0xff}) {
		text = text[2:]
		if len(text)%2 != 0 {
			return "", fmt.Errorf("invalid UTF-16 text size for %s", name)
		}
		units := make([]uint16, len(text)/2)
		for i := range units {
			units[i] = binary.BigEndian.Uint16(text[i*2:])
		}
		return strings.TrimRight(string(utf16.Decode(units)), "\x00"), nil
	}
	return strings.TrimRight(string(text), "\x00"), nil
}

func chunkTablePosition(box *MP4Box, delta, shiftFrom int64) int64 {
	pos := box.StartOffset + 16
	if shiftFrom > 0 && box.StartOffset >= shiftFrom {
		pos += delta
	}
	return pos
}

func rewriteUserData(path string, changes []userDataText) error {
	source, err := Open(path)
	if err != nil {
		return err
	}
	defer source.Close()
	_, boxes, err := source.actualRead()
	if err != nil {
		return err
	}
	moov := boxes.getBoxByPath("moov")
	udta := boxes.getBoxByPath("moov.udta")
	if moov.HeaderSize != 8 || udta.HeaderSize != 8 {
		return fmt.Errorf("64-bit moov or udta resizing is not supported")
	}

	var added bytes.Buffer
	remove := make(map[string]bool, len(changes))
	for _, change := range changes {
		remove["moov.udta.(c)"+change.atom] = true
		if change.value == "" {
			continue
		}
		if len(change.value) > math.MaxUint16 {
			return fmt.Errorf("UserData text for %s exceeds 16 bits", change.atom)
		}
		var header [12]byte
		binary.BigEndian.PutUint32(header[:4], uint32(len(header)+len(change.value)))
		copy(header[4:8], "\xa9"+change.atom)
		binary.BigEndian.PutUint16(header[8:10], uint16(len(change.value)))
		// Packed ISO 639-2 "und" identifies UTF-8 text with no specific language.
		binary.BigEndian.PutUint16(header[10:12], 21<<10|14<<5|4)
		added.Write(header[:])
		added.WriteString(change.value)
	}

	out, err := os.CreateTemp(filepath.Dir(path), ".mp4tag-udta-*")
	if err != nil {
		return err
	}
	defer os.Remove(out.Name())
	defer out.Close()
	copyRange := func(start, end int64) error {
		_, err := io.CopyN(out, io.NewSectionReader(source.f, start, end-start), end-start)
		return err
	}
	delta := int64(added.Len())
	var pos int64
	for _, box := range boxes.Boxes {
		if !remove[box.Path] || box.StartOffset < udta.StartOffset || box.EndOffset > udta.EndOffset {
			continue
		}
		if err := copyRange(pos, box.StartOffset); err != nil {
			return err
		}
		pos = box.EndOffset
		delta -= box.BoxSize
	}
	if err := copyRange(pos, udta.EndOffset); err != nil {
		return err
	}
	if _, err := added.WriteTo(out); err != nil {
		return err
	}
	if err := copyRange(udta.EndOffset, source.size); err != nil {
		return err
	}

	for _, box := range []*MP4Box{moov, udta} {
		size := box.BoxSize + delta
		if size < 8 || size > math.MaxUint32 {
			return fmt.Errorf("resized %s size %d is invalid", box.Path, size)
		}
		var value [4]byte
		binary.BigEndian.PutUint32(value[:], uint32(size))
		if _, err := out.WriteAt(value[:], box.StartOffset); err != nil {
			return err
		}
	}
	if boxes.getBoxByPath("mdat").StartOffset >= udta.EndOffset {
		if err := source.updateChunkOffsetsAt(out, boxes, delta, udta.EndOffset); err != nil {
			return err
		}
	}
	if err := out.Close(); err != nil {
		return err
	}
	if err := source.Close(); err != nil {
		return err
	}
	return moveMP4(out.Name(), path)
}
