package repository

import (
	"context"
	"errors"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/dto"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TransactionRepository interface {
		CreateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) (entity.Transaction, error)
		GetProductByID(ctx context.Context, tx *gorm.DB, productId uuid.UUID) (entity.Product, error)
		UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) error
		GetTransactionByReference(ctx context.Context, tx *gorm.DB, reference string) (entity.Transaction, error)
	}

	transactionRepository struct {
		db *gorm.DB
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) (entity.Transaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetProductByID(ctx context.Context, tx *gorm.DB, productId uuid.UUID) (entity.Product, error) {
	if tx == nil {
		tx = r.db
	}

	var product entity.Product
	if err := tx.WithContext(ctx).Where("id = ?", productId).First(&product).Error; err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, tx *gorm.DB, transactionId uuid.UUID) (entity.Transaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.Transaction
	if err := tx.WithContext(ctx).Where("id = ?", transactionId).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Transaction{}, dto.ErrTransactionNotFound
		}
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Updates(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrTransactionNotFound
		}
		return err
	}

	return nil
}

func (r *transactionRepository) DeleteTransaction(ctx context.Context, tx *gorm.DB, transactionId uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.Transaction
	if err := tx.WithContext(ctx).Where("id = ?", transactionId).Delete(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrTransactionNotFound
		}
		return err
	}

	return nil
}

func (r *transactionRepository) GetTransactionByReference(ctx context.Context, tx *gorm.DB, reference string) (entity.Transaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.Transaction
	if err := tx.WithContext(ctx).Where("reference = ?", reference).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Transaction{}, dto.ErrTransactionNotFound
		}
		return entity.Transaction{}, err
	}

	return transaction, nil
}
