<template>
  <div class="application-deployment">
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h2>应用部署</h2>
            </div>
          </template>
          
          <el-tabs v-model="activeTab">
            <!-- YAML部署 -->
            <el-tab-pane label="YAML部署" name="yaml">
              <el-upload
                class="upload-demo"
                drag
                action="/api/deploy/yaml"
                accept=".yaml,.yml"
                :on-success="handleUploadSuccess"
                :on-error="handleUploadError"
                :before-upload="beforeUpload"
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  将文件拖到此处，或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持上传 YAML 文件 (.yaml, .yml)
                  </div>
                </template>
              </el-upload>

              <div v-if="yamlContent" class="mt-4">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <h3>YAML预览</h3>
                      <el-button type="primary" @click="deployYaml">部署</el-button>
                    </div>
                  </template>
                  <pre>{{ yamlContent }}</pre>
                </el-card>
              </div>
            </el-tab-pane>

            <!-- Helm Charts部署 -->
            <el-tab-pane label="Helm Charts" name="helm">
              <el-upload
                class="upload-demo"
                drag
                action="/api/deploy/helm"
                accept=".tgz"
                :on-success="handleHelmUploadSuccess"
                :on-error="handleUploadError"
                :before-upload="beforeHelmUpload"
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  将Helm包拖到此处，或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持上传Helm Charts包 (.tgz)
                  </div>
                </template>
              </el-upload>

              <div v-if="helmChartInfo" class="mt-4">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <h3>Chart信息</h3>
                    </div>
                  </template>
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="名称">{{ helmChartInfo.name }}</el-descriptions-item>
                    <el-descriptions-item label="版本">{{ helmChartInfo.version }}</el-descriptions-item>
                    <el-descriptions-item label="描述">{{ helmChartInfo.description }}</el-descriptions-item>
                  </el-descriptions>
                  
                  <div class="mt-4">
                    <h4>配置参数</h4>
                    <el-form :model="helmValues" label-width="120px">
                      <el-form-item
                        v-for="(value, key) in helmValues"
                        :key="key"
                        :label="key"
                      >
                        <el-input v-model="helmValues[key]" />
                      </el-form-item>
                    </el-form>
                    
                    <el-button type="primary" @click="deployHelm">部署Chart</el-button>
                  </div>
                </el-card>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>

    <!-- 部署历史 -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h2>部署历史</h2>
              <el-button type="primary" @click="refreshDeployHistory">刷新</el-button>
            </div>
          </template>
          <el-table :data="deployHistory" style="width: 100%">
            <el-table-column prop="name" label="应用名称" />
            <el-table-column prop="type" label="部署类型" />
            <el-table-column prop="namespace" label="命名空间" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 'Success' ? 'success' : 'danger'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="deployTime" label="部署时间" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button-group>
                  <el-button type="primary" size="small" @click="viewDeployDetail(row)">详情</el-button>
                  <el-button type="danger" size="small" @click="deleteDeployment(row)">删除</el-button>
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
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const activeTab = ref('yaml')
const yamlContent = ref('')
const helmChartInfo = ref(null)
const helmValues = ref({})
const deployHistory = ref([])

const handleUploadSuccess = (response) => {
  yamlContent.value = response.content
  ElMessage.success('文件上传成功')
}

const handleHelmUploadSuccess = (response) => {
  helmChartInfo.value = response.chartInfo
  helmValues.value = response.defaultValues
  ElMessage.success('Helm Chart上传成功')
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

const beforeUpload = (file) => {
  const isYAML = file.name.endsWith('.yaml') || file.name.endsWith('.yml')
  if (!isYAML) {
    ElMessage.error('请上传YAML文件')
    return false
  }
  return true
}

const beforeHelmUpload = (file) => {
  const isTgz = file.name.endsWith('.tgz')
  if (!isTgz) {
    ElMessage.error('请上传Helm Chart包(.tgz)')
    return false
  }
  return true
}

const deployYaml = async () => {
  try {
    // TODO: 调用后端API部署YAML
    ElMessage.success('部署请求已提交')
  } catch (error) {
    ElMessage.error('部署失败')
  }
}

const deployHelm = async () => {
  try {
    // TODO: 调用后端API部署Helm Chart
    ElMessage.success('Helm Chart部署请求已提交')
  } catch (error) {
    ElMessage.error('部署失败')
  }
}

const refreshDeployHistory = async () => {
  try {
    // TODO: 获取部署历史
    ElMessage.success('部署历史已更新')
  } catch (error) {
    ElMessage.error('获取部署历史失败')
  }
}

const viewDeployDetail = (deployment) => {
  // TODO: 查看部署详情
}

const deleteDeployment = async (deployment) => {
  try {
    // TODO: 删除部署
    ElMessage.success('部署已删除')
  } catch (error) {
    ElMessage.error('删除失败')
  }
}
</script>

<style scoped>
.application-deployment {
  padding: 20px;
}

.mb-4 {
  margin-bottom: 20px;
}

.mt-4 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2, .card-header h3 {
  margin: 0;
}

.upload-demo {
  text-align: center;
}

pre {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
}
</style>