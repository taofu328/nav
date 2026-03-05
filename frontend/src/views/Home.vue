<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <header class="bg-white shadow-md sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center space-x-3">
            <el-icon :size="32" class="text-blue-600">
              <Link />
            </el-icon>
            <h1 class="text-xl font-bold text-gray-800">{{ siteSettings.siteTitle }}</h1>
          </div>
          <div class="flex items-center space-x-4">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索网址..."
              prefix-icon="Search"
              clearable
              class="w-64"
              @input="handleSearch"
            />
            <el-button @click="showAdminLogin" type="primary" size="small">
              <el-icon class="mr-1"><Setting /></el-icon>
              管理
            </el-button>
          </div>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="bg-white rounded-lg shadow p-4 mb-6">
        <div class="flex flex-wrap gap-2">
          <div
            @click="selectedCategory = null"
            :class="[
              'px-4 py-2 rounded-full cursor-pointer transition-all duration-300',
              selectedCategory === null 
                ? 'bg-blue-600 text-white shadow-md transform scale-105' 
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:shadow-md'
            ]"
            class="category-tag"
          >
            <span class="text-sm font-medium">全部</span>
          </div>
          <div
            v-for="category in categories"
            :key="category.id"
            @click="selectedCategory = category.id"
            :class="[
              'px-4 py-2 rounded-full cursor-pointer transition-all duration-300',
              selectedCategory === category.id 
                ? 'bg-blue-600 text-white shadow-md transform scale-105' 
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:shadow-md'
            ]"
            class="category-tag"
          >
            <span class="text-sm font-medium">{{ category.name }}</span>
          </div>
        </div>
      </div>

      <div v-if="loading" class="text-center py-12">
        <el-icon class="is-loading" :size="40">
          <Loading />
        </el-icon>
      </div>

      <div v-else-if="filteredCategories.length === 0" class="text-center py-12 text-gray-500">
        <el-icon :size="64" class="text-gray-300">
          <Document />
        </el-icon>
        <p class="mt-4">暂无网址</p>
      </div>

      <div v-else>
        <div v-for="category in filteredCategories" :key="category.id" class="mb-8">
          <div class="bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100">
              <div class="flex items-center">
                <el-icon class="text-blue-600 mr-2">
                  <Folder />
                </el-icon>
                <h2 class="text-lg font-semibold text-gray-800">{{ category.name }}</h2>
              </div>
              <span class="text-sm text-gray-400">{{ getCategoryBookmarks(category.id).length }} 个网址</span>
            </div>
            <div class="p-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                <a
                  v-for="bookmark in getCategoryBookmarks(category.id)"
                  :key="bookmark.id"
                  :href="bookmark.url"
                  target="_blank"
                  @click="incrementVisit(bookmark)"
                  class="block group"
                >
                  <div class="bg-gray-50 rounded-lg p-4 hover:bg-blue-50 hover:border-blue-200 border border-transparent hover:border transition-all">
                    <div class="flex items-start space-x-3">
                      <img
                        :src="getIconUrl(bookmark)"
                        :alt="bookmark.title"
                        class="w-12 h-12 rounded object-cover flex-shrink-0"
                        @error="handleIconError"
                      />
                      <div class="flex-1 min-w-0">
                        <h3 class="font-medium text-gray-800 group-hover:text-blue-600 truncate">
                          {{ bookmark.title }}
                        </h3>
                        <p v-if="bookmark.description" class="text-sm text-gray-500 mt-1 line-clamp-2">
                          {{ bookmark.description }}
                        </p>
                      </div>
                    </div>
                  </div>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer class="bg-white border-t border-gray-200 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex justify-between items-center">
          <div class="flex items-center space-x-3">
              <img v-if="siteSettings.siteLogo" :src="siteSettings.siteLogo" alt="Logo" class="h-8 w-auto" />
              <div class="flex flex-col">
                <p class="text-gray-500 text-sm">
                  © 2024 {{ siteSettings.siteTitle }}. All rights reserved.
                </p>
              </div>
            </div>
        </div>
      </div>
    </footer>


  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const allCategories = ref([])
