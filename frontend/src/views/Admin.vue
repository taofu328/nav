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
                    <el-dropdown-item command="batchFetchIcons">
                      <el-icon><Picture /></el-icon>
                      批量获取图标
                    </el-dropdown-item>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <el-table :data="paginatedBookmarks" stripe style="width: 100%" @selection-change="handleSelectionChange">
              <el-table-column type="selection" width="55" />
              <el-table-column prop="title" label="标题" width="150" />
              <el-table-column label="图标" width="80" align="center">
                <template #default="{ row }">
                  <img
                    :src="getIconUrl(row)"
                    :alt="row.title"
                    class="w-8 h-8 rounded object-cover"
                    @error="handleIconError"
                  />
                </template>
              </el-table-column>
              <el-table-column prop="url" label="网址" min-width="200">
                <template #default="{ row }">
                  <a :href="row.url" target="_blank" class="text-blue-600 hover:underline">
                    {{ row.url }}
                  </a>
                </template>
              </el-table-column>
              <el-table-column prop="category.name" label="分类" width="120" />
              <el-table-column prop="sort_order" label="排序" width="80" align="center" />
              <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />

              <el-table-column label="操作" width="180">
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

            <div class="flex justify-between items-center mt-4">
              <div class="text-sm text-gray-600">
                共 {{ totalItems }} 条记录，第 {{ currentPage }} / {{ totalPages }} 页
              </div>
              <div class="flex items-center space-x-2">
                <el-button
                  @click="currentPage = currentPage > 1 ? currentPage - 1 : 1"
                  :disabled="currentPage === 1"
                  size="small"
                >
                  上一页
                </el-button>
                <el-button
                  @click="currentPage = currentPage < totalPages ? currentPage + 1 : totalPages"
                  :disabled="currentPage === totalPages"
                  size="small"
                >
                  下一页
                </el-button>
                <el-select v-model="pageSize" size="small" @change="currentPage = 1">
                  <el-option label="10条/页" value="10" />
                  <el-option label="20条/页" value="20" />
                  <el-option label="50条/页" value="50" />
                </el-select>
              </div>
            </div>
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

            <el-table 
              ref="categoryTableRef"
              :data="categories" 
              stripe 
              style="width: 100%"
              row-key="id"
            >
              <el-table-column width="50" align="center">
                <template #default="{ row }">
                  <el-icon 
                    :class="['cursor-move', row.is_default ? 'text-gray-400' : 'text-blue-600']"
                    :disabled="row.is_default"
                  >
                    <Operation />
                  </el-icon>
                </template>
              </el-table-column>
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
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div class="border border-gray-200 rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-800 mb-4">
                  <el-icon class="mr-2 text-blue-600"><User /></el-icon>
                  修改用户信息
                </h3>
                <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="100px">
                  <el-form-item label="用户名" prop="username">
                    <el-input v-model="userForm.username" placeholder="请输入用户名" />
                  </el-form-item>
                  <el-form-item label="新密码" prop="newPassword">
                    <el-input
                      v-model="userForm.newPassword"
                      type="password"
                      placeholder="留空则不修改密码"
                      show-password
                    />
                  </el-form-item>
                  <el-form-item label="确认密码" prop="confirmPassword">
                    <el-input
                      v-model="userForm.confirmPassword"
                      type="password"
                      placeholder="请再次输入新密码"
                      show-password
                    />
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="updateUserInfo" :loading="updatingUser">
                      保存修改
                    </el-button>
                    <el-button @click="resetUserForm">取消</el-button>
                  </el-form-item>
                </el-form>
              </div>

              <div class="border border-gray-200 rounded-lg p-6">
                <h3 class="text-lg font-semibold text-gray-800 mb-4">
                  <el-icon class="mr-2 text-blue-600"><Setting /></el-icon>
                  修改网站信息
                </h3>
                <el-form :model="siteSettings" :rules="siteRules" ref="siteFormRef" label-width="100px">
                  <el-form-item label="网站标题" prop="siteTitle">
                    <el-input v-model="siteSettings.siteTitle" placeholder="请输入网站标题" />
                  </el-form-item>

                  <el-form-item label="网站Logo">
                    <div class="flex space-x-2 items-center">
                      <el-upload
                        class="upload-demo"
                        :action="'/api/icons/upload'"
                        :headers="{ Authorization: `Bearer ${adminToken}` }"
                        :on-success="handleLogoUpload"
                        :on-error="handleLogoUploadError"
                        :show-file-list="false"
                        accept=".png,.jpg,.jpeg,.svg,.ico"
                        name="icon"
                      >
                        <el-button type="primary">
                          <el-icon class="mr-1"><Upload /></el-icon>
                          上传Logo
                        </el-button>
                      </el-upload>
                      <el-button @click="removeLogo" v-if="siteSettings.siteLogo" type="danger">
                        <el-icon class="mr-1"><Delete /></el-icon>
                        删除Logo
                      </el-button>
                      <div v-if="siteSettings.siteLogo" class="flex-shrink-0">
                        <el-image 
                          :src="siteSettings.siteLogo" 
                          fit="cover" 
                          :preview-src-list="[siteSettings.siteLogo]" 
                          style="height: 32px; width: auto; max-width: 100px; border-radius: 4px;" 
                        />
                      </div>
                    </div>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="updateSiteSettings" :loading="updatingSettings">
                      保存设置
                    </el-button>
                  </el-form-item>
                </el-form>
              </div>
            </div>
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
          <div class="flex space-x-2">
            <el-button @click="fetchIcon" type="primary" :loading="fetchingIcon">
              <el-icon class="mr-1"><Connection /></el-icon>
              重新获取
            </el-button>
            <el-upload
              class="upload-demo"
              :action="'/api/icons/upload'"
              :headers="{ Authorization: `Bearer ${adminToken}` }"
              :on-success="handleIconUpload"
              :on-error="handleIconUploadError"
              :show-file-list="false"
              accept=".png,.jpg,.jpeg,.svg,.ico"
            >
              <el-button type="success">
                <el-icon class="mr-1"><Upload /></el-icon>
                上传图标
              </el-button>
            </el-upload>
          </div>
          <div v-if="bookmarkForm.icon" class="mt-2">
            <el-image :src="bookmarkForm.icon" fit="cover" :preview-src-list="[bookmarkForm.icon]" style="width: 48px; height: 48px;" />
          </div>
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

    <el-dialog v-model="batchFetchIconLoading" title="批量获取图标" width="600px" :close-on-click-modal="false" :close-on-press-escape="false">
      <div class="batch-progress-container">
        <div class="progress-info">
          <el-progress :percentage="Math.round((batchFetchIconProgress / batchFetchIconTotal) * 100)" :status="batchFetchIconCancelled ? 'exception' : 'success'" />
          <div class="progress-text">
            <span class="current">{{ batchFetchIconProgress }} / {{ batchFetchIconTotal }}</span>
            <span class="percentage">{{ Math.round((batchFetchIconProgress / batchFetchIconTotal) * 100) }}%</span>
          </div>
        </div>
        
        <div v-if="batchFetchIconResults.length > 0" class="results-list">
          <div class="results-header">
            <span>获取结果</span>
            <el-button size="small" @click="batchFetchIconCancelled = true" :disabled="!batchFetchIconLoading">
              取消操作
            </el-button>
          </div>
          <div class="results-content">
            <div
              v-for="result in batchFetchIconResults"
              :key="result.id"
              class="result-item"
              :class="{ success: result.status === 'success', failed: result.status === 'error' }"
            >
              <div class="result-title">{{ result.title }}</div>
              <div class="result-url">{{ result.url }}</div>
              <div class="result-status">
                <el-icon v-if="result.status === 'success'" class="success-icon"><CircleCheck /></el-icon>
                <el-icon v-else class="error-icon"><CircleClose /></el-icon>
                <span :class="result.status">{{ result.status === 'success' ? '成功' : '失败' }}</span>
              </div>
              <div v-if="result.error" class="result-error">{{ result.error }}</div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import api from '@/utils/api'
