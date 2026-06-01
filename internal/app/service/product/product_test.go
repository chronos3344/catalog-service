package mproduct

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository/mocks"
	"github.com/chronos3344/catalog-service/internal/pkg/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type createProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}
	type want struct {
		err error
	}

	categoryGUID := uuid.New()
	createError := errors.New("create failed")
	dbError := errors.New("database error")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{GUID: categoryGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, &args.req.CategoryGUID).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.Anything).
					Return(nil).
					Once()
			},
		},
		{
			name: "create error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: createError},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{GUID: categoryGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, &args.req.CategoryGUID).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.Anything).
					Return(createError). // ← ОБЯЗАТЕЛЬНО
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing Product",
					Price:        500,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{GUID: categoryGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, &args.req.CategoryGUID).
					Return([]entity.Product{{Name: "Existing Product"}}, nil).
					Once()

				// Create НЕ настраиваем
			},
		},
		{
			name: "create - list error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: dbError},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{GUID: categoryGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, &args.req.CategoryGUID).
					Return(nil, dbError).
					Once()

				// Create НЕ настраиваем
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{}, entity.ErrNotFound).
					Once()

				// List и Create НЕ настраиваем
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.Error(err)
				s.ErrorIs(err, tc.want.err)
				// При ошибке не проверяем GUID, так как сервис может возвращать заполненный продукт
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

type getByGUIDProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *getByGUIDProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDProductSuite))
}

