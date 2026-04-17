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
              <div 
                :ref="el => setCategoryRef(category.id, el)" 
                class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
                :data-category-id="category.id"
              >
                <div
                  v-for="bookmark in getCategoryBookmarks(category.id)"
                  :key="bookmark.id"
                  class="cursor-move"
                  :data-bookmark-id="bookmark.id"
                >
                  <a
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
      </div>
    </main>

    <footer class="bg-white border-t border-gray-200 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex justify-between items-center">
          <div class="flex items-center space-x-3">
              <img v-if="siteSettings.siteLogo" :src="siteSettings.siteLogo" alt="Logo" class="h-8 w-auto" />
              <div class="flex flex-col">
                <p class="text-gray-500 text-sm">
                  © {{ new Date().getFullYear() }} {{ siteSettings.siteTitle }}. All rights reserved.
                </p>
              </div>
            </div>
        </div>
      </div>
    </footer>


  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Sortable from 'sortablejs'

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

// 存储分类的ref
const categoryRefs = ref({})

// 设置分类的ref
const setCategoryRef = (categoryId, el) => {
  if (el) {
    categoryRefs.value[categoryId] = el
  } else {
    delete categoryRefs.value[categoryId]
  }
}

// 初始化SortableJS
const initSortable = () => {
  Object.keys(categoryRefs.value).forEach(categoryId => {
    const container = categoryRefs.value[categoryId]
    if (container) {
      new Sortable(container, {
        animation: 150,
        ghostClass: 'sortable-ghost',
        chosenClass: 'sortable-chosen',
        dragClass: 'sortable-drag',
        group: 'bookmarks',
        handle: '.cursor-move',
        forceFallback: true,
        fallbackClass: 'sortable-fallback',
        fallbackOnBody: true,
        fallbackTolerance: 3,
        scroll: true,
        scrollSensitivity: 30,
        scrollSpeed: 10,
        
        // 开始拖动
        onStart: function(evt) {
          // 检查用户是否已登录
          if (!isLoggedIn()) {
            evt.preventDefault()
            ElMessage.warning('请先登录后再进行排序操作')
            return
          }
          
          // 保存原始排序值
          window.originalSortOrders = {}
          bookmarks.value.forEach(bookmark => {
            window.originalSortOrders[bookmark.id] = {
              sort_order: bookmark.sort_order,
              category_id: bookmark.category_id
            }
          })
          
          // 保存拖动前的书签数组副本
          window.bookmarksBeforeDrag = [...bookmarks.value]
          
          // 添加拖动状态样式
          evt.item.classList.add('opacity-50')
        },
        
        // 结束拖动
        onEnd: async function(evt) {
          // 移除拖动状态样式
          evt.item.classList.remove('opacity-50')
          
          const oldIndex = evt.oldIndex
          const newIndex = evt.newIndex
          const oldCategoryId = parseInt(evt.from.dataset.categoryId)
          const newCategoryId = parseInt(evt.to.dataset.categoryId)
          const bookmarkId = parseInt(evt.item.dataset.bookmarkId)
          
          // 检查是否真的发生了位置变化
          if (oldIndex === newIndex && oldCategoryId === newCategoryId) return
          
          try {
            // 使用拖动前保存的书签数组副本进行操作
            let newBookmarks = window.bookmarksBeforeDrag.map(b => ({...b}))
            
            // 找到被拖动的书签
            const bookmark = newBookmarks.find(b => b.id === bookmarkId)
            if (!bookmark) {
              throw new Error('Bookmark not found')
            }
            
            if (oldCategoryId === newCategoryId) {
              // 同一分类内拖动
              // 获取该分类的所有书签
              let categoryBookmarks = newBookmarks.filter(b => b.category_id === oldCategoryId)
              
              // 在该分类的书签数组中移除被拖动的书签
              categoryBookmarks.splice(oldIndex, 1)
              
              // 在该分类的书签数组中插入到新位置
              categoryBookmarks.splice(newIndex, 0, bookmark)
              
              // 重建newBookmarks数组
              newBookmarks = []
              allCategories.value.forEach(category => {
                if (category.id === oldCategoryId) {
                  // 对于被拖动的分类，使用更新后的分类书签
                  newBookmarks.push(...categoryBookmarks)
                } else {
                  // 对于其他分类，保持原顺序
                  const categoryBookmarks = window.bookmarksBeforeDrag.filter(b => b.category_id === category.id)
                  newBookmarks.push(...categoryBookmarks.map(b => ({...b})))
                }
              })
            } else {
              // 跨分类拖动
              // 更新书签的分类
              bookmark.category_id = newCategoryId
              
              // 获取原分类的所有书签，移除被拖动的书签
              let oldCategoryBookmarks = newBookmarks.filter(b => b.category_id === oldCategoryId)
              oldCategoryBookmarks.splice(oldIndex, 1)
              
              // 获取新分类的所有书签，插入被拖动的书签
              let newCategoryBookmarks = newBookmarks.filter(b => b.category_id === newCategoryId)
              newCategoryBookmarks.splice(newIndex, 0, bookmark)
              
              // 重建newBookmarks数组
              newBookmarks = []
              allCategories.value.forEach(category => {
                if (category.id === oldCategoryId) {
                  // 对于原分类，使用更新后的分类书签
                  newBookmarks.push(...oldCategoryBookmarks)
                } else if (category.id === newCategoryId) {
                  // 对于新分类，使用更新后的分类书签
                  newBookmarks.push(...newCategoryBookmarks)
                } else {
                  // 对于其他分类，保持原顺序
                  const categoryBookmarks = window.bookmarksBeforeDrag.filter(b => b.category_id === category.id)
                  newBookmarks.push(...categoryBookmarks.map(b => ({...b})))
                }
              })
            }
            
            // 更新bookmarks数组
            bookmarks.value = newBookmarks
            
            // 重新计算所有书签的排序值
            const categorySortOrders = {}
            bookmarks.value.forEach(bookmark => {
              if (!categorySortOrders[bookmark.category_id]) {
                categorySortOrders[bookmark.category_id] = 1
              }
              bookmark.sort_order = categorySortOrders[bookmark.category_id]
              categorySortOrders[bookmark.category_id]++
            })
            
            // 保存到后端
            // 只保存排序值或分类发生变化的书签
            const bookmarksToUpdate = bookmarks.value.filter(b => {
              const original = window.originalSortOrders[b.id]
              return original && (b.sort_order !== original.sort_order || b.category_id !== original.category_id)
            })
            
            // 清理临时存储
            delete window.originalSortOrders
            delete window.bookmarksBeforeDrag
            
            if (bookmarksToUpdate.length > 0) {
              await saveBookmarkOrder(bookmarksToUpdate)
            }
            
            ElMessage.success('书签排序已更新')
          } catch (error) {
            console.error('Failed to update bookmark order:', error)
            ElMessage.error('排序更新失败，请稍后重试')
            // 清理临时存储
            delete window.originalSortOrders
            delete window.bookmarksBeforeDrag
            // 重新加载数据恢复原始状态
            await loadPublicData()
          }
        }
      })
    }
  })
}

