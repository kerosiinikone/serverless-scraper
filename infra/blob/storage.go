package blob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

func New(s *session.Session) *s3.S3 {
	return s3.New(s)
}

func NewUploader(s *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s)
}

func NewDownloader(s *session.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(s)
}

func RetrieveFiles(svc *s3.S3, downloader *s3manager.Downloader, cid string, rid string) ([]models.DataEntry, error) {
    var (
        d []models.DataEntry
        dataEntry models.DataEntry
    )

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
        Bucket: aws.String(bucketName()),
        Prefix: aws.String(pathToFile(cid, rid, "")),
    })
    if err != nil {
        return nil, err
    }
    for _, item := range resp.Contents {
        b := &aws.WriteAtBuffer{}
        _, err := downloader.Download(b, &s3.GetObjectInput{
            Bucket: aws.String(bucketName()),
            Key:    item.Key,
        })
        if err != nil {
            return nil, err
        }
        err = json.Unmarshal(b.Bytes(), &dataEntry)
        if err != nil {
            return nil, err
        }
        d = append(d, dataEntry)
    }

    return d, nil
}

func SaveFile(u *s3manager.Uploader, d models.DataEntry) error {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return err
	}
	_, err = u.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucketName()),
        Key:    aws.String(pathToFile(d.ClientID, d.RequestID, d.Post.Id)),
        Body:   bytes.NewReader(jsonData),
    })
    if err != nil {
        return err
    }
	return nil
}

func pathToFile(cId, rId, pId string) string {
    return fmt.Sprintf("client-%s/request-%s/post-%s", cId, rId, pId)
}

func bucketName() string {
    return os.Getenv("AWS_BUCKET")
}