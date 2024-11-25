package service

import (
	"context"
	"io"
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	"github.com/minio/minio-go/v7"
)

type MinioStore struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStore(conn *db.MinioConn, bucketName string) MinioStore {
	// get new minio client
	client := db.NewMinioClient(conn)

	// check if bucket exists
	ctx := context.Background()
	isExists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Error while setting up minio store: %v", err)
	}

	// if not, make one
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

func (m MinioStore) SaveFile(file *types.FileUpload) error {
	_, err := m.client.PutObject(context.Background(), m.bucketName, file.Name, file.Content, file.Size, minio.PutObjectOptions{
		ContentType: file.Headers["Content-Type"],
	})
	if err != nil {
		return err
	}

	return nil
}

func (m MinioStore) GetFile(filePath string) ([]byte, error) {
	// fetch the object
	object, err := m.client.GetObject(context.Background(), m.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	// read the object in bytes
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}
