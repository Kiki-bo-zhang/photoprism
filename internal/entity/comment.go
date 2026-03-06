package entity

import (
	"time"

	"github.com/photoprism/photoprism/pkg/rnd"
)

// Comment represents a user comment on a photo.
type Comment struct {
	ID        uint      `gorm:"primary_key" yaml:"-"`
	UID       string    `gorm:"type:VARBINARY(42);unique_index;" json:"UID" yaml:"UID"`
	PhotoUID  string    `gorm:"type:VARBINARY(42);index;" json:"PhotoUID" yaml:"PhotoUID"`
	UserUID   string    `gorm:"type:VARBINARY(42);index;" json:"UserUID" yaml:"UserUID,omitempty"`
	UserName  string    `gorm:"type:VARCHAR(200);" json:"UserName" yaml:"UserName,omitempty"`
	Content   string    `gorm:"type:VARCHAR(4000);" json:"Content" yaml:"Content"`
	Type      string    `gorm:"type:VARBINARY(8);default:'text';" json:"Type" yaml:"Type"` // text, audio, video
	MediaURL  string    `gorm:"type:VARBINARY(1024);" json:"MediaURL,omitempty" yaml:"MediaURL,omitempty"`
	CreatedAt time.Time `json:"CreatedAt" yaml:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" yaml:"UpdatedAt"`
}

// TableName returns the entity database table name.
func (Comment) TableName() string {
	return "comments"
}

// NewComment creates a new comment.
func NewComment(photoUID, userUID, userName, content string) *Comment {
	return &Comment{
		UID:      rnd.GenerateUID('c'),
		PhotoUID: photoUID,
		UserUID:  userUID,
		UserName: userName,
		Content:  content,
		Type:     "text",
	}
}

// NewMediaComment creates a new media comment (audio/video).
func NewMediaComment(photoUID, userUID, userName, commentType, mediaURL string) *Comment {
	return &Comment{
		UID:      rnd.GenerateUID('c'),
		PhotoUID: photoUID,
		UserUID:  userUID,
		UserName: userName,
		Type:     commentType,
		MediaURL: mediaURL,
	}
}

// CreateComment creates the comment in the database.
func (m *Comment) CreateComment() error {
	return Db().Create(m).Error
}

// SaveComment saves the comment in the database.
func (m *Comment) SaveComment() error {
	return Db().Save(m).Error
}

// DeleteComment deletes the comment from the database.
func (m *Comment) DeleteComment() error {
	return Db().Delete(m).Error
}

// FindCommentByUID finds a comment by its UID.
func FindCommentByUID(uid string) (*Comment, error) {
	var comment Comment
	if err := Db().Where("uid = ?", uid).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindCommentsByPhotoUID finds all comments for a photo.
func FindCommentsByPhotoUID(photoUID string) ([]Comment, error) {
	var comments []Comment
	if err := Db().Where("photo_uid = ?", photoUID).Order("created_at asc").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// DeleteCommentsByPhotoUID deletes all comments for a photo.
func DeleteCommentsByPhotoUID(photoUID string) error {
	return Db().Where("photo_uid = ?", photoUID).Delete(&Comment{}).Error
}
