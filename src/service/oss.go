package service

import (
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/yiwen-ai/yiwen-static/src/conf"
)

type OSS struct {
	bucket *oss.Bucket
}

func NewOSS() *OSS {
	cfg := conf.Config.OSS
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		panic(err)
	}

	return &OSS{
		bucket: bucket,
	}
}

func (s *OSS) GetFile(ctx context.Context, objectKey string) ([]byte, error) {
	r, err := s.bucket.GetObject(objectKey)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return io.ReadAll(r)
}
