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

// Get user data
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		loggedInUsername, _ := c.Get("username") // Fetching authorized username

		username := c.DefaultQuery("name", loggedInUsername.(string)) // If no username is provided, then it will return authorized user's own details

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		queryForUserDataFetch := "SELECT username, firstname, lastname, email, premium FROM users WHERE username = ?"

		var userdata models.SavedUser

		if err := db.QueryRowContext(ctx, queryForUserDataFetch, username).Scan(&userdata.UserName, &userdata.FirstName, &userdata.LastName, &userdata.Email, &userdata.IsPremium); err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"username": userdata.UserName, "fullname": userdata.FirstName + " " + userdata.LastName, "email": userdata.Email, "premium": userdata.IsPremium})
	}
}

// http://localhost:5000/user/edit?username=pramit&email=pramit@gmail.com&firstname=Pramit&lastname=Mondal
// Edit user account details using query
func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, _ := c.Get("userid")

		username := c.Query("username")
		email := c.Query("email")
		firstname := c.Query("firstname")
		lastname := c.Query("lastname")

		if username != "" {
			// Check whether the username is already taken or not
			sqlQueryToCheckUsernameAvailablity := "SELECT user_id FROM users WHERE username = ?"
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			var savedUserid string
			err := db.QueryRowContext(ctx, sqlQueryToCheckUsernameAvailablity, username).Scan(&savedUserid)
			if err == sql.ErrNoRows {
				sqlQueryToUpdateUsername := "UPDATE users SET username = ? WHERE user_id = ?"
				statement, err := db.PrepareContext(ctx, sqlQueryToUpdateUsername)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				res, err := statement.ExecContext(ctx, username, userid)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				rows, err := res.RowsAffected()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"rows": rows})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "username is already taken"})
				return
			}
		}

		if email != "" {
			// Check whether the email is already taken or not
			sqlQueryToCheckEmailAvailablity := "SELECT user_id FROM users WHERE email = ?"
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			var savedUserid string
			err := db.QueryRowContext(ctx, sqlQueryToCheckEmailAvailablity, email).Scan(&savedUserid)
			if err == sql.ErrNoRows {
				sqlQueryToUpdateEmail := "UPDATE users SET email = ? WHERE user_id = ?"
				statement, err := db.PrepareContext(ctx, sqlQueryToUpdateEmail)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				res, err := statement.ExecContext(ctx, email, userid)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				rows, err := res.RowsAffected()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"rows": rows})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "email is already taken"})
				return
			}
		}

		if firstname != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			sqlQueryToUpdateFirstName := "UPDATE users SET firstname = ? WHERE user_id = ?"
			statement, err := db.PrepareContext(ctx, sqlQueryToUpdateFirstName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			res, err := statement.ExecContext(ctx, firstname, userid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			rows, err := res.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"rows": rows})

		}

		if lastname != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			sqlQueryToUpdateLastName := "UPDATE users SET lastname = ? WHERE user_id = ?"
			statement, err := db.PrepareContext(ctx, sqlQueryToUpdateLastName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			res, err := statement.ExecContext(ctx, lastname, userid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			rows, err := res.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"rows": rows})

		}
	}
}

// Delete Account
func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, _ := c.Get("userid")

		sqlQueryToDeleteRows := "DELETE FROM users WHERE user_id = ?"

		res, err := db.Exec(sqlQueryToDeleteRows, userid)
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
		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
	}
}

// Purchase Premium / Get Premium
// http://localhost:5000/user/purchase_premium?year=2
func PurchasePremium() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, _ := c.Get("userid")
		year := c.DefaultQuery("year", "1")

		yearValue, err := strconv.ParseInt(year, 0, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if yearValue <= 0 || yearValue > 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		// Query to check whether the user has premium membership or not
		queryToCheckAlreadyPremium := "SELECT premium, premiumExpiry FROM users WHERE user_id = ?"

		var isPremiumMember string
		var premiumExpiryDate string

		if err := db.QueryRow(queryToCheckAlreadyPremium, userid).Scan(&isPremiumMember, &premiumExpiryDate); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "no user found"})
				return
			}
		}

		// Query to update premium membership in DB
		queryToAddPremium := "UPDATE users SET premium = 'true', premiumPurchase = ?, premiumExpiry = ? WHERE user_id = ?"

		Purchased, _ := time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		Expiry, _ := time.Parse(time.RFC1123, time.Now().Add(time.Hour*24*365*time.Duration(yearValue)).Format(time.RFC1123))

		if isPremiumMember == "false" {
			result1, err1 := db.Exec(queryToAddPremium, Purchased, Expiry, userid)

			if err1 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
				return
			}

			rows1, err := result1.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			fmt.Println(rows1)

			var unit string
			if yearValue == 1 {
				unit = " year"
			} else {
				unit = " years"
			}
			c.JSON(http.StatusOK, gin.H{"message": "Purchased subscription for " + year + unit})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Your already have a premium subscription", "Expiry": premiumExpiryDate})
			return
		}
	}
}
