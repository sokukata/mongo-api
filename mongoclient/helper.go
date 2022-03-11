package mongoclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	// SERVER           = "mongodb://127.0.0.1:27017"
	// DATABASE         = "dataimpact"
	// USERS_COLLECTION = "users"

	DIRECTORY = "output"

	HELPER = `Usage:
POST /login
	body (json):
		{
			"id":<Mongo username>,
			"password":<Mongo password>
		}
	Description:
		Log in to mongo Db to use API
	
POST /add/users
	body (json):
		[
			{
				id: ...,
				password: ...,
				data: ...,
				[...],
			},
			{
				id: ...,
				password: ...,
				data: ...,
				[...]
			}
		]
	Description:
		Add a list of users to Mongodb and create a file wth user 'data'

DELETE /delete/user/{id}
	args:
		id: user id to delete
	Description:
		Delete user by {id} in mongodb and filesystem

GET /users/list
	Description:	
		Get list of user from Moongodb

GET /user/{id}
	args:
		id: user id to read
	Description:	
		Get a user by {id}

PUT /user/{id}
	args:
		id: user id to read
	Body:
		{
			id: ...,
			password: ...,
			data: ...,
			[...],
		}
	Description:	
		Update user by {id} in Mongodb and filesystem
	`
)

var (
	Server       = "mongodb://127.0.0.1:27017"
	DatabaseName = "dataimpact"
	Collection   = "users"
)

func MongoURI(uri string, user string, password string) (string, error) {
	var strScheme string
	if strings.HasPrefix(uri, connstring.SchemeMongoDBSRV+"://") {
		strScheme = connstring.SchemeMongoDBSRV
		// remove the scheme
		uri = uri[len(connstring.SchemeMongoDBSRV)+3:]
	} else if strings.HasPrefix(uri, connstring.SchemeMongoDB+"://") {
		strScheme = connstring.SchemeMongoDB
		// remove the scheme
		uri = uri[len(connstring.SchemeMongoDB)+3:]
	} else {
		// if neither mongodb:// nor mongodb+srv://, we consider the default as
		//   mongodb://
		strScheme = connstring.SchemeMongoDB
	}

	// Is there a <something>@... ? If yes, this is an error
	if idx := strings.Index(uri, "@"); idx != -1 {
		err := fmt.Errorf("Mongodb URI should not include credentials")
		return "", err
	}

	// And rebuild URI from pieces
	if user == "" {
		// case of empty user
		return "", fmt.Errorf("No user set")
	}
	ret := fmt.Sprintf("%s://%s:%s@%s",
		strScheme,
		url.QueryEscape(user),
		url.QueryEscape(password),
		uri)
	return ret, nil
}

// func unauthorized(w http.ResponseWriter, r *http.Request) {
// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// }

// func AuthenticateMiddleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	if database == nil {
// 		return unauthorized
// 	}
// 	return f
// }

func AuthenticateMiddleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if database == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

func getId(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError)+": "+err.Error(), http.StatusInternalServerError)
}

func notFoundError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusNotFound)+": "+err.Error(), http.StatusNotFound)
}

func badRequestError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusBadRequest)+": "+err.Error(), http.StatusBadRequest)
}
