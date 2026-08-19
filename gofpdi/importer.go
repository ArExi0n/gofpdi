package gofpdi

type Importer struct {
	sourceFile string
	readers    map[string]*PDFreader
	writers    map[string]*PDFwriter
}
