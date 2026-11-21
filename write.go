package mp4tag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const BufSize = 4096*1024

func overwriteTags(mergedTags, tags *MP4Tags, delStrings []string) *MP4Tags{
	if containsStr(delStrings, "alltags") {
		mergedPics := mergedTags.Pictures
		mergedTags = &MP4Tags{}
		mergedTags.Pictures = mergedPics
	} else if containsStr(delStrings, "allcustom") {
		mergedTags.Custom = map[string]string{}
	}
	if containsStr(delStrings, "allothercustom") {
		mergedTags.OtherCustom = map[string][]string{}
	}

	if mergedTags.Custom == nil {
		mergedTags.Custom = map[string]string{}
	}
	if mergedTags.OtherCustom == nil {
		mergedTags.OtherCustom = map[string][]string{}
	}

	if containsStr(delStrings,  "album") {
		mergedTags.Album = ""
	}

	if containsStr(delStrings,  "albumartist") {
		mergedTags.AlbumArtist = ""
		mergedTags.AlbumArtists = nil
	}

	if containsStr(delStrings,  "albumartistsort") {
		mergedTags.AlbumArtistSort = ""
		mergedTags.AlbumArtistSorts = nil
	}

	if containsStr(delStrings,  "albumsort") {
		mergedTags.AlbumSort = ""
	}

	if containsStr(delStrings,  "artist") {
		mergedTags.Artist = ""
		mergedTags.Artists = nil
	}

	if containsStr(delStrings,  "artistsort") {
		mergedTags.ArtistSort = ""
		mergedTags.ArtistSorts = nil
	}

	if containsStr(delStrings,  "bpm") {
		mergedTags.BPM = 0
	}

	if containsStr(delStrings,  "comment") {
		mergedTags.Comment = ""
	}

	if containsStr(delStrings, "compilation") {
		mergedTags.Compilation = nil
	}

	if containsStr(delStrings,  "composer") {
		mergedTags.Composer = ""
		mergedTags.Composers = nil
	}

	if containsStr(delStrings,  "composersort") {
		mergedTags.ComposerSort = ""
		mergedTags.ComposerSorts = nil
	}

	if containsStr(delStrings,  "conductor") {
		mergedTags.Conductor = ""
		mergedTags.Conductors = nil
	}

	if containsStr(delStrings,  "copyright") {
		mergedTags.Copyright = ""
	}

	if containsStr(delStrings,  "customgenre") {
		mergedTags.CustomGenre = ""
		mergedTags.CustomGenres = nil
	}

	if containsStr(delStrings,  "date") {
		mergedTags.Date = ""
	}

	if containsStr(delStrings,  "description") {
		mergedTags.Description = ""
	}

	if containsStr(delStrings,  "director") {
		mergedTags.Director = ""
		mergedTags.Directors = nil
	}

	if containsStr(delStrings,  "discnumber") || containsStr(delStrings,  "disknumber") {
		mergedTags.DiscNumber = 0
	}

	if containsStr(delStrings,  "disctotal") || containsStr(delStrings,  "disktotal") {
		mergedTags.DiscTotal = 0
	}

	if containsStr(delStrings,  "genre") {
		mergedTags.Genre = GenreNone
	}

	if containsStr(delStrings,  "itunesadvisory") {
		mergedTags.ItunesAdvisory = ItunesAdvisoryNone
	}

	if containsStr(delStrings,  "itunesalbumid") {
		mergedTags.ItunesAlbumID = 0
	}	

	if containsStr(delStrings, "itunescatalogid") {
		mergedTags.ItunesCatalogID = 0
	}
	if containsStr(delStrings, "itunesstorefrontid") || containsStr(delStrings, "itunescountryid") {
		mergedTags.ItunesStorefrontID = 0
	}
	if containsStr(delStrings, "releasedate") {
		mergedTags.ReleaseDate = ""
	}
	if containsStr(delStrings, "xid") {
		mergedTags.XID = ""
	}
	if containsStr(delStrings, "itunesgapless") {
		mergedTags.ItunesGapless = nil
	}
	if containsStr(delStrings, "itunesgenreid") {
		mergedTags.ItunesGenreID = 0
		mergedTags.ItunesGenreIDs = nil
	}
	if containsStr(delStrings, "itunesmediatype") {
		mergedTags.ItunesMediaType = nil
	}

	if containsStr(delStrings,  "itunesartistid") {
		mergedTags.ItunesArtistID = 0
		mergedTags.ItunesArtistIDs = nil
	}

	if containsStr(delStrings,  "lyrics") {
		mergedTags.Lyrics = ""
	}

	if containsStr(delStrings,  "narrator") {
		mergedTags.Narrator = ""
	}

	if containsStr(delStrings,  "publisher") {
		mergedTags.Publisher = ""
	}

	if containsStr(delStrings,  "title") {
		mergedTags.Title = ""
	}

	if containsStr(delStrings,  "titlesort") {
		mergedTags.TitleSort = ""
	}

	if containsStr(delStrings,  "tracknumber") {
		mergedTags.TrackNumber = 0
	}

	if containsStr(delStrings,  "tracktotal") {
		mergedTags.TrackTotal = 0
	}

	if containsStr(delStrings,  "year") {
		mergedTags.Year = 0
	}

	if containsStr(delStrings, "allpictures") {
		mergedTags.Pictures = []*MP4Picture{}
	}

	if tags.Album != "" {
		mergedTags.Album = tags.Album
	}

	if tags.AlbumSort != "" {
		mergedTags.AlbumSort = tags.AlbumSort
	}

	if tags.AlbumArtists != nil {
		mergedTags.AlbumArtist = ""
		mergedTags.AlbumArtists = tags.AlbumArtists
	} else if tags.AlbumArtist != "" {
		mergedTags.AlbumArtist = tags.AlbumArtist
		mergedTags.AlbumArtists = nil
	}

	if tags.AlbumArtistSorts != nil {
		mergedTags.AlbumArtistSort = ""
		mergedTags.AlbumArtistSorts = tags.AlbumArtistSorts
	} else if tags.AlbumArtistSort != "" {
		mergedTags.AlbumArtistSort = tags.AlbumArtistSort
		mergedTags.AlbumArtistSorts = nil
	}

	if tags.Artists != nil {
		mergedTags.Artist = ""
		mergedTags.Artists = tags.Artists
	} else if tags.Artist != "" {
		mergedTags.Artist = tags.Artist
		mergedTags.Artists = nil
	}

	if tags.ArtistSorts != nil {
		mergedTags.ArtistSort = ""
		mergedTags.ArtistSorts = tags.ArtistSorts
	} else if tags.ArtistSort != "" {
		mergedTags.ArtistSort = tags.ArtistSort
		mergedTags.ArtistSorts = nil
	}

	if tags.BPM > 0 {
		mergedTags.BPM = tags.BPM
	}

	if tags.Comment != "" {
		mergedTags.Comment = tags.Comment
	}

	if tags.Compilation != nil {
		mergedTags.Compilation = tags.Compilation
	}

	if tags.Composers != nil {
		mergedTags.Composer = ""
		mergedTags.Composers = tags.Composers
	} else if tags.Composer != "" {
		mergedTags.Composer = tags.Composer
		mergedTags.Composers = nil
	}

	if tags.ComposerSorts != nil {
		mergedTags.ComposerSort = ""
		mergedTags.ComposerSorts = tags.ComposerSorts
	} else if tags.ComposerSort != "" {
		mergedTags.ComposerSort = tags.ComposerSort
		mergedTags.ComposerSorts = nil
	}
	mergedRoles := mergedTags.nativeCreditLists()
	for i, role := range tags.nativeCreditLists() {
		if *role.values != nil {
			*mergedRoles[i].values = *role.values
		}
	}
	if tags.Conductors != nil {
		mergedTags.Conductor = ""
	} else if tags.Conductor != "" {
		mergedTags.Conductor = tags.Conductor
		mergedTags.Conductors = nil
	}

	if tags.Copyright != "" {
		mergedTags.Copyright = tags.Copyright
	}

	if tags.CustomGenres != nil {
		mergedTags.CustomGenre = ""
		mergedTags.CustomGenres = tags.CustomGenres
	} else if tags.CustomGenre != "" {
		mergedTags.CustomGenre = tags.CustomGenre
		mergedTags.CustomGenres = nil
	}

	if tags.Date != "" {
		mergedTags.Date = tags.Date
	}

	if tags.Description != "" {
		mergedTags.Description = tags.Description
	}

	if tags.Directors != nil {
		mergedTags.Director = ""
		mergedTags.Directors = tags.Directors
	} else if tags.Director != "" {
		mergedTags.Director = tags.Director
		mergedTags.Directors = nil
	}

	if tags.DiscNumber > 0 {
		mergedTags.DiscNumber = tags.DiscNumber
	}

	if tags.DiscTotal > 0 {
		mergedTags.DiscTotal = tags.DiscTotal
	}

	if tags.ItunesAdvisory != ItunesAdvisoryNone {
		mergedTags.ItunesAdvisory = tags.ItunesAdvisory
	}

	if tags.ItunesAlbumID > 0 {
		mergedTags.ItunesAlbumID = tags.ItunesAlbumID
	}

	if containsStr(delStrings, "itunescomposerid") {
		mergedTags.ItunesComposerIDs = nil
	}
	if tags.ItunesArtistIDs != nil {
		mergedTags.ItunesArtistIDs = tags.ItunesArtistIDs
		mergedTags.ItunesArtistID = 0
	} else if tags.ItunesArtistID > 0 {
		mergedTags.ItunesArtistID = tags.ItunesArtistID
		mergedTags.ItunesArtistIDs = nil
	}
	if tags.ItunesComposerIDs != nil {
		mergedTags.ItunesComposerIDs = tags.ItunesComposerIDs
	}
	if tags.ItunesCatalogID != 0 {
		mergedTags.ItunesCatalogID = tags.ItunesCatalogID
	}
	if tags.ItunesStorefrontID != 0 {
		mergedTags.ItunesStorefrontID = tags.ItunesStorefrontID
	}
	if tags.ReleaseDate != "" {
		mergedTags.ReleaseDate = tags.ReleaseDate
	}
	if tags.XID != "" {
		mergedTags.XID = tags.XID
	}
	if tags.ItunesGapless != nil {
		mergedTags.ItunesGapless = tags.ItunesGapless
	}
	if tags.ItunesGenreIDs != nil {
		mergedTags.ItunesGenreIDs = tags.ItunesGenreIDs
		mergedTags.ItunesGenreID = 0
	} else if tags.ItunesGenreID != 0 {
		mergedTags.ItunesGenreID = tags.ItunesGenreID
		mergedTags.ItunesGenreIDs = nil
	}
	if tags.ItunesMediaType != nil {
		mergedTags.ItunesMediaType = tags.ItunesMediaType
	}

	if tags.Lyrics != "" {
		mergedTags.Lyrics = tags.Lyrics
	}

	if tags.Narrator != "" {
		mergedTags.Narrator = tags.Narrator
	}

	if tags.Publisher != "" {
		mergedTags.Publisher = tags.Publisher
	}

	if tags.Title != "" {
		mergedTags.Title = tags.Title
	}

	if tags.TitleSort != "" {
		mergedTags.TitleSort = tags.TitleSort
	}

	if tags.TrackNumber > 0 {
		mergedTags.TrackNumber = tags.TrackNumber
	}

	if tags.TrackTotal > 0 {
		mergedTags.TrackTotal = tags.TrackTotal
	}

	if tags.Year > 0 {
		mergedTags.Year = tags.Year
	}	

	if tags.Genre != GenreNone {
		mergedTags.Genre = tags.Genre
	}

	for k, v := range tags.Custom {
		if containsStr(delStrings, "custom:"+k) {
			continue
		}
		if v != "" {
			mergedTags.Custom[k] = v
		}
	}

	for k, v := range tags.OtherCustom {
		if containsStr(delStrings, "custom:"+k) || len(v) < 1 {
			continue
		}
		_, ok := mergedTags.OtherCustom[k]
		if ok {
			mergedTags.OtherCustom[k] = append(mergedTags.OtherCustom[k], v...)
		} else {
			mergedTags.OtherCustom[k] = v
		}

	}

	var filteredPics []*MP4Picture

	for idx, p := range mergedTags.Pictures {
		if !containsStr(delStrings,  fmt.Sprintf("picture:%d", idx+1)) {
			filteredPics = append(filteredPics, p)
		}
	}

	for _, p := range tags.Pictures {
		filteredPics = append(filteredPics, p)
	}

	mergedTags.Pictures = filteredPics
	return mergedTags
}

