package mailo

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var dec = &mime.WordDecoder{
	CharsetReader: charsetReader,
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	// for iPhone
	if strings.EqualFold(charset, "CP932") {
		charset = "Shift_JIS"
	}
	enc, err := ianaindex.IANA.Encoding(charset)
	if err != nil {
		return nil, fmt.Errorf("charset not supported: %q: %v", charset, err)
	}

	return transform.NewReader(input, enc.NewDecoder()), nil
}

func DecodeHeader(header string) (string, error) {
	return dec.DecodeHeader(header)
}

func ParseAddress(address string) (*mail.Address, error) {
	return (&mail.AddressParser{WordDecoder: dec}).Parse(address)
}

func ParseAddressList(list string) ([]*mail.Address, error) {
	return (&mail.AddressParser{WordDecoder: dec}).ParseList(list)
}

func ReadBody(msg *mail.Message) (b []byte, err error) {
	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain; charset=us-ascii"
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse content-type")
	}
	if !strings.HasPrefix(mediaType, "text/") {
		return nil, fmt.Errorf("unsupported media type: %q", mediaType)
	}

	r := msg.Body
	encoding := msg.Header.Get("Content-Transfer-Encoding")
	if encoding == "" {
		encoding = "7BIT"
	}
	switch {
	case strings.EqualFold(encoding, "quoted-printable"):
		r = quotedprintable.NewReader(r)
	case strings.EqualFold(encoding, "base64"):
		r = base64.NewDecoder(base64.StdEncoding, r)
		defer func() {
			if err == io.ErrUnexpectedEOF {
				err = nil
			}
		}()
	case strings.EqualFold(encoding, "7bit"), strings.EqualFold(encoding, "8bit"):
	default:
		return nil, fmt.Errorf("unsupported encoding: %q", encoding)
	}

	if cs, ok := params["charset"]; ok {
		cr, err := charsetReader(cs, r)
		if err != nil {
			return nil, err
		}
		r = cr
	}
	return ioutil.ReadAll(r)
}

type Message struct {
	Date    time.Time
	From    *mail.Address
	To      []*mail.Address
	Subject string

	Text io.Reader
	HTML io.Reader

	Resources   map[string]io.Reader
	Attachments map[string]io.Reader
}

func ReadMessage(r io.Reader) (*Message, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}

	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain; charset=us-ascii"
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse content-type")
	}

	switch mediaType {
	case "multipart/mixed":
		return readMixed(msg.Body, params["boundary"])
	case "multipart/related":
		return readRelated(msg.Body, params["boundary"])
	case "multipart/alternative":
		return readAlternative(msg.Body, params["boundary"])
	case "text/html":

	case "text/plain":
	}
}

func readMixed(r io.Reader, boundary string) (*Message, error) {
	mp := multipart.NewReader(r, boundary)
	for {
		part, err := mp.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
}

func readRelated(r io.Reader, boundary string) (*Message, error) {

}

func readAlternative(r io.Reader, boundary string) (*Message, error) {

}

func ReadText(msg *mail.Message) (b []byte, err error) {
	encoding := msg.Header.Get("Content-Transfer-Encoding")
	if encoding == "" {
		encoding = "7BIT"
	}
	r := msg.Body
	switch {
	case strings.EqualFold(encoding, "quoted-printable"):
		r = quotedprintable.NewReader(r)
	case strings.EqualFold(encoding, "base64"):
		r = base64.NewDecoder(base64.StdEncoding, r)
		defer func() {
			if err == io.ErrUnexpectedEOF {
				err = nil
			}
		}()
	case strings.EqualFold(encoding, "7bit"), strings.EqualFold(encoding, "8bit"):
	default:
		return nil, fmt.Errorf("unsupported encoding: %q", encoding)
	}

	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain; charset=us-ascii"
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse content-type")
	}
	if cs, ok := params["charset"]; ok {
		cr, err := charsetReader(cs, r)
		if err != nil {
			return nil, err
		}
		r = cr
	}
	return nil, nil
}
