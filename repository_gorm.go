package option

import (
	"context"
	"errors"

	"github.com/KOMKZ/go-yogan-domain-option/model"
	"gorm.io/gorm"
)

// GORMRepository GORM 配置仓储实现
type GORMRepository struct {
	db *gorm.DB
}

// NewGORMRepository 创建配置仓储
func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{db: db}
}

// FindAll 获取所有配置
func (r *GORMRepository) FindAll(ctx context.Context) ([]model.Option, error) {
	var options []model.Option
	err := r.db.WithContext(ctx).Order("group_type ASC, `key` ASC").Find(&options).Error
	return options, err
}

// FindByID 根据ID获取配置
func (r *GORMRepository) FindByID(ctx context.Context, id uint) (*model.Option, error) {
	var option model.Option
	err := r.db.WithContext(ctx).First(&option, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &option, nil
}

// FindByKey 根据Key获取配置
func (r *GORMRepository) FindByKey(ctx context.Context, key string) (*model.Option, error) {
	var option model.Option
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&option).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &option, nil
}

// FindByGroupType 根据分组类型获取配置列表
func (r *GORMRepository) FindByGroupType(ctx context.Context, groupType string) ([]model.Option, error) {
	var options []model.Option
	err := r.db.WithContext(ctx).Where("group_type = ?", groupType).Order("`key` ASC").Find(&options).Error
	return options, err
}

// ExistsByKey 检查Key是否存在
func (r *GORMRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Option{}).Where("`key` = ?", key).Count(&count).Error
	return count > 0, err
}

// Create 创建配置
func (r *GORMRepository) Create(ctx context.Context, option *model.Option) error {
	return r.db.WithContext(ctx).Create(option).Error
}

// Update 更新配置
func (r *GORMRepository) Update(ctx context.Context, option *model.Option) error {
	return r.db.WithContext(ctx).Save(option).Error
}

// DeleteByID 根据ID删除配置
func (r *GORMRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Option{}, id).Error
}

// DeleteByKey 根据Key删除配置
func (r *GORMRepository) DeleteByKey(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("`key` = ?", key).Delete(&model.Option{}).Error
}