import Sortable from 'sortablejs'

const router = useRouter()

const activeTab = ref('bookmarks')
const loading = ref(false)
const saving = ref(false)
const batchLoading = ref(false)
const fetchingIcon = ref(false)
const batchFetchIconLoading = ref(false)
const batchFetchIconProgress = ref(0)
const batchFetchIconTotal = ref(0)
const batchFetchIconResults = ref([])
const batchFetchIconCancelled = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const categories = ref([])
const bookmarks = ref([])
const adminUser = ref(null)
const selectedBookmarks = ref([])
const adminToken = ref('')
const categoryTableRef = ref(null)

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
  sort_order: 99
})

const categoryForm = reactive({
  name: '',
  description: '',
  sort_order: 99
})

const settings = reactive({
  siteTitle: 'Van Nav',
  siteDescription: '个人收藏网址导航',
  siteLogo: '',
  pageSize: 20
})

const userForm = reactive({
  username: '',
  newPassword: '',
  confirmPassword: ''
})

const siteSettings = reactive({
  siteTitle: 'Van Nav',
  siteLogo: ''
})

const userFormRef = ref(null)
const siteFormRef = ref(null)
const updatingUser = ref(false)
const updatingSettings = ref(false)

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

const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  newPassword: [
    { min: 6, message: '密码长度至少为 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { 
      validator: (rule, value, callback) => {
        if (userForm.newPassword && value !== userForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const siteRules = {
  siteTitle: [
    { required: true, message: '请输入网站标题', trigger: 'blur' }
  ]
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

const paginatedBookmarks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredBookmarks.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredBookmarks.value.length / pageSize.value)
})

const totalItems = computed(() => {
  return filteredBookmarks.value.length
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
  // 筛选时重置到第一页
  currentPage.value = 1
}

const handleCategoryClear = () => {
  console.log('Category cleared')
  // 清除筛选时重置到第一页
  currentPage.value = 1
}

// 监听搜索关键词变化，重置到第一页
watch(searchKeyword, () => {
  currentPage.value = 1
})

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
  } else if (command === 'batchFetchIcons') {
    await batchFetchIcons()
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

const batchFetchIcons = async () => {
  if (selectedBookmarks.value.length === 0) {
    ElMessage.warning('请先选择要获取图标的网址')
    return
  }

  batchFetchIconLoading.value = true
  batchFetchIconProgress.value = 0
  batchFetchIconTotal.value = selectedBookmarks.value.length
  batchFetchIconResults.value = []
  batchFetchIconCancelled.value = false

  // 创建一个固定长度的数组副本，避免循环过程中数组长度变化
  const bookmarksToProcess = [...selectedBookmarks.value]
  const total = bookmarksToProcess.length

  let successCount = 0
  let failCount = 0

  for (let i = 0; i < total; i++) {
    if (batchFetchIconCancelled.value) {
      ElMessage.info('批量获取图标已取消')
      break
    }

    const bookmark = bookmarksToProcess[i]
    try {
      const response = await api.get(`/icons?url=${encodeURIComponent(bookmark.url)}`)
      if (response.icon) {
        const oldIcon = bookmark.icon
        
        if (oldIcon && oldIcon !== response.icon) {
          try {
            await api.delete(`/icons?url=${encodeURIComponent(oldIcon)}`)
          } catch (error) {
            console.error('Failed to delete old icon:', error)
          }
        }
        
        await api.put(`/bookmarks/${bookmark.id}`, { 
          title: bookmark.title,
          url: bookmark.url,
          category_id: bookmark.category_id,
          description: bookmark.description,
          sort_order: bookmark.sort_order,
          icon: response.icon 
        })
        successCount++
        batchFetchIconResults.value.push({
          id: bookmark.id,
          title: bookmark.title,
          url: bookmark.url,
          status: 'success'
        })
      } else {
        failCount++
        batchFetchIconResults.value.push({
          id: bookmark.id,
          title: bookmark.title,
          url: bookmark.url,
          status: 'failed',
          error: 'Failed to fetch icon'
        })
      }
    } catch (error) {
      failCount++
      batchFetchIconResults.value.push({
        id: bookmark.id,
        title: bookmark.title,
        url: bookmark.url,
        status: 'error',
        error: error.message
      })
    }

    batchFetchIconProgress.value = i + 1
    
    if ((i + 1) % 5 === 0) {
      await new Promise(resolve => setTimeout(resolve, 1000))
    }
  }

  batchFetchIconLoading.value = false
  selectedBookmarks.value = []

  ElNotification({
    title: '批量获取图标完成',
    message: `成功获取 ${successCount} 个图标，失败 ${failCount} 个`,
    type: 'success',
    duration: 3000
  })

  await loadData()
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

    const data = await api.get('/export')


    if (!data || typeof data !== 'object') {
      console.error('响应数据无效:', data)
      ElMessage.error('导出失败：服务器返回的数据格式不正确')
      return
    }

    
    // 兼容大小写字段名
    const categories = data.Categories || data.categories || []
    const bookmarks = data.Bookmarks || data.bookmarks || []
    
    if (categories.length === 0 && bookmarks.length === 0) {
      console.error('响应数据为空:', data)
      ElMessage.error('导出失败：服务器返回的数据为空')
      return
    }
    

    
    // 转换为正确的字段名
    const exportData = {
      Categories: categories,
      Bookmarks: bookmarks
    }
    
    const jsonData = JSON.stringify(exportData, null, 2)

    const blob = new Blob([jsonData], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `nav-backup-${new Date().toISOString().split('T')[0]}.json`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出数据异常:', error)
    ElMessage.error('导出失败：' + (error.message || '未知错误'))
  }
}

const importData = async (file) => {
  try {
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const data = JSON.parse(e.target.result)
        
        // 构建导入请求，默认使用覆盖策略
        const importRequest = {
          data: data,
          conflict: 'overwrite'
        }
        
        const response = await api.post('/import', importRequest)
        ElMessage.success(`导入成功！导入了 ${response.categories} 个分类和 ${response.bookmarks} 个书签`)
        await loadData()
      } catch (error) {
        console.error('导入失败:', error)
        ElMessage.error('导入失败：' + (error.message || '文件格式错误'))
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

const updateUserInfo = async () => {
  if (!userFormRef.value) return
  
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      updatingUser.value = true
      try {
        const response = await api.put('/admin/user', {
          username: userForm.username,
          new_password: userForm.newPassword
        })
        
        if (response.user) {
          localStorage.setItem('admin_user', JSON.stringify(response.user))
          adminUser.value = response.user
        }
        
        ElMessage.success('用户信息更新成功')
        resetUserForm()
      } catch (error) {
        console.error('Failed to update user info:', error)
        ElMessage.error('更新失败：' + (error.message || '未知错误'))
      } finally {
        updatingUser.value = false
      }
    }
  })
}

const updateSiteSettings = async () => {
  if (!siteFormRef.value) return
  
  await siteFormRef.value.validate(async (valid) => {
    if (valid) {
      updatingSettings.value = true
      try {
        const response = await api.put('/admin/settings', {
          site_title: siteSettings.siteTitle,
          site_logo: siteSettings.siteLogo
        })
        
        localStorage.setItem('site_settings', JSON.stringify({
          siteTitle: siteSettings.siteTitle,
          siteLogo: siteSettings.siteLogo
        }))
        
        // 更新页面标题
        if (siteSettings.siteTitle) {
          document.title = siteSettings.siteTitle
        }
        
        ElMessage.success('网站设置更新成功')
      } catch (error) {
        console.error('Failed to update site settings:', error)
        ElMessage.error('更新失败：' + (error.message || '未知错误'))
      } finally {
        updatingSettings.value = false
      }
    }
  })
}

const resetUserForm = () => {
  if (adminUser.value) {
    userForm.username = adminUser.value.username
  }
  userForm.newPassword = ''
  userForm.confirmPassword = ''
  userFormRef.value?.resetFields()
}

const handleLogoUpload = (response) => {
  if (response.icon) {
    siteSettings.siteLogo = response.icon + '?t=' + Date.now()
    // 更新favicon
    const favicon = document.getElementById('favicon')
    if (favicon) {
      favicon.href = siteSettings.siteLogo
    }
    ElMessage.success('Logo上传成功')
  } else {
    ElMessage.error('Logo上传失败，请稍后重试')
  }
}

const handleLogoUploadError = (error) => {
  console.error('Logo upload error:', error)
  ElMessage.error('Logo上传失败，请稍后重试')
}

const removeLogo = async () => {
  try {
    await ElMessageBox.confirm('确定要删除网站Logo吗？', '确认删除', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    
    if (siteSettings.siteLogo) {
      try {
        await api.delete(`/icons?url=${encodeURIComponent(siteSettings.siteLogo)}`)
      } catch (error) {
        console.error('Failed to delete logo:', error)
      }
    }
    
    siteSettings.siteLogo = ''
    // 恢复默认favicon
    const favicon = document.getElementById('favicon')
    if (favicon) {
      favicon.href = '/api/icons/default.svg'
    }
    ElMessage.success('Logo删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove logo:', error)
    }
  }
}

const resetBookmarkForm = () => {
  Object.assign(bookmarkForm, {
    title: '',
    url: '',
    category_id: null,
    icon: '',
    description: '',
    sort_order: 99
  })
  bookmarkFormRef.value?.resetFields()
}

const resetCategoryForm = () => {
  Object.assign(categoryForm, {
    name: '',
    description: '',
    sort_order: 99
  })
  categoryFormRef.value?.resetFields()
}

const fetchIcon = async () => {
  if (!bookmarkForm.url) {
    ElMessage.warning('请先输入网址')
    return
  }
  
  // 保存旧的图标路径
  const oldIcon = bookmarkForm.icon
  
  fetchingIcon.value = true
  try {
    const response = await api.get(`/icons?url=${encodeURIComponent(bookmarkForm.url)}`)
    if (response.icon) {
      // 添加时间戳参数，避免浏览器缓存旧图标
      bookmarkForm.icon = response.icon + '?t=' + Date.now()
      ElMessage.success('图标获取成功')
      
      // 如果有旧图标且与新图标不同，删除旧图标
      if (oldIcon && oldIcon !== response.icon) {
        try {
          await api.delete(`/icons?url=${encodeURIComponent(oldIcon)}`)
        } catch (error) {
          console.error('Failed to delete old icon:', error)
        }
      }
    } else {
      ElMessage.error('图标获取失败，请稍后重试')
    }
  } catch (error) {
    console.error('Failed to fetch icon:', error)
    ElMessage.error('图标获取失败，请稍后重试')
  } finally {
    fetchingIcon.value = false
  }
}

const handleIconUpload = (response) => {
  if (response.icon) {
    // 添加时间戳参数，避免浏览器缓存旧图标
    bookmarkForm.icon = response.icon + '?t=' + Date.now()
    ElMessage.success('图标上传成功')
  } else {
    ElMessage.error('图标上传失败，请稍后重试')
  }
}

const handleIconUploadError = (error) => {
  console.error('Icon upload error:', error)
  ElMessage.error('图标上传失败，请稍后重试')
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

// 初始化SortableJS
const initSortable = () => {
  if (!categoryTableRef.value) return
  
  const table = categoryTableRef.value
  const tbody = table.$el.querySelector('.el-table__body-wrapper tbody')
  
  if (!tbody) return
  
  // 销毁之前的Sortable实例
  if (window.categorySortable) {
    window.categorySortable.destroy()
  }
  
  // 保存原始排序信息
  const originalSortOrder = new Map()
  
  // 初始化SortableJS
  window.categorySortable = Sortable.create(tbody, {
    animation: 150, // 动画持续时间
    ghostClass: 'sortable-ghost', // 拖动时的样式
    chosenClass: 'sortable-chosen', // 选中时的样式
    dragClass: 'sortable-drag', // 拖动中的样式
    handle: '.cursor-move', // 拖动手柄
    forceFallback: true, // 强制使用fallback模式
    fallbackClass: 'sortable-fallback', // fallback模式的样式
    fallbackOnBody: true, // fallback元素添加到body
    fallbackTolerance: 3, // 鼠标移动多少像素才开始拖动
    scroll: true, // 允许滚动
    scrollSensitivity: 30, // 滚动灵敏度
    scrollSpeed: 10, // 滚动速度
    
    // 开始拖动
    onStart: function(evt) {
      const draggedRow = categories.value[evt.oldIndex]
      // 禁止拖动默认分类
      if (draggedRow.is_default) {
        evt.preventDefault()
        ElMessage.warning('默认分类不可拖动')
        return
      }
      // 保存原始排序信息
      categories.value.forEach(category => {
        originalSortOrder.set(category.id, category.sort_order)
      })
      // 添加拖动状态样式
      evt.item.classList.add('opacity-50')
    },
    
    // 结束拖动
    onEnd: async function(evt) {
      // 移除拖动状态样式
      evt.item.classList.remove('opacity-50')
      
      const oldIndex = evt.oldIndex
      const newIndex = evt.newIndex
      
      // 检查是否真的发生了位置变化
      if (oldIndex === newIndex) return
      
      try {
        // 获取拖动的分类和目标位置的分类
        const draggedCategory = categories.value[oldIndex]
        const targetCategory = categories.value[newIndex]
        
        // 检查目标分类是否是默认分类
        if (targetCategory.is_default) {
          // 恢复原始位置
          categories.value.splice(oldIndex, 0, categories.value.splice(newIndex, 1)[0])
          ElMessage.warning('默认分类不可作为目标位置')
          return
        }
        
        // 交换位置
        categories.value.splice(newIndex, 0, categories.value.splice(oldIndex, 1)[0])
        
        // 重新计算排序值
        categories.value.forEach((category, index) => {
          category.sort_order = index + 1
        })
        
        // 找出排序值发生变化的分类
        const changedCategories = categories.value.filter(category => {
          return originalSortOrder.get(category.id) !== category.sort_order
        })
        
        // 保存到后端（只更新排序值发生变化的分类）
        for (const category of changedCategories) {
          await api.put(`/categories/${category.id}`, { sort_order: category.sort_order })
        }
        
        ElMessage.success('分类排序已更新')
      } catch (error) {
        console.error('Failed to update category sort order:', error)
        ElMessage.error('排序更新失败，请稍后重试')
        // 重新加载数据恢复原始状态
        await loadData()
      }
    }
  })
}

// 拖动排序相关方法
const draggingRow = ref(null)

const handleRowDragStart = (row, column, event) => {
  console.log("handleRowDragStart")
  // 禁止拖动默认分类
  if (row.is_default) {
    event.preventDefault()
    return
  }
  draggingRow.value = row
  event.target.classList.add('opacity-50')
}

const handleRowDragEnd = (row, column, event) => {
  draggingRow.value = null
  event.target.classList.remove('opacity-50')
}

const handleRowDrop = async (draggedRow, targetRow, column, event) => {
  console.log("handleRowDrop")
  // 禁止拖动默认分类或拖动到默认分类
  if (draggedRow.is_default || targetRow.is_default) {
    ElMessage.warning('默认分类不可拖动')
    return
  }

  try {
    // 交换排序值
    const tempSort = draggedRow.sort_order
    draggedRow.sort_order = targetRow.sort_order
    targetRow.sort_order = tempSort

    // 保存到后端
    await Promise.all([
      api.put(`/categories/${draggedRow.id}`, { sort_order: draggedRow.sort_order }),
      api.put(`/categories/${targetRow.id}`, { sort_order: targetRow.sort_order })
    ])

    // 重新加载数据以确保排序正确
    await loadData()
    ElMessage.success('分类排序已更新')
  } catch (error) {
    console.error('Failed to update category sort order:', error)
    ElMessage.error('排序更新失败，请稍后重试')
    // 重新加载数据恢复原始状态
    await loadData()
  }
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

const loadData = async () => {
  loading.value = true
  try {
    const [categoriesRes, bookmarksRes] = await Promise.all([
      api.get('/categories'),
      api.get('/bookmarks')
    ])
    categories.value = categoriesRes
    bookmarks.value = bookmarksRes
  } catch (error) {
    console.error('Failed to load data:', error)
  } finally {
    loading.value = false
    // 初始化SortableJS
    nextTick(() => {
      initSortable()
    })
  }
}

onMounted(async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    router.push('/admin/login')
    return
  }
  adminToken.value = token
  
  const userData = localStorage.getItem('admin_user')
  if (userData) {
    adminUser.value = JSON.parse(userData)
    userForm.username = adminUser.value.username
  }
  
  // 从后端 API 获取网站设置
  try {
    const response = await api.get('/admin/settings')
    if (response.settings) {
      const savedSettings = response.settings
      Object.assign(settings, {
        siteTitle: savedSettings.site_title || 'Van Nav',
        siteDescription: savedSettings.site_description || '个人收藏网址导航',
        siteLogo: savedSettings.site_logo || '',
        pageSize: savedSettings.page_size || 20
      })
      Object.assign(siteSettings, {
        siteTitle: savedSettings.site_title || 'Van Nav',
        siteLogo: savedSettings.site_logo || ''
      })
      
      // 将设置保存到 localStorage 以便快速加载
      localStorage.setItem('site_settings', JSON.stringify({
        siteTitle: savedSettings.site_title || 'Van Nav',
        siteLogo: savedSettings.site_logo || ''
      }))
      
      // 更新页面标题
      if (siteSettings.siteTitle) {
        document.title = siteSettings.siteTitle
      }
    }
  } catch (error) {
    console.error('Failed to load site settings:', error)
  }
  
  await loadData()
})
</script>

<style scoped>
/* 动画效果 */
.sortable-ghost {
  opacity: 0.5;
  background: #f0f0f0;
}

.sortable-chosen {
  background: #e6f7ff;
}

.sortable-drag {
  opacity: 0.8;
  background: #e6f7ff;
}

.batch-progress-container {
  padding: 20px;
}

.progress-info {
  margin-bottom: 20px;
}

.progress-text {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

.results-list {
  border-top: 1px solid #e8e8e8;
  padding-top: 20px;
  margin-top: 20px;
  max-height: 300px;
  overflow-y: auto;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: bold;
}

.result-item {
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: 1px solid #e8e8e8;
}

.result-item.success {
  border-color: #52c41a;
  background-color: #f6ffed;
}

.result-item.failed {
  border-color: #ff4d4f;
  background-color: #fff1f0;
}

.result-title {
  font-weight: bold;
  margin-bottom: 5px;
}

.result-url {
  font-size: 12px;
  color: #666;
  margin-bottom: 5px;
  word-break: break-all;
}

.result-status {
  display: flex;
  align-items: center;
  font-size: 12px;
  margin-bottom: 5px;
}

.result-status .success {
  color: #52c41a;
}

.result-status .error {
  color: #ff4d4f;
}

.success-icon {
  color: #52c41a;
  margin-right: 5px;
}

.error-icon {
  color: #ff4d4f;
  margin-right: 5px;
}

.result-error {
  font-size: 12px;
  color: #ff4d4f;
  margin-top: 5px;
  word-break: break-all;
}
</style>