package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"regexp"
	"test/storage/proto"
)

var (
	bucketName = "images"
	thumbName = "thumbs"
	location = "ru-east-1"
	endpoint = "minio:9000"
	accessKeyID = "XKTWRN4QJNPEM35M2MXR"
	secretAccessKey = "gb1AU5YzuOhWgaVRE7jqurezp4bK7tRSdU9EG07Z"
)

func NewService() *service {
	s := new(service)

	// Initi minio storage to store images
	// ToDO: refactor after deploy
	//endpoint = os.Getenv("")
	//accessKeyID = os.Getenv("")
	//secretAccessKey = os.Getenv("")

	// Create storage client
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.Fatalln(err)
	}

	// Trying to create bucket for images
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	if !exists {
		err := minioClient.MakeBucket(bucketName, location)
		if err != nil {
			log.Fatalln(err)
		}

		err = minioClient.SetBucketPolicy(bucketName, `{"Version": "2012-10-17","Statement": [{"Action":["s3:GetBucketLocation"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::images"],"Sid":""},{"Action":["s3:ListBucket"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::images"],"Sid":""},{"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::images/*"],"Sid":""}]}`)
		// Set remote access to read bucket files
		if err != nil {
			log.Fatalln(err)
		}
	}

	// trying to create thumb bucket
	ok, err := minioClient.BucketExists(thumbName)
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		err := minioClient.MakeBucket(thumbName, location)
		if err != nil {
			log.Fatalln(err)
		}
		err = minioClient.SetBucketPolicy(thumbName, `{"Version": "2012-10-17","Statement": [{"Action":["s3:GetBucketLocation"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::thumbs"],"Sid":""},{"Action":["s3:ListBucket"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::thumbs"],"Sid":""},{"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::thumbs/*"],"Sid":""}]}`)
		if err != nil {
			log.Fatalln(err)
		}
	}

	s.minioClient = minioClient
	return s
}

type service struct {
	minioClient *minio.Client
}

func (s *service) Save (ctx context.Context, req *proto.File, rsp *proto.Response) error {
	_, err := s.minioClient.PutObject(
		bucketName,
		req.Name,
		bytes.NewReader(req.Content),
		req.Size,
		minio.PutObjectOptions{
			ContentType: req.Type,
		},
	)
	if err != nil {
		return err
	}
	//fmt.Println(n)

	// Resize image
	img, _, err := image.Decode(bytes.NewReader(req.Content))
	if err == nil {
		thumbnail := resize.Thumbnail(100, 100, img, resize.Lanczos3)

		// Get image type
		r, err := regexp.Compile("image/([a-z]+)")
		if err != nil {
			return err
		}
		match := r.FindStringSubmatch(req.Type)
		if len(match) == 2 {
			// Convert image -> bytes
			var buff bytes.Buffer
			switch match[1] {
				case "png":
					png.Encode(&buff, thumbnail)
					break
				case "jpg", "jpeg":
					jpeg.Encode(&buff, thumbnail, nil)
					break
				default:
					fmt.Println("is not an image")

			}
			// Upload resized thumbnail
			if buff.Len() > 0 {
				_, err := s.minioClient.PutObject(
					thumbName,
					req.Name,
					bytes.NewReader(buff.Bytes()),
					int64(buff.Len()),
					minio.PutObjectOptions{
						ContentType: req.Type,
					},
				)
				if err != nil {
					return err
				}
			}
		}
	}
	rsp.Msg = "Success"
	return nil
}
