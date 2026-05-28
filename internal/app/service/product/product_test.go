package mproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository/mocks"
	"github.com/chronos3344/catalog-service/internal/pkg/testutil"
	"github.com/gofrs/uuid"
	uuid2 "github.com/google/uuid"
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

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}
	type want struct {
		err error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

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
					CategoryGUID: uuid2.UUID(categoryGUID),
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

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{GUID: uuid2.UUID(categoryGUID)}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing Product",
					Price:        500,
					CategoryGUID: uuid2.UUID(categoryGUID),
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

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{{Name: "Existing Product"}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGUID: uuid2.UUID(categoryGUID),
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

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, args.req.CategoryGUID).
					Return(entity.Category{}, entity.ErrNotFound).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Empty(result.GUID)
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
	// TODO:
	// Метод GetByGUID — простая обёртка над репозиторием.
	// Создайте табличные тесты с двумя кейсами:
	// 1. "success" — репозиторий вернул продукт, ошибки нет.
	// 2. "not found" — репозиторий вернул entity.ErrNotFound.
	//
	// InsideTx здесь не используется — метод вызывает
	// репозиторий напрямую. Мок нужен только для productRepo.GetByGUID.
	productGUID := uuid.Must(uuid.NewV4())
	expectedProduct := entity.Product{
		GUID:         uuid2.UUID(productGUID),
		Name:         "Test Product",
		Description:  testutil.PtrString("A test product description"),
		Price:        1000,
		CategoryGUID: uuid2.UUID(uuid.Must(uuid.NewV4())),
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
					Once()
			},
		},
		{
			name:    "not found",
			guid:    uuid.Must(uuid.NewV4()),
			wantErr: entity.ErrNotFound,
			prepare: func(guid uuid.UUID) {
				s.productRepo.EXPECT().
					GetByGUID(s.ctx, guid).
					Return(entity.Product{}, entity.ErrNotFound).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.guid)

			result, err := s.srv.Get(s.ctx, uuid2.UUID(tc.guid))

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
	productGUID := uuid.Must(uuid.NewV4())
	existingProduct := entity.Product{
		GUID: uuid2.UUID(productGUID),
		Name: "Test Product",
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
					Delete(s.ctx, guid).
					Return(nil).
					Once()
			},
		},
		{
			name:    "not found",
			guid:    uuid.Must(uuid.NewV4()),
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

				// Delete не вызывается — не настраиваем мок
			},
		},
		{
			name:    "delete error",
			guid:    productGUID,
			wantErr: errors.New("database error"),
			prepare: func(guid uuid.UUID) {
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
					Delete(s.ctx, guid).
					Return(errors.New("database error")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.guid)

			err := s.srv.Delete(s.ctx, uuid2.UUID(tc.guid))

			if tc.wantErr != nil {
				s.Error(err)
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
	categoryGUID := uuid.Must(uuid.NewV4())

	expectedProducts := []entity.Product{
		{
			GUID:         uuid2.UUID(uuid.Must(uuid.NewV4())),
			Name:         "Product 1",
			Description:  testutil.PtrString("Description 1"),
			Price:        1000,
			CategoryGUID: uuid2.UUID(categoryGUID),
		},
		{
			GUID:         uuid2.UUID(uuid.Must(uuid.NewV4())),
			Name:         "Product 2",
			Description:  testutil.PtrString("Description 2"),
			Price:        2000,
			CategoryGUID: uuid2.UUID(categoryGUID),
		},
	}

	testCases := []struct {
		name         string
		categoryGUID *uuid.UUID
		minPrice     *int
		maxPrice     *int
		wantErr      error
		wantLen      int
		wantProducts []entity.Product
		prepare      func(categoryGUID *uuid.UUID, minPrice *int, maxPrice *int)
	}{
		{
			name:         "success",
			categoryGUID: &categoryGUID,
			minPrice:     testutil.PtrInt(500),
			maxPrice:     testutil.PtrInt(3000),
			wantErr:      nil,
			wantLen:      2,
			wantProducts: expectedProducts,
			prepare: func(categoryGUID *uuid.UUID, minPrice *int, maxPrice *int) {
				s.productRepo.EXPECT().
					List(s.ctx, nil, categoryGUID, minPrice, maxPrice).
					Return(expectedProducts, nil).
					Once()
			},
		},
		{
			name:         "empty result",
			categoryGUID: &categoryGUID,
			minPrice:     nil,
			maxPrice:     nil,
			wantErr:      nil,
			wantLen:      0,
			wantProducts: []entity.Product{},
			prepare: func(categoryGUID *uuid.UUID, minPrice *int, maxPrice *int) {
				s.productRepo.EXPECT().
					List(s.ctx, nil, categoryGUID, minPrice, maxPrice).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.categoryGUID, tc.minPrice, tc.maxPrice)

			result, err := s.srv.List(s.ctx, tc.categoryGUID, tc.minPrice, tc.maxPrice)

			if tc.wantErr != nil {
				s.Error(err)
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
	productGUID := uuid.Must(uuid.NewV4())
	categoryGUID := uuid.Must(uuid.NewV4())
	newCategoryGUID := uuid.Must(uuid.NewV4())

	existingProduct := entity.Product{
		GUID:         uuid2.UUID(productGUID),
		Name:         "Old Name",
		Description:  testutil.PtrString("Old Description"),
		Price:        1000,
		CategoryGUID: uuid2.UUID(categoryGUID),
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
				CategoryGUID: (*uuid2.UUID)(&newCategoryGUID),
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
					List(s.ctx, req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, *req.CategoryGUID).
					Return(entity.Category{GUID: uuid2.UUID(newCategoryGUID)}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.GUID == guid &&
							p.Name == *req.Name &&
							p.Description == req.Description &&
							p.Price == *req.Price &&
							p.CategoryGUID == *req.CategoryGUID
					})).
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
					List(s.ctx, req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				// CategoryGUID не передан, поэтому GetByGUID category не вызывается
				// Price = 0, но это может быть валидным значением

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.GUID == guid &&
							p.Name == *req.Name &&
							p.Description == existingProduct.Description &&
							p.Price == existingProduct.Price &&
							p.CategoryGUID == existingProduct.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			guid: uuid.Must(uuid.NewV4()),
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

				// List, GetByGUID category и Update не вызываются
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

				// Возвращаем продукт с другим GUID, но таким же именем
				otherProduct := entity.Product{
					GUID: uuid2.UUID(uuid.Must(uuid.NewV4())),
					Name: *req.Name,
				}
				s.productRepo.EXPECT().
					List(s.ctx, req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{otherProduct}, nil).
					Once()

				// Update не вызывается из-за ошибки
			},
		},
		{
			name: "category not found",
			guid: productGUID,
			req: entity.RequestProductUpdate{
				Name:         testutil.PtrString("New Name"),
				CategoryGUID: (*uuid2.UUID)(&newCategoryGUID),
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

				s.productRepo.EXPECT().
					List(s.ctx, req.Name, (*uuid.UUID)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, *req.CategoryGUID).
					Return(entity.Category{}, entity.ErrNotFound).
					Once()

				// Update не вызывается из-за ошибки
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.guid, tc.req)

			err := s.srv.Update(s.ctx, tc.guid, tc.req)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
			}
		})
	}
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
