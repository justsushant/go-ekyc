package controller

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type MinioStore struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStore(client *minio.Client, bucketName string) MinioStore {
	ctx := context.Background()
	isExists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Error while setting up minio store: %v", err)
	}

	if !isExists {
		client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		log.Println("Bucket created successfully")
	} else {
		log.Println("Bucket already exists")
	}

	return MinioStore{
		client:     client,
		bucketName: bucketName,
	}
}

func (m MinioStore) SaveFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(context.Background(), m.bucketName, fileHeader.Filename, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return err
	}

	return nil
}
