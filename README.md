```sql
CREATE TABLE users (
	user_id VARCHAR(40) NOT NULL,
  	username VARCHAR(30) NOT NULL UNIQUE,
    firstname VARCHAR(40) NOT NULL,
  	lastname VARCHAR(40) NOT NULL,
    email VARCHAR(40) NOT NULL UNIQUE,
    hash_password VARCHAR(100) NOT NULL,
  	premium VARCHAR(5) NOT NULL DEFAULT 'false',
  	premiumPurchase VARCHAR(30),
  	premiumExpiry VARCHAR(30),
    last_login VARCHAR(30) NOT NULL,
    created_at VARCHAR(30) NOT NULL,
    PRIMARY KEY (user_id)
);


CREATE TABLE posts (
	post_id INT NOT NULL UNIQUE AUTO_INCREMENT,
  	text_data VARCHAR(10000) NOT NULL,
  	image_url VARCHAR(100),
    posted_by VARCHAR(40) NOT NULL,
  	posted_on VARCHAR(30) NOT NULL,
  	PRIMARY KEY (post_id),
    FOREIGN KEY (posted_by) REFERENCES users(user_id) ON DELETE CASCADE
);


CREATE TABLE upvotes (
	upvote_id INT NOT NULL UNIQUE AUTO_INCREMENT,
  	post_id INT NOT NULL,
  	upvoted_by VARCHAR(40) NOT NULL,
  	PRIMARY KEY (upvote_id),
  	FOREIGN KEY (upvoted_by) REFERENCES users(user_id) ON DELETE CASCADE,
  	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);


CREATE TABLE comments (
	comment_id INT NOT NULL UNIQUE AUTO_INCREMENT,
  	post_id INT NOT NULL,
  	text_data VARCHAR(1000) NOT NULL,
  	commented_by VARCHAR(40) NOT NULL,
  	commented_on VARCHAR (30) NOT NULL,
  	PRIMARY KEY (comment_id),
  	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
  	FOREIGN KEY (commented_by) REFERENCES users(user_id) ON DELETE CASCADE
);

```