func putI16BE(n int16) []byte {
	buf := make([]byte, 2)

	binary.BigEndian.PutUint16(buf, uint16(n))
	return buf
}

func putI32BE(n int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf
}

func (mp4 MP4) updateChunkOffsets(outF *os.File, boxes MP4Boxes, oldIlistSize, newIlistSize int64) error {
	return mp4.updateChunkOffsetsAt(outF, boxes, newIlistSize-oldIlistSize, 0)
}

func (mp4 MP4) updateChunkOffsetsAt(outF *os.File, boxes MP4Boxes, delta, shiftFrom int64) error {
	if delta == 0 {
		return nil
	}

	stcos := boxes.getBoxesByPath("moov.trak.mdia.minf.stbl.stco")
	co64s := boxes.getBoxesByPath("moov.trak.mdia.minf.stbl.co64")
	if len(stcos) == 0 && len(co64s) == 0 {
		return &ErrBoxNotPresent{Msg: "moov.trak.mdia.minf.stbl.stco or co64 box not present"}
	}

	for _, stco := range stcos {
		if err := mp4.updateStcoChunkOffsets(outF, stco, delta, shiftFrom); err != nil {
			return fmt.Errorf("update stco offsets: %w", err)
		}
	}
	for _, co64 := range co64s {
		if err := mp4.updateCo64ChunkOffsets(outF, co64, delta, shiftFrom); err != nil {
			return fmt.Errorf("update co64 offsets: %w", err)
		}
	}
	return nil
}

