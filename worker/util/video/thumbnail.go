package video

import (
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ExtractScreenshotToFile(
	inPath string,
	outPath string,
	pos float64,
	quality int,
	height int,
) error {
	err := ExtractToFile(inPath, outPath, pos, ffmpeg.KwArgs{
		"c:v":     "libwebp",
		"vframes": "1",
		"quality": fmt.Sprintf("%d", quality),
		"vf":      fmt.Sprintf("scale=-2:%d", height),
	})
	if err != nil {
		return fmt.Errorf("unable to extract screenshot: %w", err)
	}
	return nil
}

func ExtractAnimateScreenshotToFile(inPath string,
	outPath string,
	pos float64,
	seconds int64,
	loop int,
	quality int,
	height int,
) error {
	err := ExtractToFile(inPath, outPath, pos, ffmpeg.KwArgs{
		"c:v":     "libwebp",
		"t":       fmt.Sprintf("%d", seconds),
		"loop":    fmt.Sprintf("%d", loop),
		"quality": fmt.Sprintf("%d", quality),
		"vf":      fmt.Sprintf("fps=5,scale=-2:%d", height),
	})
	if err != nil {
		return fmt.Errorf("unable to extract screenshot: %w", err)
	}

	return nil
}

func ExtractToFile(inPath string,
	outPath string,
	pos float64,
	args ffmpeg.KwArgs,
) error {
	err := ffmpeg.
		Input(inPath, ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", pos),
		}).
		Output(outPath, args).
		OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to extract screenshot: %w", err)
	}
	return nil
}
