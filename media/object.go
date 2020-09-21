package media

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BUCKET_LOCAL  = 0
	BUCKET_AWS_S3 = 1

	FILETYPE_IMAGE   = 0
	FILETYPE_VIDEO   = 1
	FILETYPE_UNKNOWN = 2
)

// FileModel struct for holding th file metadata
type FileModel struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"_id_user,omitempty" bson:"_id_user"`

	FileType     string `json:"file_type,omitempty" bson:"file_type"`
	OriginalName string `json:"original_name,omitempty" bson:"original_name"`
	Size         int64  `json:"size,omitempty" bson:"size"`

	Privacy    int8              `json:"privacy,omitempty" bson:"privacy"`
	Bucket     int8              `json:"bucket,omitempty" bson:"bucket"`
	BucketMeta map[string]string `json:"bucket_meta,omitempty" bson:"bucket_meta"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type Media interface {
	Save()

	MoveMediaSafe()
	Optimize()
}
