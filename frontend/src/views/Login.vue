<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md p-8">
      <div class="text-center mb-8">
        <el-icon :size="48" class="text-blue-600">
          <Lock />
        </el-icon>
        <h1 class="text-2xl font-bold text-gray-800 mt-4">后台管理登录</h1>
        <p class="text-gray-500 mt-2">请输入管理员密码以访问后台管理系统</p>
      </div>

      <el-form :model="loginForm" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="管理员用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="管理员密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="w-full"
            native-type="submit"
          >
            登录后台
          </el-button>
        </el-form-item>
      </el-form>

      <div class="text-center mt-6">
        <el-button @click="goHome" link class="text-gray-500">
          <el-icon class="mr-1"><House /></el-icon>
          返回前台
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await fetch('/api/admin/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(loginForm)
        })
        
        if (response.ok) {
          const data = await response.json()
          localStorage.setItem('admin_token', data.token)
          localStorage.setItem('admin_user', JSON.stringify(data.user))
          // 移除登录成功的提示，避免与Admin.vue中的提示重复
          // ElMessage.success('登录成功')
          router.push('/admin')
        } else {
          ElMessage.error('登录失败，请检查用户名和密码')
        }
      } catch (error) {
        console.error('Login failed:', error)
        ElMessage.error('登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const goHome = () => {
  router.push('/')
}
</script>
