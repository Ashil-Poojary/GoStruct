package repository

import (
	"github.com/ashil-poojary/gostruct/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Default_DB *gorm.DB
	Replica_DB *gorm.DB
}

func NewOrderRepository(defaultDb, replicaDb *gorm.DB) *OrderRepository {
	return &OrderRepository{Default_DB: defaultDb, Replica_DB: replicaDb}
}

// --- Order Operations ---

func (r *OrderRepository) GetAll() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.Replica_DB.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetByID(id int64) (*models.Order, error) {
	var order models.Order
	if err := r.Default_DB.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.Default_DB.Create(order).Error
}

func (r *OrderRepository) UpdateStatus(id int64, status string) error {
	return r.Default_DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// --- OrderReturn Operations ---

func (r *OrderRepository) CreateReturn(returnOrder *models.OrderReturn) error {
	return r.Default_DB.Create(returnOrder).Error
}

func (r *OrderRepository) GetAllReturns() ([]*models.OrderReturn, error) {
	var returns []*models.OrderReturn
	if err := r.Replica_DB.Find(&returns).Error; err != nil {
		return nil, err
	}
	return returns, nil
}

func (r *OrderRepository) GetReturnByID(id int64) (*models.OrderReturn, error) {
	var ret models.OrderReturn
	if err := r.Replica_DB.First(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (r *OrderRepository) UpdateReturn(returnOrder *models.OrderReturn) error {
	return r.Default_DB.Save(returnOrder).Error
}

func (r *OrderRepository) DeleteReturn(id int64) error {
	return r.Default_DB.Delete(&models.OrderReturn{}, id).Error
}
