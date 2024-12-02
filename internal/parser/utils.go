package parser

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

func (p *Parser) IsTextFile(file *os.File) (bool, error) {
	// Read first 512 bytes for MIME detection
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return false, err
	}
	header = header[:n]

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return false, err
	}

	// Check MIME type
	mimeType := http.DetectContentType(header)
	isText := false
	for textType := range p.textMimeTypes {
		if strings.HasPrefix(mimeType, textType) {
			isText = true
			break
		}
	}
	if !isText {
		return false, nil
	}

	// Validate UTF-8 using a buffered reader
	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if r == utf8.RuneError {
			return false, nil
		}
	}

	// Reset file pointer again
	_, err = file.Seek(0, 0)
	return true, err
}