func (mp4 MP4) chunkOffsetCount(box *MP4Box) (int64, error) {
	if box == nil {
		return 0, &ErrBoxNotPresent{Msg: "chunk offset box not present"}
	}
	if _, err := mp4.f.Seek(box.StartOffset+12, io.SeekStart); err != nil {
		return 0, err
	}
	var raw [4]byte
	if _, err := io.ReadFull(mp4.f, raw[:]); err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint32(raw[:])), nil
}

func (mp4 MP4) updateStcoChunkOffsets(outF *os.File, box *MP4Box, delta, shiftFrom int64) error {
	count, err := mp4.chunkOffsetCount(box)
	if err != nil {
		return err
	}
	if box.BoxSize != count*4+16 {
		return &ErrInvalidStcoSize{}
	}

	if _, err := mp4.f.Seek(box.StartOffset+16, io.SeekStart); err != nil {
		return err
	}
	if _, err := outF.Seek(chunkTablePosition(box, delta, shiftFrom), io.SeekStart); err != nil {
		return err
	}

	const maxUint32 = uint64(^uint32(0))
	var raw [4]byte
	for i := int64(0); i < count; i++ {
		if _, err := io.ReadFull(mp4.f, raw[:]); err != nil {
			return err
		}
		offset := uint64(binary.BigEndian.Uint32(raw[:]))
		adjusted, err := adjustChunkOffset(offset, delta)
		if err != nil {
			return err
		}
		if adjusted > maxUint32 {
			return fmt.Errorf("offset %d exceeds 32-bit stco range", adjusted)
		}
		binary.BigEndian.PutUint32(raw[:], uint32(adjusted))
		if _, err := outF.Write(raw[:]); err != nil {
			return err
		}
	}
	return nil
}

