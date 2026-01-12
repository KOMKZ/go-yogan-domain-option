package option

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-option/model"
)

// Repository 配置仓储接口
type Repository interface {
	FindAll(ctx context.Context) ([]model.Option, error)
	FindByID(ctx context.Context, id uint) (*model.Option, error)
	FindByKey(ctx context.Context, key string) (*model.Option, error)
	FindByGroupType(ctx context.Context, groupType string) ([]model.Option, error)
	ExistsByKey(ctx context.Context, key string) (bool, error)
	Create(ctx context.Context, option *model.Option) error
	Update(ctx context.Context, option *model.Option) error
	DeleteByID(ctx context.Context, id uint) error
	DeleteByKey(ctx context.Context, key string) error
}
