package repositories

import (
	"log"
	"po-backend/models"

	"gorm.io/gorm"
)

type HashtagRepository struct {
	DB *gorm.DB
}

func NewHashtagRepository(db *gorm.DB) *HashtagRepository {
	return &HashtagRepository{DB: db}
}

func (r *HashtagRepository) FindOrCreate(name string) (*models.Hashtag, error) {
	var hashtag models.Hashtag
	err := r.DB.Where("name = ?", name).First(&hashtag).Error
	if err == nil {
		return &hashtag, nil
	}
	hashtag = models.Hashtag{Name: name}
	err = r.DB.Create(&hashtag).Error
	return &hashtag, err
}

func (r *HashtagRepository) SyncPostHashtags(postID uint, tagNames []string) error {
	r.DB.Where("post_id = ?", postID).Delete(&models.PostHashtag{})

	for _, name := range tagNames {
		hashtag, err := r.FindOrCreate(name)
		if err != nil {
			return err
		}
		ph := models.PostHashtag{PostID: postID, HashtagID: hashtag.ID}
		if err := r.DB.Create(&ph).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *HashtagRepository) GetPostsByHashtag(tag string, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Preload("Reposts").
		Joins("JOIN post_hashtags ON post_hashtags.post_id = posts.id").
		Joins("JOIN hashtags ON hashtags.id = post_hashtags.hashtag_id").
		Where("hashtags.name = ?", tag).
		Order("posts.created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *HashtagRepository) GetTrending(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	rows, err := r.DB.Table("post_hashtags").
		Select("hashtags.id, hashtags.name, COUNT(post_hashtags.post_id) as post_count").
		Joins("JOIN hashtags ON hashtags.id = post_hashtags.hashtag_id").
		Group("hashtags.id, hashtags.name").
		Order("post_count DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	for rows.Next() {
		var id uint
		var name string
		var postCount int64
		if err := rows.Scan(&id, &name, &postCount); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":        id,
			"name":      name,
			"postCount": postCount,
		})
	}
	return results, nil
}
