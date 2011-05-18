package filesystem

const (
	_ = iota
	// Entries with this type incapsulates playable track objects.
	TypeTrack = iota
	// Entries with this type incapsulates directory objects.
	TypeDirectory = iota
)

// Entry represents any filesystem entry object.
// For example, it can be playable track, directory, etc.
type Entry struct {
	// Type of the Entry.
	t int
	// Incapsulated object (e. g. Track or Directory).
	item interface{}
}

// NewEntry returns newly created and initialized Entry object.
// t is the type of the incapsulated object.
// item is incapsulated object itself.
func NewEntry(t int, item interface{}) *Entry {
	return &Entry{t, item}
}

// Type returns type of the object incapsulated in the Entry.
func (e *Entry) Type() int {
	return e.t
}

// TypeString returns string representaion of the type.
func (e *Entry) TypeString() string {
	typeDescriptions := map[int]string{
		TypeTrack:     "TRACK",
		TypeDirectory: "DIRECTORY",
	}

	// XXX: Check to be sure.
	str, ok := typeDescriptions[e.t]
	if !ok {
		panic("Unknown type recived")
	}

	return str
}

// Track returns Track object incapsulated by Entry.
// Before calling this method you should be sure that Type returns
// TypeTrack, otherwise panic will happend.
func (e *Entry) Track() *Track {
	if e.t != TypeTrack {
		panic("Entry doesn't incapsulate Track object")
	}

	return e.item.(*Track)
}

// Directory returns Directory object incapsulated by Entry.
// Before calling this method you should be sure that Entry incapsulates
// Directory object.
func (e *Entry) Directory() *Directory {
	if e.t != TypeDirectory {
		panic("Entry doesn't incapsulate Directory object")
	}

	return e.item.(*Directory)
}
