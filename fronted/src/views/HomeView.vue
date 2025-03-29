<script setup>
// 无需导入不存在的组件
</script>

<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="text-center mb-16">
      <h1 class="text-4xl font-bold text-gray-800 mb-4 sm:text-5xl">欢迎使用<span class="text-blue-600">知识库</span></h1>
      <p class="text-lg text-gray-600 max-w-2xl mx-auto mt-4">
        一个简单易用的个人知识管理工具，帮助您整理、分类和管理您的笔记和知识。
      </p>
      <div class="mt-8">
        <router-link to="/notes" class="btn btn-primary px-8 py-3 text-lg shadow-lg hover:shadow-xl">开始使用</router-link>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3">Markdown支持</h3>
        <p class="text-gray-600">
          使用Markdown语法编写笔记，支持标题、列表、代码块等格式，让您的笔记更加丰富多彩。
        </p>
      </div>

      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3">标签管理</h3>
        <p class="text-gray-600">
          通过标签对笔记进行分类，方便查找和整理您的知识，让信息管理更加条理清晰。
        </p>
      </div>

      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3">全文搜索</h3>
        <p class="text-gray-600">
          快速搜索您的笔记内容，找到您需要的信息，高效获取知识，提升工作学习效率。
        </p>
      </div>
    </div>

    <div class="bg-blue-50 rounded-xl p-8 shadow-inner">
      <div class="max-w-4xl mx-auto">
        <h2 class="text-2xl font-bold text-center mb-6 text-gray-800">功能亮点</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="flex items-start">
            <div class="bg-blue-100 rounded-full p-2 mr-4 text-blue-600">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900">响应式设计</h3>
              <p class="text-gray-600 mt-1">自适应PC和移动端，随时随地管理您的笔记</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 rounded-full p-2 mr-4 text-blue-600">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900">实时预览</h3>
              <p class="text-gray-600 mt-1">编辑Markdown时实时查看渲染效果</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 rounded-full p-2 mr-4 text-blue-600">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900">分类管理</h3>
              <p class="text-gray-600 mt-1">通过标签系统对笔记进行分类整理</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 rounded-full p-2 mr-4 text-blue-600">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900">代码高亮</h3>
              <p class="text-gray-600 mt-1">支持多种编程语言的代码块语法高亮</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
