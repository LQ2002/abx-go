package abx

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"testing"
)

type Package struct {
	Name     string `xml:"name,attr"`
	CodePath string `xml:"codePath,attr"`
	UserID   int    `xml:"userId,attr"`
}

func TestParsePackagesXML(t *testing.T) {

	data, err := os.ReadFile("packages.xml")
	if err != nil {
		t.Fatal(err)
	}

	reader, ok := NewReader(bytes.NewReader(data))
	if !ok {
		t.Fatal("file is not ABX format")
	}

	decoder := xml.NewTokenDecoder(reader)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		switch se := tok.(type) {

		case xml.StartElement:
			if se.Name.Local == "package" {

				var p Package
				if err := decoder.DecodeElement(&p, &se); err != nil {
					t.Fatal(err)
				}

				t.Logf("package=%s path=%s uid=%d",
					p.Name, p.CodePath, p.UserID)
			}
		}
	}
}
