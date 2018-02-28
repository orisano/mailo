package jpmail_test

import (
	"net/mail"
	"strings"
	"testing"

	"github.com/orisano/jpmail"
)

func mustMessage(msg *mail.Message, err error) *mail.Message {
	if err != nil {
		panic(err)
	}
	return msg
}

func message(s string) *mail.Message {
	return mustMessage(mail.ReadMessage(strings.NewReader(s)))
}

func TestReadBody(t *testing.T) {
	tests := []struct {
		msg      *mail.Message
		expected string
	}{
		{
			msg: message(`To: Another Gopher <to@example.com>
Subject: Gophers at Gophercon
Date: Mon, 23 Jun 2015 11:40:36 -0400
From: Gopher <from@example.com>
MIME-Version: 1.0
Content-Type: text/plain; charset="UTF-8"

Message body
`),
			expected: "Message body\n",
		},
		{
			msg: message(`To: Another Gopher <to@example.com>
Subject: =?ISO-2022-JP?B?GyRCIVolRiU5JUg0RDYtIVslNSUkJUg5OT83JCw0ME47JDckXiQ3JD8bKEI=?=
Date: Tue, 15 Sep 2015 16:17:23 +0000
From: Gopher <from@example.com>
MIME-Version: 1.0
Content-Type: text/plain; charset=ISO-2022-JP

$B%5%$%H$r99?7$7$?>uBV$KJ]$D$3$H$O%;%-%e%j%F%#$K$H$C$F=EMW$G$9!#$=$l$O$^$?!"$"$J$?$H$"$J$?$NFI<T$K$H$C$F%$%s%?!<%M%C%H$r$h$j0BA4$J>l=j$K$9$k$3$H$G$b$"$j$^$9!#(B
`),
			expected: "サイトを更新した状態に保つことはセキュリティにとって重要です。それはまた、あなたとあなたの読者にとってインターネットをより安全な場所にすることでもあります。\n",
		},
		{
			msg: message(`To: Another Gopher <to@example.com>
Subject: =?ISO-2022-JP?B?GyRCIVolRiU5JUg0RDYtIVslNSUkJUg5OT83JCw0ME47JDckXiQ3JD8bKEI=?=
Date: Tue, 15 Sep 2015 16:17:23 +0000
From: Gopher <from@example.com>
MIME-Version: 1.0
Content-Type: text/plain; charset=ISO-2022-JP
Content-Transfer-Encoding: quoted-printable

=1B=24B=255=25=24=25H=24r99=3F7=247=24=3F=3EuBV=24KJ=5D=24D=243=24H=24O=25=3B=25=2D=25e=25j=25F=25=23=24K=24H=24C=24F=3DEMW=24G=249=21=23=24=3D=24l=24O=24=5E=24=3F=21=22=24=22=24J=24=3F=24H=24=22=24J=24=3F=24NFI=3CT=24K=24H=24C=24F=25=24=25s=25=3F=21=3C=25M=25C=25H=24r=24h=24j0BA4=24J=3El=3Dj=24K=249=24k=243=24H=24G=24b=24=22=24j=24=5E=249=21=23=1B=28B
`),
			expected: "サイトを更新した状態に保つことはセキュリティにとって重要です。それはまた、あなたとあなたの読者にとってインターネットをより安全な場所にすることでもあります。\n",
		},
		{
			msg: message(`MIME-Version: 1.0
Content-Type: text/plain; charset="utf-8"
Content-Transfer-Encoding: base64
To: Another Gopher <to@example.com>
Subject: Gophers at Gophercon
Date: Mon, 23 Jun 2015 11:40:36 -0400
From: Gopher <from@example.com>

44K144Kk44OI44KS5pu05paw44GX44Gf54q25oWL44Gr5L+d44Gk44GT44Go44Gv44K744Kt
44Ol44Oq44OG44Kj44Gr44Go44Gj44Gm6YeN6KaB44Gn44GZ44CC44Gd44KM44Gv44G+44Gf
44CB44GC44Gq44Gf44Go44GC44Gq44Gf44Gu6Kqt6ICF44Gr44Go44Gj44Gm44Kk44Oz44K/
44O844ON44OD44OI44KS44KI44KK5a6J5YWo44Gq5aC05omA44Gr44GZ44KL44GT44Go44Gn
44KC44GC44KK44G+44GZ44CCCg==%
`),
			expected: "サイトを更新した状態に保つことはセキュリティにとって重要です。それはまた、あなたとあなたの読者にとってインターネットをより安全な場所にすることでもあります。\n",
		},
		{
			msg: message(`To: Another Gopher <to@example.com>
Subject: Gophers at Gophercon
From: Gopher <from@example.com>
Date: Fri, 18 Sep 2015 17:51:01 +0900
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: quoted-printable
Content-Disposition: inline
MIME-Version: 1.0

=E3=82=B5=E3=82=A4=E3=83=88=E3=82=92=E6=9B=B4=E6=96=B0=E3=81=97=E3=81=9F=
=E7=8A=B6=E6=85=8B=E3=81=AB=E4=BF=9D=E3=81=A4=E3=81=93=E3=81=A8=E3=81=AF=
=E3=82=BB=E3=82=AD=E3=83=A5=E3=83=AA=E3=83=86=E3=82=A3=E3=81=AB=E3=81=A8=
=E3=81=A3=E3=81=A6=E9=87=8D=E8=A6=81=E3=81=A7=E3=81=99=E3=80=82=E3=81=9D=
=E3=82=8C=E3=81=AF=E3=81=BE=E3=81=9F=E3=80=81=E3=81=82=E3=81=AA=E3=81=9F=
=E3=81=A8=E3=81=82=E3=81=AA=E3=81=9F=E3=81=AE=E8=AA=AD=E8=80=85=E3=81=AB=
=E3=81=A8=E3=81=A3=E3=81=A6=E3=82=A4=E3=83=B3=E3=82=BF=E3=83=BC=E3=83=8D=
=E3=83=83=E3=83=88=E3=82=92=E3=82=88=E3=82=8A=E5=AE=89=E5=85=A8=E3=81=AA=
=E5=A0=B4=E6=89=80=E3=81=AB=E3=81=99=E3=82=8B=E3=81=93=E3=81=A8=E3=81=A7=
=E3=82=82=E3=81=82=E3=82=8A=E3=81=BE=E3=81=99=E3=80=82
`),
			expected: "サイトを更新した状態に保つことはセキュリティにとって重要です。それはまた、あなたとあなたの読者にとってインターネットをより安全な場所にすることでもあります。\n",
		},
	}

	for _, test := range tests {
		b, err := jpmail.ReadBody(test.msg)
		if err != nil {
			t.Error(err)
			continue
		}
		if got := string(b); got != test.expected {
			t.Errorf("unexpected message body. expected: %q, but got: %q", test.expected, got)
		}
	}
}
