<template>
  <div class="min-h-screen bg-gray-100">
    <header class="bg-white shadow-md">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center space-x-3">
            <el-icon :size="32" class="text-blue-600">
              <Link />
            </el-icon>
            <h1 class="text-xl font-bold text-gray-800">后台管理</h1>
          </div>
          <div class="flex items-center space-x-4">
            <el-button @click="goHome" type="primary" plain>
              <el-icon class="mr-1"><House /></el-icon>
              返回前台
            </el-button>
            <el-dropdown @command="handleCommand">
              <el-button>
                <el-icon class="mr-1"><User /></el-icon>
                {{ adminUser?.username || '管理员' }}
                <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="网址管理" name="bookmarks">
          <div class="bg-white rounded-lg shadow-sm p-6">
            <div class="flex justify-between items-center mb-6">
              <div class="flex items-center space-x-4">
                <el-select
                  v-model="selectedCategory"
                  placeholder="选择分类"
                  clearable
                  class="w-48"
                  @change="handleCategoryChange"
                  @clear="handleCategoryClear"
                >
                  <el-option label="全部" value="" />
                  <el-option
                    v-for="category in categories"
                    :key="category.id"
                    :label="category.name"
                    :value="category.id"
                  />
                </el-select>
                <el-input
                  v-model="searchKeyword"
                  placeholder="搜索网址..."
                  prefix-icon="Search"
                  clearable
                  class="w-64"
                />
              </div>
              <div class="flex items-center space-x-2">
                <el-button @click="showBookmarkDialog" type="primary">
                  <el-icon class="mr-1"><Plus /></el-icon>
                  添加网址
                </el-button>
                <el-dropdown @command="handleBatchCommand" :disabled="selectedBookmarks.length === 0">
                  <el-button type="primary">
                    <el-icon class="mr-1"><Operation /></el-icon>
                    批量操作
                    <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-item command="batchDelete">
                      <el-icon><Delete /></el-icon>
                      批量删除
                    </el-dropdown-item>
                    <el-dropdown-item command="batchTest">
                      <el-icon><Connection /></el-icon>
                      批量测试
                    </el-dropdown-item>
                    <el-dropdown-item command="batchDeduplicate">
                      <el-icon><DocumentCopy /></el-icon>
                      批量去重
                    </el-dropdown-item>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <el-table :data="filteredBookmarks" stripe style="width: 100%" @selection-change="handleSelectionChange">
              <el-table-column type="selection" width="55" />
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="url" label="网址" min-width="250">
                <template #default="{ row }">
                  <a :href="row.url" target="_blank" class="text-blue-600 hover:underline">
                    {{ row.url }}
                  </a>
                </template>
              </el-table-column>
              <el-table-column prop="category.name" label="分类" width="120" />
              <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
              <el-table-column prop="visit_count" label="访问次数" width="100" align="center" />
              <el-table-column label="操作" width="180" fixed="right">
                <template #default="{ row }">
                  <el-button @click="editBookmark(row)" type="primary" size="small">
                    编辑
                  </el-button>
                  <el-button @click="deleteBookmark(row)" type="danger" size="small">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="分类管理" name="categories">
          <div class="bg-white rounded-lg shadow-sm p-6">
            <div class="flex justify-between items-center mb-6">
              <h3 class="text-lg font-semibold text-gray-800">分类列表</h3>
              <el-button @click="showCategoryDialog" type="primary">
                <el-icon class="mr-1"><Plus /></el-icon>
                添加分类
              </el-button>
            </div>

            <el-table :data="categories" stripe style="width: 100%">
              <el-table-column prop="name" label="分类名称" width="200">
                <template #default="{ row }">
                  <span v-if="row.is_default" class="text-blue-600 font-semibold">
                    <el-icon class="mr-1"><Star /></el-icon>
                    {{ row.name }}
                  </span>
                  <span v-else>{{ row.name }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="描述" min-width="250" />
              <el-table-column prop="sort_order" label="排序" width="100" align="center" />
              <el-table-column label="网址数量" width="120" align="center">
                <template #default="{ row }">
                  {{ getCategoryBookmarkCount(row.id) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180" fixed="right">
                <template #default="{ row }">
                  <el-button @click="editCategory(row)" type="primary" size="small">
                    编辑
                  </el-button>
                  <el-button 
                    @click="deleteCategory(row)" 
                    type="danger" 
                    size="small"
                    :disabled="row.is_default"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="数据管理" name="data">
          <div class="bg-white rounded-lg shadow-sm p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="border border-gray-200 rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-800 mb-4">
                  <el-icon class="mr-2 text-blue-600"><Download /></el-icon>
                  导出数据
                </h3>
                <p class="text-gray-600 mb-4">
                  将所有分类和网址导出为JSON文件，可用于数据备份或迁移。
                </p>
                <el-button @click="exportData" type="success" class="w-full">
                  <el-icon class="mr-2"><Download /></el-icon>
                  导出数据
                </el-button>
              </div>

              <div class="border border-gray-200 rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-800 mb-4">
                  <el-icon class="mr-2 text-blue-600"><Upload /></el-icon>
                  导入数据
                </h3>
                <p class="text-gray-600 mb-4">
                  从JSON文件导入分类和网址数据。注意：导入会添加到现有数据中。
                </p>
                <el-upload
                  :auto-upload="false"
                  :on-change="importData"
                  :show-file-list="false"
                  accept=".json"
                  class="w-full"
                >
                  <el-button type="warning" class="w-full">
                    <el-icon class="mr-2"><Upload /></el-icon>
                    选择文件导入
                  </el-button>
                </el-upload>
              </div>
            </div>

            <div class="mt-6 border-t border-gray-200 pt-6">
              <h3 class="text-lg font-semibold text-gray-800 mb-4">
                <el-icon class="mr-2 text-red-600"><Warning /></el-icon>
                危险操作
              </h3>
              <p class="text-gray-600 mb-4">
                清空所有数据，此操作不可恢复！
              </p>
              <el-button @click="clearAllData" type="danger">
                <el-icon class="mr-2"><Delete /></el-icon>
                清空所有数据
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="系统设置" name="settings">
          <div class="bg-white rounded-lg shadow-sm p-6">
            <el-form :model="settings" label-width="120px">
              <el-form-item label="网站标题">
                <el-input v-model="settings.siteTitle" placeholder="Van Nav" />
              </el-form-item>
              <el-form-item label="网站描述">
                <el-input
                  v-model="settings.siteDescription"
                  type="textarea"
                  :rows="3"
                  placeholder="个人收藏网址导航"
                />
              </el-form-item>
              <el-form-item label="每页显示">
                <el-input-number v-model="settings.pageSize" :min="10" :max="100" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveSettings">
                  保存设置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </main>

    <el-dialog v-model="bookmarkDialogVisible" :title="editingBookmark ? '编辑网址' : '添加网址'" width="600px">
      <el-form :model="bookmarkForm" :rules="bookmarkRules" ref="bookmarkFormRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="bookmarkForm.title" placeholder="网址标题" />
        </el-form-item>
        <el-form-item label="网址" prop="url">
          <el-input v-model="bookmarkForm.url" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="bookmarkForm.category_id" placeholder="选择分类（可选）" clearable class="w-full">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="bookmarkForm.icon" placeholder="图标URL（可选）" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="bookmarkForm.description"
            type="textarea"
            :rows="3"
            placeholder="网址描述（可选）"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="bookmarkForm.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bookmarkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveBookmark" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="categoryDialogVisible" :title="editingCategory ? '编辑分类' : '添加分类'" width="500px">
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
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCategory" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import api from '@/utils/api'

const router = useRouter()

const activeTab = ref('bookmarks')
const loading = ref(false)
const saving = ref(false)
const batchLoading = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const categories = ref([])
const bookmarks = ref([])
const adminUser = ref(null)
const selectedBookmarks = ref([])

const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const bookmarkDialogVisible = ref(false)
const categoryDialogVisible = ref(false)
const editingBookmark = ref(null)
const editingCategory = ref(null)

const bookmarkFormRef = ref(null)
const categoryFormRef = ref(null)

const bookmarkForm = reactive({
  title: '',
  url: '',
  category_id: null,
  icon: '',
  description: '',
  sort_order: 0
})

const categoryForm = reactive({
  name: '',
  description: '',
  sort_order: 0
})

const settings = reactive({
  siteTitle: 'Van Nav',
  siteDescription: '个人收藏网址导航',
  pageSize: 20
})

const bookmarkRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入网址', trigger: 'blur' },
    { type: 'url', message: '请输入正确的网址格式', trigger: 'blur' }
  ]
}

const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const filteredBookmarks = computed(() => {
  let result = bookmarks.value
  
  if (selectedCategory.value !== null && selectedCategory.value !== undefined && selectedCategory.value !== '') {
    result = result.filter(bookmark => bookmark.category_id === selectedCategory.value)
  }
  
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(bookmark => 
      bookmark.title.toLowerCase().includes(keyword) ||
      bookmark.url.toLowerCase().includes(keyword) ||
      (bookmark.description && bookmark.description.toLowerCase().includes(keyword))
    )
  }
  
  return result
})

