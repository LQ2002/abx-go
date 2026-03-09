package abx

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"testing"
)

const maxUnsignedShort = 65535

func TestParsePackagesXMLDemo(t *testing.T) {
	abxData := buildPackagesDemoABX(t)

	reader, ok := NewReader(bytes.NewReader(abxData))
	if !ok {
		t.Fatal("expected ABX reader to be created")
	}

	decoder := xml.NewTokenDecoder(reader)

	var got struct {
		XMLName  xml.Name `xml:"packages"`
		Packages []struct {
			Name    string `xml:"name,attr"`
			Version int    `xml:"version,attr"`
			Enabled bool   `xml:"enabled,attr"`
		} `xml:"package"`
	}
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode packages.xml ABX: %v", err)
	}

	if got.XMLName.Local != "packages" {
		t.Fatalf("unexpected root tag: %q", got.XMLName.Local)
	}
	if len(got.Packages) != 1 {
		t.Fatalf("unexpected package count: %d", len(got.Packages))
	}

	pkg := got.Packages[0]
	if pkg.Name != "com.demo.app" {
		t.Fatalf("unexpected package name: %q", pkg.Name)
	}
	if pkg.Version != 123 {
		t.Fatalf("unexpected package version: %d", pkg.Version)
	}
	if !pkg.Enabled {
		t.Fatal("expected package enabled=true")
	}
}

func buildPackagesDemoABX(t *testing.T) []byte {
	t.Helper()

	var out bytes.Buffer
	out.Write([]byte{0x41, 0x42, 0x58, 0x00})

	out.WriteByte(0x00)
	out.WriteByte(0x02)
	writeInternedUTF(t, &out, "packages")

	out.WriteByte(0x02)
	writeInternedUTF(t, &out, "package")

	out.WriteByte(0x2f)
	writeInternedUTF(t, &out, "name")
	writeUTF(t, &out, "com.demo.app")

	out.WriteByte(0x6f)
	writeInternedUTF(t, &out, "version")
	if err := binary.Write(&out, binary.BigEndian, int32(123)); err != nil {
		t.Fatalf("write version: %v", err)
	}

	out.WriteByte(0xcf)
	writeInternedUTF(t, &out, "enabled")

	out.WriteByte(0x03)
	writeInternedUTF(t, &out, "package")

	out.WriteByte(0x03)
	writeInternedUTF(t, &out, "packages")

	out.WriteByte(0x01)
	return out.Bytes()
}

func writeInternedUTF(t *testing.T, out *bytes.Buffer, value string) {
	t.Helper()
	if err := binary.Write(out, binary.BigEndian, uint16(maxUnsignedShort)); err != nil {
		t.Fatalf("write string ref marker: %v", err)
	}
	writeUTF(t, out, value)
}

func writeUTF(t *testing.T, out *bytes.Buffer, value string) {
	t.Helper()
	if err := binary.Write(out, binary.BigEndian, uint16(len(value))); err != nil {
		t.Fatalf("write utf length: %v", err)
	}
	if _, err := out.WriteString(value); err != nil {
		t.Fatalf("write utf value: %v", err)
	}
}
