package gofpdi

import (
	"bufio"
	"bytes"
	"os"

	"github.com/pkg/errors"
)

type PdfWriter struct {
	f       *os.File
	w       *bufio.Writer
	r       *PDFreader
	k       float64
	tpls    []*PdfTemplate
	m       int
	n       int
	offsets map[int]int
	offset  int
	result  map[int]string
	// keep track of which objects have already been writen
	obj_stack       map[int]*PdfValue
	don_obj_stack   map[int]*PdfValue
	written_obj     map[*PdfObjectId][]byte
	written_obj_pos map[*PdfObjectId]map[int]string
	current_obj     *PdfObject
	current_obj_id  int
	tpl_id_offset   int
	use_hash        bool
}

type PdfObjectId struct {
	id   int
	hash string
}

type PdfObject struct {
	id     *PdfObjectId
	buffer *bytes.Buffer
}

func (this *PdfWriter) SetTpIdOffset(n int) {
	this.tpl_id_offset = n
}

func (this *PdfWriter) Init() {
	this.k = 1
	this.obj_stack = make(map[int]*PdfValue, 0)
	this.don_obj_stack = make(map[int]*PdfValue, 0)
	this.tpls = make([]*PdfTemplate, 0)
	this.written_obj = make(map[*PdfObjectId][]byte, 0)
	this.written_obj_pos = make(map[*PdfObjectId]map[int]string, 0)
	this.current_obj = new(PdfObject)
}

func (this *PdfWriter) SetUseHash(b bool) {
	this.use_hash = b
}

func (this *PdfWriter) SetNextObjectId(id int) {
	this.n = id - 1
}

func NewPdfWriter(filename string) (*PdfWriter, error) {
	writer := &PdfWriter{}
	writer.Init()

	if filename != "" {
		var err error
		f, err := os.Create(filename)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to create filename: "+filename)
		}
		writer.f = f
		writer.w = bufio.NewWriter(f)
	}
	return writer, nil
}

// done with parsing. Now, create template
type PdfTemplate struct {
	Id        int
	Reader    *PDFreader
	Resources *PdfValue
	Buffer    string
	Box       map[string]float64
	Boxes     map[string]map[string]float64
	X         float64
	Y         float64
	W         float64
	H         float64
	Rotation  int
	N         int
}

func (this *PdfWriter) GetImportedObjects() map[*PdfObjectId][]byte {
	return this.written_obj
}

// For each object (uniquely identified by a sha1 hash), return
// the positions of each hash within the object, to be replaced
// with pdf object ids (integers)
func (this *PdfWriter) GetImportedObjHashPos() map[*PdfObjectId]map[int]string {
	return this.written_obj_pos
}

func (this *PdfWriter) ClearImportedObjects() {
	this.written_obj = make(map[*PdfObjectId][]byte, 0)
}

// Create a Pdftemplate object from a page number and a box number
func (this *PdfWriter) ImportPage(reader *PDFreader, pageno int, boxname string) (int, error) {
	var err error

	// set default scale to 1
	this.k = 1

	// Get all the page boxes
	pageBoxes, err := reader.getPageBoxes(1, this.k)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to get page boxes")
	}

	return len(this.tpls) - 1, nil
}
