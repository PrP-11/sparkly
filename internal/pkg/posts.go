package pkg

import "time"

// type PostInteractionAction string

// const (
// 	Like          PostInteractionAction = "LIKE"
// 	Unlike        PostInteractionAction = "UNLIKE"
// 	Comment       PostInteractionAction = "COMMENT"
// 	DeleteComment PostInteractionAction = "DELETE_COMMENT"
// 	Share         PostInteractionAction = "SHARE"
// )

type PostInteraction struct {
	UserID    string    `json:"userId" bson:"userId"`
	PostID    string    `json:"postId" bson:"postId"`
	Action    string    `json:"action" bson:"action"`
	CommentID string    `json:"commentId,omitempty" bson:"commentId,omitempty"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type PopularPost struct {
	PostID string `json:"postId"`
	Count  int64  `json:"count"`
}

const CachekeyPopularPosts = "popular_posts_%v"
