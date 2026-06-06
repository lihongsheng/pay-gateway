<template>
  <div class="welcome-container">
    <!-- 欢迎横幅 -->
    <div class="welcome-banner">
      <h1 class="welcome-title">欢迎回来，{{ username }}！</h1>
      <p class="welcome-date">今天是 {{ currentDate }}，祝您工作愉快</p>
    </div>
  </div>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue'
import { useUserStore } from '@/pinia/modules/user'

// 用户信息
const userStore = useUserStore()
const username = computed(() => userStore.userInfo?.nickName || userStore.userInfo?.userName || '管理员')

// 当前日期
const currentDate = ref('')

// 格式化当前日期
const formatCurrentDate = () => {
  const date = new Date()
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const weekday = weekdays[date.getDay()]

  currentDate.value = `${year}年${month}月${day}日 ${weekday}`
}

onMounted(() => {
  formatCurrentDate()
})

defineOptions({
  name: 'Welcome'
})
</script>

<style scoped>
.welcome-container {
  width: 100%;
  min-height: 100%;
  background: transparent; /* 透明背景，跟随主体页面背景 */
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.welcome-banner {
  text-align: center;
  max-width: 800px;
  width: 100%;
  background: white;
  border-radius: 12px;
  padding: 48px 32px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  animation: fadeIn 0.5s ease-in-out;
}

.welcome-title {
  font-size: 2.5rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 16px;
  line-height: 1.3;
}

.welcome-date {
  font-size: 1.25rem;
  color: #64748b;
  margin: 0;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .welcome-banner {
    padding: 32px 20px;
  }

  .welcome-title {
    font-size: 2rem;
  }

  .welcome-date {
    font-size: 1.125rem;
  }
}

@media (max-width: 480px) {
  .welcome-title {
    font-size: 1.5rem;
  }

  .welcome-date {
    font-size: 1rem;
  }
}
</style>