package posts

import (
	"context"
	"fmt"
	"sort"

	"prp.com/sparkly/internal/pkg"
)

func (service *service) incrementPostsCache(ctx context.Context, postID string, increment float64) {

	for key := range pkg.AnalyticsTimeFrames {
		service.incrementPostsCacheByDuration(ctx, key, postID, increment)
	}

}

func (service *service) incrementPostsCacheByDuration(ctx context.Context, key, postID string, increment float64) {

	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, key)
	service.cacheService.ZIncrBy(ctx, cacheKey, postID, increment)
	service.cacheService.Expire(ctx, cacheKey, 2*pkg.AnalyticsTimeFrames[key])

}

func (service *service) getPostsCache(ctx context.Context, key string, limit int) ([]pkg.PopularPost, error) {

	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, key)
	data, err := service.cacheService.ZRevRangeWithScores(ctx, cacheKey, int64(limit))
	if err != nil {
		return nil, err
	}

	posts := make([]pkg.PopularPost, 0)

	for key, value := range data {
		posts = append(posts, pkg.PopularPost{
			PostID: key,
			Count:  int64(value),
		})
	}

	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].Count > posts[j].Count
	})

	return posts, nil

}

func (service *service) BackfillPolularPosts(ctx context.Context) {

	for key, duration := range pkg.AnalyticsTimeFrames {

		posts, err := service.getPopularPostsFromDB(ctx, duration, 10)
		if err != nil {
			continue
		}

		service.cacheService.DeleteKeys(ctx, fmt.Sprintf(pkg.CachekeyPopularPosts, key))

		for _, post := range posts {
			service.incrementPostsCacheByDuration(ctx, key, post.PostID, float64(post.Count))
		}

	}

}
