package cache

import (
	"io"
	"net/http"
)

type CachedHttpClient struct {
	client *http.Client
	cache  *HttpCache
}

func NewCachedHttpClient() *CachedHttpClient {
	return &CachedHttpClient{
		client: http.DefaultClient,
		cache:  GetHttpCacheInstance(),
	}
}

func (c *CachedHttpClient) GetWithCache(url string) ([]byte, error) {
	if value, exists := c.cache.Get(url); exists {
		return value, nil
	}

	res, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	c.cache.Add(url, body)

	return body, nil
}
