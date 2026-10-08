package util

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"uuid"

	"github.com/minio/minio-go/v7"
)

var (
	bucket                    = ""
	minioClient *minio.Client = nil
)

func CreateAssetKey(extension string) string {
	u := uuid.NewV7()
	return fmt.Sprintf("public/%s.%s", u.String(), extension)
}

func PutObject(
	ctx context.Context,
	key string,
	reader io.Reader,
	length int64,
	options minio.PutObjectOptions,
) (info minio.UploadInfo, err error) {
	slog.Info("put object to S3", "key", key)
	return minioClient.PutObject(
		ctx, bucket,
		key,
		reader,
		length,
		options,
	)
}
