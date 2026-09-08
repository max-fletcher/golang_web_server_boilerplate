package formatters

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

// NOTE: Exporting type fields
// Turns out any type that has fields that are not exported as capital case(i.e first char of name is capital) is
// not exported in golang. So if you used id instead of ID, the returned struct(and consequently the JSON) will be missing that field.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NOTE: Exporting functions from a folder(a.k.a a package as different folders are each considered a package in go)
// Turns out any exported function from a folder/package needs to be pascal-case(or at least starts with a capital case)
// else, the function is not made available when imported and used. This is why the convention for writing golang often
// is to use pascal-case.

// This is a function that converts struct keys from pascal-case to use camel-case. This is so we get camel-cased JSON.
// Remember that database.go has the type for the user being fetched from database and that is being passed here.
func DatabaseUserToUser(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func DatabaseUsersToUsers(dbUsers []db.User) []User {
	users := []User{}

	// var description *string // a var containing a pointer to a string
	for _, dbUser := range dbUsers {
		// if dbUser.Description.Valid {
		// 	description = &dbUser.Description.String
		// }
		users = append(users, User{
			ID:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
			// Avatar: dbUser.Avatar,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		})

		// This also works if you want to reuse functionality i.e replace above block with the line below
		// posts = append(posts, DatabasePostsToPost(dbFeed))
	}

	return users
}

type Post struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	// Although we are passing an sql.NullString inside dbPosts, we are doing this since passing a pointer will cause either
	// the string or the null value(the value that the pointer is pointing to) to be inside the description field when it is returned as json.
	// This is how marshalling to json in go works i.e if a pointer points to null, it will return null in that json field,
	// and if a pointer points to a string, it will return string in that json field. Otherwise, if we directly used post.Description,
	// due to the, dbFeed.description struct containing nested fields(sql.NullString obj) it will be marshalled to
	// "description": { "String" : "Some des", Valid : true }
	// Description   *string   `json:"description"`
	Content   *string   `json:"content"`
	Photo     *string   `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabasePostToPost(dbPost db.Post) Post {
	// var description *string // a var containing a pointer to a string
	// if dbPost.Content.Valid {
	// 	description = &dbPost.Description.String
	// }

	// *IMPORTANT: This is how you dead with nullable fields.(Sources: posts/structs.go(especially CreatePostRequest and CreatePostInput),
	// posts/services.go and posts/formatters.go)
	var photo *string
	if dbPost.Photo.Valid {
		photo = &dbPost.Photo.String
	}
	var content *string
	if dbPost.Content.Valid {
		photo = &dbPost.Content.String
	}

	return Post{
		ID:        dbPost.ID,
		Title:     dbPost.Title,
		Content:   content,
		Photo:     photo,
		CreatedAt: dbPost.CreatedAt,
		UpdatedAt: dbPost.UpdatedAt,
	}
}

func DatabasePostsToPosts(dbPosts []db.Post) []Post {
	posts := []Post{}
	var photo *string
	var content *string

	// var description *string // a var containing a pointer to a string
	for _, dbPost := range dbPosts {
		// if dbPost.Description.Valid {
		// 	description = &dbPost.Description.String
		// }

		if dbPost.Photo.Valid {
			photo = &dbPost.Photo.String
		}
		if dbPost.Content.Valid {
			photo = &dbPost.Content.String
		}

		posts = append(posts, Post{
			ID:        dbPost.ID,
			Title:     dbPost.Title,
			Content:   content,
			Photo:     photo,
			CreatedAt: dbPost.CreatedAt,
			UpdatedAt: dbPost.UpdatedAt,
		})
	}

	return posts
}

type CreatedByUser struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
type PostWithUser struct {
	ID        uuid.UUID     `json:"id"`
	Title     string        `json:"title"`
	Content   *string       `json:"content"`
	Photo     *string       `json:"photo"`
	User      CreatedByUser `json:"user"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func DatabasePostWUserToPostWUser(dbPostWUser db.GetPostWithUserByIdRow) PostWithUser {
	// *IMPORTANT: This is how you dead with nullable fields.(Sources: posts/structs.go(especially CreatePostRequest and CreatePostInput),
	// posts/services.go and posts/formatters.go)
	var photo *string
	if dbPostWUser.Photo.Valid {
		photo = &dbPostWUser.Photo.String
	}
	var content *string
	if dbPostWUser.Content.Valid {
		photo = &dbPostWUser.Content.String
	}

	return PostWithUser{
		ID:      dbPostWUser.ID,
		Title:   dbPostWUser.Title,
		Content: content,
		Photo:   photo,
		User: CreatedByUser{
			ID:    dbPostWUser.UserID,
			Name:  dbPostWUser.UserName,
			Email: dbPostWUser.UserEmail,
		},
		CreatedAt: dbPostWUser.PostCreatedAt,
		UpdatedAt: dbPostWUser.PostUpdatedAt,
	}
}

func DatabasePostsWUserToPostsWUser(dbPostsWithUser []db.GetPostsWithUserRow) []PostWithUser {
	posts := []PostWithUser{}
	var photo *string
	var content *string

	for _, dbPostWithUser := range dbPostsWithUser {
		if dbPostWithUser.Photo.Valid {
			photo = &dbPostWithUser.Photo.String
		}
		if dbPostWithUser.Content.Valid {
			photo = &dbPostWithUser.Content.String
		}

		posts = append(posts, PostWithUser{
			ID:      dbPostWithUser.ID,
			Title:   dbPostWithUser.Title,
			Content: content,
			Photo:   photo,
			User: CreatedByUser{
				ID:    dbPostWithUser.UserID,
				Name:  dbPostWithUser.UserName,
				Email: dbPostWithUser.UserEmail,
			},
			CreatedAt: dbPostWithUser.PostCreatedAt,
			UpdatedAt: dbPostWithUser.PostUpdatedAt,
		})
	}

	return posts
}

func StringPointerToNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *value,
		Valid:  true,
	}
}
