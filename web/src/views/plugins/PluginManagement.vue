<template>
  <div class="plugin-management">
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h2>插件管理</h2>
              <el-button type="primary" @click="refreshPlugins">刷新</el-button>
            </div>
          </template>
          
          <!-- 插件市场 -->
          <el-tabs v-model="activeTab">
            <el-tab-pane label="已安装插件" name="installed">
              <el-table :data="installedPlugins" style="width: 100%">
                <el-table-column prop="name" label="插件名称" />
                <el-table-column prop="version" label="版本" width="100" />
                <el-table-column prop="status" label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'Active' ? 'success' : 'warning'">
                      {{ row.status }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="description" label="描述" />
                <el-table-column label="操作" width="250">
                  <template #default="{ row }">
                    <el-button-group>
                      <el-button 
                        type="primary" 
                        size="small"
                        @click="configurePlugin(row)"
                      >配置</el-button>
                      <el-button 
                        :type="row.status === 'Active' ? 'warning' : 'success'"
                        size="small"
                        @click="togglePluginStatus(row)"
                      >{{ row.status === 'Active' ? '停用' : '启用' }}</el-button>
                      <el-button 
                        type="danger" 
                        size="small"
                        @click="uninstallPlugin(row)"
                      >卸载</el-button>
                    </el-button-group>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <el-tab-pane label="插件市场" name="marketplace">
              <el-row :gutter="20">
                <el-col :span="8" v-for="plugin in marketplacePlugins" :key="plugin.id">
                  <el-card class="plugin-card mb-4">
                    <template #header>
                      <div class="plugin-card-header">
                        <h3>{{ plugin.name }}</h3>
                        <el-tag size="small">v{{ plugin.version }}</el-tag>
                      </div>
                    </template>
                    <p class="plugin-description">{{ plugin.description }}</p>
                    <div class="plugin-footer">
                      <el-button 
                        type="primary" 
                        :loading="plugin.installing"
                        @click="installPlugin(plugin)"
                      >安装</el-button>
                      <el-button type="info" @click="viewPluginDetail(plugin)">详情</el-button>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>

    <!-- 插件配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      :title="`配置插件: ${currentPlugin?.name || ''}`"
      width="50%"
    >
      <el-form
        v-if="currentPlugin"
        :model="pluginConfig"
        label-width="120px"
      >
        <el-form-item
          v-for="(field, key) in currentPlugin.configFields"
          :key="key"
          :label="field.label"
        >
          <el-input
            v-if="field.type === 'string'"
            v-model="pluginConfig[key]"
            :placeholder="field.placeholder"
          />
          <el-switch
            v-else-if="field.type === 'boolean'"
            v-model="pluginConfig[key]"
          />
          <el-input-number
            v-else-if="field.type === 'number'"
            v-model="pluginConfig[key]"
            :min="field.min"
            :max="field.max"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="configDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="savePluginConfig">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('installed')
const installedPlugins = ref([])
const marketplacePlugins = ref([])
const configDialogVisible = ref(false)
const currentPlugin = ref(null)
const pluginConfig = ref({})

const refreshPlugins = async () => {
  try {
    // TODO: 获取已安装插件和市场插件列表
    ElMessage.success('插件列表已更新')
  } catch (error) {
    ElMessage.error('获取插件列表失败')
  }
}

const configurePlugin = (plugin) => {
  currentPlugin.value = plugin
  pluginConfig.value = { ...plugin.config }
  configDialogVisible.value = true
}

const savePluginConfig = async () => {
  try {
    // TODO: 保存插件配置
    configDialogVisible.value = false
    ElMessage.success('插件配置已保存')
  } catch (error) {
    ElMessage.error('保存配置失败')
  }
}

const togglePluginStatus = async (plugin) => {
  try {
    // TODO: 切换插件状态
    ElMessage.success(`插件已${plugin.status === 'Active' ? '停用' : '启用'}`)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const uninstallPlugin = async (plugin) => {
  try {
    // TODO: 卸载插件
    ElMessage.success('插件已卸载')
  } catch (error) {
    ElMessage.error('卸载失败')
  }
}

const installPlugin = async (plugin) => {
  try {
    plugin.installing = true
    // TODO: 安装插件
    ElMessage.success('插件安装成功')
  } catch (error) {
    ElMessage.error('安装失败')
  } finally {
    plugin.installing = false
  }
}

const viewPluginDetail = (plugin) => {
  // TODO: 查看插件详情
}
</script>

<style scoped>
.plugin-management {
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

.plugin-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plugin-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.plugin-card-header h3 {
  margin: 0;
}

.plugin-description {
  flex-grow: 1;
  margin: 10px 0;
  color: #666;
}

.plugin-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 15px;
}
</style>