package config

import "github.com/aws/aws-sdk-go/aws/credentials"

// TODO move to some env based thing. Fine for now.
var (

	// APP_PORT port on which app listens to
	APP_PORT = 8001

	// LOCAL_FOLDER port on which app has tosearch for media
	LOCAL_FOLDER = "./uploads/"

	// AWSCredentials : AWS credentials
	AWSCredentials = credentials.NewStaticCredentials(
		"AKIAQKMXTOR2KBPU3TLH",
		"VQ+9S7pqbj6eaJL2xwl8+rkWnLM/qZs1GrNYtF2N",
		"", // a token will be created when the session it's used.
	)
	// AWSRegion : AWS Location
	AWSRegion = "ap-south-1"

	// AWSBuckets : AWS Bucket list
	AWSBuckets = map[string]string{
		"media_store": "uat-jw-storage",
	}

	// VIMEO_OAuthToken vimeo Oauth token
	VIMEO_OAuthToken = "ef99c2e4eb68eebca1f4fce3c62742b6"
	// VIMEO_Folders : VIMEO_Folders flolder list
	VIMEO_Folders = map[string]string{
		"generic": "2690221",
	}
)
