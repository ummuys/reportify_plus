package miniocli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type minioCli struct {
	cli    *minio.Client
	pcli   *minio.Client
	logger zerolog.Logger
	bucket string
}

func NewMinIOCli(baseLogger zerolog.Logger) (MinIOClient, error) {
	cfg, err := config.ParseMinIOCliConfig()
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "minio").Logger()

	cli, err := minio.New(cfg.DockerEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
		Region: "us-east-1",
	})

	if err != nil {
		return nil, err
	}

	pcli, err := minio.New(cfg.PublicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
		Region: "us-east-1",
	})

	if err != nil {
		return nil, err
	}

	return &minioCli{
		cli:    cli,
		pcli:   pcli,
		logger: logger,
	}, nil
}

func (c *minioCli) UploadAndPresign(ctx context.Context, in dto.PutReportIn) (string, error) {
	c.logger.Debug().Str("evt", "call LoadFile").Msg("")

	fctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, err := c.cli.BucketExists(fctx, in.Bucket)
	if err != nil {
		return "", err
	}

	if !exists {
		err = c.cli.MakeBucket(fctx, in.Bucket, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			return "", err
		}
	}

	_, err = c.cli.PutObject(fctx, in.Bucket, in.ObjectName, in.Reader, -1, minio.PutObjectOptions{
		ContentType: in.ContentType,
	})

	if err != nil {
		return "", err
	}

	q := make(url.Values)
	q.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%s"`, in.FileName))

	u, err := c.pcli.PresignedGetObject(fctx, in.Bucket, in.ObjectName, in.Expire, q)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (c *minioCli) DeleteExpiredFiles(ctx context.Context, in dto.DeleteExpiredFilesParams) error {

	objectsCh := make(chan minio.ObjectInfo, 100)
	go func() {
		defer close(objectsCh)
		for _, name := range in.Names {
			objectsCh <- minio.ObjectInfo{Key: name}
		}
	}()

	var errs []error
	for err := range c.cli.RemoveObjects(ctx, c.bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		errs = append(errs, err.Err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
