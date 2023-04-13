// Code generated by MockGen. DO NOT EDIT.
// Source: services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	models "crud-books/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTokener is a mock of Tokener interface.
type MockTokener struct {
	ctrl     *gomock.Controller
	recorder *MockTokenerMockRecorder
}

// MockTokenerMockRecorder is the mock recorder for MockTokener.
type MockTokenerMockRecorder struct {
	mock *MockTokener
}

// NewMockTokener creates a new mock instance.
func NewMockTokener(ctrl *gomock.Controller) *MockTokener {
	mock := &MockTokener{ctrl: ctrl}
	mock.recorder = &MockTokenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokener) EXPECT() *MockTokenerMockRecorder {
	return m.recorder
}

// GenerateTokens mocks base method.
func (m *MockTokener) GenerateTokens(userId string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokens", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokens indicates an expected call of GenerateTokens.
func (mr *MockTokenerMockRecorder) GenerateTokens(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokens", reflect.TypeOf((*MockTokener)(nil).GenerateTokens), userId)
}

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// CreateBook mocks base method.
func (m *MockDB) CreateBook(title, description, fileToken, emailOwner string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBook", title, description, fileToken, emailOwner)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBook indicates an expected call of CreateBook.
func (mr *MockDBMockRecorder) CreateBook(title, description, fileToken, emailOwner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBook", reflect.TypeOf((*MockDB)(nil).CreateBook), title, description, fileToken, emailOwner)
}

// CreateUser mocks base method.
func (m *MockDB) CreateUser(email, passwordHash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", email, passwordHash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockDBMockRecorder) CreateUser(email, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDB)(nil).CreateUser), email, passwordHash)
}

// DeleteBook mocks base method.
func (m *MockDB) DeleteBook(tokenBook string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBook", tokenBook)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBook indicates an expected call of DeleteBook.
func (mr *MockDBMockRecorder) DeleteBook(tokenBook interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBook", reflect.TypeOf((*MockDB)(nil).DeleteBook), tokenBook)
}

// GetBook mocks base method.
func (m *MockDB) GetBook(bookToken string) (*models.BookData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBook", bookToken)
	ret0, _ := ret[0].(*models.BookData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBook indicates an expected call of GetBook.
func (mr *MockDBMockRecorder) GetBook(bookToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBook", reflect.TypeOf((*MockDB)(nil).GetBook), bookToken)
}

// GetFileData mocks base method.
func (m *MockDB) GetFileData(fileToken string) (*models.FileData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileData", fileToken)
	ret0, _ := ret[0].(*models.FileData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileData indicates an expected call of GetFileData.
func (mr *MockDBMockRecorder) GetFileData(fileToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileData", reflect.TypeOf((*MockDB)(nil).GetFileData), fileToken)
}

// GetListBooksOfUser mocks base method.
func (m *MockDB) GetListBooksOfUser(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListBooksOfUser", paramsOfBooks)
	ret0, _ := ret[0].(*[]models.BookData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListBooksOfUser indicates an expected call of GetListBooksOfUser.
func (mr *MockDBMockRecorder) GetListBooksOfUser(paramsOfBooks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListBooksOfUser", reflect.TypeOf((*MockDB)(nil).GetListBooksOfUser), paramsOfBooks)
}

// GetListBooksPublic mocks base method.
func (m *MockDB) GetListBooksPublic(paramsOfBooks *models.ValidateDataInGetLists) (*[]models.BookData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListBooksPublic", paramsOfBooks)
	ret0, _ := ret[0].(*[]models.BookData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListBooksPublic indicates an expected call of GetListBooksPublic.
func (mr *MockDBMockRecorder) GetListBooksPublic(paramsOfBooks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListBooksPublic", reflect.TypeOf((*MockDB)(nil).GetListBooksPublic), paramsOfBooks)
}

// GetUserData mocks base method.
func (m *MockDB) GetUserData(email string) (*models.UserData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", email)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserData indicates an expected call of GetUserData.
func (mr *MockDBMockRecorder) GetUserData(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockDB)(nil).GetUserData), email)
}

// GetUserDataByInsertedId mocks base method.
func (m *MockDB) GetUserDataByInsertedId(userId string) (*models.UserData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDataByInsertedId", userId)
	ret0, _ := ret[0].(*models.UserData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDataByInsertedId indicates an expected call of GetUserDataByInsertedId.
func (mr *MockDBMockRecorder) GetUserDataByInsertedId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDataByInsertedId", reflect.TypeOf((*MockDB)(nil).GetUserDataByInsertedId), userId)
}

// UpdateBook mocks base method.
func (m *MockDB) UpdateBook(bookId string, updater models.BookDataUpdater) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBook", bookId, updater)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBook indicates an expected call of UpdateBook.
func (mr *MockDBMockRecorder) UpdateBook(bookId, updater interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBook", reflect.TypeOf((*MockDB)(nil).UpdateBook), bookId, updater)
}

// UploadFileData mocks base method.
func (m *MockDB) UploadFileData(fileToken, downloadPage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileData", fileToken, downloadPage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFileData indicates an expected call of UploadFileData.
func (mr *MockDBMockRecorder) UploadFileData(fileToken, downloadPage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileData", reflect.TypeOf((*MockDB)(nil).UploadFileData), fileToken, downloadPage)
}

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// DeleteFile mocks base method.
func (m *MockStorager) DeleteFile(fileToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", fileToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockStoragerMockRecorder) DeleteFile(fileToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockStorager)(nil).DeleteFile), fileToken)
}

// UploadFile mocks base method.
func (m *MockStorager) UploadFile(file []byte, isTest bool) (*models.UploadFileReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", file, isTest)
	ret0, _ := ret[0].(*models.UploadFileReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockStoragerMockRecorder) UploadFile(file, isTest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockStorager)(nil).UploadFile), file, isTest)
}

// MockHasher is a mock of Hasher interface.
type MockHasher struct {
	ctrl     *gomock.Controller
	recorder *MockHasherMockRecorder
}

// MockHasherMockRecorder is the mock recorder for MockHasher.
type MockHasherMockRecorder struct {
	mock *MockHasher
}

// NewMockHasher creates a new mock instance.
func NewMockHasher(ctrl *gomock.Controller) *MockHasher {
	mock := &MockHasher{ctrl: ctrl}
	mock.recorder = &MockHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHasher) EXPECT() *MockHasherMockRecorder {
	return m.recorder
}

// CompareHashWithPassword mocks base method.
func (m *MockHasher) CompareHashWithPassword(password, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareHashWithPassword", password, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompareHashWithPassword indicates an expected call of CompareHashWithPassword.
func (mr *MockHasherMockRecorder) CompareHashWithPassword(password, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareHashWithPassword", reflect.TypeOf((*MockHasher)(nil).CompareHashWithPassword), password, hash)
}

// GetNewHash mocks base method.
func (m *MockHasher) GetNewHash(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewHash", password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewHash indicates an expected call of GetNewHash.
func (mr *MockHasherMockRecorder) GetNewHash(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewHash", reflect.TypeOf((*MockHasher)(nil).GetNewHash), password)
}