func (mp4 MP4) updateCo64ChunkOffsets(outF *os.File, box *MP4Box, delta, shiftFrom int64) error {
	count, err := mp4.chunkOffsetCount(box)
	if err != nil {
		return err
	}
	if box.BoxSize != count*8+16 {
		return fmt.Errorf("co64 size is invalid")
	}

	if _, err := mp4.f.Seek(box.StartOffset+16, io.SeekStart); err != nil {
		return err
	}
	if _, err := outF.Seek(chunkTablePosition(box, delta, shiftFrom), io.SeekStart); err != nil {
		return err
	}

	var raw [8]byte
	for i := int64(0); i < count; i++ {
		if _, err := io.ReadFull(mp4.f, raw[:]); err != nil {
			return err
		}
		offset := binary.BigEndian.Uint64(raw[:])
		adjusted, err := adjustChunkOffset(offset, delta)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint64(raw[:], adjusted)
		if _, err := outF.Write(raw[:]); err != nil {
			return err
		}
	}
	return nil
}

func adjustChunkOffset(offset uint64, delta int64) (uint64, error) {
	if delta >= 0 {
		add := uint64(delta)
		if offset > ^uint64(0)-add {
			return 0, fmt.Errorf("chunk offset overflows uint64")
		}
		return offset + add, nil
	}

	sub := uint64(-delta)
	if offset < sub {
		return 0, fmt.Errorf("chunk offset %d is smaller than adjustment %d", offset, sub)
	}
	return offset - sub, nil
}

func (mp4 MP4) readToOffset(f *os.File, startOffset int64) error {
	_, err := mp4.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	buf := make([]byte, BufSize)
	var totalRead int64

	for {
		read, err := mp4.f.Read(buf)
		if err != nil { 
			if err == io.EOF {
                break
            }
            return err
		}
		readI64 := int64(read)
		totalRead += readI64
		if totalRead > startOffset {
			_, err = f.Write(buf[:readI64+startOffset-totalRead])
			if err != nil {
				return err
			}
			break
		}
		_, err = f.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}


func writeRegular(f *os.File, boxName, val string, prefix bool) error {
	// boxSize := utf8.RuneCountInString(val) + 24
	valBytes := []byte(val)
	boxSize := len(valBytes) + 24
	boxSizeI32 := int32(boxSize)
	boxSizeBytes := putI32BE(boxSizeI32)
	_, err := f.Write(boxSizeBytes)
	if err != nil {
		return err
	}
	if prefix {
		_, err = f.Write([]byte{0xA9})
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString(boxName)
	if err != nil {
		return err
	}
	boxSizeBytes = putI32BE(boxSizeI32-8)
	_, err = f.Write(boxSizeBytes)
	if err != nil {
		return err
	}	
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}

	_, err = f.Write(
		[]byte{0x0, 0x0, 0x0, 0x01, 0x0, 0x0, 0x0, 0x0})
	if err != nil {
		return err
	}

	_, err = f.Write(valBytes)
	return err
}

func writeTextValues(f *os.File, name string, values []string) error {
	size := uint64(8)
	for _, value := range values {
		size += uint64(len(value)) + 16
	}
	if size > math.MaxUint32 {
		return fmt.Errorf("text metadata item %q exceeds the 32-bit box size", name)
	}

	if _, err := f.Write(putI32BE(int32(size))); err != nil {
		return err
	}
	if _, err := f.WriteString(name); err != nil {
		return err
	}
	// Text items store each value in its own UTF-8 data box.
	for _, value := range values {
		if _, err := f.Write(putI32BE(int32(len(value) + 16))); err != nil {
			return err
		}
		if _, err := f.WriteString("data\x00\x00\x00\x01\x00\x00\x00\x00"); err != nil {
			return err
		}
		if _, err := f.WriteString(value); err != nil {
			return err
		}
	}
	return nil
}

func writeGenre(f *os.File, genre Genre) error {
	_, err := f.Write([]byte{0x0, 0x0, 0x0, 0x1A})
	if err != nil {
		return err
	}
	_, err = f.WriteString("gnre")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x12})
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}
	_, err = f.Write(bytes.Repeat([]byte{0x0}, 9))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{byte(genre)})
	return err
}

func writeTrknDisc(f *os.File, n, total int16, isTrkn bool) error {
	var boxSize int32 = 30
	if n < 0 {
		n = 0
	}
	if total < 0 {
		total = 0
	}
	if isTrkn {
		boxSize += 2
	}
	boxSizeBytes := putI32BE(boxSize)
	_, err := f.Write(boxSizeBytes)
	if err != nil {
		return err
	}
	if isTrkn {
		_, err = f.WriteString("trkn")
	} else {
		_, err = f.WriteString("disk")
	}
	if err != nil {
		return err
	}
	boxSizeBytes = putI32BE(boxSize-8)
	if err != nil {
		return err
	}
	_, err = f.Write(boxSizeBytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")	
	if err != nil {
		return err
	}
	_, err = f.Write(bytes.Repeat([]byte{0x0}, 10))
	if err != nil {
		return err
	}

	nBytes := putI16BE(n)
	_, err = f.Write(nBytes)
	if err != nil {
		return err
	}
	totalBytes := putI16BE(total)
	_, err = f.Write(totalBytes)
	if err != nil {
		return err
	}
	if isTrkn {
		_, err = f.Write([]byte{0x0, 0x0})
		return err
	}
	return nil
}

func writeBPM(f *os.File, bpm int16) error {
	_, err := f.Write([]byte{0x0, 0x0, 0x0, 0x1A})
	if err != nil {
		return err
	}
	_, err = f.WriteString("tmpo")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x12})
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}
	_, err = f.Write(
		[]byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0})
	if err != nil {
		return err
	}
	bpmBytes := putI16BE(bpm)
	_, err = f.Write(bpmBytes)
	return err
}

