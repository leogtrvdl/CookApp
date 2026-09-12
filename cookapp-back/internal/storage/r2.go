package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type R2Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewR2Storage() (*R2Storage, error) {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("R2_BUCKET")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       "auto",
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	})

	return &R2Storage{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

func (r *R2Storage) UploadFile(file multipart.File, filename string, folder string) (string, error) {

	id := uuid.New()
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("%s/%s%s", folder, id.String(), ext)

	contentType, err := detectContentType(file)
	if err != nil {
		return "", fmt.Errorf("erreur détection type fichier: %w", err)
	}
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", fmt.Errorf("mauvais type de fichier")
	}
	_, err = r.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("erreur upload R2: %w", err)
	}

	return fmt.Sprintf("%s/%s", r.publicURL, key), nil
}

func (r *R2Storage) DeleteFile(fileURL string) error {
	key := strings.TrimPrefix(fileURL, r.publicURL + "/")

	_, err := r.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}

func detectContentType(file multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	contentType := http.DetectContentType(buf[:n])

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return contentType, nil
}
