package storages

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

// PublicObject minio storage
type PublicObject struct {
	BaseObject
	Location string
	URL      string
}

func (PublicObject) bucketName() string {
	return "public"
}

func (p *PublicObject) generatePublicURL() string {
	u, err := url.Parse(p.FileName)
	if err != nil {
		log.Panicf("error in PublicObject.generatePublicURL: %v \n", err)
	}

	return variables.MinioConfig.URIEndpoint + "/" + path.Join(
		p.bucketName(),
		strings.Join([]string{"desatanjungbunga", variables.Mode}, "-"),
		p.PathName,
		u.String(),
	)
}

func (p *PublicObject) generateObjectName() string {
	return path.Join(
		strings.Join([]string{"desatanjungbunga", variables.Mode}, "-"),
		p.PathName,
		p.FileName,
	)
}

// Upload file to public bucket
// rewrite public URL to this variable
func (p *PublicObject) Upload(ctx context.Context) (*minio.UploadInfo, error) {
	client := services.MinioClient
	info, err := client.PutObject(
		ctx,
		p.bucketName(),
		p.generateObjectName(),
		bytes.NewReader(p.File),
		p.Size,
		minio.PutObjectOptions{ContentType: http.DetectContentType(p.File)},
	)
	if err != nil {
		return nil, err
	}

	p.URL = p.generatePublicURL()
	return &info, nil
}
