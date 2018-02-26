package jpmail

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"net/textproto"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var addrParser = &mail.AddressParser{
	WordDecoder: &mime.WordDecoder{
		CharsetReader: func(charset string, input io.Reader) (io.Reader, error) {
			switch {
			case strings.EqualFold(charset, "iso-2022-jp"):
				return transform.NewReader(input, japanese.ISO2022JP.NewDecoder()), nil
			case strings.EqualFold(charset, "Windows-31J"), strings.EqualFold(charset, "Shift_JIS"):
				return transform.NewReader(input, japanese.ShiftJIS.NewDecoder()), nil
			default:
				return nil, fmt.Errorf("charset not supported: %q", charset)
			}
		},
	},
}

func ReadMessage(r io.Reader) (*mail.Message, error) {
	tp := textproto.NewReader(bufio.NewReader(r))

	hdr, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	return &mail.Message{
		Header: mail.Header(hdr),
		Body:   tp.R,
	}, nil
}
