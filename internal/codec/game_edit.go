package codec

import (
	"fmt"
	"sort"
)

type gameGap struct {
	tail, after []byte
}

type gameSpan struct {
	start, end int64
	index      int
}

// EncodeGameDirectoryEdit rebuilds an indexed edit from desired Entries and
// original Raw/Base. Unclaimed original bytes remain after their preceding
// physical payload; retainedTail preserves a replaced compressed stream's gap.
// Payload starts, including compressed terminal boundaries, are four-aligned.
func EncodeGameDirectoryEdit(directory GameDirectory, index int, retainedTail []byte, base int32) ([]byte, error) {
	count, directoryEnd, err := validateGameDirectory(directory.Raw, directory.Base)
	if err != nil {
		return nil, err
	}
	if len(directory.Entries) != count || index < 0 || index >= count {
		return nil, fmt.Errorf("game directory edit: entry count or index mismatch")
	}
	spans := make([]gameSpan, 0, count)
	for i := 0; i < count; i++ {
		entry := readGameEntry(directory.Raw, i)
		start, end, err := gameEntryRange(directory.Raw, directory.Base, count, directoryEnd, i, entry)
		if err != nil {
			return nil, err
		}
		if end > start {
			spans = append(spans, gameSpan{start, end, i})
		}
	}
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].start != spans[j].start {
			return spans[i].start < spans[j].start
		}
		return spans[i].end < spans[j].end
	})
	gaps := make([]gameGap, count)
	gaps[index].tail = retainedTail
	var leading []byte
	cursor, previous := directoryEnd, -1
	for _, span := range spans {
		if span.start > cursor {
			gap := directory.Raw[int(cursor):int(span.start):int(span.start)]
			if previous < 0 {
				leading = gap
			} else {
				gaps[previous].after = gap
			}
		}
		if span.end > cursor {
			cursor, previous = span.end, span.index
		}
	}
	if cursor < int64(len(directory.Raw)) {
		gap := directory.Raw[int(cursor):len(directory.Raw):len(directory.Raw)]
		if previous < 0 {
			leading = gap
		} else {
			gaps[previous].after = gap
		}
	}
	return encodeGameDirectory(directory.Entries, base, leading, gaps)
}
