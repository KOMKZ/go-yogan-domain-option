package option

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KOMKZ/go-yogan-domain-option/model"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

// Service 配置服务
type Service struct {
	repo   Repository
	logger *logger.CtxZapLogger
}

// NewService 创建配置服务
func NewService(repo Repository, log *logger.CtxZapLogger) *Service {
	return &Service{
		repo:   repo,
		logger: log,
	}
}

// GetAll 获取所有配置
func (s *Service) GetAll(ctx context.Context) ([]model.Option, error) {
	options, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.ErrorCtx(ctx, "获取所有配置失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	return options, nil
}

// GetByID 根据ID获取配置
func (s *Service) GetByID(ctx context.Context, id uint) (*model.Option, error) {
	option, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorCtx(ctx, "根据ID获取配置失败", zap.Uint("id", id), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if option == nil {
		return nil, ErrNotFound.WithMsgf("配置不存在: id=%d", id)
	}
	return option, nil
}

// GetByKey 根据Key获取配置
func (s *Service) GetByKey(ctx context.Context, key string) (*model.Option, error) {
	option, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		s.logger.ErrorCtx(ctx, "根据Key获取配置失败", zap.String("key", key), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if option == nil {
		return nil, ErrNotFound.WithMsgf("配置不存在: key=%s", key)
	}
	return option, nil
}

// GetValue 获取配置值
func (s *Service) GetValue(ctx context.Context, key string) (string, error) {
	option, err := s.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}
	return option.Value, nil
}

// GetByGroupType 根据分组类型获取配置列表
func (s *Service) GetByGroupType(ctx context.Context, groupType string) ([]model.Option, error) {
	options, err := s.repo.FindByGroupType(ctx, groupType)
	if err != nil {
		s.logger.ErrorCtx(ctx, "根据分组类型获取配置失败", zap.String("groupType", groupType), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	return options, nil
}

// CreateInput 创建配置输入
type CreateInput struct {
	Key             string
	Value           string
	GroupType       string
	Component       string
	ComponentParams string
}

// Create 创建配置
func (s *Service) Create(ctx context.Context, input *CreateInput) (*model.Option, error) {
	// 检查Key是否已存在
	exists, err := s.repo.ExistsByKey(ctx, input.Key)
	if err != nil {
		s.logger.ErrorCtx(ctx, "检查Key是否存在失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if exists {
		return nil, ErrKeyExists.WithMsgf("配置键已存在: key=%s", input.Key)
	}

	// 解析 componentParams
	var componentParams model.JSONValue
	if input.ComponentParams != "" {
		if err := parseJSON(input.ComponentParams, &componentParams); err != nil {
			s.logger.WarnCtx(ctx, "解析componentParams失败", zap.String("key", input.Key), zap.Error(err))
			// 不报错，使用空值
		}
	}

	option := &model.Option{
		Key:             input.Key,
		Value:           input.Value,
		GroupType:       input.GroupType,
		Component:       input.Component,
		ComponentParams: componentParams,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, option); err != nil {
		s.logger.ErrorCtx(ctx, "创建配置失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "创建配置成功", zap.String("key", option.Key), zap.Uint("id", option.ID))
	return option, nil
}

// UpdateInput 更新配置输入
type UpdateInput struct {
	Value           *string
	GroupType       *string
	Component       *string
	ComponentParams *string
}

// Update 更新配置
func (s *Service) Update(ctx context.Context, id uint, input *UpdateInput) (*model.Option, error) {
	option, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Value != nil {
		option.Value = *input.Value
	}
	if input.GroupType != nil {
		option.GroupType = *input.GroupType
	}
	if input.Component != nil {
		option.Component = *input.Component
	}
	if input.ComponentParams != nil {
		var componentParams model.JSONValue
		if *input.ComponentParams != "" {
			if err := parseJSON(*input.ComponentParams, &componentParams); err != nil {
				s.logger.WarnCtx(ctx, "解析componentParams失败", zap.String("key", option.Key), zap.Error(err))
			}
		}
		option.ComponentParams = componentParams
	}
	option.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, option); err != nil {
		s.logger.ErrorCtx(ctx, "更新配置失败", zap.Uint("id", id), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "更新配置成功", zap.String("key", option.Key), zap.Uint("id", option.ID))
	return option, nil
}

// UpdateByKey 根据Key更新配置值
func (s *Service) UpdateByKey(ctx context.Context, key string, value string) (*model.Option, error) {
	option, err := s.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	option.Value = value
	option.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, option); err != nil {
		s.logger.ErrorCtx(ctx, "更新配置值失败", zap.String("key", key), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "更新配置值成功", zap.String("key", key))
	return option, nil
}

// Delete 删除配置
func (s *Service) Delete(ctx context.Context, id uint) error {
	option, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		s.logger.ErrorCtx(ctx, "删除配置失败", zap.Uint("id", id), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "删除配置成功", zap.String("key", option.Key), zap.Uint("id", id))
	return nil
}

// DeleteByKey 根据Key删除配置
func (s *Service) DeleteByKey(ctx context.Context, key string) error {
	option, err := s.GetByKey(ctx, key)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteByKey(ctx, key); err != nil {
		s.logger.ErrorCtx(ctx, "删除配置失败", zap.String("key", key), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "删除配置成功", zap.String("key", key), zap.Uint("id", option.ID))
	return nil
}

// parseJSON 解析JSON字符串
func parseJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}
