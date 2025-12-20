package miniocli

import (
	"context"
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
	cli        *minio.Client
	pcli       *minio.Client
	logger     zerolog.Logger
	bucketName string
}

func NewMinIOCli(baseLogger zerolog.Logger) (MinIOClient, error) {
	cfg, err := config.ParseMinIOCliConfig()
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "minio").Logger()

	cli, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
		Region: "us-east-1",
	})

	if err != nil {
		return nil, err
	}

	pcli, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
		Region: "us-east-1",
	})

	if err != nil {
		return nil, err
	}

	return &minioCli{
		cli:        cli,
		pcli:       pcli,
		logger:     logger,
		bucketName: "report",
	}, nil
}

func (c *minioCli) UploadAndPresign(ctx context.Context, in dto.PutReportIn) (string, error) {
	c.logger.Debug().Str("evt", "call LoadFile").Msg("")

	fctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.cli.PutObject(fctx, in.Bucket, in.ObjectName, in.Reader, -1, minio.PutObjectOptions{
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
