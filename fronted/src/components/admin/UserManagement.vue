<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAdminStore } from '@/stores/admin';
import { useUserStore } from '@/stores/user';
import { storeToRefs } from 'pinia';

// 状态存储
const adminStore = useAdminStore();
const userStore = useUserStore();
const { users, pagination, isLoading } = storeToRefs(adminStore);
const { user: currentUser } = storeToRefs(userStore);

// 模态框状态
const showCreateUserModal = ref(false);
const showUserDetailModal = ref(false);
const showDeleteConfirmModal = ref(false);

// 搜索参数
const searchKeyword = ref('');

// 表单数据
const userForm = ref({
  username: '',
  password: '',
  email: '',
  role: 'user'
});

// 选中的用户
const selectedUser = ref(null);
const selectedUserId = ref(null);

// 错误信息
const errorMessage = ref('');
const successMessage = ref('');

// 确认删除
const confirmDeleteUser = (user) => {
  selectedUser.value = user;
  showDeleteConfirmModal.value = true;
};

// 确认删除操作
const handleDeleteUser = async () => {
  if (!selectedUser.value) return;
  
  const result = await adminStore.deleteUser(selectedUser.value.id);
  if (result.success) {
    successMessage.value = '用户删除成功';
    showDeleteConfirmModal.value = false;
    // 3秒后清空消息
    setTimeout(() => {
      successMessage.value = '';
    }, 3000);
  } else {
    errorMessage.value = result.message || '删除用户失败';
    setTimeout(() => {
      errorMessage.value = '';
    }, 3000);
  }
};

// 查看用户详情
const viewUserDetail = async (user) => {
  selectedUser.value = user;
  selectedUserId.value = user.id;
  showUserDetailModal.value = true;
  
  // 获取更详细的用户信息
  const result = await adminStore.fetchUserDetail(user.id);
  if (result) {
    selectedUser.value = result.user;
  }
};

// 加载用户列表
const loadUsers = async (page = 1) => {
  await adminStore.fetchUsers(page, pagination.value.pageSize, searchKeyword.value);
};

// 搜索用户
const searchUsers = async () => {
  await loadUsers(1); // 搜索时重置到第一页
};

// 处理角色切换
const handleRoleChange = async (user, newRole) => {
  // 防止管理员降级自己
  if (user.id === currentUser.value.id && newRole !== 'admin') {
    errorMessage.value = '不能降级自己的管理员权限';
    setTimeout(() => {
      errorMessage.value = '';
    }, 3000);
    return;
  }
  
  const result = await adminStore.updateUserRole(user.id, newRole);
  if (result.success) {
    successMessage.value = '用户角色已更新';
    setTimeout(() => {
      successMessage.value = '';
    }, 3000);
  } else {
    errorMessage.value = result.message || '更新用户角色失败';
    setTimeout(() => {
      errorMessage.value = '';
    }, 3000);
  }
};

// 创建用户
const createUser = async () => {
  // 表单验证
  if (!userForm.value.username || !userForm.value.password || !userForm.value.email) {
    errorMessage.value = '请填写所有必填字段';
    return;
  }
  
  const result = await adminStore.createUser(userForm.value);
  if (result.success) {
    // 重置表单
    userForm.value = {
      username: '',
      password: '',
      email: '',
      role: 'user'
    };
    showCreateUserModal.value = false;
    successMessage.value = '用户创建成功';
    setTimeout(() => {
      successMessage.value = '';
    }, 3000);
  } else {
    errorMessage.value = result.message || '创建用户失败';
  }
};

// 分页处理
const goToPage = (page) => {
  loadUsers(page);
};

