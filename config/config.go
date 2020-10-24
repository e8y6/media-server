package config

import "github.com/aws/aws-sdk-go/aws/credentials"

// TODO move to some env based thing. Fine for now.
var (

	// APP_PORT port on which app listens to external apps
	APP_PORT_EXTERNAL = 8001

	// APP_PORT port on which app listens to the internal apps
	APP_PORT_INTERNAL = 8002

	// DATABASE_NAME name
	DATABASE_NAME = "johny_walker"

	// DATABASE_CONNECTION_URI string
	DATABASE_CONNECTION_URI = "mongodb://localhost:27017/"

	// LOCAL_FOLDER port on which app has tosearch for media
	LOCAL_FOLDER = "./persist/uploads/"

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
		"media_store": "prod-jw-storage",
	}

	// VIMEO_OAuthToken vimeo Oauth token
	VIMEO_OAuthToken = "0eb8d07493ce6d7131e4087161109f37"
	// VIMEO_Folders : VIMEO_Folders flolder list
	VIMEO_Folders = map[string]string{
		"generic": "2690221",
	}

	// CF_STREAM_TOKEN Token
	CF_STREAM_TOKEN = "c1757d9a32a006edc2b3a90444f5749da225c"
	CF_ACCOUNT_ID   = "c2cea9aa148e4a8aaf4021d9a8f99a3b"
	CF_EMAIL        = "jobinrjohnson@gmail.com"
)
