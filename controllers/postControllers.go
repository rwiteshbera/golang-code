package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/orbit/models"
)

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logged in userid
		userid, _ := c.Get("userid")

		var post models.Post

		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		queryToUploadPost := "INSERT INTO posts (text_data, image_url, posted_by, posted_on) VALUES(?, ?, ?, ?)"

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		statement, err := db.PrepareContext(ctx, queryToUploadPost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		postedOn, _ := time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		res, err := statement.ExecContext(ctx, post.Text_Data, post.ImageURL, userid, postedOn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(rows)
		c.JSON(http.StatusOK, gin.H{"message": "Posted Successfully"})
	}
}

// Upvote a post
// http://localhost:5000/upvote?post_id=4
func UpVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logged in userid
		userid, _ := c.Get("userid")
		post_id := c.Query("post_id")

		queryToUpvote := "UPDATE posts SET upvote = upvote + 1 WHERE post_id = ? AND posted_by = ?"

		res, err := db.Exec(queryToUpvote, post_id, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(rows)
		c.JSON(http.StatusOK, gin.H{"message": "Upvoted"})
	}
}

// Downvote a post
// http://localhost:5000/downvote?post_id=3
func DownVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logged in userid
		userid, _ := c.Get("userid")
		post_id := c.Query("post_id")

		queryToUpvote := "UPDATE posts SET downvote = downvote + 1 WHERE post_id = ? AND posted_by = ?"

		res, err := db.Exec(queryToUpvote, post_id, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(rows)
		c.JSON(http.StatusOK, gin.H{"message": "Downvoted"})
	}
}
