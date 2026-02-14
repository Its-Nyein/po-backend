package repositories

import (
	"po-backend/models"
	"time"

	"gorm.io/gorm"
)

type StoryRepository struct {
	DB *gorm.DB
}

func NewStoryRepository(db *gorm.DB) *StoryRepository {
	return &StoryRepository{DB: db}
}

func (r *StoryRepository) Create(story *models.Story) error {
	return r.DB.Create(story).Error
}

func (r *StoryRepository) Delete(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("story_id = ?", id).Delete(&models.StoryView{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Story{}, id).Error
	})
}

func (r *StoryRepository) GetByID(id uint) (*models.Story, error) {
	var story models.Story
	err := r.DB.Preload("User").Preload("Views.Viewer").First(&story, id).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func (r *StoryRepository) IsOwner(storyID, userID uint) bool {
	var story models.Story
	err := r.DB.Where("id = ? AND user_id = ?", storyID, userID).First(&story).Error
	return err == nil
}

func (r *StoryRepository) GetFeedStories(viewerID uint, followingIDs []uint, mutualIDs []uint) ([]models.Story, error) {
	var stories []models.Story
	now := time.Now()

	query := r.DB.Preload("User").Preload("Views").
		Where("expires_at > ?", now).
		Order("created_at DESC")

	if len(mutualIDs) > 0 {
		query = query.Where(
			"user_id = ? OR privacy = 'public' OR (privacy = 'friends' AND user_id IN ?)",
			viewerID, mutualIDs,
		)
	} else {
		query = query.Where("user_id = ? OR privacy = 'public'", viewerID)
	}

	err := query.Find(&stories).Error
	return stories, err
}

func (r *StoryRepository) GetUserStories(targetUserID, viewerID uint, mutualIDs []uint) ([]models.Story, error) {
	var stories []models.Story
	now := time.Now()

	query := r.DB.Preload("User").Preload("Views").
		Where("user_id = ? AND expires_at > ?", targetUserID, now).
		Order("created_at ASC")

	if targetUserID != viewerID {
		isMutual := false
		for _, id := range mutualIDs {
			if id == targetUserID {
				isMutual = true
				break
			}
		}
		if isMutual {
			query = query.Where("privacy IN ?", []string{"public", "friends"})
		} else {
			query = query.Where("privacy = 'public'")
		}
	}

	err := query.Find(&stories).Error
	return stories, err
}

func (r *StoryRepository) CreateView(storyID, viewerID uint) (*models.StoryView, error) {
	view := &models.StoryView{StoryID: storyID, ViewerID: viewerID}
	result := r.DB.Where("story_id = ? AND viewer_id = ?", storyID, viewerID).FirstOrCreate(view)
	return view, result.Error
}

func (r *StoryRepository) GetViewers(storyID uint) ([]models.StoryView, error) {
	var views []models.StoryView
	err := r.DB.Preload("Viewer").Where("story_id = ?", storyID).
		Order("created_at DESC").Find(&views).Error
	return views, err
}

func (r *StoryRepository) DeleteExpired() error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("story_id IN (?)",
			tx.Model(&models.Story{}).Select("id").Where("expires_at < ?", time.Now()),
		).Delete(&models.StoryView{}).Error; err != nil {
			return err
		}
		return tx.Where("expires_at < ?", time.Now()).Delete(&models.Story{}).Error
	})
}
