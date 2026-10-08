package util

type UtilOptions struct {
	S3Bucket      string
	MongoDatabase string
}

func Init(options UtilOptions) {
	bucket = options.S3Bucket
}
