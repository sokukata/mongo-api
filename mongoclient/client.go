package mongoclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sokukata/mongo-api/convert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Usage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, HELPER)
}

type Credential struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

var (
	database *mongo.Database
)

// LOGIN

func login(w http.ResponseWriter, cred Credential) {
	auth := &options.Credential{
		Username:    cred.Id,
		Password:    cred.Password,
		PasswordSet: true,
		AuthSource:  "admin",
	}

	strURI, err := MongoURI(Server, cred.Id, cred.Password)
	if err != nil {
		badRequestError(w, err)
		return
	}

	opts := options.Client()
	opts = opts.ApplyURI(strURI)
	opts = opts.SetAuth(*auth)

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		internalError(w, err)
		return
	}
	// Check credentials
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		badRequestError(w, err)
		return
	}
	database = mongoClient.Database(DatabaseName)
}

func Login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var cred Credential
	json.Unmarshal(reqBody, &cred)
	login(w, cred)
}

// CREATE

func addUsers(w http.ResponseWriter, mUsers []map[string]interface{}) {
	collection := database.Collection(Collection)
	ctx := context.Background()

	var docs []interface{}
	ch := make(chan bson.M)
	errc := make(chan error)
	for _, m := range mUsers {
		go func(mUser map[string]interface{}) {
			opts := options.FindOne().SetProjection(bson.D{{"id", 1}})
			res := collection.FindOne(ctx, bson.M{"id": mUser["id"]}, opts)
			err := res.Err()
			if err == nil {
				// Ignore user already in base
				ch <- nil
				return
			}
			if err != mongo.ErrNoDocuments {
				errc <- err
				return
			}

			convRec := convert.ConvertMap(mUser, false)
			hash, err := bcrypt.GenerateFromPassword([]byte(convRec["password"].(string)), 10 /* Default cost*/)
			if err != nil {
				errc <- err
				return
			}
			convRec["password"] = string(hash)

			err = os.WriteFile(filepath.Join(DIRECTORY, convRec["id"].(string)), []byte(convRec["data"].(string)), 0644)
			if err != nil {
				errc <- err
				return
			}
			ch <- convRec
		}(m)
	}
	//count := 0
	var err error
	for i := 0; i < len(mUsers); i++ {
		select {
		case err := <-errc:
			internalError(w, err)
			return
		case cr := <-ch:
			if cr != nil {
				docs = append(docs, cr)
			}
		}
	}
	if len(docs) == 0 {
		// InsertMany need at least one element
		return
	}
	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		internalError(w, err)
		return
	}
}

func AddUsers(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	mUsers := make([]map[string]interface{}, 0)
	json.Unmarshal(reqBody, &mUsers)
	addUsers(w, mUsers)
}

// DELETE

func deleteUser(w http.ResponseWriter, id string) {
	collection := database.Collection(Collection)
	ctx := context.Background()

	path := filepath.Join(DIRECTORY, id)
	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
		if err != nil {
			internalError(w, err)
			return
		}
	}

	deleteRes, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		internalError(w, err)
		return
	}
	if deleteRes.DeletedCount == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := getId(r)
	deleteUser(w, id)
}

// READ

func getUsersList(w http.ResponseWriter) {
	collection := database.Collection(Collection)
	ctx := context.Background()
	cur, err := collection.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		internalError(w, err)
		return
	}
	records := make([]map[string]interface{}, 0)
	// get at most rowcount records
	for cur.Next(ctx) {
		var itemBson bson.D
		err = cur.Decode(&itemBson)
		if err != nil {
			internalError(w, err)
			return
		}

		// Convert from mongo document
		rec, err := convert.ConvertDoc(itemBson)
		if err != nil {
			internalError(w, err)
			return
		}
		records = append(records, rec)
	}
	json.NewEncoder(w).Encode(records)
}

func GetUsersList(w http.ResponseWriter, r *http.Request) {
	getUsersList(w)
}

func getUser(w http.ResponseWriter, id string) {
	collection := database.Collection(Collection)
	ctx := context.Background()
	res := collection.FindOne(ctx, bson.M{"id": id})
	err := res.Err()
	if err == mongo.ErrNoDocuments {
		notFoundError(w, err)
		return
	}
	if err != nil {
		internalError(w, err)
		return
	}

	var itemBson bson.D
	err = res.Decode(&itemBson)
	if err != nil {
		internalError(w, err)
		return
	}

	rec, err := convert.ConvertDoc(itemBson)
	if err != nil {
		internalError(w, err)
		return
	}
	json.NewEncoder(w).Encode(rec)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := getId(r)
	getUser(w, id)
}

// UPDATE

func updateUser(w http.ResponseWriter, id string, mUser map[string]interface{}) {
	collection := database.Collection(Collection)
	ctx := context.Background()

	convRec := convert.ConvertMap(mUser, true)
	updateRes, err := collection.UpdateOne(ctx, bson.M{"id": id}, convRec)
	if err != nil {
		internalError(w, err)
		return
	}
	if updateRes.MatchedCount == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, ok := convRec["$set"].(bson.M)["data"]
	if !ok {
		return
	}

	path := filepath.Join(DIRECTORY, id)
	err = os.WriteFile(path, []byte(data.(string)), 0644)
	if err != nil {
		internalError(w, err)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	reqBody, _ := ioutil.ReadAll(r.Body)
	mUser := make(map[string]interface{}, 0)
	json.Unmarshal(reqBody, &mUser)
	updateUser(w, id, mUser)
}
