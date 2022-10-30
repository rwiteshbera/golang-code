package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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
func UpvotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get loggedin userid
		userid, _ := c.Get("userid")
		postid := c.Query("post")

		postid_int, err := strconv.ParseUint(postid, 0, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//  Check if user already upvoted the post or not?
		// condition 1 : If upvoted -> Remove Upvote
		// condition 2 : If not -> Insert Upvote

		queryToCheckUpvote := "SELECT upvote_id FROM upvotes WHERE post_id = ? AND upvoted_by = ?"

		var upvote_id int

		err1 := db.QueryRow(queryToCheckUpvote, postid_int, userid).Scan(&upvote_id)
		if err1 == sql.ErrNoRows {
			queryToUpvote := "INSERT upvotes (post_id, upvoted_by) VALUES(?, ?)"
			res, err := db.Exec(queryToUpvote, postid_int, userid)
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
			c.JSON(http.StatusOK, gin.H{"message": "upvoted"})
		} else {
			queryToRemoveUpvote := "DELETE FROM upvotes WHERE post_id = ? AND upvoted_by = ?"
			res, err := db.Exec(queryToRemoveUpvote, postid_int, userid)
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
			c.JSON(http.StatusOK, gin.H{"message": "remove upvote"})
		}
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get loggedin userid
		userid, _ := c.Get("userid")
		postid := c.Query("post")

		postid_int, err := strconv.ParseUint(postid, 0, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		queryToRemoveUpvote := "DELETE FROM posts WHERE post_id = ? AND posted_by = ?"
		res, err := db.Exec(queryToRemoveUpvote, postid_int, userid)
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
		c.JSON(http.StatusOK, gin.H{"message": "Deleted post successfully"})
	}
}
