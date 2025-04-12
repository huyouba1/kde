<template>
  <div class="cluster-management">
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h2>集群概览</h2>
              <el-button type="primary" @click="refreshClusterStatus">刷新状态</el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="节点总数" :value="clusterStats.totalNodes">
                <template #suffix>
                  <el-icon><Connection /></el-icon>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="6">
              <el-statistic title="Pod总数" :value="clusterStats.totalPods">
                <template #suffix>
                  <el-icon><Box /></el-icon>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="6">
              <el-statistic title="命名空间数" :value="clusterStats.totalNamespaces">
                <template #suffix>
                  <el-icon><Files /></el-icon>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="6">
              <el-statistic title="部署数" :value="clusterStats.totalDeployments">
                <template #suffix>
                  <el-icon><Cpu /></el-icon>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h2>节点列表</h2>
              <el-button type="success" @click="addNode">添加节点</el-button>
            </div>
          </template>
          <el-table :data="nodes" style="width: 100%">
            <el-table-column prop="name" label="节点名称" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 'Ready' ? 'success' : 'danger'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="role" label="角色" />
            <el-table-column prop="ip" label="IP地址" />
            <el-table-column prop="cpu" label="CPU使用率">
              <template #default="{ row }">
                <el-progress :percentage="row.cpu" />
              </template>
            </el-table-column>
            <el-table-column prop="memory" label="内存使用率">
              <template #default="{ row }">
                <el-progress :percentage="row.memory" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button-group>
                  <el-button type="primary" size="small" @click="viewNodeDetail(row)">详情</el-button>
                  <el-button type="warning" size="small" @click="drainNode(row)">维护</el-button>
                  <el-button type="danger" size="small" @click="removeNode(row)">删除</el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection, Box, Files, Cpu } from '@element-plus/icons-vue'

const clusterStats = ref({
  totalNodes: 0,
  totalPods: 0,
  totalNamespaces: 0,
  totalDeployments: 0
})

const nodes = ref([])

const refreshClusterStatus = async () => {
  try {
    // TODO: 调用后端API获取集群状态
    ElMessage.success('集群状态已更新')
  } catch (error) {
    ElMessage.error('获取集群状态失败')
  }
}

const addNode = () => {
  // TODO: 实现添加节点逻辑
}

const viewNodeDetail = (node) => {
  // TODO: 实现查看节点详情逻辑
}

const drainNode = (node) => {
  // TODO: 实现节点维护模式逻辑
}

const removeNode = (node) => {
  // TODO: 实现删除节点逻辑
}

onMounted(() => {
  refreshClusterStatus()
})
</script>

<style scoped>
.cluster-management {
  padding: 20px;
}

.mb-4 {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
}

.el-statistic {
  text-align: center;
}
</style>