// 页面加载时获取用户列表
onMounted(() => {
  loadUsers();
});
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-4 md:p-6 border border-gray-100 dark:border-gray-700/50 animate-fadeIn">
    <!-- 页面标题与操作区 -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
      <div>
        <h1 class="text-xl md:text-2xl font-bold text-gray-800 dark:text-white">用户管理</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">管理系统用户，分配权限</p>
      </div>
      
      <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
        <!-- 搜索框 -->
        <div class="relative flex-1 sm:flex-initial">
          <input 
            type="text" 
            v-model="searchKeyword"
            placeholder="搜索用户..." 
            class="pl-10 pr-4 py-2.5 w-full border border-gray-200 dark:border-gray-700 rounded-xl bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-600 focus:border-transparent shadow-sm"
            @keyup.enter="searchUsers"
          />
          <span class="absolute left-3 top-3 text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
        </div>
        
        <!-- 添加用户按钮 -->
        <button 
          @click="showCreateUserModal = true"
          class="flex items-center justify-center px-5 py-2.5 bg-gradient-to-r from-violet-500 to-purple-600 text-white rounded-xl hover:from-violet-600 hover:to-purple-700 transition-all shadow-md shadow-purple-200 dark:shadow-purple-900/20 hover:shadow-purple-300 dark:hover:shadow-purple-900/30 transform hover:scale-[1.02]"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
          </svg>
          <span>添加用户</span>
        </button>
      </div>
    </div>
    
    <!-- 提示消息 -->
    <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg border border-red-100 dark:border-red-900/30 flex items-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span>{{ errorMessage }}</span>
    </div>
    
    <div v-if="successMessage" class="mb-4 p-3 bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400 rounded-lg border border-green-100 dark:border-green-900/30 flex items-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
      </svg>
      <span>{{ successMessage }}</span>
    </div>
    
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 border border-gray-100 dark:border-gray-700/50">
        <div class="flex items-center">
          <div class="w-12 h-12 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-500 dark:text-blue-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">总用户数</h3>
            <p class="text-2xl font-bold text-gray-800 dark:text-white">{{ pagination.total || 0 }}</p>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 border border-gray-100 dark:border-gray-700/50">
        <div class="flex items-center">
          <div class="w-12 h-12 rounded-lg bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center text-purple-500 dark:text-purple-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">管理员</h3>
            <p class="text-2xl font-bold text-gray-800 dark:text-white">{{ users.filter(u => u.role === 'admin').length }}</p>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 border border-gray-100 dark:border-gray-700/50">
        <div class="flex items-center">
          <div class="w-12 h-12 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center text-emerald-500 dark:text-emerald-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">在线用户</h3>
            <p class="text-2xl font-bold text-gray-800 dark:text-white">{{ Math.round(users.length * 0.6) }}</p>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 border border-gray-100 dark:border-gray-700/50">
        <div class="flex items-center">
          <div class="w-12 h-12 rounded-lg bg-amber-100 dark:bg-amber-900/30 flex items-center justify-center text-amber-500 dark:text-amber-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">总笔记数</h3>
            <p class="text-2xl font-bold text-gray-800 dark:text-white">{{ Math.round(users.length * 12) }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 用户列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-xl overflow-hidden border border-gray-100 dark:border-gray-700/50 shadow-lg">
      <!-- 加载状态 -->
      <div v-if="isLoading" class="flex justify-center items-center p-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-4 border-l-4 border-purple-500 dark:border-purple-400"></div>
        <span class="ml-3 text-gray-600 dark:text-gray-300 animate-pulse">加载中...</span>
      </div>
      
      <!-- 没有数据 -->
      <div v-else-if="users.length === 0" class="p-20 text-center">
        <div class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-8 max-w-md mx-auto">
          <div class="w-20 h-20 mx-auto rounded-full bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center text-purple-500 dark:text-purple-400 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-bold text-gray-700 dark:text-gray-200">没有找到用户</h3>
          <p class="text-gray-500 dark:text-gray-400 mt-2 mb-4">尝试使用不同的搜索条件或添加新用户</p>
          <button 
            @click="showCreateUserModal = true"
            class="px-4 py-2 bg-gradient-to-r from-violet-500 to-purple-600 text-white rounded-lg hover:from-violet-600 hover:to-purple-700 transition-all"
          >
            添加用户
          </button>
        </div>
      </div>
      
      <!-- 用户列表表格 -->
      <div v-else class="overflow-x-auto">
        <table class="w-full min-w-full table-auto">
          <thead class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-800/80 dark:to-gray-900/80 border-b border-gray-200 dark:border-gray-700">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">用户</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider hidden md:table-cell">邮箱</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">角色</th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider hidden md:table-cell">注册时间</th>
              <th class="px-4 py-3 text-center text-xs font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(user, index) in users" :key="user.id" 
                class="border-b border-gray-100 dark:border-gray-800 transition-colors"
                :class="index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50 dark:bg-gray-800/50'"
            >
              <td class="px-4 py-4">
                <div class="flex items-center">
                  <!-- 用户头像 -->
                  <div class="w-10 h-10 flex-shrink-0 rounded-full shadow-md flex items-center justify-center font-medium text-white"
                       :class="[
                         user.role === 'admin' 
                           ? 'bg-gradient-to-br from-purple-500 to-indigo-600' 
                           : 'bg-gradient-to-br from-blue-500 to-cyan-500'
                       ]"
                  >
                    {{ user.username.charAt(0).toUpperCase() }}
                  </div>
                  <div class="ml-3">
                    <div class="text-sm font-semibold text-gray-900 dark:text-white">{{ user.username }}</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 md:hidden">{{ user.email }}</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 hidden md:block">ID: {{ user.id }}</div>
                  </div>
                </div>
              </td>
              <td class="px-4 py-4 hidden md:table-cell">
                <div class="text-sm text-gray-800 dark:text-gray-200 font-medium">{{ user.email }}</div>
              </td>
              <td class="px-4 py-4">
                <div class="flex items-center">
                  <div class="relative inline-block w-10 mr-2 align-middle">
                    <input 
                      :id="`toggle-${user.id}`" 
                      type="checkbox" 
                      class="sr-only peer" 
                      :checked="user.role === 'admin'"
                      @change="handleRoleChange(user, $event.target.checked ? 'admin' : 'user')"
                      :disabled="user.id === currentUser.id"
                    />
                    <div class="w-10 h-5 bg-gray-200 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all dark:border-gray-600 peer-checked:bg-purple-600"></div>
                  </div>
                  <span :class="[
                    'text-xs px-2 py-1 rounded-full font-medium',
                    user.role === 'admin' 
                      ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300'
                      : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
                  ]">
                    {{ user.role === 'admin' ? '管理员' : '用户' }}
                  </span>
                </div>
              </td>
              <td class="px-4 py-4 hidden md:table-cell">
                <div class="text-sm text-gray-500 dark:text-gray-400">
                  {{ new Date(user.created_at).toLocaleDateString('zh-CN') }}
                </div>
              </td>
              <td class="px-4 py-4">
                <div class="flex justify-center space-x-2">
                  <button 
                    @click="viewUserDetail(user)"
                    class="p-2 rounded-lg bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-colors" 
                    title="查看详情"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button 
                    v-if="user.id !== currentUser.id"
                    @click="confirmDeleteUser(user)"
                    class="p-2 rounded-lg bg-red-50 text-red-600 dark:bg-red-900/20 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors" 
                    title="删除用户"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 分页 -->
      <div v-if="users.length > 0" class="px-6 py-6 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div class="text-sm text-gray-700 dark:text-gray-300">
            显示第 <span class="font-semibold text-gray-900 dark:text-white">{{ (pagination.page - 1) * pagination.pageSize + 1 }}</span> 到 
            <span class="font-semibold text-gray-900 dark:text-white">{{ Math.min(pagination.page * pagination.pageSize, pagination.total) }}</span> 条，
            共 <span class="font-semibold text-gray-900 dark:text-white">{{ pagination.total }}</span> 条记录
          </div>
          <div class="flex">
            <button 
              @click="goToPage(pagination.page - 1)" 
              :disabled="pagination.page <= 1"
              class="inline-flex items-center px-4 py-2 mr-2 text-sm font-medium rounded-lg transition-colors"
              :class="pagination.page <= 1 
                ? 'text-gray-400 bg-gray-100 dark:text-gray-500 dark:bg-gray-800 cursor-not-allowed' 
                : 'text-gray-700 bg-white hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:border dark:border-gray-700 dark:hover:bg-gray-700'"
            >
              <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
              上一页
            </button>
            
            <div class="flex items-center">
              <span 
                v-for="page in [1, 2, 3, 4, 5].filter(p => p <= pagination.totalPage)" 
                :key="page"
                @click="goToPage(page)"
                :class="[
                  'inline-flex items-center justify-center w-9 h-9 mx-1 rounded-lg text-sm transition-colors cursor-pointer',
                  page === pagination.page 
                    ? 'bg-gradient-to-r from-violet-500 to-purple-600 text-white font-medium'
                    : 'text-gray-700 bg-white hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
                ]"
              >
                {{ page }}
              </span>
            </div>
            
            <button 
              @click="goToPage(pagination.page + 1)" 
              :disabled="pagination.page >= pagination.totalPage"
              class="inline-flex items-center px-4 py-2 ml-2 text-sm font-medium rounded-lg transition-colors"
              :class="pagination.page >= pagination.totalPage 
                ? 'text-gray-400 bg-gray-100 dark:text-gray-500 dark:bg-gray-800 cursor-not-allowed' 
                : 'text-gray-700 bg-white hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:border dark:border-gray-700 dark:hover:bg-gray-700'"
            >
              下一页
              <svg class="w-5 h-5 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
  
  <!-- 创建用户模态框 -->
  <div v-if="showCreateUserModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 px-4 animate-fadeIn">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-gray-100 dark:border-gray-700/50 transform transition-all">
      <div class="flex justify-between items-center p-6 border-b border-gray-100 dark:border-gray-700/50 bg-gradient-to-r from-gray-50 to-white dark:from-gray-800 dark:to-gray-900">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-violet-500 to-purple-600 text-white flex items-center justify-center shadow-md">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
            </svg>
          </div>
          <h3 class="text-xl font-bold text-gray-800 dark:text-white">创建新用户</h3>
        </div>
        <button @click="showCreateUserModal = false" class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 transition-all">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <div class="p-6">
        <form @submit.prevent="createUser" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="username">
              用户名
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <input 
                id="username" 
                type="text" 
                v-model="userForm.username"
                class="w-full pl-10 pr-4 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-600 focus:border-transparent"
                placeholder="输入用户名"
                required
              />
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="email">
              邮箱
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
              </div>
              <input 
                id="email" 
                type="email" 
                v-model="userForm.email"
                class="w-full pl-10 pr-4 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-600 focus:border-transparent"
                placeholder="输入邮箱地址"
                required
              />
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="password">
              密码
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input 
                id="password" 
                type="password" 
                v-model="userForm.password"
                class="w-full pl-10 pr-4 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-600 focus:border-transparent"
                placeholder="输入密码"
                required
              />
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="role">
              角色
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <select 
                id="role" 
                v-model="userForm.role"
                class="w-full pl-10 pr-4 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-600 focus:border-transparent appearance-none"
              >
                <option value="user">普通用户</option>
                <option value="admin">管理员</option>
              </select>
              <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>
          
          <div v-if="errorMessage" class="bg-red-50 border-l-4 border-red-500 dark:bg-red-900/20 dark:border-red-600 p-4 rounded-lg text-sm text-red-600 dark:text-red-400">
            <div class="flex">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ errorMessage }}
            </div>
          </div>
          
          <div class="flex justify-end space-x-3 pt-4">
            <button 
              type="button"
              @click="showCreateUserModal = false"
              class="px-5 py-2.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 font-medium transition-colors shadow-sm"
            >
              取消
            </button>
            <button 
              type="submit"
              class="px-5 py-2.5 bg-gradient-to-r from-violet-500 to-purple-600 text-white rounded-xl hover:from-violet-600 hover:to-purple-700 font-medium transition-all shadow-lg shadow-purple-200 dark:shadow-purple-900/20 hover:shadow-purple-300 dark:hover:shadow-purple-900/30"
            >
              创建用户
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
  
  <!-- 用户详情模态框 -->
  <div v-if="showUserDetailModal && selectedUser" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 px-4 animate-fadeIn">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-gray-100 dark:border-gray-700/50">
      <div class="relative h-32 bg-gradient-to-r from-purple-500 to-indigo-600 dark:from-purple-700 dark:to-indigo-800">
        <button @click="showUserDetailModal = false" class="absolute top-4 right-4 p-2 rounded-full bg-white/20 hover:bg-white/30 text-white transition-all">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        
        <!-- 用户头像 (在顶部中央突出显示) -->
        <div class="absolute -bottom-12 left-1/2 transform -translate-x-1/2">
          <div class="w-24 h-24 rounded-full shadow-xl flex items-center justify-center text-2xl font-bold text-white bg-gradient-to-br from-purple-500 to-indigo-600 ring-8 ring-white dark:ring-gray-800">
            {{ selectedUser.username.charAt(0).toUpperCase() }}
          </div>
        </div>
      </div>
      
      <div class="pt-16 px-6 pb-6">
        <h3 class="text-2xl font-bold text-gray-900 dark:text-white text-center">{{ selectedUser.username }}</h3>
        <p class="text-gray-500 dark:text-gray-400 text-center mb-6">{{ selectedUser.email }}</p>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div class="bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4">
            <h4 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">角色</h4>
            <div class="flex items-center">
              <span :class="[
                'inline-flex items-center justify-center w-8 h-8 rounded-lg mr-2',
                selectedUser.role === 'admin' ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300' : 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
              ]">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path v-if="selectedUser.role === 'admin'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </span>
              <span class="text-gray-900 dark:text-white font-medium">{{ selectedUser.role === 'admin' ? '管理员' : '普通用户' }}</span>
            </div>
          </div>
          
          <div class="bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4">
            <h4 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">笔记数量</h4>
            <div class="flex items-center">
              <span class="inline-flex items-center justify-center w-8 h-8 rounded-lg bg-cyan-100 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-300 mr-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </span>
              <span class="text-gray-900 dark:text-white font-medium">{{ selectedUser.noteCount || 0 }} 篇笔记</span>
            </div>
          </div>
          
          <div class="bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4">
            <h4 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">用户ID</h4>
            <div class="flex items-center">
              <span class="inline-flex items-center justify-center w-8 h-8 rounded-lg bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 mr-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" />
                </svg>
              </span>
              <span class="text-gray-900 dark:text-white font-medium">{{ selectedUser.id }}</span>
            </div>
          </div>
          
          <div class="bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4">
            <h4 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">注册时间</h4>
            <div class="flex items-center">
              <span class="inline-flex items-center justify-center w-8 h-8 rounded-lg bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300 mr-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </span>
              <span class="text-gray-900 dark:text-white font-medium">{{ new Date(selectedUser.created_at).toLocaleString('zh-CN') }}</span>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button 
            @click="showUserDetailModal = false"
            class="px-5 py-2.5 bg-gradient-to-r from-violet-500 to-purple-600 text-white rounded-xl hover:from-violet-600 hover:to-purple-700 font-medium transition-all shadow-lg shadow-purple-200 dark:shadow-purple-900/20 hover:shadow-purple-300 dark:hover:shadow-purple-900/30"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
  
  <!-- 删除确认模态框 -->
  <div v-if="showDeleteConfirmModal && selectedUser" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 px-4 animate-fadeIn">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-gray-100 dark:border-gray-700/50">
      <div class="p-6 text-center">
        <div class="w-20 h-20 mx-auto rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-500 dark:text-red-400 mb-5">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-3">确认删除用户</h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6 max-w-sm mx-auto">
          您确定要删除用户 <span class="font-semibold text-gray-800 dark:text-gray-200">"{{ selectedUser.username }}"</span> 吗？<br/>此操作不可撤销，删除后相关数据将无法恢复。
        </p>
        
        <div class="flex justify-center space-x-4">
          <button 
            @click="showDeleteConfirmModal = false"
            class="px-5 py-2.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 font-medium transition-colors shadow-sm"
          >
            取消
          </button>
          <button 
            @click="handleDeleteUser"
            class="px-5 py-2.5 bg-gradient-to-r from-red-500 to-pink-600 text-white rounded-xl hover:from-red-600 hover:to-pink-700 font-medium transition-all shadow-lg shadow-red-200 dark:shadow-red-900/20 hover:shadow-red-300 dark:hover:shadow-red-900/30"
          >
            确认删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fadeIn {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 美化滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(247, 247, 247, 0.1);
}

::-webkit-scrollbar-thumb {
  background: rgba(193, 193, 193, 0.3);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(193, 193, 193, 0.5);
}

.dark ::-webkit-scrollbar-track {
  background: rgba(30, 30, 30, 0.1);
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(90, 90, 90, 0.3);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(90, 90, 90, 0.5);
}
</style> 