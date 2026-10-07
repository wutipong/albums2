package media

import (
	"slices"
	"strings"
)

var animationExts = []string{
	".gif",
	".webp",
}

func HasAnimationExt(ext string) bool {
	return slices.Contains(animationExts, strings.ToLower(ext))
}
