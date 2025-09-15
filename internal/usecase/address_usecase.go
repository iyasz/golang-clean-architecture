package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iyasz/golang-clean-architecture/internal/entity"
	"github.com/iyasz/golang-clean-architecture/internal/helper"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/iyasz/golang-clean-architecture/internal/model/converter"
	"github.com/iyasz/golang-clean-architecture/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AddressUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	AddressRepository *repository.AddressRepository
	ContactRepository *repository.ContactRepository
}

func NewAddressUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, addressRepository *repository.AddressRepository, 
	contactRepository *repository.ContactRepository) *AddressUseCase { 
		return &AddressUseCase{
			DB: db,
			Log: log,
			Validate: validate,
			AddressRepository: addressRepository,
			ContactRepository: contactRepository,
		}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, helper.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, helper.ErrNotFound
	}

	address := &entity.Address{
		ID:         uuid.NewString(),
		ContactId:  contact.ID,
		Street:     request.Street,
		City:       request.City,
		Province:   request.Province,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return nil, helper.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, helper.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Update(ctx context.Context, request *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, helper.ErrBadRequest
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, helper.ErrNotFound
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, contact.ID); err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, helper.ErrNotFound
	}

	address.Street = request.Street
	address.City = request.City
	address.Province = request.Province
	address.PostalCode = request.PostalCode
	address.Country = request.Country

	if err := c.AddressRepository.Update(tx, address); err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return nil, helper.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, helper.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}