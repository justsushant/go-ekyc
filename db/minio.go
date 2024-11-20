package db

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConn struct {
	Endpoint string
	User     string
	Password string
	Ssl      bool
}

func NewMinioClient(conn *MinioConn) *minio.Client {
	client, err := minio.New(conn.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conn.User, conn.Password, ""),
		Secure: conn.Ssl,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Minio client connected")
	return client
}
