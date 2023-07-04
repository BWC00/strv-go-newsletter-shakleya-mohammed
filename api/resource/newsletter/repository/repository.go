package repository

import (
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/api/resource/newsletter"
)

type Repository struct {
	postgresDB *gorm.DB
}

func NewRepository(postgresDB *gorm.DB) *Repository {
	return &Repository{
		postgresDB: postgresDB,
	}
}

func (r *Repository) ListNewsletters(userId uint32) (newsletter.Newsletters, error) {
	newsletters := make([]*newsletter.Newsletter, 0)

	if err := r.postgresDB.Find(&newsletters, "editor_id = ?", userId).Error; err != nil {
		return nil, err
	}

	if len(newsletters) == 0 {
		return nil, nil
	}

	return newsletters, nil
}

func (r *Repository) CreateNewsletter(newsletter *newsletter.Newsletter) (*newsletter.Newsletter, error) {
	if err := r.postgresDB.Create(newsletter).Error; err != nil {
		return nil, err
	}

	return newsletter, nil
}

func (r *Repository) ReadNewsletter(newsletterId uint32, userId uint32) (*newsletter.Newsletter, error) {
	newsletter := newsletter.Newsletter{}
	if err := r.postgresDB.Where("id = ? AND editor_id = ?", newsletterId, userId).First(&newsletter).Error; err != nil {
		return nil, err
	}

	return &newsletter, nil
}

func (r *Repository) UpdateNewsletter(newsletterInfo *newsletter.Newsletter) error {
	if _, err := r.ReadNewsletter(newsletterInfo.ID, newsletterInfo.EditorId); err != nil {
		return err
	}
	if err := r.postgresDB.Model(&newsletter.Newsletter{}).Where("id = ? AND editor_id = ?", newsletterInfo.ID, newsletterInfo.EditorId).Updates(&newsletter.Newsletter{Name: newsletterInfo.Name, Description: newsletterInfo.Description}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteNewsletter(newsletterId uint32, userId uint32) error {
	if err := r.postgresDB.Where("id = ? AND editor_id = ?", newsletterId, userId).Delete(&newsletter.Newsletter{}).Error; err != nil {
		return err
	}

	return nil
}
