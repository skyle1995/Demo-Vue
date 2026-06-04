<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from "vue";
import { getSystemInfo, getLoginLogs } from "@/api/admin/dashboard";
import dayjs from "dayjs";
import { ElMessage } from "element-plus";

defineOptions({
  name: "Index"
});

const systemInfo = ref({
  version: "",
  mode: false,
  db_type: "",
  uptime: "",
  uptime_seconds: 0
});

const loginLogs = ref([]);
const totalLogs = ref(0);
const loading = ref(false);
const pagination = ref({
  currentPage: 1,
  pageSize: 30
});

const emptyFormatter = (row: any, column: any, cellValue: any) => {
  return (cellValue === "" || cellValue === null || cellValue === undefined) ? "-" : cellValue;
};

let timer: number | null = null;

const fetchSystemInfo = async () => {
  try {
    const res = await getSystemInfo();
    if (res.code === 0) {
      systemInfo.value = res.data;
    }
  } catch (error: any) {
    console.error(
      "Failed to fetch system info",
      error.response?.status,
      error.message,
      error
    );
    ElMessage.error(error.response?.data?.message || "获取系统信息失败");
  }
};

const fetchLoginLogs = async () => {
  loading.value = true;
  try {
    const res = await getLoginLogs({
      page: pagination.value.currentPage,
      limit: pagination.value.pageSize
    });
    if (res.code === 0) {
      loginLogs.value = res.data.list || [];
      totalLogs.value = res.data.total || 0;
    }
  } catch (error: any) {
    console.error(
      "Failed to fetch login logs",
      error.response?.status,
      error.message,
      error
    );
    ElMessage.error(error.response?.data?.message || "获取登录日志失败");
  } finally {
    loading.value = false;
  }
};

const handleSizeChange = (val: number) => {
  pagination.value.pageSize = val;
  fetchLoginLogs();
};

const handleCurrentChange = (val: number) => {
  pagination.value.currentPage = val;
  fetchLoginLogs();
};

const formatDate = (dateStr: string) => {
  return dayjs(dateStr).format("YYYY-MM-DD HH:mm:ss");
};

const formattedUptime = computed(() => {
  const totalSeconds = systemInfo.value.uptime_seconds;
  if (!totalSeconds) return systemInfo.value.uptime || "未知";

  const days = Math.floor(totalSeconds / 86400);
  const hours = Math.floor((totalSeconds % 86400) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (days > 0) {
    return `${days}天 ${hours}小时 ${minutes}分钟 ${seconds}秒`;
  } else if (hours > 0) {
    return `${hours}小时 ${minutes}分钟 ${seconds}秒`;
  } else if (minutes > 0) {
    return `${minutes}分钟 ${seconds}秒`;
  } else {
    return `${seconds}秒`;
  }
});

onMounted(() => {
  fetchSystemInfo();
  fetchLoginLogs();

  timer = window.setInterval(() => {
    if (systemInfo.value.uptime_seconds > 0) {
      systemInfo.value.uptime_seconds++;
    }
  }, 1000);
});

onUnmounted(() => {
  if (timer) {
    clearInterval(timer);
  }
});
</script>

<template>
  <div>
    <el-row :gutter="20" class="mt-4">
      <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="mb-4">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          <div class="flex justify-around items-center h-[120px]">
            <div class="text-center">
              <div class="text-[var(--el-text-color-secondary)] text-sm mb-2">
                项目版本
              </div>
              <div class="text-2xl font-bold text-[var(--el-color-primary)]">
                {{ systemInfo.version || "-" }}
              </div>
            </div>
            <div class="text-center">
              <div class="text-[var(--el-text-color-secondary)] text-sm mb-2">
                运行环境
              </div>
              <div class="text-2xl font-bold text-[var(--el-color-success)]">
                {{ systemInfo.mode ? "开发环境" : "生产环境" }}
              </div>
            </div>
            <div class="text-center">
              <div class="text-[var(--el-text-color-secondary)] text-sm mb-2">
                数据库
              </div>
              <div class="text-2xl font-bold text-[var(--el-color-info)]">
                {{ systemInfo.db_type || "-" }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="mb-4">
          <template #header>
            <div class="card-header">
              <span>运行状态</span>
            </div>
          </template>
          <div
            class="uptime-display flex justify-center items-center h-[120px]"
          >
            <div class="text-center">
              <div class="text-[var(--el-text-color-secondary)] text-sm mb-2">
                系统已运行时间
              </div>
              <div
                class="text-3xl font-bold text-[var(--el-text-color-primary)]"
              >
                {{ formattedUptime }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-4">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近登录日志</span>
            </div>
          </template>
          <el-table
            v-loading="loading"
            :data="loginLogs"
            style="width: 100%"
            table-layout="auto"
            show-overflow-tooltip
            border
          >
            <el-table-column prop="created_at" label="登录时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" width="150" :formatter="emptyFormatter" />
            <el-table-column prop="ip" label="登录IP" width="200" :formatter="emptyFormatter" />
            <el-table-column prop="status" label="登录状态" width="100" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? "成功" : "失败" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="message" label="登录信息" width="240" :formatter="emptyFormatter" />
            <el-table-column prop="user_agent" label="浏览器类型" min-width="320" show-overflow-tooltip :formatter="emptyFormatter" />
          </el-table>
          <div class="flex mt-4 w-full overflow-x-auto">
            <div class="ml-auto shrink-0">
              <el-pagination
                v-model:current-page="pagination.currentPage"
                v-model:page-size="pagination.pageSize"
                :page-sizes="[10, 20, 30, 50, 100, 200, 500, 1000]"
                :background="true"
                layout="total, sizes, prev, pager, next, jumper"
                :total="totalLogs"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.card-header {
  font-weight: bold;
}

.uptime-display {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
}
</style>