const getCategoryBookmarkCount = (categoryId) => {
  return bookmarks.value.filter(b => b.category_id === categoryId).length
}

const showBookmarkDialog = () => {
  editingBookmark.value = null
  resetBookmarkForm()
  bookmarkDialogVisible.value = true
}

const showCategoryDialog = () => {
  editingCategory.value = null
  resetCategoryForm()
  categoryDialogVisible.value = true
}

const editBookmark = (bookmark) => {
  editingBookmark.value = bookmark
  Object.assign(bookmarkForm, {
    title: bookmark.title,
    url: bookmark.url,
    category_id: bookmark.category_id,
    icon: bookmark.icon || '',
    description: bookmark.description || '',
    sort_order: bookmark.sort_order
  })
  bookmarkDialogVisible.value = true
}

const editCategory = (category) => {
  editingCategory.value = category
  Object.assign(categoryForm, {
    name: category.name,
    description: category.description || '',
    sort_order: category.sort_order
  })
  categoryDialogVisible.value = true
}

const handleCategoryChange = () => {
  console.log('Category changed to:', selectedCategory.value)
}

const handleCategoryClear = () => {
  console.log('Category cleared')
}

const handleSelectionChange = (selection) => {
  selectedBookmarks.value = selection
}

const toggleSelectAll = () => {
  if (selectedBookmarks.value.length > 0 && selectedBookmarks.value.length === filteredBookmarks.data.length) {
    selectedBookmarks.value = []
  } else {
    selectedBookmarks.value = [...filteredBookmarks.data]
  }
}

