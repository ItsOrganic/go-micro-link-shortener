package mongo

import (
	"context"
	"time"

	"github.com/ItsOrganic/go-micro-link-shortener/shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
type mongoRepo struct {
    client *mongo.Client
    database string
    timeout time.Duration
}
func newMongoClient(mongoURL string,mongoTimeout int)(*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
    if err != nil {
        return nil, err
    }
    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        return nil, err
    }

    return client, err
}

func NewMongoRepo(mongoURL string, mongo string, mongoTimeout int) (shortener.RedirectRepo, error) {
    repo := &mongoRepo{
        timeout: time.Duration(mongoTimeout)*time.Second,
        database: mongo,
    }
    client, err := newMongoClient(mongoURL, mongoTimeout)
    if err != nil {
        return nil, errors.Wrap(err, "repo.newMongoClient")
    }
    repo.client = client
    return repo, nil
}
func(r *mongoRepo) Find(code string) (*shortener.Redirect, error){
    ctx, cancel := context.WithTimeout(context.Background(),r.timeout)
    defer cancel()
    redirect := &shortener.Redirect{}
    collection := r.client.Database(r.database).Collection("redirects")
    filter := bson.M{"code":code}
    err := collection.FindOne(ctx, filter).Decode(&redirect)
    if err != nil {
        if err == mongo.ErrNoDocuments{
            return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repo.Redirect.Find")
        }
        return nil, errors.Wrap(shortener.ErrRedirectInvalid, "repo.Redirect.Find")
    }
    return redirect, nil
}

func(r *mongoRepo) Store(redirect *shortener.Redirect) (error){
    ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
    defer cancel()
    collection := r.client.Database(r.database).Collection("redirects")
    _, err := collection.InsertOne(
    ctx,
    bson.M{
        "code": redirect.Code,
        "url": redirect.URL,
        "created_at": redirect.CreatedAt,
    },
    )
    if err != nil {
        return errors.Wrap(err, "repo.Redirect.Store")
    }
    return nil
}




