package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"uuid"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	bucket               = ""
	client *minio.Client = nil
)

func Init(endpoint string, accessKeyId string, secret string) error {
	endpoint, secure, err := GetMinioEndpoint(endpoint)
	if err != nil {
		slog.Error("unable to parse endpoint", "error", err)
		return fmt.Errorf("unable to parse endpoint")
	}

	client, err = minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyId, secret, ""),
		Secure:       secure,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return fmt.Errorf("unable to create minio client: %w", err)
	}

	return nil
}

func GetMinioEndpoint(input string) (endpoint string, secure bool, err error) {
	endpointUrl, err := url.Parse(input)
	if err != nil {
		err = fmt.Errorf("unable to parse input url")
		return
	}

	scheme := endpointUrl.Scheme

	if scheme == "https" {
		secure = true
	} else {
		secure = false
	}
	endpoint = input
	endpoint = strings.TrimPrefix(endpoint, scheme)
	endpoint = strings.TrimPrefix(endpoint, "://")

	return
}

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
	return client.PutObject(
		ctx, bucket,
		key,
		reader,
		length,
		options,
	)
}

func GetObject(ctx context.Context,
	key string,
	options minio.GetObjectOptions,
) (object *minio.Object, err error) {
	slog.Info("get object from S3", "key", key)

	return client.GetObject(
		ctx,
		bucket,
		key,
		options,
	)
}