const handleBatchCommand = async (command) => {
  if (command === 'batchDelete') {
    await batchDeleteBookmarks()
  } else if (command === 'batchTest') {
    await batchTestBookmarks()
  } else if (command === 'batchDeduplicate') {
    await batchDeduplicateBookmarks()
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  selectedBookmarks.value = []
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
        bookmarkDialogVisible.value = false
        resetBookmarkForm()
        await loadData()
      } catch (error) {
        console.error('Failed to save bookmark:', error)
      } finally {
        saving.value = false
      }
    }
  })
}

const deleteBookmark = async (bookmark) => {
  try {
    await ElMessageBox.confirm(`确定要删除网址"${bookmark.title}"吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    await api.delete(`/bookmarks/${bookmark.id}`)
    ElMessage.success('删除成功')
    await loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete bookmark:', error)
    }
  }
}

const batchDeleteBookmarks = async () => {
  if (selectedBookmarks.value.length === 0) {
    ElMessage.warning('请先选择要删除的网址')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedBookmarks.value.length} 个网址吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })

    batchLoading.value = true
    let successCount = 0
    let failCount = 0

    for (const bookmark of selectedBookmarks.value) {
      try {
        await api.delete(`/bookmarks/${bookmark.id}`)
        successCount++
      } catch (error) {
        failCount++
        console.error('Failed to delete bookmark:', error)
      }
    }

    batchLoading.value = false
    selectedBookmarks.value = []

    ElNotification({
      title: '批量删除完成',
      message: `成功删除 ${successCount} 个网址，失败 ${failCount} 个`,
      type: 'success',
      duration: 3000
    })

    await loadData()
  } catch (error) {
    batchLoading.value = false
    console.error('Batch delete failed:', error)
  }
}

const batchTestBookmarks = async () => {
  if (selectedBookmarks.value.length === 0) {
    ElMessage.warning('请先选择要测试的网址')
    return
  }

  batchLoading.value = true
  let successCount = 0
  let failCount = 0
  const results = []

  for (const bookmark of selectedBookmarks.value) {
    try {
      const response = await fetch(bookmark.url, {
        method: 'HEAD',
        mode: 'no-cors'
      })
      results.push({
        id: bookmark.id,
        title: bookmark.title,
        url: bookmark.url,
        status: response.ok ? 'success' : 'failed',
        statusCode: response.status
      })
      if (response.ok) {
        successCount++
      } else {
        failCount++
      }
    } catch (error) {
      failCount++
      results.push({
        id: bookmark.id,
        title: bookmark.title,
        url: bookmark.url,
        status: 'error',
        error: error.message
      })
    }
  }
  batchLoading.value = false
  selectedBookmarks.value = []

  ElMessage.success(`测试完成：成功 ${successCount} 个，失败 ${failCount} 个`)
  
  console.table(results)
}

const batchDeduplicateBookmarks = async () => {
  if (selectedBookmarks.value.length === 0) {
    ElMessage.warning('请先选择要去重的网址')
    return
  }

  const urlMap = new Map()
  const duplicates = []
  const uniqueBookmarks = []
  const toDelete = []

  for (const bookmark of selectedBookmarks.value) {
    const normalizedUrl = normalizeUrl(bookmark.url)
    
    if (urlMap.has(normalizedUrl)) {
      const existing = urlMap.get(normalizedUrl)
      if (existing.created_at < bookmark.created_at) {
        toDelete.push(bookmark.id)
      } else {
        toDelete.push(existing.id)
        urlMap.set(normalizedUrl, bookmark)
        uniqueBookmarks.push(bookmark)
      }
    } else {
      urlMap.set(normalizedUrl, bookmark)
      uniqueBookmarks.push(bookmark)
    }
  }

  if (toDelete.length === 0) {
    ElMessage.info('没有发现重复网址')
    return
  }

  try {
    await ElMessageBox.confirm(`发现 ${toDelete.length} 个重复网址，保留创建时间最早的记录，确定要去重吗？`, '确认去重', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })

    batchLoading.value = true
    let deletedCount = 0
    let processedCount = 0
    const total = toDelete.length

    for (const id of toDelete) {
      try {
        await api.delete(`/bookmarks/${id}`)
        deletedCount++
        processedCount++
        
        if (processedCount % 5 === 0) {
          ElMessage.info(`正在处理 ${processedCount}/${total} 条数据...`)
        }
      } catch (error) {
        console.error('Failed to delete duplicate:', error)
      }
    }

    batchLoading.value = false
    selectedBookmarks.value = []
    isAllSelected.value = false

    ElNotification({
      title: '批量去重完成',
      message: `成功删除 ${deletedCount} 个重复网址`,
      type: 'success',
      duration: 3000
    })

    await loadData()
  } catch (error) {
    batchLoading.value = false
    console.error('Batch deduplicate failed:', error)
  }
}

function normalizeUrl(url) {
  try {
    const urlObj = new URL(url)
    let normalized = urlObj.protocol + '//' + urlObj.hostname + urlObj.pathname
    const query = urlObj.search
    if (query) {
      const params = new URLSearchParams(query)
      params.sort()
      normalized += '?' + params.toString()
    }
    return normalized.toLowerCase()
  } catch {
    return url.toLowerCase()
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
        categoryDialogVisible.value = false
        resetCategoryForm()
        await loadData()
      } catch (error) {
        console.error('Failed to save category:', error)
      } finally {
        saving.value = false
      }
    }
  })
}

const deleteCategory = async (category) => {
  try {
    const message = category.is_default 
      ? '默认分类"未分类"不能删除' 
      : `确定要删除分类"${category.name}"吗？该分类下的网址将自动迁移到"未分类"。`
    
    await ElMessageBox.confirm(message, '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    await api.delete(`/categories/${category.id}`)
    ElMessage.success('删除成功')
    
    if (selectedBookmarks.value.length > 0) {
      selectedBookmarks.value = []
      ElNotification({
        title: '提示',
        message: `分类"${category.name}"已删除，选中的网址已自动取消`,
        type: 'info'
      })
    }
    
    await loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete category:', error)
    }
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
    ElMessage.error('导出失败')
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
        await loadData()
      } catch (error) {
        ElMessage.error('导入失败：文件格式错误')
      }
    }
    reader.readAsText(file.raw)
  } catch (error) {
    console.error('Failed to import data:', error)
  }
}

const clearAllData = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有数据吗？此操作不可恢复！', '危险操作', {
      type: 'error',
      confirmButtonText: '确定清空',
      cancelButtonText: '取消'
    })
    await api.delete('/admin/clear-all')
    ElMessage.success('清空成功')
    await loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to clear data:', error)
    }
  }
}

const saveSettings = () => {
  localStorage.setItem('site_settings', JSON.stringify(settings))
  ElMessage.success('设置已保存')
}

const resetBookmarkForm = () => {
  Object.assign(bookmarkForm, {
    title: '',
    url: '',
    category_id: null,
    icon: '',
    description: '',
    sort_order: 0
  })
  bookmarkFormRef.value?.resetFields()
}

const resetCategoryForm = () => {
  Object.assign(categoryForm, {
    name: '',
    description: '',
    sort_order: 0
  })
  categoryFormRef.value?.resetFields()
}

const loadData = async () => {
  loading.value = true
  try {
    const [categoriesRes, bookmarksRes] = await Promise.all([
      api.get('/categories'),
      api.get('/bookmarks')
    ])
    categories.value = categoriesRes
    bookmarks.value = bookmarksRes
    total.value = bookmarksRes.length
    currentPage.value = 1
  } catch (error) {
    console.error('Failed to load data:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const goHome = () => {
  router.push('/')
}

const handleCommand = (command) => {
  if (command === 'logout') {
    localStorage.removeItem('admin_token')
    adminUser.value = null
    router.push('/')
  }
}

onMounted(() => {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    router.push('/admin/login')
    return
  }
  
  const userData = localStorage.getItem('admin_user')
  if (userData) {
    adminUser.value = JSON.parse(userData)
  }
  
  const settingsData = localStorage.getItem('site_settings')
  if (settingsData) {
    Object.assign(settings, JSON.parse(settingsData))
  }
  
  loadData()
})
</script>
