package mime_composer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

type MimeComposer interface {
	From(email string, displayName string) MimeComposer
	To(email string, displayName string) MimeComposer
	Cc(email string, displayName string) MimeComposer
	Bcc(email string) MimeComposer
	Subject(subject string) MimeComposer
	Compose() string
	AddMultipartAlternative(contentType string, body string) MimeComposer
}

type mimeComposerImpl struct {
	from         person
	to           []person
	cc           []person
	bcc          []person
	subject      string
	multiPartAlt []multiPart
}
type person struct {
	email       string
	displayName string
}
type multiPart struct {
	contentType string
	body        string
}

func New() MimeComposer {
	return &mimeComposerImpl{}
}

func (impl *mimeComposerImpl) From(email string, displayName string) MimeComposer {
	impl.from = person{
		email:       email,
		displayName: displayName,
	}
	return impl
}

func (impl *mimeComposerImpl) To(email string, displayName string) MimeComposer {
	impl.to = append(impl.to, person{
		email:       email,
		displayName: displayName,
	})
	return impl
}

func (impl *mimeComposerImpl) Cc(email string, displayName string) MimeComposer {
	impl.cc = append(impl.cc, person{
		email:       email,
		displayName: displayName,
	})
	return impl
}

func (impl *mimeComposerImpl) Bcc(email string) MimeComposer {
	impl.bcc = append(impl.bcc, person{
		email: email,
	})
	return impl
}

func (impl *mimeComposerImpl) Subject(subject string) MimeComposer {
	impl.subject = subject
	return impl
}

func (impl *mimeComposerImpl) AddMultipartAlternative(contentType string, body string) MimeComposer {
	impl.multiPartAlt = append(impl.multiPartAlt, multiPart{
		contentType: contentType,
		body:        body,
	})
	return impl
}

func QuotedPrintable(input string) string {
	byteBuffer := new(bytes.Buffer)
	writer := quotedprintable.NewWriter(byteBuffer)
	writer.Write([]byte(input))
	writer.Close()
	return byteBuffer.String()
}

func q(input string) string {
	return mime.QEncoding.Encode("utf-8", input)
}

func b(input string) string {
	return mime.BEncoding.Encode("utf-8", input)
}
func b64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func personRemap(in []person, fn func(person) string) []string {
	out := make([]string, len(in))

	for i := 0; i < len(in); i++ {
		out[i] = fn(in[i])
	}

	return out
}

func (impl *mimeComposerImpl) Compose() string {
	body := &bytes.Buffer{}

	// Write mail header
	fmt.Fprintf(body, "From: %s <%s>\r\n", q(impl.from.displayName), impl.from.email)
	fmt.Fprintf(body, "To: %s\r\n", strings.Join(personRemap(impl.to, func(p person) string {
		return fmt.Sprintf("%s <%s>", q(p.displayName), p.email)
	}), ","))
	fmt.Fprintf(body, "Subject: %s\r\n", strings.Join(strings.Fields(b(impl.subject)), "\r\n\t"))

	mw := multipart.NewWriter(body)

	fmt.Fprintf(body, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(body, "Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", mw.Boundary())

	for _, part := range impl.multiPartAlt {
		w, err := mw.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {part.contentType},
			"Content-Transfer-Encoding": {"quoted-printable"},
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, QuotedPrintable(part.body))
	}

	mw.Close()

	return body.String()
}
