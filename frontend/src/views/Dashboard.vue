<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <el-icon :size="32" class="text-blue-600">
              <Link />
            </el-icon>
            <h1 class="ml-3 text-xl font-bold text-gray-800">网址导航</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-gray-600">{{ authStore.user?.username }}</span>
            <el-button @click="handleLogout" type="danger" size="small">
              退出登录
            </el-button>
          </div>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex flex-col gap-6">
        <div class="bg-white rounded-lg shadow p-4">
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

        <div class="flex flex-col lg:flex-row gap-6">
          <aside class="lg:w-64 flex-shrink-0">
            <div class="bg-white rounded-lg shadow p-4">
              <div class="flex justify-between items-center mb-4">
                <h2 class="text-lg font-semibold text-gray-800">数据管理</h2>
                <el-button @click="showCategoryDialog = true" type="primary" size="small" circle>
                  <el-icon><Plus /></el-icon>
                </el-button>
              </div>
              <div class="space-y-2">
                <el-button @click="exportData" class="w-full mb-2" type="success">
                  <el-icon class="mr-2"><Download /></el-icon>
                  导出数据
                </el-button>
                <el-upload
                  :auto-upload="false"
                  :on-change="importData"
                  :show-file-list="false"
                  accept=".json"
                >
                  <el-button class="w-full" type="warning">
                    <el-icon class="mr-2"><Upload /></el-icon>
                    导入数据
                  </el-button>
                </el-upload>
              </div>
            </div>
          </aside>

        <div class="flex-1">
          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
              <h2 class="text-xl font-semibold text-gray-800">
                {{ selectedCategory ? getCategoryName(selectedCategory) : '全部网址' }}
              </h2>
              <div class="flex gap-3 w-full sm:w-auto">
                <el-input
                  v-model="searchKeyword"
                  placeholder="搜索网址..."
                  prefix-icon="Search"
                  clearable
                  class="flex-1 sm:w-64"
                />
                <el-button @click="showBookmarkDialog = true" type="primary">
                  <el-icon><Plus /></el-icon>
                  添加网址
                </el-button>
              </div>
            </div>

            <div v-if="loading" class="text-center py-12">
              <el-icon class="is-loading" :size="40">
                <Loading />
              </el-icon>
            </div>

            <div v-else-if="filteredBookmarks.length === 0" class="text-center py-12 text-gray-500">
              <el-icon :size="64" class="text-gray-300">
                <Document />
              </el-icon>
              <p class="mt-4">暂无网址，点击上方按钮添加</p>
            </div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              <div
                v-for="bookmark in filteredBookmarks"
                :key="bookmark.id"
                class="border border-gray-200 rounded-lg p-4 hover:shadow-lg transition-shadow"
              >
                <div class="flex justify-between items-start">
                  <a
                    :href="bookmark.url"
                    target="_blank"
                    class="flex-1"
                    @click="incrementVisit(bookmark)"
                  >
                    <h3 class="font-semibold text-gray-800 hover:text-blue-600 truncate">
                      {{ bookmark.title }}
                    </h3>
                    <p class="text-sm text-gray-500 mt-1 truncate">{{ bookmark.url }}</p>
                  </a>
                  <el-dropdown @command="(cmd) => handleBookmarkAction(cmd, bookmark)">
                    <el-icon class="cursor-pointer hover:text-blue-600">
                      <MoreFilled />
                    </el-icon>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="edit">编辑</el-dropdown-item>
                        <el-dropdown-item command="delete">删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
                <p v-if="bookmark.description" class="text-sm text-gray-600 mt-2 line-clamp-2">
                  {{ bookmark.description }}
                </p>
                <div class="flex items-center justify-between mt-3 text-xs text-gray-400">
                  <span>{{ getCategoryName(bookmark.category_id) }}</span>
                  <span>访问 {{ bookmark.visit_count }} 次</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <el-dialog v-model="showBookmarkDialog" :title="editingBookmark ? '编辑网址' : '添加网址'" width="500px">
      <el-form :model="bookmarkForm" :rules="bookmarkRules" ref="bookmarkFormRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="bookmarkForm.title" placeholder="网址标题" />
        </el-form-item>
        <el-form-item label="网址" prop="url">
          <el-input v-model="bookmarkForm.url" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="bookmarkForm.category_id" placeholder="选择分类" class="w-full">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="bookmarkForm.description"
            type="textarea"
            :rows="3"
            placeholder="网址描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBookmarkDialog = false">取消</el-button>
        <el-button type="primary" @click="saveBookmark" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showCategoryDialog" :title="editingCategory ? '编辑分类' : '添加分类'" width="400px">
      <el-form :model="categoryForm" :rules="categoryRules" ref="categoryFormRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="categoryForm.name" placeholder="分类名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            :rows="3"
            placeholder="分类描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCategoryDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCategory" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const saving = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const showBookmarkDialog = ref(false)
const showCategoryDialog = ref(false)
const editingBookmark = ref(null)
const editingCategory = ref(null)

const categories = ref([])
const bookmarks = ref([])

