package config

import "github.com/aws/aws-sdk-go/aws/credentials"

var (
	// AWSCredentials : AWS credentials
	AWSCredentials = credentials.NewStaticCredentials(
		"AKIAQKMXTOR2KBPU3TLH",
		"VQ+9S7pqbj6eaJL2xwl8+rkWnLM/qZs1GrNYtF2Nd",
		"", // a token will be created when the session it's used.
	)
	// AWSRegion : AWS Location
	AWSRegion = "ap-south-1"

	// AWSBuckets : AWS Bucket list
	AWSBuckets = map[string]string{
		"media_store": "uat-jw-storage",
	}
)
