package charset

import (
	"fmt"

	zlog "github.com/rs/zerolog/log"
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
		zlog.Warn().Str("text", str).Str("error", err.Error()).Msg("failed to detect encoding.")

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
		zlog.Warn().Str("text", str).Str("error", err.Error()).Msg("failed to decode string.")
		return str, err
	}

	return decodedStr, nil
}

func Decode(str string) string {
	enc, _, err := DetectCharSet(str)
	if err != nil {
		zlog.Warn().Str("error", err.Error()).Msg("unable to detect encoding")
		return str
	}

	decodedStr, err := DecodeString(str, enc)
	if err != nil {
		zlog.Warn().Str("error", err.Error()).Msg("unable to decode the string")
		return str
	}

	return decodedStr
}
