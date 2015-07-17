package exif

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	err := exif.Open("_examples/resources/test.jpg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("----- Open")
	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}

func TestWriteAndParse(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	file, err := os.Open("_examples/resources/test.jpg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	defer file.Close()

	_, err = io.Copy(exif, file)

	if err != nil && err != FoundExifInData {
		t.Fatalf("Error: %s", err.Error())
	}

	err = exif.Parse()

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("----- Write and Parse")
	for key, val := range exif.Tags {
		fmt.Printf("%s: %s\n", key, val)
	}
}

func TestGetLongitude(t *testing.T) {
	exif := New()
	err := exif.Open("_examples/resources/testlocation.jpg")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	longitude, ok := exif.Tags["GPSLongitude"]
	if !ok {
		t.Fatalf("Error: Tag \"Longitude\" could not be found")
	}

	if longitude != "131,  0, 55.2063" {
		t.FailNow()
	}

}

func TestGetLatitude(t *testing.T) {
	exif := New()
	err := exif.Open("_examples/resources/testlocation.jpg")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	latitude, ok := exif.Tags["GPSLatitude"]
	if !ok {
		t.Fatalf("Error: Tag \"Latitude\" could not be found")
	}

	if latitude != "25, 21, 32.6101" {
		t.Fatalf("Error:\n Expected 25, 21, 32.6101\n Found: %s", latitude)
	}
}

func TestMakerNoteAndUserComment(t *testing.T) {
	exif := New()

	// http://www.exif.org/samples/fujifilm-mx1700.jpg
	err := exif.Open("_examples/resources/IMAG0001.jpeg")

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	fmt.Println("----- Maker Note")
	hasMakerNote := false
	for key, val := range exif.Tags {
		if key != "MakerNote" {
			continue
		}
		hasMakerNote = true
		// make sure it's base64 encoded
		// fmt.Printf("%s: %s\n", key, val)
		decodedBytes, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		if len(decodedBytes) != 2804 {
			t.Fatalf("Error: length of decoded data should be 2804: %d", len(decodedBytes))
		}
	}
	if !hasMakerNote {
		t.Fatal("should have maker note")
	}
	fmt.Println("----- User Comment")
	hasUserComment := false
	for key, val := range exif.Tags {
		if key != "UserComment" {
			continue
		}
		hasUserComment = true
		decodedBytes, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		if len(decodedBytes) != 20 { // "ASCII\0\0\0user comment"
			t.Fatalf("Error: length of decoded data should be 20: %d", len(decodedBytes))
		}
	}
	if !hasUserComment {
		t.Fatal("should have user comment")
	}
}
