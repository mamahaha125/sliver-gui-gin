package session

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// 使用 Cookie 保存 session
func EnableCookieSession(key string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(key))
	store.Options(sessions.Options{Path: "/", MaxAge: 24 * 3600})
	return sessions.Sessions("_SESSION", store)
}

// 使用 Redis 保存 session
func EnableRedisSession(key string) gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte(key))
	store.Options(sessions.Options{Path: "/", MaxAge: 6 * 3600})
	return sessions.Sessions("_SESSION", store)
}

// 使用 内存 保存 session
func EnableMemorySession(key string) gin.HandlerFunc {
	store := memstore.NewStore([]byte(key))
	store.Options(sessions.Options{Path: "/", MaxAge: 6 * 3600})
	return sessions.Sessions("_SESSION", store)
}

func Set(c *gin.Context, key string, value interface{}) (err error) {
	s := sessions.Default(c)
	s.Set(key, value)
	err = s.Save()
	if err != nil {
		return err
	}
	return nil
}

func Get(c *gin.Context, key string) interface{} {
	s := sessions.Default(c)
	return s.Get(key)
}

func Del(c *gin.Context, key string) error {
	s := sessions.Default(c)
	s.Delete(key)
	return s.Save()
}

func GetSessionId(c *gin.Context) (int64, error) {
	s := sessions.Default(c)
	auth, ok := s.Get("uid").(uint)
	if !ok {
		return 0, errors.New("无用户session")
	}
	return int64(auth), nil
}