const bookmarks = ref([])
const siteSettings = reactive({
  siteTitle: 'Van Nav',
  siteLogo: ''
})

// 只显示有有效网址的分类
const categories = computed(() => {
  return allCategories.value.filter(category => {
    const categoryBookmarks = bookmarks.value.filter(b => b.category_id === category.id)
    return categoryBookmarks.length > 0
  })
})



const filteredCategories = computed(() => {
  let result = categories.value
  
  // 按选中的分类过滤
  if (selectedCategory.value) {
    result = result.filter(category => category.id === selectedCategory.value)
  }
  
  // 按搜索关键词过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(category => {
      const categoryBookmarks = getCategoryBookmarks(category.id)
      return categoryBookmarks.some(bookmark => 
        bookmark.title.toLowerCase().includes(keyword) ||
        bookmark.url.toLowerCase().includes(keyword) ||
        (bookmark.description && bookmark.description.toLowerCase().includes(keyword))
      )
    })
  }
  
  return result
})

const getCategoryBookmarks = (categoryId) => {
  return bookmarks.value.filter(b => b.category_id === categoryId)
}

const getIconUrl = (bookmark) => {
  if (bookmark.icon && bookmark.icon.startsWith('/api/icons/')) {
    return bookmark.icon
  }
  if (bookmark.icon) {
    return bookmark.icon
  }
  return '/api/icons/default.svg'
}

const handleIconError = (event) => {
  event.target.src = '/api/icons/default.svg'
}

const handleSearch = () => {
}

const showAdminLogin = () => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    router.push('/admin')
  } else {
    router.push('/login')
  }
}



const incrementVisit = async (bookmark) => {
  try {
    await fetch(`/api/bookmarks/${bookmark.id}/visit`, { method: 'POST' })
  } catch (error) {
    console.error('Failed to increment visit:', error)
  }
}

const loadPublicData = async () => {
  loading.value = true
  try {
    const [categoriesRes, bookmarksRes, settingsRes] = await Promise.all([
      fetch('/api/public/categories'),
      fetch('/api/public/bookmarks'),
      fetch('/api/admin/settings')
    ])
    
    allCategories.value = await categoriesRes.json()
    bookmarks.value = await bookmarksRes.json()
    
    // 从后端 API 获取网站设置
    if (settingsRes.ok) {
      const settingsData = await settingsRes.json()
      if (settingsData.settings) {
        const savedSettings = settingsData.settings
        Object.assign(siteSettings, {
          siteTitle: savedSettings.site_title || 'Van Nav',
          siteLogo: savedSettings.site_logo || ''
        })
        
        // 将设置保存到 localStorage 以便快速加载
        localStorage.setItem('site_settings', JSON.stringify({
          siteTitle: savedSettings.site_title || 'Van Nav',
          siteLogo: savedSettings.site_logo || ''
        }))
      }
    } else {
      // 如果 API 请求失败，尝试从 localStorage 加载
      const settingsData = localStorage.getItem('site_settings')
      if (settingsData) {
        const savedSettings = JSON.parse(settingsData)
        Object.assign(siteSettings, {
          siteTitle: savedSettings.siteTitle || 'Van Nav',
          siteLogo: savedSettings.siteLogo || ''
        })
      }
    }
  } catch (error) {
    console.error('Failed to load data:', error)
    ElMessage.error('加载数据失败')
    
    // 如果 API 请求失败，尝试从 localStorage 加载
    const settingsData = localStorage.getItem('site_settings')
    if (settingsData) {
      const savedSettings = JSON.parse(settingsData)
      Object.assign(siteSettings, {
        siteTitle: savedSettings.siteTitle || 'Van Nav',
        siteLogo: savedSettings.siteLogo || ''
      })
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadPublicData()
})
</script>

<style scoped>
.category-tag {
  transition: all 0.3s ease;
  user-select: none;
}

.category-tag:hover {
  transform: translateY(-2px);
}

.category-tag.active {
  transform: scale(1.05);
  box-shadow: 0 4px 6px -1px rgba(37, 99, 235, 0.1);
}
</style>