func writeCompilation(f *os.File) error {
	// cpil contains a one-byte integer in its data box.
	_, err := f.WriteString("\x00\x00\x00\x19cpil\x00\x00\x00\x11data\x00\x00\x00\x15\x00\x00\x00\x00\x01")
	return err
}

func writeAdvisory(f *os.File, advisory ItunesAdvisory) error {
	_, err := f.Write([]byte{0x0, 0x0, 0x0, 0x19})
	if err != nil {
		return err
	}
	_, err = f.WriteString("rtng")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x11})
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0})
	if err != nil {
		return err
	}	
	_, err = f.Write([]byte{byte(advisory)})
	return err
}

func writeItunesInteger(f *os.File, name string, value []byte) error {
	var header [24]byte
	binary.BigEndian.PutUint32(header[0:4], uint32(len(header)+len(value)))
	copy(header[4:8], name)
	binary.BigEndian.PutUint32(header[8:12], uint32(16+len(value)))
	copy(header[12:16], "data")
	binary.BigEndian.PutUint32(header[16:20], 21) // iTunes integer data type.
	if _, err := f.Write(header[:]); err != nil {
		return err
	}
	_, err := f.Write(value)
	return err
}

func writeItunesAlbumID(f *os.File, albumID int32) error {
	_, err := f.Write([]byte{0x0, 0x0, 0x0, 0x20})
	if err != nil {
		return err
	}	
	_, err = f.WriteString("plID")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x18})
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}
	_, err = f.Write(
		[]byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	if err != nil {
		return err
	}
	albumIDBytes := putI32BE(albumID)
	_, err = f.Write(albumIDBytes)
	return err
}

func writeItunesIDs(f *os.File, name string, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	var data bytes.Buffer
	for _, id := range ids {
		if id == 0 || id > math.MaxInt64 {
			return fmt.Errorf("invalid iTunes ID for %s: %d", name, id)
		}
		width := 4
		if id > math.MaxInt32 {
			width = 8
		}
		var value [8]byte
		binary.BigEndian.PutUint64(value[:], id)
		var header [16]byte
		binary.BigEndian.PutUint32(header[:4], uint32(16+width))
		copy(header[4:8], "data")
		binary.BigEndian.PutUint32(header[8:12], 21)
		data.Write(header[:])
		data.Write(value[8-width:])
	}
	var header [8]byte
	binary.BigEndian.PutUint32(header[:4], uint32(len(header)+data.Len()))
	copy(header[4:], name)
	if _, err := f.Write(header[:]); err != nil {
		return err
	}
	_, err := data.WriteTo(f)
	return err
}

