package jpmail

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var addrParser = &mail.AddressParser{
	WordDecoder: &mime.WordDecoder{
		CharsetReader: charsetReader,
	},
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	if strings.EqualFold(charset, "CP932") {
		charset = "Shift_JIS"
	}
	enc, err := ianaindex.IANA.Encoding(charset)
	if err != nil {
		return nil, fmt.Errorf("charset not supported: %q: %v", charset, err)
	}

	return transform.NewReader(input, enc.NewDecoder()), nil
}

func ParseAddress(address string) (*mail.Address, error) {
	return addrParser.Parse(address)
}

func ParseAddressList(list string) ([]*mail.Address, error) {
	return addrParser.ParseList(list)
}

func ReadBody(msg *mail.Message) ([]byte, error) {
	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	if mediaType != "text/plain" {
		return nil, fmt.Errorf("unsupported media type: %q", mediaType)
	}

	r := msg.Body
	encoding := msg.Header.Get("Content-Transfer-Encoding")
	if encoding != "" {
		switch {
		case strings.EqualFold(encoding, "quoted-printable"):
			r = quotedprintable.NewReader(r)
		case strings.EqualFold(encoding, "base64"):
			r = base64.NewDecoder(base64.StdEncoding, r)
		case strings.EqualFold(encoding, "7bit"), strings.EqualFold(encoding, "8bit"):
		default:
			return nil, fmt.Errorf("unsupported encoding")
		}
	}

	if cs, ok := params["charset"]; ok {
		rr, err := charsetReader(cs, r)
		if err != nil {
			return nil, err
		}
		r = rr
	}

	return ioutil.ReadAll(r)
}
