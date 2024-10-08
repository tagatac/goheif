// Package heic2jpg provides convenience functions for converting an HEIC image
// to a JPG image.
package heic2jpg

import (
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/tagatac/goheif"
)

//go:generate mockgen -destination=mock_heic2jpg/mock_heic2jpg.go github.com/tagatac/goheif/heic2jpg Converter

type (
	// Converter provides a function for converting an HEIC image to a JPG image.
	Converter interface {
		// HEIC2JPG converts a specified HEIC image to JPG and writes it to dst.
		HEIC2JPG(src, dst string) error
	}
	converter struct{}
)

// NewConverter returns a Converter.
func NewConverter() Converter {
	return converter{}
}

func (converter) HEIC2JPG(src, dst string) error {
	fi, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "open HEIC file %q", src)
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		log.Println(errors.Wrapf(err, "WARN: no EXIF from HEIC file %q", src))
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		return errors.Wrapf(err, "decode HEIC file %q", src)
	}

	fo, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		errors.Wrapf(err, "create output file %q", dst)
	}
	defer fo.Close()

	w, err := newWriterExif(fo, exif)
	if err != nil {
		return errors.Wrap(err, "add EXIF data to image writer")
	}
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return errors.Wrapf(err, "encode JPG file %q", dst)
	}

	return nil
}

// skip writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}
