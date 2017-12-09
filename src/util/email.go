package util

import (
	"net/mail"
	"os"

	"email"
)

func ServerAlert() (err error) {
	subject := "server is down"
	body := `This email is system generated, please do not reply to it.`

	add, err := mail.ParseAddress("ADMIN@heshi.com")
	if err != nil {
		return
	}
	m := email.NewMessage(subject, body)
	m.From = *add
	m.To = []string{"cunying.huang@hpe.com", "huang372923@hotmail.com"}

	mlog := "server.log"
	if PathExists(mlog) {
		size, err := GetFileSize(mlog)
		if err != nil {
			return err
		}
		if size > 2<<10 {
			buf := make([]byte, 2<<10)
			f, err := os.Open(mlog)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.ReadAt(buf, size-2<<10)
			if err != nil {
				return err
			}

			m.Attachments[mlog] = &email.Attachment{
				Filename: mlog,
				Data:     buf,
				Inline:   true,
			}
		} else {
			if err = m.Inline(mlog); err != nil {
				return err
			}
		}
	}

	return email.Send("smtp3.hp.com:25", nil, m)
}

func CurrencyExchangeAlert() {
	subject := "exchange fluctuation is too high"
	body := `This email is system generated, please do not reply to it.
		exchange fluctuation is too high.`

	add, err := mail.ParseAddress("ADMIN@heshi.com")
	if err != nil {
		return
	}
	mCurrency := email.NewMessage(subject, body)
	mCurrency.From = *add
	mCurrency.To = []string{"cunying.huang@hpe.com", "huang372923@hotmail.com"}
	email.Send("smtp3.hp.com:25", nil, mCurrency)
}

func FailToGetCurrencyExchangeAlert() {
	subject := "FAIL TO GET CURRENCY EXCHANGE RATE FROM UI"
	body := `PLEASE MANUAL UPDATE TODAY CURRENCY EXCHANGE RATE AND CHECK.`

	add, err := mail.ParseAddress("ADMIN@heshi.com")
	if err != nil {
		return
	}
	mCurrency := email.NewMessage(subject, body)
	mCurrency.From = *add
	mCurrency.To = []string{"cunying.huang@hpe.com", "huang372923@hotmail.com"}
	email.Send("smtp3.hp.com:25", nil, mCurrency)
}
