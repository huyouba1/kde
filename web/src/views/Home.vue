<template>
  <el-container class="layout-container">
    <el-header>
      <el-menu mode="horizontal" :router="true">
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/cluster">集群管理</el-menu-item>
        <el-menu-item index="/deploy">应用部署</el-menu-item>
        <el-menu-item index="/plugins">插件管理</el-menu-item>
      </el-menu>
    </el-header>
    
    <el-main>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="feature-card">
            <template #header>
              <div class="card-header">
                <h3>集群管理</h3>
              </div>
            </template>
            <div class="card-content">
              <p>管理和监控您的Kubernetes集群</p>
              <el-button type="primary" @click="$router.push('/cluster')">进入管理</el-button>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="feature-card">
            <template #header>
              <div class="card-header">
                <h3>应用部署</h3>
              </div>
            </template>
            <div class="card-content">
              <p>部署和管理您的应用</p>
              <el-button type="primary" @click="$router.push('/deploy')">开始部署</el-button>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="feature-card">
            <template #header>
              <div class="card-header">
                <h3>插件管理</h3>
              </div>
            </template>
            <div class="card-content">
              <p>管理系统插件和扩展</p>
              <el-button type="primary" @click="$router.push('/plugins')">查看插件</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script setup>
import { ElContainer, ElHeader, ElMain, ElMenu, ElMenuItem, ElRow, ElCol, ElCard, ElButton } from 'element-plus'
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-header {
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
}

.feature-card {
  height: 100%;
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-content {
  text-align: center;
}

.card-content p {
  margin: 10px 0 20px;
  color: #666;
}
</style>