func (s *getByGUIDProductSuite) TestGetByGUID() {
	productGUID := uuid.New()
	expectedProduct := entity.Product{
		GUID:         productGUID,
		Name:         "Test Product",
		Description:  testutil.PtrString("A test product description"),
		Price:        1000,
		CategoryGUID: uuid.New(),
	}

	testCases := []struct {
		name    string
		guid    uuid.UUID
		wantErr error
		prepare func(guid uuid.UUID)
	}{
		{
			name:    "success",
			guid:    productGUID,
			wantErr: nil,
			prepare: func(guid uuid.UUID) {
				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(expectedProduct, nil).
					Maybe()
			},
		},
		{
			name:    "not found",
			guid:    uuid.New(),
			wantErr: entity.ErrNotFound,
			prepare: func(guid uuid.UUID) {
				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(entity.Product{}, entity.ErrNotFound).
					Maybe()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.guid)

			result, err := s.srv.Get(s.ctx, tc.guid)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.Equal(expectedProduct.GUID, result.GUID)
				s.Equal(expectedProduct.Name, result.Name)
				s.Equal(expectedProduct.Description, result.Description)
				s.Equal(expectedProduct.Price, result.Price)
				s.Equal(expectedProduct.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

type deleteProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {
	productGUID := uuid.New()
	existingProduct := entity.Product{
		GUID: productGUID,
		Name: "Test Product",
	}

	dbError := errors.New("database error")

	testCases := []struct {
		name    string
		guid    uuid.UUID
		wantErr error
		prepare func(guid uuid.UUID)
	}{
		{
			name:    "success",
			guid:    productGUID,
			wantErr: nil,
			prepare: func(guid uuid.UUID) {
				// Настраиваем InsideTx так, чтобы он выполнил функцию и вернул nil
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						// Выполняем функцию транзакции
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, guid).
					Return(nil).
					Once()
			},
		},
		{
			name:    "not found",
			guid:    uuid.New(),
			wantErr: entity.ErrNotFound,
			prepare: func(guid uuid.UUID) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(entity.Product{}, entity.ErrNotFound).
					Once()

				// Delete не должен вызываться
			},
		},
		{
			name:    "delete error",
			guid:    productGUID,
			wantErr: dbError,
			prepare: func(guid uuid.UUID) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						// Выполняем функцию и возвращаем ошибку, которую вернет Delete
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, guid).
					Return(dbError).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.guid)

			err := s.srv.Delete(s.ctx, tc.guid)

			if tc.wantErr != nil {
				s.Error(err, "Expected error but got nil")
				if tc.name == "delete error" {
					s.Equal(dbError, err, "Expected database error")
				}
				if tc.wantErr == entity.ErrNotFound {
					s.ErrorIs(err, entity.ErrNotFound)
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

type listProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func (s *listProductSuite) TestList() {
	expectedProducts := []entity.Product{
		{
			GUID:         uuid.New(),
			Name:         "Product 1",
			Description:  testutil.PtrString("Description 1"),
			Price:        1000,
			CategoryGUID: uuid.New(),
		},
		{
			GUID:         uuid.New(),
			Name:         "Product 2",
			Description:  testutil.PtrString("Description 2"),
			Price:        2000,
			CategoryGUID: uuid.New(),
		},
	}

	testCases := []struct {
		name         string
		wantErr      error
		wantLen      int
		wantProducts []entity.Product
		prepare      func()
	}{
		{
			name:         "success",
			wantErr:      nil,
			wantLen:      2,
			wantProducts: expectedProducts,
			prepare: func() {
				// Сервис вызывает List с nil, nil
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), (*uuid.UUID)(nil)).
					Return(expectedProducts, nil).
					Once()
			},
		},
		{
			name:         "empty result",
			wantErr:      nil,
			wantLen:      0,
			wantProducts: []entity.Product{},
			prepare: func() {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
		{
			name:    "repository error",
			wantErr: errors.New("database error"),
			wantLen: 0,
			prepare: func() {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), (*uuid.UUID)(nil)).
					Return(nil, errors.New("database error")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare()

			result, err := s.srv.List(s.ctx)

			if tc.wantErr != nil {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Len(result, tc.wantLen)

				for i, product := range result {
					s.Equal(tc.wantProducts[i].Name, product.Name)
					s.Equal(tc.wantProducts[i].Price, product.Price)
				}
			}
		})
	}
}

type updateProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	productGUID := uuid.New()
	categoryGUID := uuid.New()
	newCategoryGUID := uuid.New()

	dbError := errors.New("database error")

	existingProduct := entity.Product{
		GUID:         productGUID,
		Name:         "Old Name",
		Description:  testutil.PtrString("Old Description"),
		Price:        1000,
		CategoryGUID: categoryGUID,
	}

	testCases := []struct {
		name    string
		guid    uuid.UUID
		req     entity.RequestProductUpdate
		wantErr error
		prepare func(guid uuid.UUID, req entity.RequestProductUpdate)
	}{
		{
			name: "full update",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name:         testutil.PtrString("New Name"),
				Description:  testutil.PtrString("New Description"),
				Price:        testutil.PtrFloat64(2000),
				CategoryGUID: &newCategoryGUID,
			},
			wantErr: nil,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, newCategoryGUID).
					Return(entity.Category{GUID: newCategoryGUID}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, mock.Anything, mock.Anything).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.Anything).
					Return(nil).
					Once()
			},
		},
		{
			name: "partial update - name only",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name: testutil.PtrString("New Name Only"),
			},
			wantErr: nil,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, mock.Anything, mock.Anything).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.Anything).
					Return(nil).
					Once()
			},
		},

		{
			name: "update - list error",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name: testutil.PtrString("New Name"),
			},
			wantErr: dbError,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						mock.MatchedBy(func(name *string) bool {
							return name != nil && *name == "New Name"
						}),
						mock.MatchedBy(func(catGUID *uuid.UUID) bool {
							return catGUID != nil && *catGUID == categoryGUID
						}),
					).
					Return(nil, dbError). // Возвращаем ошибку базы данных
					Once()

				// Update НЕ должен вызываться
				s.productRepo.EXPECT().
					Update(s.ctx, mock.Anything).
					Times(0)
			},
		},

		{
			name: "partial update - price only",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Price: testutil.PtrFloat64(3000),
			},
			wantErr: nil,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				// ВАЖНО: при обновлении только цены, проверка уникальности все равно выполняется
				// и должна вернуть список с текущим продуктом
				s.productRepo.EXPECT().
					List(s.ctx, mock.Anything, mock.Anything).
					Return([]entity.Product{existingProduct}, nil). // Возвращаем текущий продукт
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.Anything).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			guid: uuid.New(),
			req: entity.RequestProductUpdate{
				Name: testutil.PtrString("New Name"),
			},
			wantErr: entity.ErrNotFound,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(entity.Product{}, entity.ErrNotFound).
					Once()

			},
		},
		{
			name: "duplicate name",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name: testutil.PtrString("Duplicate Name"),
			},
			wantErr: entity.ErrAlreadyExists,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				// Возвращаем дубликат с другим GUID
				s.productRepo.EXPECT().
					List(s.ctx, mock.Anything, mock.Anything).
					Return([]entity.Product{{
						GUID:         uuid.New(),
						Name:         "Duplicate Name",
						CategoryGUID: categoryGUID,
					}}, nil).
					Once()

			},
		},
		{
			name: "category not found",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name:         testutil.PtrString("New Name"),
				CategoryGUID: &newCategoryGUID,
			},
			wantErr: entity.ErrNotFound,
			prepare: func(guid uuid.UUID, req entity.RequestProductUpdate) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(existingProduct, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, newCategoryGUID).
					Return(entity.Category{}, entity.ErrNotFound).
					Once()

			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.productRepo.ExpectedCalls = nil
			s.productRepo.Calls = nil
			s.categoryRepo.ExpectedCalls = nil
			s.categoryRepo.Calls = nil

			tc.prepare(tc.guid, tc.req)

			result, err := s.srv.Update(s.ctx, tc.guid, tc.req)

			if tc.wantErr != nil {
				s.Error(err)
				s.ErrorIs(err, tc.wantErr)
				s.Equal(uuid.Nil, result.GUID)
			} else {
				s.NoError(err)
				s.NotEqual(uuid.Nil, result.GUID)
				s.Equal(tc.guid, result.GUID)
			}
		})
	}

}

func TestNewService(t *testing.T) {
	productRepo := mocks.NewMockProduct(t)
	categoryRepo := mocks.NewMockCategory(t)

	service := NewService(productRepo, categoryRepo)

	assert.NotNil(t, service)

	// Проверяем, что возвращается правильный тип
	srvImpl, ok := service.(*srv)
	assert.True(t, ok)
	assert.Equal(t, productRepo, srvImpl.repoProduct)
	assert.Equal(t, categoryRepo, srvImpl.repoCategory)
}
