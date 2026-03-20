package tests

import (
	"bytes"
	"encoding/json"
	"haveYouWorkedOutToday/controllers"
	"haveYouWorkedOutToday/global"
	"haveYouWorkedOutToday/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	setupTestDB()

	// 测试正常场景：创建文章成功
	t.Run("正常场景：创建文章成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 先创建一个用户
		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		global.Db.Create(&user)

		// 创建文章数据
		article := models.Article{
			Title:   "测试文章",
			Content: "测试内容",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(article)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("username", "testuser") // 设置用户认证信息

		// 调用创建文章函数
		controllers.CreateArticle(c)

		// 验证响应
		assert.Equal(t, http.StatusCreated, w.Code)

		// 验证响应体包含文章信息
		var response models.Article
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "测试文章", response.Title)
		assert.Equal(t, "测试内容", response.Content)
	})

	// 测试异常场景：未认证用户
	t.Run("异常场景：未认证用户", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建文章数据
		article := models.Article{
			Title:   "测试文章",
			Content: "测试内容",
		}

		// 将数据转换为JSON
		jsonData, err := json.Marshal(article)
		assert.NoError(t, err)

		// 创建请求
		req, err := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// 不设置用户认证信息

		// 调用创建文章函数
		controllers.CreateArticle(c)

		// 验证响应
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// 测试异常场景：请求体格式错误
	t.Run("异常场景：请求体格式错误", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建格式错误的请求体
		invalidData := []byte(`{"title": "测试文章"}`) // 缺少content字段

		// 创建请求
		req, err := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(invalidData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("username", "testuser") // 设置用户认证信息

		// 调用创建文章函数
		controllers.CreateArticle(c)

		// 验证响应
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllArticles(t *testing.T) {
	setupTestDB()

	// 测试正常场景：获取所有文章成功
	t.Run("正常场景：获取所有文章成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建测试文章
		articles := []models.Article{
			{Title: "文章1", Content: "内容1"},
			{Title: "文章2", Content: "内容2"},
		}
		global.Db.Create(&articles)

		// 创建请求
		req, err := http.NewRequest("GET", "/api/articles", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用获取所有文章函数
		controllers.GetAllArticles(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体包含文章列表
		var response []models.Article
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	// 测试边界场景：没有文章时返回空列表
	t.Run("边界场景：没有文章时返回空列表", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 确保数据库中没有文章
		global.Db.Where("1 = 1").Delete(&models.Article{})

		// 创建请求
		req, err := http.NewRequest("GET", "/api/articles", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 调用获取所有文章函数
		controllers.GetAllArticles(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体为空列表
		var response []models.Article
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response)
	})
}

func TestGetArticleByID(t *testing.T) {
	setupTestDB()

	// 测试正常场景：获取文章详情成功
	t.Run("正常场景：获取文章详情成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建测试文章
		article := models.Article{
			Title:   "测试文章",
			Content: "测试内容",
		}
		global.Db.Create(&article)

		// 创建请求
		req, err := http.NewRequest("GET", "/api/articles/1", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}} // 设置文章ID参数

		// 调用获取文章详情函数
		controllers.GetArticleByID(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体包含文章信息
		var response models.Article
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "测试文章", response.Title)
		assert.Equal(t, "测试内容", response.Content)
	})

	// 测试异常场景：文章不存在
	t.Run("异常场景：文章不存在", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建请求
		req, err := http.NewRequest("GET", "/api/articles/999", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "999"}} // 设置不存在的文章ID

		// 调用获取文章详情函数
		controllers.GetArticleByID(c)

		// 验证响应
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetArticleByUser(t *testing.T) {
	setupTestDB()

	// 测试正常场景：获取用户自己的文章成功
	t.Run("正常场景：获取用户自己的文章成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 先创建一个用户
		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		global.Db.Create(&user)

		// 创建该用户的文章
		articles := []models.Article{
			{Title: "我的文章1", Content: "内容1", UserID: user.ID},
			{Title: "我的文章2", Content: "内容2", UserID: user.ID},
		}
		global.Db.Create(&articles)

		// 创建请求
		req, err := http.NewRequest("GET", "/api/my/articles", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("username", "testuser") // 设置用户认证信息

		// 调用获取用户文章函数
		controllers.GetArticleByUser(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体包含用户的文章列表
		var response []models.Article
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	// 测试异常场景：未认证用户
	t.Run("异常场景：未认证用户", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建请求
		req, err := http.NewRequest("GET", "/api/my/articles", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// 不设置用户认证信息

		// 调用获取用户文章函数
		controllers.GetArticleByUser(c)

		// 验证响应
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestDeleteArticle(t *testing.T) {
	setupTestDB()

	// 测试正常场景：删除文章成功
	t.Run("正常场景：删除文章成功", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建测试文章
		article := models.Article{
			Title:   "测试文章",
			Content: "测试内容",
		}
		global.Db.Create(&article)

		// 创建请求
		req, err := http.NewRequest("DELETE", "/api/articles/1", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}} // 设置文章ID参数

		// 调用删除文章函数
		controllers.DeleteArticle(c)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		// 验证响应体包含成功消息
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Article deleted", response["message"])
	})

	// 测试异常场景：文章不存在
	t.Run("异常场景：文章不存在", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// 创建请求
		req, err := http.NewRequest("DELETE", "/api/articles/999", nil)
		assert.NoError(t, err)

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 创建Gin上下文
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "999"}} // 设置不存在的文章ID

		// 调用删除文章函数
		controllers.DeleteArticle(c)

		// 验证响应
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
