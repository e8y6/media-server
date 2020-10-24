package s3

import (
	"../../../config"
	"../../../misc/exceptions"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Delete(objectID string, bucket string) {

	sess, err := session.NewSession(&aws.Config{
		Credentials: config.AWSCredentials,
		Region:      &config.AWSRegion,
	})

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(objectID)})
	if err != nil {
		panic(exceptions.Exception{
			Cause: "S3: Error ocurred while deleting " + objectID + " form " + bucket,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectID),
	})

	if err != nil {
		panic(exceptions.Exception{
			Cause: "S3: Error ocurred while deleting " + objectID + " form " + bucket,
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

}
