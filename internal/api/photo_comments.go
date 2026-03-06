package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/i18n"
	"github.com/photoprism/photoprism/pkg/txt"
)

// GetPhotoComments returns all comments for a photo.
//
//	@Summary		returns all comments for a photo
//	@Id				GetPhotoComments
//	@Tags			Photos
//	@Produce		json
//	@Success		200				{array}	entity.Comment
//	@Failure		401,403,404,500	{object}	i18n.Response
//	@Router			/api/v1/photos/:uid/comments [get]
func GetPhotoComments(router *gin.RouterGroup) {
	router.GET("/photos/:uid/comments", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionView)

		if s.Abort(c) {
			return
		}

		photoUID := clean.UID(c.Param("uid"))
		if photoUID == "" {
			AbortBadRequest(c)
			return
		}

		// Verify photo exists
		photo, err := query.PhotoByUID(photoUID)
		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		// Check permissions
		if photo.PhotoPrivate && !acl.Events.Allow(acl.ResourcePhotos, s.GetUserRole(), acl.ActionView) {
			AbortEntityNotFound(c)
			return
		}

		comments, err := entity.FindCommentsByPhotoUID(photoUID)
		if err != nil {
			AbortUnexpected(c)
			return
		}

		c.JSON(http.StatusOK, comments)
	})
}

// AddPhotoComment adds a new comment to a photo.
//
//	@Summary		adds a new comment to a photo
//	@Id				AddPhotoComment
//	@Tags			Photos
//	@Produce		json
//	@Success		200				{object}	entity.Comment
//	@Failure		400,401,403,404,500	{object}	i18n.Response
//	@Router			/api/v1/photos/:uid/comments [post]
func AddPhotoComment(router *gin.RouterGroup) {
	router.POST("/photos/:uid/comments", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionUpdate)

		if s.Abort(c) {
			return
		}

		photoUID := clean.UID(c.Param("uid"))
		if photoUID == "" {
			AbortBadRequest(c)
			return
		}

		// Verify photo exists
		photo, err := query.PhotoByUID(photoUID)
		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		// Check permissions
		if photo.PhotoPrivate && !acl.Events.Allow(acl.ResourcePhotos, s.GetUserRole(), acl.ActionUpdate) {
			AbortEntityNotFound(c)
			return
		}

		var f form.Comment

		if err := c.BindJSON(&f); err != nil {
			AbortBadRequest(c)
			return
		}

		user := s.User()
		comment := entity.NewComment(photoUID, user.UserUID, user.UserName, f.Content)
		
		if f.Type != "" {
			comment.Type = f.Type
		}
		if f.MediaURL != "" {
			comment.MediaURL = f.MediaURL
		}

		if err := comment.CreateComment(); err != nil {
			AbortUnexpected(c)
			return
		}

		c.JSON(http.StatusOK, comment)
	})
}

// DeletePhotoComment deletes a comment from a photo.
//
//	@Summary		deletes a comment from a photo
//	@Id				DeletePhotoComment
//	@Tags			Photos
//	@Produce		json
//	@Success		200				{object}	entity.Comment
//	@Failure		401,403,404,500	{object}	i18n.Response
//	@Router			/api/v1/photos/:uid/comments/:commentuid [delete]
func DeletePhotoComment(router *gin.RouterGroup) {
	router.DELETE("/photos/:uid/comments/:commentuid", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionDelete)

		if s.Abort(c) {
			return
		}

		photoUID := clean.UID(c.Param("uid"))
		commentUID := clean.UID(c.Param("commentuid"))
		
		if photoUID == "" || commentUID == "" {
			AbortBadRequest(c)
			return
		}

		comment, err := entity.FindCommentByUID(commentUID)
		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		// Verify the comment belongs to the photo
		if comment.PhotoUID != photoUID {
			AbortEntityNotFound(c)
			return
		}

		// Check permissions - only allow delete if user is comment author or has admin role
		user := s.User()
		if comment.UserUID != user.UserUID && !acl.Events.Allow(acl.ResourcePhotos, s.GetUserRole(), acl.ActionDelete) {
			AbortForbidden(c)
			return
		}

		if err := comment.DeleteComment(); err != nil {
			AbortUnexpected(c)
			return
		}

		c.JSON(http.StatusOK, comment)
	})
}
