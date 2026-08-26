package formatters

import (
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

// NOTE: Exporting type fields
// Turns out any type that has fields that are not exported as pascal-case(or at least starts with a capital case) is
// not exported. So if you used id instead of ID, the returned struct(and consequently the JSON) will be missing that field.
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

type Post struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	// Although we are passing an sql.NullString inside dbPosts, we are doing this since passing a pointer will cause either
	// the string or the null value(the value that the pointer is pointing to) to be inside the description field when it is returned as json.
	// This is how marshalling to json in go works i.e if a pointer points to null, it will return null in that json field,
	// and if a pointer points to a string, it will return string in that json field. Otherwise, if we directly used post.Description,
	// due to the, dbFeed.description struct containing nested fields(sql.NullString obj) it will be marshalled to
	// "description": { "String" : "Some des", Valid : true }
	Content     *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// func DatabasePostsToPost(dbPost database.Post) Post {
// 	var description *string // a var containing a pointer to a string
// 	if dbPost.Description.Valid {
// 		description = &dbPost.Description.String
// 	}

// 	return Post{
// 		ID:          dbPost.ID,
// 		FeedID:      dbPost.FeedID,
// 		Title:       dbPost.Title,
// 		Description: description,
// 		PublishedAt: dbPost.PublishedAt,
// 		Url:         dbPost.Url,
// 		CreatedAt:   dbPost.CreatedAt,
// 		UpdatedAt:   dbPost.UpdatedAt,
// 	}
// }

// func DatabasePostsToPosts(dbPosts []database.Post) []Post {
// 	posts := []Post{}

// 	var description *string // a var containing a pointer to a string
// 	for _, dbPost := range dbPosts {
// 		if dbPost.Description.Valid {
// 			description = &dbPost.Description.String
// 		}
// 		posts = append(posts, Post{
// 			ID:          dbPost.ID,
// 			FeedID:      dbPost.FeedID,
// 			Title:       dbPost.Title,
// 			Description: description,
// 			PublishedAt: dbPost.PublishedAt,
// 			Url:         dbPost.Url,
// 			CreatedAt:   dbPost.CreatedAt,
// 			UpdatedAt:   dbPost.UpdatedAt,
// 		})

// 		// This also works if you want to reuse functionality i.e replace above block with the line below
// 		// posts = append(posts, DatabasePostsToPost(dbFeed))
// 	}

// 	return posts
// }
