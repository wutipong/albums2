package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"uuid"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	zlog "github.com/rs/zerolog/log"
)

var (
	bucket               = ""
	client *minio.Client = nil
)

func Init(
	endpoint string,
	accessKeyId string,
	secret string,
	bucketName string,
	region string,
) error {
	zlog.Info().
		Str("endpoint", endpoint).
		Str("accessKeyId", accessKeyId).
		Str("secret", secret[:2]+strings.Repeat("*", len(secret)-2)).
		Str("bucketName", bucketName).
		Str("region", region).
		Msg("using S3 credential")

	endpoint, secure, err := GetMinioEndpoint(endpoint)
	if err != nil {
		zlog.Error().Err(err).Msg("unable to parse endpoint")
		return fmt.Errorf("unable to parse endpoint")
	}

	client, err = minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyId, secret, ""),
		Secure:       secure,
		BucketLookup: minio.BucketLookupPath,
		Region:       region,
	})
	if err != nil {
		return fmt.Errorf("unable to create minio client: %w", err)
	}

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list buckets: %w", err)
	}

	for _, b := range buckets {
		zlog.Debug().Str("name", b.Name).Time("created", b.CreationDate).Msg("found bucket")
	}

	bucket = bucketName

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

func CreateAssetKey(collectionId, albumId, extension string) string {
	u := uuid.NewV4()

	keyStr := strings.ReplaceAll(u.String(), "-", "")

	return fmt.Sprintf("_cache/%s/%s/%s%s",
		collectionId,
		albumId,
		keyStr,
		extension)
}

func CreateAssetKeyPrefix(collectionId, albumId string) string {
	return fmt.Sprintf("_cache/%s/%s",
		collectionId,
		albumId,
	)
}

func PutObject(
	ctx context.Context,
	key string,
	reader io.Reader,
	length int64,
	options minio.PutObjectOptions,
) (info minio.UploadInfo, err error) {
	zlog.Info().Str("key", key).Msg("put object to S3")
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
	zlog.Info().Str("key", key).Msg("get object from S3")

	return client.GetObject(
		ctx,
		bucket,
		key,
		options,
	)
}

func RemoveObject(ctx context.Context, key string, options minio.RemoveObjectOptions) error {
	zlog.Info().Str("key", key).Msg("remove object from S3")
	return client.RemoveObject(ctx,
		bucket,
		key,
		options,
	)
}

func Client() *minio.Client {
	return client
}

func Bucket() string {
	return bucket
}
