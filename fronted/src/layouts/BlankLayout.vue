<script setup>
// 空白布局组件，不包含头部和底部
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full">
      <main class="w-full">
        <router-view />
      </main>
    </div>
  </div>
</template> 