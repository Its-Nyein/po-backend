package repositories

import (
	"po-backend/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) GetAll(limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Preload("QuotedPost.User").
		Order("created_at DESC").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	r.PopulateQuoteCounts(posts)
	return posts, nil
}

func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Preload("QuotedPost.User").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	r.PopulateQuoteCounts([]models.Post{post})
	return &post, nil
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) GetByUserIDs(userIDs []uint, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Likes").
		Preload("Likes").
		Preload("QuotedPost.User").
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	r.PopulateQuoteCounts(posts)
	return posts, nil
}

func (r *PostRepository) Update(id uint, content string) (*models.Post, error) {
	if err := r.DB.Model(&models.Post{}).Where("id = ?", id).Update("content", content).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *PostRepository) IsOwner(postID, userID uint) bool {
	var post models.Post
	err := r.DB.Where("id = ? AND user_id = ?", postID, userID).First(&post).Error
	return err == nil
}

func (r *PostRepository) PopulateQuoteCounts(posts []models.Post) {
	if len(posts) == 0 {
		return
	}
	ids := make([]uint, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}

	type result struct {
		QuotedPostID uint
		Count        int
	}
	var results []result
	r.DB.Model(&models.Post{}).
		Select("quoted_post_id, count(*) as count").
		Where("quoted_post_id IN ?", ids).
		Group("quoted_post_id").
		Scan(&results)

	counts := make(map[uint]int, len(results))
	for _, r := range results {
		counts[r.QuotedPostID] = r.Count
	}
	for i := range posts {
		posts[i].QuoteCount = counts[posts[i].ID]
	}
}