// 检查用户是否已登录
const isLoggedIn = () => {
  return localStorage.getItem('admin_token') !== null
}

// 保存书签排序
const saveBookmarkOrder = async (bookmarksToUpdate) => {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    throw new Error('请先登录后再进行排序操作')
  }
  
  for (const bookmark of bookmarksToUpdate) {
    try {
      await fetch(`/api/bookmarks/${bookmark.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          title: bookmark.title,
          url: bookmark.url,
          category_id: bookmark.category_id,
          description: bookmark.description,
          sort_order: bookmark.sort_order,
          icon: bookmark.icon
        })
      })
    } catch (error) {
      console.error('Failed to save bookmark order:', error)
      throw error
    }
  }
}

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
    // 初始化SortableJS
    nextTick(() => {
      initSortable()
    })
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

/* 拖放样式 */
.sortable-ghost {
  opacity: 0.5;
  background: #f0f0f0;
  border: 2px dashed #ccc;
  border-radius: 8px;
}

.sortable-chosen {
  background: #e6f7ff;
  border: 2px solid #1890ff;
  border-radius: 8px;
}

.sortable-drag {
  opacity: 0.8;
  background: #e6f7ff;
  border: 2px solid #1890ff;
  border-radius: 8px;
  z-index: 9999;
}

/* 拖动时的光标样式 */
.cursor-move {
  cursor: move;
  user-select: none;
}

/* 拖动时的视觉反馈 */
.cursor-move:active {
  cursor: grabbing;
}
</style>
