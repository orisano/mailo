package mailo

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"strings"

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
		return nil, fmt.Errorf("unsupported charset: %q: %v", charset, err)
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
