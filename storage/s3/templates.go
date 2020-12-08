package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s storage) GetHTMLTemplate(userID int64, templateName string) (*s3.GetObjectOutput, error) {
	return nil, nil
}
