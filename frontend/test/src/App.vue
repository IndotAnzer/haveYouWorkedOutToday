<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from './store'

const router = useRouter()
const userStore = useUserStore()

const logout = () => {
  userStore.clearToken()
  router.push('/login')
}
</script>

<template>
  <div class="app">
    <header class="header">
      <div class="logo">
        <span class="logo-icon">💪</span>
        <h1>今天你练了吗</h1>
      </div>
      <nav class="nav">
        <router-link to="/" class="nav-link">
          <span class="nav-icon">🏠</span>
          <span>首页</span>
        </router-link>
        <router-link v-if="userStore.isLoggedIn" to="/my-articles" class="nav-link">
          <span class="nav-icon">📝</span>
          <span>我的记录</span>
        </router-link>
        <template v-if="userStore.isLoggedIn">
          <button @click="logout" class="logout-btn">
            <span class="nav-icon">🚪</span>
            <span>退出</span>
          </button>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link">
            <span class="nav-icon">🔐</span>
            <span>登录</span>
          </router-link>
          <router-link to="/register" class="register-btn">
            <span class="nav-icon">✨</span>
            <span>注册</span>
          </router-link>
        </template>
      </nav>
    </header>
    <main class="main">
      <router-view></router-view>
    </main>
    <footer class="footer">
      <p>💪 坚持训练，遇见更好的自己</p>
    </footer>
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 32px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.header h1 {
  font-size: 22px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 700;
  margin: 0;
}

.nav {
  display: flex;
  gap: 8px;
  align-items: center;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
  color: #555;
  font-weight: 500;
  padding: 10px 16px;
  border-radius: 12px;
  transition: all 0.3s ease;
  font-size: 14px;
}

.nav-link:hover {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  color: #667eea;
  transform: translateY(-2px);
}

.nav-link.router-link-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.nav-icon {
  font-size: 16px;
}

.register-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.3);
}

.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 107, 107, 0.4);
}

.main {
  flex: 1;
  max-width: 1000px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
}

.footer {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 16px;
  text-align: center;
  color: #888;
  font-size: 13px;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
  }

  .nav {
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
  }

  .nav-link, .register-btn, .logout-btn {
    padding: 8px 14px;
    font-size: 13px;
  }

  .header h1 {
    font-size: 18px;
  }

  .main {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .nav-link span:last-child,
  .register-btn span:last-child,
  .logout-btn span:last-child {
    display: none;
  }

  .nav-link, .register-btn, .logout-btn {
    padding: 10px;
  }
}
</style>
