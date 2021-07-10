package mimemail

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"os"
	"strings"
)

type Parser interface {
	Parts() []Part
	Subject() string
	Boundary() string
	MediaType() string
}

type Part struct {
	ContentType string
	Body        string
}

type parserImpl struct {
	boundary  string
	mediaType string
	subject   string
	parts     []Part
}

func New(fileName string) Parser {
	fileReader, err := os.Open(fileName) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	info := &parserImpl{}

	err = info.ParseMessage(fileReader)

	if err != nil {
		log.Fatal(err)
	}

	return info
}

func (impl *parserImpl) Parts() []Part {
	return impl.parts
}

func (impl *parserImpl) Boundary() string {
	return impl.boundary
}

func (impl *parserImpl) MediaType() string {
	return impl.mediaType
}

func (impl *parserImpl) Subject() string {
	return impl.subject
}

func (impl *parserImpl) ParseMessage(r io.Reader) error {
	msg, err := mail.ReadMessage(bufioReader(r))
	if err != nil {
		return err
	}

	for _, values := range msg.Header {
		for idx, val := range values {
			values[idx] = decodeRFC2047(val)
		}
	}
	return impl.parseMessageWithHeader(msg)
}

func (impl *parserImpl) parseMessageWithHeader(msg *mail.Message) error {
	bufferedReader := contentReader(msg.Header, msg.Body)

	var err error
	var mediaType string
	var mediaTypeParams map[string]string

	impl.subject = msg.Header.Get("Subject")

	if contentType := msg.Header.Get("Content-Type"); len(contentType) > 0 {
		mediaType, mediaTypeParams, err = mime.ParseMediaType(contentType)

		impl.mediaType = mediaType

		//fmt.Println("Parse Message With Header: ", mediaType)
		if err != nil {
			return err
		}
	} // Lack of contentType is not a problem

	// Can only have one of the following: Parts, SubMessage, or Body
	if strings.HasPrefix(mediaType, "multipart") {
		impl.boundary = mediaTypeParams["boundary"]

		parts, err := impl.readParts(bufferedReader, impl.boundary)

		if err == nil {
			if parts != nil {
				return nil
			}
		}

	} else if strings.HasPrefix(mediaType, "message") {
		return impl.ParseMessage(bufferedReader)
	}

	return nil
}

func (impl *parserImpl) readParts(bodyReader io.Reader, boundary string) (io.Reader, error) {

	multipartReader := multipart.NewReader(bodyReader, boundary)

	for part, partErr := multipartReader.NextRawPart(); partErr != io.EOF; part, partErr = multipartReader.NextRawPart() {
		if partErr != nil && partErr != io.EOF || part == nil {
			return nil, partErr
		}

		b, err := ioutil.ReadAll(mimeContentReader(part.Header, part))
		if err != nil {
			log.Println(err)
		} else {
			impl.parts = append(impl.parts, Part{
				ContentType: part.Header.Get("Content-Type"),
				Body:        string(b),
			})
		}
	}

	return nil, errors.New("not found")
}

func contentReader(headers mail.Header, bodyReader io.Reader) *bufio.Reader {
	if headers.Get("Content-Transfer-Encoding") == "quoted-printable" {
		return bufioReader(quotedprintable.NewReader(bodyReader))
	}
	if headers.Get("Content-Transfer-Encoding") == "base64" {
		return bufioReader(base64.NewDecoder(base64.StdEncoding, bodyReader))
	}
	return bufioReader(bodyReader)
}

func mimeContentReader(headers textproto.MIMEHeader, bodyReader io.Reader) *bufio.Reader {
	if headers.Get("Content-Transfer-Encoding") == "quoted-printable" {
		// headers.Del("Content-Transfer-Encoding")
		return bufioReader(quotedprintable.NewReader(bodyReader))
	}
	if headers.Get("Content-Transfer-Encoding") == "base64" {
		// headers.Del("Content-Transfer-Encoding")
		return bufioReader(base64.NewDecoder(base64.StdEncoding, bodyReader))
	}
	return bufioReader(bodyReader)
}

// bufioReader ...
func bufioReader(r io.Reader) *bufio.Reader {
	if bufferedReader, ok := r.(*bufio.Reader); ok {
		return bufferedReader
	}
	return bufio.NewReader(r)
}

// decodeRFC2047 ...
func decodeRFC2047(s string) string {
	// GO 1.5 does not decode headers, but this may change in future releases...
	decoded, err := (&mime.WordDecoder{}).DecodeHeader(s)
	if err != nil || len(decoded) == 0 {
		return s
	}
	return decoded
}
