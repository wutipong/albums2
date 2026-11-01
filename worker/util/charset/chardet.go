package charset

import (
	"fmt"
	"log/slog"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

var charsetMap = map[string]string{
	"GB-18030": "GB18030",
}

func DetectCharSet(str string) (encoding encoding.Encoding, charset string, err error) {
	detector := chardet.NewTextDetector()

	chardetResult, err := detector.DetectBest([]byte(str))
	if err != nil {
		err = fmt.Errorf("failed to detect charset: %w", err)
		slog.Warn(
			"failed to detect encoding.",
			slog.String("text", str),
			slog.String("error", err.Error()),
		)

		return
	}
	charset, ok := charsetMap[chardetResult.Charset]
	if !ok {
		charset = chardetResult.Charset
	}

	encoding, err = ianaindex.IANA.Encoding(charset)
	return
}

func DecodeString(str string, enc encoding.Encoding) (string, error) {
	decoder := enc.NewDecoder()
	decodedStr, err := decoder.String(str)
	if err != nil {
		err = fmt.Errorf("failed to decode string: %w", err)
		slog.Warn(
			"failed to decode string.",
			slog.String("text", str),
			slog.String("error", err.Error()),
		)
		return str, err
	}

	return decodedStr, nil
}

func Decode(str string) string {
	enc, _, err := DetectCharSet(str)
	if err != nil {
		slog.Warn("unable to detect encoding", slog.String("error", err.Error()))
		return str
	}

	decodedStr, err := DecodeString(str, enc)
	if err != nil {
		slog.Warn("unable to decode the string", slog.String("error", err.Error()))
		return str
	}

	return decodedStr
}