func writeCustom(f *os.File, name, value string, upper bool, others map[string][]string) error {	
	valueBytes := []byte(value)
	valueSize := len(valueBytes)

	if upper {
		name = strings.ToUpper(name)
	}
	nameUpperBytes := []byte(name)
	nameSize := len(nameUpperBytes)

	totalSize := nameSize+valueSize+64
	otherCust, _ := others[name]

	for _, other := range otherCust {
		totalSize += len([]byte(other)) + 16
	}

	// 48
	sizeBytes := putI32BE(int32(totalSize))
	_, err := f.Write(sizeBytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("----")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{0x0, 0x0, 0x0, 0x1C})
	if err != nil {
		return err
	}	
	_, err = f.WriteString("mean")
	if err != nil {
		return err
	}
	_, err = f.Write(bytes.Repeat([]byte{0x0}, 4))
	if err != nil {
		return err
	}
	_, err = f.WriteString("com.apple.iTunes")
	if err != nil {
		return err
	}
	sizeBytes = putI32BE(int32(nameSize)+12)
	_, err = f.Write(sizeBytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("name")
	if err != nil {
		return err
	}
	_, err = f.Write(bytes.Repeat([]byte{0x0}, 4))
	if err != nil {
		return err
	}
	_, err = f.Write(nameUpperBytes)
	if err != nil {
		return err
	}
	sizeBytes = putI32BE(int32(valueSize)+16)
	_, err = f.Write(sizeBytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("data")
	if err != nil {
		return err
	}
	_, err = f.Write(
		[]byte{0x0, 0x0, 0x0, 0x01, 0x0, 0x0, 0x0, 0x0})
	if err != nil {
		return err
	}
	_, err = f.Write(valueBytes)
	if err != nil {
		return err
	}


	for _, v := range otherCust {
		valueBytes = []byte(v)
		valueSize = len(valueBytes)
		sizeBytes = putI32BE(int32(valueSize)+16)
		_, err = f.Write(sizeBytes)
		if err != nil {
			return err
		}
		_, err = f.WriteString("data")
		if err != nil {
			return err
		}
		_, err = f.Write(
			[]byte{0x0, 0x0, 0x0, 0x01, 0x0, 0x0, 0x0, 0x0})
		if err != nil {
			return err
		}
		_, err = f.Write(valueBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPicFormat(imageType ImageType, magic []byte) uint8 {
	if imageType == ImageTypeAuto {
		if bytes.Equal(magic, []byte{0x89, 0x50, 0x4E, 0x47}) {
			return 0xE
		}
	}
	if imageType == ImageTypePNG {
		return 0xE
	}
	return 0x0D
}

func writePics(f *os.File, pics []*MP4Picture) error {
	var boxSize int32 = 8
	for _, pic := range pics {
		dataSize := len(pic.Data)
		if dataSize < 1 {
			continue
		}
		boxSize += int32(dataSize + 16)
	}
	if boxSize == 8 {
		return nil
	}

	boxSizeBytes := putI32BE(boxSize)
	_, err := f.Write(boxSizeBytes)
	if err != nil {
		return err
	}
	_, err = f.WriteString("covr")
	if err != nil {
		return err
	}

	for _, pic := range pics {
		dataSize := len(pic.Data)
		if dataSize < 1 {
			continue
		}
		boxSizeBytes = putI32BE(int32(dataSize+16))
		_, err = f.Write(boxSizeBytes)
		if err != nil {
			return err
		}
		_, err = f.WriteString("data")
		if err != nil {
			return err
		}

		format := getPicFormat(pic.Format, pic.Data[:4])
		_, err = f.Write([]byte{0x0, 0x0, 0x0, format, 0x0, 0x0, 0x0, 0x0})
		if err != nil {
			return err
		}
		_, err = f.Write(pic.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func resizeBoxes(f *os.File, boxes MP4Boxes, oldIlstSize, newIlistSize int64) error {
	moov := boxes.getBoxByPath("moov")
	udta := boxes.getBoxByPath("moov.udta")
	meta := boxes.getBoxByPath("moov.udta.meta")

	// Sizes are rewritten in place as 32-bit fields, which only works for
	// boxes with a normal 8-byte header whose new size still fits.
	for _, resized := range []struct {
		box     *MP4Box
		newSize int64
	}{
		{nil, newIlistSize},
		{moov, moov.BoxSize - oldIlstSize + newIlistSize},
		{udta, udta.BoxSize - oldIlstSize + newIlistSize},
		{meta, meta.BoxSize - oldIlstSize + newIlistSize},
	} {
		if resized.box != nil && resized.box.HeaderSize != 8 {
			return fmt.Errorf("%s box uses a 64-bit size header; resizing it is not supported", resized.box.Path)
		}
		if resized.newSize > math.MaxUint32 {
			return fmt.Errorf("resized box size %d exceeds 32 bits", resized.newSize)
		}
	}

	sizeBytes := putI32BE(int32(newIlistSize))
	_, err := f.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = f.Seek(moov.StartOffset, io.SeekStart)
	if err != nil {
		return err
	}
	newMoovSize := moov.BoxSize - oldIlstSize + newIlistSize
	sizeBytes = putI32BE(int32(newMoovSize))
	_, err = f.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = f.Seek(udta.StartOffset, io.SeekStart)
	if err != nil {
		return err
	}
	newUdtaSize := udta.BoxSize - oldIlstSize + newIlistSize
	sizeBytes = putI32BE(int32(newUdtaSize))
	_, err = f.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = f.Seek(meta.StartOffset, io.SeekStart)
	if err != nil {
		return err
	}
	newMetaSize := meta.BoxSize - oldIlstSize + newIlistSize
	sizeBytes = putI32BE(int32(newMetaSize))
	_, err = f.Write(sizeBytes)
	return err
}

func (mp4 MP4) writeRemaining(f *os.File) error {
	buf := make([]byte, BufSize)
	for {
		read, err := mp4.f.Read(buf)
		if err != nil {
			if err == io.EOF {
                break
            }
			return err
		}
		if read < BufSize {
			_, err := f.Write(buf[:read])
			if err != nil {
				return err
			}
		} else {
			_, err = f.Write(buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (mp4 MP4) writeTags(boxes MP4Boxes, tags *MP4Tags, tempPath string) error {
	ilst := boxes.getBoxByPath("moov.udta.meta.ilst")
    oldIlstSize := ilst.BoxSize
    f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
    	return err
    }
    defer f.Close()
    err = mp4.readToOffset(f, ilst.StartOffset)
    if err != nil {
    	return err
    }
    ilstStartOffset, err := getPos(f)
    if err != nil {
    	return err
    }

   	_, err = f.Write(bytes.Repeat([]byte{0x0}, 4))
	if err != nil {
		return err
	}
	_, err = f.WriteString("ilst")
	if err != nil {
		return err
	}
	if tags.Title != "" {
		err = writeRegular(f, "nam", tags.Title, true)
		if err != nil {
			return err
		}
	}
	if tags.TitleSort != "" {
		err = writeRegular(f, "sonm", tags.TitleSort, false)
		if err != nil {
			return err
		}
	}
	if tags.Album != "" {
		err = writeRegular(f, "alb", tags.Album, true)
		if err != nil {
			return err
		}
	}
	if tags.AlbumSort != "" {
		err = writeRegular(f, "soal", tags.AlbumSort, false)
		if err != nil {
			return err
		}
	}

	if len(tags.AlbumArtists) > 0 {
		if err := writeTextValues(f, "aART", tags.AlbumArtists); err != nil {
			return err
		}
	} else if tags.AlbumArtist != "" {
		err = writeRegular(f, "aART", tags.AlbumArtist, false)
		if err != nil {
			return err
		}
	}

	if len(tags.AlbumArtistSorts) > 0 {
		if err := writeTextValues(f, "soaa", tags.AlbumArtistSorts); err != nil {
			return err
		}
	} else if tags.AlbumArtistSort != "" {
		err = writeRegular(f, "soaa", tags.AlbumArtistSort, false)
		if err != nil {
			return err
		}
	}

	if len(tags.Artists) > 0 {
		if err := writeTextValues(f, "\xa9ART", tags.Artists); err != nil {
			return err
		}
	} else if tags.Artist != "" {
		err = writeRegular(f, "ART", tags.Artist, true)
		if err != nil {
			return err
		}
	}

	if len(tags.ArtistSorts) > 0 {
		if err := writeTextValues(f, "soar", tags.ArtistSorts); err != nil {
			return err
		}
	} else if tags.ArtistSort != "" {
		err = writeRegular(f, "soar", tags.ArtistSort, false)
		if err != nil {
			return err
		}
	}

	if tags.Comment != "" {
		err = writeRegular(f, "cmt", tags.Comment, true)
		if err != nil {
			return err
		}
	}

	if len(tags.Directors) > 0 {
		if err := writeTextValues(f, "\xa9dir", tags.Directors); err != nil {
			return err
		}
	} else if tags.Director != "" {
		if err := writeRegular(f, "dir", tags.Director, true); err != nil {
			return err
		}
	}

	if len(tags.Composers) > 0 {
		if err := writeTextValues(f, "\xa9wrt", tags.Composers); err != nil {
			return err
		}
	} else if tags.Composer != "" {
		err = writeRegular(f, "wrt", tags.Composer, true)
		if err != nil {
			return err
		}
	}

	if len(tags.ComposerSorts) > 0 {
		if err := writeTextValues(f, "soco", tags.ComposerSorts); err != nil {
			return err
		}
	} else if tags.ComposerSort != "" {
		err = writeRegular(f, "soco", tags.ComposerSort, false)
		if err != nil {
			return err
		}
	}

	if tags.Copyright != "" {
		err = writeRegular(f, "cprt", tags.Copyright, false)
		if err != nil {
			return err
		}
	}

	if tags.Lyrics != "" {
		err = writeRegular(f, "lyr", tags.Lyrics, true)
		if err != nil {
			return err
		}
	}

	if len(tags.CustomGenres) != 0 {
		if err := writeTextValues(f, "\xa9gen", tags.CustomGenres); err != nil {
			return err
		}
	} else if tags.CustomGenre != "" {
		if err := writeRegular(f, "gen", tags.CustomGenre, true); err != nil {
			return err
		}
	}

	if tags.Description != "" {
		err = writeRegular(f, "desc", tags.Description, false)
		if err != nil {
			return err
		}
	}

	if tags.Publisher != "" {
		err = writeRegular(f, "pub", tags.Publisher, true)
		if err != nil {
			return err
		}
	}	

	if len(tags.Conductors) == 0 && tags.Conductor != "" {
		if err := writeRegular(f, "con", tags.Conductor, true); err != nil {
			return err
		}
	}
	for _, role := range tags.nativeCreditLists() {
		if len(*role.values) != 0 {
			if err := writeTextValues(f, role.atom, *role.values); err != nil {
				return err
			}
		}
	}

	if tags.Compilation != nil && *tags.Compilation {
		if err := writeCompilation(f); err != nil {
			return err
		}
	}

	if tags.ItunesAdvisory != ItunesAdvisoryNone {
		err = writeAdvisory(f, tags.ItunesAdvisory)
		if err != nil {
			return err
		}
	}

	if tags.ItunesAlbumID > 0  {
		err = writeItunesAlbumID(f, tags.ItunesAlbumID)
		if err != nil {
			return err
		}
	}

	artistIDs := tags.ItunesArtistIDs
	if len(artistIDs) == 0 && tags.ItunesArtistID > 0 {
		artistIDs = []uint64{uint64(tags.ItunesArtistID)}
	}
	if err := writeItunesIDs(f, "atID", artistIDs); err != nil {
		return err
	}
	if err := writeItunesIDs(f, "cmID", tags.ItunesComposerIDs); err != nil {
		return err
	}
	if tags.ItunesCatalogID != 0 {
		var catalogID [4]byte
		binary.BigEndian.PutUint32(catalogID[:], tags.ItunesCatalogID)
		if err := writeItunesInteger(f, "cnID", catalogID[:]); err != nil {
			return err
		}
	}
	if tags.ItunesStorefrontID != 0 {
		var storefrontID [4]byte
		binary.BigEndian.PutUint32(storefrontID[:], tags.ItunesStorefrontID)
		if err := writeItunesInteger(f, "sfID", storefrontID[:]); err != nil {
			return err
		}
	}
	if tags.ReleaseDate != "" {
		if err := writeRegular(f, "rldt", tags.ReleaseDate, false); err != nil {
			return err
		}
	}
	if tags.XID != "" {
		if err := writeRegular(f, "xid ", tags.XID, false); err != nil {
			return err
		}
	}
	if tags.ItunesGapless != nil {
		var value byte
		if *tags.ItunesGapless {
			value = 1
		}
		if err := writeItunesInteger(f, "pgap", []byte{value}); err != nil {
			return err
		}
	}
	genreIDs := tags.ItunesGenreIDs
	if len(genreIDs) == 0 && tags.ItunesGenreID != 0 {
		genreIDs = []uint64{uint64(tags.ItunesGenreID)}
	}
	if err := writeItunesIDs(f, "geID", genreIDs); err != nil {
		return err
	}
	if tags.ItunesMediaType != nil {
		if err := writeItunesInteger(f, "stik", []byte{byte(*tags.ItunesMediaType)}); err != nil {
			return err
		}
	}

	if tags.TrackNumber > 0 || tags.TrackTotal > 0 {
		err = writeTrknDisc(f, tags.TrackNumber, tags.TrackTotal, true)
		if err != nil {
			return err
		}		
	}

	if tags.DiscNumber > 0 || tags.DiscTotal > 0 {
		err = writeTrknDisc(f, tags.DiscNumber, tags.DiscTotal, false)
		if err != nil {
			return err
		}		
	}
	
	if tags.BPM > 0 {
		err = writeBPM(f, tags.BPM)
		if err != nil {
			return err
		}		
	}

	if tags.Year > 0 {
		err = writeRegular(f, "day", strconv.Itoa(int(tags.Year)), true)
		if err != nil {
			return err
		}
	} else if tags.Date != "" {
		err = writeRegular(f, "day", tags.Date, true)
		if err != nil {
			return err
		}	
	}

	if tags.Genre != GenreNone {
		err = writeGenre(f, tags.Genre)
		if err != nil {
			return err
		}		
	}
	
	for k, v := range tags.Custom {
		err = writeCustom(f, k, v, mp4.upperCustom, tags.OtherCustom)
		if err != nil {
			return err
		}
	}

	err = writePics(f, tags.Pictures)
	if err != nil {
		return err
	}

	newIlstEndOffset, err := getPos(f)	
	if err != nil {
		return err
	}
	newIlstSize := newIlstEndOffset - ilstStartOffset
	_, err = f.Seek(ilstStartOffset, io.SeekStart)
	if err != nil {
		return err
	}

	err = resizeBoxes(f, boxes, oldIlstSize, newIlstSize)
	if err != nil {
		return err
	}

	mdat := boxes.getBoxByPath("mdat")
	if mdat.StartOffset > ilstStartOffset && oldIlstSize != newIlstSize {
		err = mp4.updateChunkOffsets(f, boxes, oldIlstSize, newIlstSize)
		if err != nil {
			return err
		}	
	}

	_, err = f.Seek(newIlstEndOffset, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = mp4.f.Seek(ilst.EndOffset, io.SeekStart)
	if err != nil {
		return err
	}
	err = mp4.writeRemaining(f)
	return err
}

func (mp4 *MP4) replaceCustomValues(tags *MP4Tags, replacements map[string][]string) error {
	seen := make(map[string]bool, len(replacements))
	for key, values := range replacements {
		canonical := key
		if mp4.upperCustom {
			canonical = strings.ToUpper(key)
		}
		if seen[canonical] {
			return fmt.Errorf("duplicate replacement for custom tag %q", canonical)
		}
		seen[canonical] = true

		// The legacy inputs may have supplied a differently cased spelling.
		matches := func(name string) bool {
			if mp4.upperCustom {
				return strings.ToUpper(name) == canonical
			}
			return name == canonical
		}
		for name := range tags.Custom {
			if matches(name) {
				delete(tags.Custom, name)
			}
		}
		for name := range tags.OtherCustom {
			if matches(name) {
				delete(tags.OtherCustom, name)
			}
		}
		if len(values) != 0 {
			tags.Custom[canonical] = values[0]
			tags.OtherCustom[canonical] = values[1:]
		}
	}
	return nil
}

func (mp4 *MP4) actualWrite(tags *MP4Tags, _delStrings []string) error {
	delStrings := strArrToLower(_delStrings)

	mergedTags, boxes, err := mp4.actualRead()
	if err != nil {
		return err
	}
	if boxes.getBoxByPath("moov.udta.meta.ilst") == nil {
		return &ErrBoxNotPresent{Msg: "ilst box not present, implement me"}
	}
	userData := userDataChanges(mergedTags, tags, delStrings)
	mergedTags = overwriteTags(mergedTags, tags, delStrings)
	if err := mp4.replaceCustomValues(mergedTags, tags.ReplaceCustomValues); err != nil {
		return err
	}
   	tempPath := getTempPath(mp4.path)
   	err = mp4.writeTags(boxes, mergedTags, tempPath)
	if err != nil {
		return err
	}
	if len(userData) != 0 {
		if err := rewriteUserData(tempPath, userData); err != nil {
			_ = os.Remove(tempPath)
			return err
		}
	}
	mp4.Close()
	err = moveMP4(tempPath, mp4.path)
	if err != nil {
		return err
	}

	m, err := Open(mp4.path)
	if err != nil {
		return err
	}
	mp4.f = m.f
	mp4.size = m.size
	return nil
}

