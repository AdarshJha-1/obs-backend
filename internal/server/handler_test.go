package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"obs/internal/database"
	"obs/internal/models"
	"obs/internal/server"
	"obs/internal/types"
	"obs/internal/utils"
)

// MockDB implements database.Service
type MockDB struct {
	mock.Mock
}

func (m *MockDB) MigrateSchema() {
	m.Called()
}

func (m *MockDB) Health() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetUser(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockDB) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetBlogs() ([]models.Blog, error) {
	args := m.Called()
	return args.Get(0).([]models.Blog), args.Error(1)
}

func (m *MockDB) GetBlog(id uint) (*models.Blog, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Blog), args.Error(1)
}

func (m *MockDB) CreateBlog(blog *models.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockDB) DeleteBlog(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) UpdateBlog(blog *models.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}
func setupRouter(db database.Service) (*gin.Engine, *server.Server) {
	s := server.NewServer // Corrected line
	router := s.RegisterRoutes().(*gin.Engine)
	return router, s
}

func TestGetUser(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	user := &models.User{ID: 1, Username: "testuser"}
	mockDB.On("GetUser", uint(1)).Return(user, nil)

	req, _ := http.NewRequest("GET", "/api/user/u/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, user.Username, response["users"].(map[string]interface{})["username"])
	mockDB.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	users := []models.User{{ID: 1, Username: "testuser"}}
	mockDB.On("GetUsers").Return(users, nil)

	req, _ := http.NewRequest("GET", "/api/user/all", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser", response["users"].([]interface{})[0].(map[string]interface{})["username"])
	mockDB.AssertExpectations(t)
}

func TestRegisterUser(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password"}
	userJSON, _ := json.Marshal(user)

	mockDB.On("GetUserByEmail", user.Email).Return(nil, nil)
	mockDB.On("CreateUser", mock.Anything).Return(nil)

	req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDB.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	login := types.SignInModel{Identifier: "test@example.com", Password: "password"}
	loginJSON, _ := json.Marshal(login)

	hashedPassword, _ := utils.HashPassword("password")                                      // Handle both return values
	user := &models.User{ID: 1, Email: "test@example.com", Password: string(hashedPassword)} // Convert []byte to string
	mockDB.On("GetUserByEmail", login.Identifier).Return(user, nil)

	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}
func TestDeleteUserById(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	mockDB.On("DeleteUser", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/user/u/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}

func TestUpdateUserById(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	user := models.User{Username: "updateduser"}
	userJSON, _ := json.Marshal(user)

	mockDB.On("UpdateUser", mock.Anything).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/user/u/1", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}

func TestGetAllBlogs(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	blogs := []models.Blog{{ID: 1, Title: "Test Blog", Content: "Test Content"}}
	mockDB.On("GetBlogs").Return(blogs, nil)

	req, _ := http.NewRequest("GET", "/api/blog/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Blog",
		response["blogs"].([]interface{})[0].(map[string]interface{})["title"])
	mockDB.AssertExpectations(t)
}

func TestGetBlogByID(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	blog := &models.Blog{ID: 1, Title: "Test Blog"}
	mockDB.On("GetBlog", uint(1)).Return(blog, nil)

	req, _ := http.NewRequest("GET", "/api/blog/b/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Blog", response["blog"].(map[string]interface{})["title"])
	mockDB.AssertExpectations(t)
}

func TestCreateNewBlog(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	blog := models.Blog{Title: "New Blog", Content: "New Content"}
	blogJSON, _ := json.Marshal(blog)

	mockDB.On("CreateBlog", mock.Anything).Return(nil)

	req, _ := http.NewRequest("POST", "/api/blog/", bytes.NewBuffer(blogJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDB.AssertExpectations(t)
}

func TestDeleteBlogByID(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	mockDB.On("DeleteBlog", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/blog/b/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockDB.AssertExpectations(t)
}

func TestUpdateBlog(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	blog := models.Blog{Title: "Updated Blog", Content: "Updated Content"}
	blogJSON, _ := json.Marshal(blog)

	mockDB.On("UpdateBlog", mock.Anything).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/blog/b/1", bytes.NewBuffer(blogJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}

func TestHelloWorldHandler(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Hello World", response["message"])
}

func TestHealthHandler(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	healthResponse := map[string]string{"status": "ok"}
	mockDB.On("Health").Return(healthResponse)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	mockDB.AssertExpectations(t)
}

func TestRegisterRoutesCORS(t *testing.T) {
	mockDB := new(MockDB)
	router, _ := setupRouter(mockDB)

	req, _ := http.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
