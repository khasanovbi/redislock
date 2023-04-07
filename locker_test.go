package redislock

import (
	"context"
	"testing"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var redisOptions = &redis.Options{
	Addr: "localhost:6380",
}

type RedisTestSuite struct {
	suite.Suite
	client redis.UniversalClient
}

func (suite *RedisTestSuite) requireFlush() {
	err := suite.client.FlushAll(context.Background()).Err()
	suite.Require().NoError(err)
}

func (suite *RedisTestSuite) SetupSuite() {
	if !wait.New().Do([]string{redisOptions.Addr}) {
		panic("can't await redis port")
	}

	suite.client = redis.NewClient(redisOptions)
}

func (suite *RedisTestSuite) TearDownSuite() {
	suite.requireFlush()
}

func (suite *RedisTestSuite) SetupTest() {
	suite.requireFlush()
}

func (suite *RedisTestSuite) TestSuccess() {
	locker, err := New(suite.client, WithTTL(time.Hour), WithRefreshPeriod(time.Millisecond))
	suite.Require().NoError(err)

	lock, _, err := locker.Lock(context.Background(), "key")
	suite.Require().NoError(err)

	err = lock.Unlock()
	suite.Require().NoError(err)
}

func (suite *RedisTestSuite) TestAlreadyLocked() {
	key := "key"
	locker, err := New(suite.client, WithTTL(time.Hour), WithRefreshPeriod(time.Millisecond))
	suite.Require().NoError(err)

	_, _, err = locker.Lock(context.Background(), key)
	suite.Require().NoError(err)

	_, _, err = locker.Lock(context.Background(), key)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, ErrAlreadyLocked)
}

func (suite *RedisTestSuite) TestUnlockError() {
	key := "key"

	locker, err := New(suite.client, WithTTL(time.Hour), WithRefreshPeriod(time.Millisecond))
	suite.Require().NoError(err)

	// NOTE: Force redislock not to set deadline to context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	lock, ctx, err := locker.Lock(ctx, key)
	suite.Require().NoError(err)

	err = lock.Unlock()
	suite.Require().NoError(err)

	suite.Require().ErrorIs(ctx.Err(), context.Canceled)

	err = lock.Unlock()
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, redislock.ErrLockNotHeld)
}

func TestRedisTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RedisTestSuite))
}

func TestNew_InvalidOptions(t *testing.T) {
	t.Parallel()

	_, err := New(nil, WithTTL(time.Second), WithRefreshPeriod(time.Hour))
	require.Error(t, err)
	require.ErrorIs(t, err, errSmallTTL)
}

func TestLocker_Lock_InvalidOptions(t *testing.T) {
	t.Parallel()

	locker, err := New(nil)
	require.NoError(t, err)

	_, _, err = locker.Lock(context.Background(), "key", WithTTL(time.Second), WithRefreshPeriod(time.Hour))
	require.Error(t, err)
	require.ErrorIs(t, err, errSmallTTL)
}
