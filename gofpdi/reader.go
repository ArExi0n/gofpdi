package gofpdi

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/errors"
)

type PDFreader struct {
	availableBoxes []string
	stack          []string
	trailer        *PdfValue
	catalog        *PdfValue
	pages          []*PdfValue
	xrefPos        int
	xref           map[int]map[int]int
	xrefStream     map[int][2]int
	f              io.ReadSeeker
	nbytes         int64
	sourceFile     string
}

func NewPdfReaderFromStream(rs io.ReadSeeker) (*PDFreader, error) {
	length, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to determine stream length")
	}

	_, err = rs.Seek(0, io.SeekStart)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to rewind stream")
	}

	parser := &PDFreader{f: rs, nbytes: length}
	if err := parser.init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize parser")
	}
	if err := parser.read(); err != nil {
		return nil, errors.Wrap(err, "Failed to read pdf from stream")
	}
	return parser, nil
}

func NewPdfReader(filename string) (*PDFreader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open file")
	}

	info, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain file information")
	}

	parser := &PDFreader{
		f:          f,
		sourceFile: filename,
		nbytes:     info.Size(),
	}

	if err = parser.init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize parser")
	}
	if err = parser.read(); err != nil {
		return nil, errors.Wrap(err, "failed to read pdf")
	}

	return parser, nil
}

func (this *PDFreader) init() error {
	this.availableBoxes = []string{
		"/MediaBox",
		"/CropBox",
		"/BleedBox",
		"/TrimBox",
		"/ArtBox",
	}
	this.xref = make(map[int]map[int]int)
	this.xrefStream = make(map[int][2]int)

	err := this.read()
	if err != nil {
		return errors.Wrap(err, "Failed to read pdf")
	}

	return nil
}

// stub so the file compiles
func (this *PDFreader) read() error {
	return nil
}

type PdfValue struct {
	Type       int
	String     string
	Token      string
	Int        int
	Real       float64
	Bool       bool
	Dictionary map[string]*PdfValue
	Array      []*PdfValue
	Id         int
	NewId      int
	Gen        int
	Value      *PdfValue
	Stream     *PdfValue
	Bytes      []byte
}

// jump over comments
func (this *PDFreader) skipComments(r *bufio.Reader) error {
	var err error
	var b byte

	for {
		b, err = r.ReadByte()
		if err != nil {
			return errors.Wrap(err, "Failed to ReadByte while skipping comments")
		}

		if b == '\n' || b == '\r' {
			if b == '\r' {
				// Peek and see if next char is \n
				b2, err := r.ReadByte()
				if err != nil {
					return errors.Wrap(err, "Failed to read byte")
				}
				if b2 != '\n' {
					r.UnreadByte()
				}
			}
			break
		}
	}
	return nil
}

// advance reader so that whitespase is ignored
func (this *PDFreader) skipWhitespace(r *bufio.Reader) error {
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to read byte")
		}

		if b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			continue
		} else {
			r.UnreadByte()
			break
		}
	}

	return nil
}

// read a token
func (this *PDFreader) readToken(r *bufio.Reader) (string, error) {
	var err error

	// If there is a token available on the stack, pop it out and returns it.
	if len(this.stack) > 0 {
		var popped string
		popped, this.stack = this.stack[len(this.stack)-1], this.stack[:len(this.stack)-1]
		return popped, nil
	}

	err = this.skipWhitespace(r)
	if err != nil {
		return "", errors.Wrap(err, "Failed to skip the whitespace")
	}
	return "", nil
}
