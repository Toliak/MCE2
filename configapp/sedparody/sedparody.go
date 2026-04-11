package sedparody

import (
	"bufio"
	"os"
	"regexp"
)

type ReplacerReader interface {
	HasNext() bool
	Next() (string, error)
}

type ReplacerWriter interface {
	WriteString(str string) (int, error)
}

type Replacer struct {
	// Line by line reader or the whole file reader
	reader ReplacerReader
	writer ReplacerWriter
}

type scannerReplacerReader struct {
    scanner *bufio.Scanner
}

func (s *scannerReplacerReader) HasNext() bool {
    return s.scanner.Scan()
}

func (s *scannerReplacerReader) Next() (string, error) {
    return s.scanner.Text(), nil
}

func ScannerToReplacerReader(scanner *bufio.Scanner) ReplacerReader {
	return &scannerReplacerReader{scanner: scanner}
}

type fullFileReader struct {
	path string
	alreadyRead bool
}

func (s *fullFileReader) HasNext() bool {
    return !s.alreadyRead
}

func (s *fullFileReader) Next() (string, error) {
	s.alreadyRead = true
    b, err := os.ReadFile(s.path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func FullFileReader(path string) ReplacerReader {
	return &fullFileReader{
		path: path,
		alreadyRead: false,
	}
}

type ReplacerOption func(*Replacer)

func NewReplacer(
	reader ReplacerReader,
	// writer ReplacerWriter,
	opts ...ReplacerOption,
) *Replacer {
	r := &Replacer{
		reader: reader,
		// writer: writer,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (f *Replacer) Replace(re *regexp.Regexp, replaceBy string, maxTimes int) ([]string, int, error) {
	replaceCtr := 0
	blocks := make([]string, 0)

	for f.reader.HasNext() {
		textBlock, err := f.reader.Next()
		if err != nil {
			return nil, 0, err
		}
		if maxTimes != 0 && replaceCtr < maxTimes {
			// fmt.Printf("Replaced string: '%s'\n", textBlock)

			textBlockNew := re.ReplaceAllString(textBlock, replaceBy)

			// TODO: I would like to have something like [regexp.Regexp.ReplaceAllString]
			// but with the boolean flag (replaced or not)
			if textBlockNew != textBlock {
				replaceCtr += 1
			}
			textBlock = textBlockNew
		}
        blocks = append(blocks, textBlock)
    }

	return blocks, replaceCtr, nil
}
