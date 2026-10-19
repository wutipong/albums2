package service

import (
	"context"
	"fmt"
	"os"

	"github.com/wutipong/albums2/worker/service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetImageSpec(ctx context.Context) (
	resp *pb.GetImageSpecResponse, err error,
) {
	conn, err := grpc.NewClient(
		os.Getenv("MLSERVICE_ADDRESS"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		err = fmt.Errorf("failed to create grpc client: %w", err)
		return
	}

	defer conn.Close()

	client := pb.NewEncodingServiceClient(conn)
	resp, err = client.GetImageSpec(ctx, &pb.GetImageSpecRequest{})

	if err != nil {
		err = fmt.Errorf("failed to get image spec: %w", err)
	}

	return
}

func EncodeText(ctx context.Context, input string) (
	resp *pb.EncodeResponse, err error,
) {
	conn, err := grpc.NewClient(
		os.Getenv("MLSERVICE_ADDRESS"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		err = fmt.Errorf("failed to create grpc client: %w", err)
		return
	}

	defer conn.Close()

	client := pb.NewEncodingServiceClient(conn)
	resp, err = client.EncodeText(ctx, &pb.EncodeTextRequest{
		Input: input,
	})

	if err != nil {
		err = fmt.Errorf("failed to encode text: %w", err)
	}

	return
}

func EncodeImage(ctx context.Context, input []byte) (
	resp *pb.EncodeResponse, err error,
) {
	conn, err := grpc.NewClient(
		os.Getenv("MLSERVICE_ADDRESS"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		err = fmt.Errorf("failed to create grpc client: %w", err)
		return
	}

	defer conn.Close()

	client := pb.NewEncodingServiceClient(conn)
	resp, err = client.EncodeImage(ctx, &pb.EncodeImageRequest{
		Image: input,
	})

	if err != nil {
		err = fmt.Errorf("failed to encode image: %w", err)
	}

	return
}
