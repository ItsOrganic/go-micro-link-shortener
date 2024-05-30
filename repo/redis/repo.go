package redis

import (
	"fmt"
	"strconv"

	"github.com/ItsOrganic/go-micro-link-shortener/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)
type redisRepo struct {
    client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error){
    opts, err := redis.ParseURL(redisURL)
    if err != nil{
        return nil, err
    }
    client := redis.NewClient(opts)
    _, err = client.Ping().Result()
    if err != nil{
        return nil, err
    }
    return client, nil
}
func NewRedisRepo(redisURL string) (shortener.RedirectRepo, error){
    repo := &redisRepo{}
        client, err := newRedisClient(redisURL)
        if err != nil {
            return nil, errors.Wrap(err, "repository.NewRedisRepo")
        }
        repo.client = client
        return repo, nil
    }
func(r *redisRepo) generateKey(code string) string{
    return fmt.Sprintf("redirect:%s", code)
}

func(r *redisRepo) Find(code string) (*shortener.Redirect, error) {
    redirect := &shortener.Redirect{}
    key := r.generateKey(code)
    data, err := r.client.HGetAll(key).Result()
    if err != nil {
        return nil, errors.Wrap(err, "repository.Redirect.Find")
    }
    if len(data) <= 0 {
        return nil, errors.Wrap(err, "repository.Redirect.Find")
    }
    createdAt, err := strconv.ParseInt(data["created_at"],10,64)
    if err != nil {
        return nil, errors.Wrap(err, "repository.Redirect.Find")
    }
    redirect.Code = data["code"]
    redirect.URL = data["url"]
    redirect.CreatedAt = createdAt
    return redirect, nil
}

func (r *redisRepo) Store(redirect *shortener.Redirect) error {
    key := r.generateKey(redirect.Code)
    data := map[string]interface{}{
        "code": redirect.Code,
        "url": redirect.URL,
        "created_at": redirect.CreatedAt,
    }
    _, err := r.client.HMSet(key, data).Result()
    if err != nil {
        return errors.Wrap(err, "repo.Redirect.Store")
    }
    return nil
}