const bookmarkFormRef = ref(null)
const categoryFormRef = ref(null)

const bookmarkForm = reactive({
  title: '',
  url: '',
  category_id: null,
  description: ''
})

const categoryForm = reactive({
  name: '',
  description: ''
})

const bookmarkRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入网址', trigger: 'blur' },
    { type: 'url', message: '请输入正确的网址格式', trigger: 'blur' }
  ],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
}

const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const filteredBookmarks = computed(() => {
  let result = bookmarks.value
  
  if (selectedCategory.value) {
    result = result.filter(b => b.category_id === selectedCategory.value)
  }
  
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(b => 
      b.title.toLowerCase().includes(keyword) || 
      b.url.toLowerCase().includes(keyword) ||
      (b.description && b.description.toLowerCase().includes(keyword))
    )
  }
  
  return result
})

const getCategoryName = (categoryId) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : '未分类'
}

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await api.get('/categories')
    categories.value = response.data
  } catch (error) {
    console.error('Failed to load categories:', error)
  } finally {
    loading.value = false
  }
}

const loadBookmarks = async () => {
  loading.value = true
  try {
    const response = await api.get('/bookmarks')
    bookmarks.value = response.data
  } catch (error) {
    console.error('Failed to load bookmarks:', error)
  } finally {
    loading.value = false
  }
}

const saveBookmark = async () => {
  if (!bookmarkFormRef.value) return
  
  await bookmarkFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingBookmark.value) {
          await api.put(`/bookmarks/${editingBookmark.value.id}`, bookmarkForm)
          ElMessage.success('更新成功')
        } else {
          await api.post('/bookmarks', bookmarkForm)
          ElMessage.success('添加成功')
        }
        showBookmarkDialog.value = false
        resetBookmarkForm()
        await loadBookmarks()
      } catch (error) {
        console.error('Failed to save bookmark:', error)
      } finally {
        saving.value = false
      }
    }
  })
}

const handleBookmarkAction = async (command, bookmark) => {
  if (command === 'edit') {
    editingBookmark.value = bookmark
    Object.assign(bookmarkForm, {
      title: bookmark.title,
      url: bookmark.url,
      category_id: bookmark.category_id,
      description: bookmark.description || ''
    })
    showBookmarkDialog.value = true
  } else if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除这个网址吗？', '确认删除', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })
      await api.delete(`/bookmarks/${bookmark.id}`)
      ElMessage.success('删除成功')
      await loadBookmarks()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('Failed to delete bookmark:', error)
      }
    }
  }
}

const saveCategory = async () => {
  if (!categoryFormRef.value) return
  
  await categoryFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingCategory.value) {
          await api.put(`/categories/${editingCategory.value.id}`, categoryForm)
          ElMessage.success('更新成功')
        } else {
          await api.post('/categories', categoryForm)
          ElMessage.success('添加成功')
        }
        showCategoryDialog.value = false
        resetCategoryForm()
        await loadCategories()
      } catch (error) {
        console.error('Failed to save category:', error)
      } finally {
        saving.value = false
      }
    }
  })
}

const handleCategoryAction = async (command, category) => {
  if (command === 'edit') {
    editingCategory.value = category
    Object.assign(categoryForm, {
      name: category.name,
      description: category.description || ''
    })
    showCategoryDialog.value = true
  } else if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除这个分类吗？分类下的网址也会被删除。', '确认删除', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })
      await api.delete(`/categories/${category.id}`)
      ElMessage.success('删除成功')
      if (selectedCategory.value === category.id) {
        selectedCategory.value = null
      }
      await loadCategories()
      await loadBookmarks()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('Failed to delete category:', error)
      }
    }
  }
}

const incrementVisit = async (bookmark) => {
  try {
    await api.post(`/bookmarks/${bookmark.id}/visit`)
    bookmark.visit_count++
  } catch (error) {
    console.error('Failed to increment visit count:', error)
  }
}

const exportData = async () => {
  try {
    const response = await api.get('/export')
    const data = JSON.stringify(response.data, null, 2)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `nav-backup-${new Date().toISOString().split('T')[0]}.json`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('Failed to export data:', error)
  }
}

const importData = async (file) => {
  try {
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const data = JSON.parse(e.target.result)
        await api.post('/import', data)
        ElMessage.success('导入成功')
        await loadCategories()
        await loadBookmarks()
      } catch (error) {
        ElMessage.error('导入失败：文件格式错误')
      }
    }
    reader.readAsText(file.raw)
  } catch (error) {
    console.error('Failed to import data:', error)
  }
}

const resetBookmarkForm = () => {
  editingBookmark.value = null
  Object.assign(bookmarkForm, {
    title: '',
    url: '',
    category_id: null,
    description: ''
  })
  bookmarkFormRef.value?.resetFields()
}

const resetCategoryForm = () => {
  editingCategory.value = null
  Object.assign(categoryForm, {
    name: '',
    description: ''
  })
  categoryFormRef.value?.resetFields()
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadCategories()
  loadBookmarks()
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
