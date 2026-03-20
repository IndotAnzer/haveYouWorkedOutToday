package tests

import (
	"bytes"
	"encoding/json"
	"haveYouWorkedOutToday/controllers"
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"haveYouWorkedOutToday/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	// 自动迁移模型结构
	err = db.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{}, &models.Reply{}, &models.FitnessAction{}, &models.FitnessActionGroup{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	global.Db = db
}

func TestRegister(t *testing.T) {
	setupTestDB()

	// 测试正常场景：注册成功
	t.Run("正常场景：注册成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建测试数据
		user := models.User{
			Username: "testuser",
			Password: "password123",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用注册函数
		controllers.Register(c)

		// 验证响应
		assert.Equal(t, http.StatusCreated, w.Code)

		// 验证响应体包含token
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
	})

	// 测试异常场景：用户名已存在
	t.Run("异常场景：用户名已存在", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 先创建一个用户
		existingUser := models.User{
			Username: "existinguser",
			Password: "password123",
		}
		global.Db.Create(&existingUser)

		// 尝试注册相同的用户名
		user := models.User{
			Username: "existinguser",
			Password: "password123",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用注册函数
		controllers.Register(c)

		// 验证响应
		assert.Equal(t, http.StatusBadRequest, w.Code) // 修复后应该返回400，用户名已存在
	})

	// 测试异常场景：请求体格式错误
	t.Run("异常场景：请求体格式错误", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建格式错误的请求体
		invalidData := []byte(`{"username": "testuser"}`) // 缺少password字段

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(invalidData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用注册函数
		controllers.Register(c)

		// 验证响应
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	setupTestDB()

	// 先创建一个测试用户
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := models.User{
		Username: "testuser",
		Password: hashedPassword,
	}
	global.Db.Create(&testUser)

	// 测试正常场景：登录成功
	t.Run("正常场景：登录成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建登录数据
		loginData := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "testuser",
			Password: "password123",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(loginData)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用登录函数
		controllers.Login(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体包含token和user信息
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")
	})

	// 测试异常场景：用户名不存在
	t.Run("异常场景：用户名不存在", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建登录数据
		loginData := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "nonexistent",
			Password: "password123",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(loginData)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用登录函数
		controllers.Login(c)

		// 验证响应
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// 测试异常场景：密码错误
	t.Run("异常场景：密码错误", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建登录数据
		loginData := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "testuser",
			Password: "wrongpassword",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(loginData)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用登录函数
		controllers.Login(c)

		// 验证响应
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// 测试异常场景：请求体格式错误
	t.Run("异常场景：请求体格式错误", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建格式错误的请求体
		invalidData := []byte(`{"username": "testuser"}`) // 缺少password字段

		// 创建请求
		req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(invalidData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用登录函数
		controllers.Login(c)

		// 验证响应
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
