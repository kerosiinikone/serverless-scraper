package infra

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	defaultRegion = "eu-north-1"
)

func New(sessOpt *session.Options) *session.Session {
    if sessOpt == nil {
        sessOpt = &session.Options{
            Config: aws.Config{
                Region: aws.String(getRegion()),
            },
            SharedConfigState: session.SharedConfigEnable,
        }
    }
    return session.Must(session.NewSessionWithOptions(*sessOpt)) 
}

func getRegion() string {
    region := os.Getenv("AWS_REGION")
    if region == "" {
        region = defaultRegion
    }
    return region
}