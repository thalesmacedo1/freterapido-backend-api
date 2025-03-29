// redis_client_test.go

package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestSetCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	expiration := 5 * time.Minute

	mock.ExpectSet(key, value, expiration).SetVal("OK")

	err := redisClient.SetCache(ctx, key, value, expiration)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	expiration := 5 * time.Minute

	mock.ExpectSet(key, value, expiration).SetErr(errors.New("set error"))

	err := redisClient.SetCache(ctx, key, value, expiration)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set cache for key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"
	expectedValue := "testValue"

	mock.ExpectGet(key).SetVal(expectedValue)

	value, err := redisClient.GetCache(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCacheNotExists(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "nonExistentKey"

	mock.ExpectGet(key).RedisNil()

	value, err := redisClient.GetCache(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache key "+key+" does not exist")
	assert.Equal(t, "", value)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"

	mock.ExpectGet(key).SetErr(errors.New("get error"))

	value, err := redisClient.GetCache(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get cache for key")
	assert.Equal(t, "", value)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"

	mock.ExpectDel(key).SetVal(1)

	err := redisClient.DeleteCache(ctx, key)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"

	mock.ExpectDel(key).SetErr(errors.New("delete error"))

	err := redisClient.DeleteCache(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete cache for key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	keyExists := "existingKey"
	keyNotExists := "nonExistingKey"

	// Caso onde a chave existe
	mock.ExpectExists(keyExists).SetVal(1)
	exists, err := redisClient.ExistsCache(ctx, keyExists)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Caso onde a chave não existe
	mock.ExpectExists(keyNotExists).SetVal(0)
	exists, err = redisClient.ExistsCache(ctx, keyNotExists)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsCacheError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "testKey"

	mock.ExpectExists(key).SetErr(errors.New("exists error"))

	exists, err := redisClient.ExistsCache(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check existence of cache key")
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetCacheJSON(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "jsonKey"
	value := map[string]interface{}{
		"field1": "value1",
		"field2": 2,
	}
	expiration := 10 * time.Minute

	marshaledValue, _ := json.Marshal(value)

	mock.ExpectSet(key, marshaledValue, expiration).SetVal("OK")

	err := redisClient.SetCacheJSON(ctx, key, value, expiration)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetCacheJSONMarshalError(t *testing.T) {
	redisClient := &RedisClient{} // O client não será usado neste teste

	ctx := context.Background()
	key := "jsonKey"
	value := make(chan int) // Canais não podem ser marshalizados para JSON
	expiration := 10 * time.Minute

	err := redisClient.SetCacheJSON(ctx, key, value, expiration)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal value for key")
}

func TestGetCacheJSON(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "jsonKey"
	expectedValue := map[string]interface{}{
		"field1": "value1",
		"field2": 2.0,
	}
	marshaledValue, _ := json.Marshal(expectedValue)

	mock.ExpectGet(key).SetVal(string(marshaledValue))

	var result map[string]interface{}
	err := redisClient.GetCacheJSON(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCacheJSONGetError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "jsonKey"

	mock.ExpectGet(key).SetErr(errors.New("get error"))

	var result map[string]interface{}
	err := redisClient.GetCacheJSON(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get cache for key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCacheJSONUnmarshalError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	ctx := context.Background()
	key := "jsonKey"
	invalidJSON := "invalid json"

	mock.ExpectGet(key).SetVal(invalidJSON)

	var result map[string]interface{}
	err := redisClient.GetCacheJSON(ctx, key, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal cache data for key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClose(t *testing.T) {
	db, _ := redismock.NewClientMock()
	defer db.Close()

	redisClient := &RedisClient{client: db}

	err := redisClient.Close()
	assert.NoError(t, err)
}
