<script setup lang="ts">
import { getConfig } from "@/config";
import {
  Monitor,
  User,
  Lock,
  DataLine,
  Setting,
  Grid
} from "@element-plus/icons-vue";

defineOptions({
  name: "HomeMain"
});

const features = [
  {
    title: "基础后台框架",
    icon: Monitor,
    description:
      "提供标准化的前后端分层结构、路由组织与页面骨架，适合作为企业后台、工具平台或业务管理系统的起步模板，便于在现有基础上快速二次开发。",
    color: "#409eff"
  },
  {
    title: "用户与权限基础",
    icon: User,
    description:
      "内置管理员登录、个人资料、公开设置与基础接口能力，可作为后续扩展 RBAC、菜单控制与多角色体系的基础骨架，降低项目冷启动成本。",
    color: "#67c23a"
  },
  {
    title: "日志与安全能力",
    icon: Lock,
    description:
      "已集成登录日志、操作日志、验证码与会话安全等基础能力，为通用后台提供最常见的安全防护与审计入口，适合后续按业务需要逐步增强。",
    color: "#e6a23c"
  },
  {
    title: "数据面板预留",
    icon: DataLine,
    description:
      "当前卡片作为仪表盘与业务统计能力的占位区域，可在后续根据项目需求接入订单、用户、设备或其他核心指标，逐步形成完整的运营视图。",
    color: "#f56c6c"
  },
  {
    title: "系统配置预留",
    icon: Setting,
    description:
      "当前卡片作为系统设置、主题配置、站点信息与运行参数的展示预留位，适合在后续扩展更多基础能力时继续沿用同一门户结构与视觉风格。",
    color: "#909399"
  },
  {
    title: "功能模块预留",
    icon: Grid,
    description:
      "当前项目暂未定义更多门户描述项，因此保留空白能力位用于后续模块接入。未来可根据真实业务逐步替换为订单、资源、流程或其他核心功能入口。",
    color: "#8e44ad"
  }
];
</script>

<template>
  <div class="home-container">
    <div class="hero-section">
      <h1 class="hero-title">{{ getConfig()?.Title || "加载中..." }}</h1>
      <p class="hero-subtitle">
        {{
          getConfig()?.Description ||
          "Demo 是一款现代化的企业级管理系统基础框架，提供基础的 Web 管理界面、权限骨架与 API 服务支持。"
        }}
      </p>
    </div>

    <div class="features-section">
      <h2 class="section-title">核心功能特性</h2>
      <div class="features-grid">
        <el-card
          v-for="(item, index) in features"
          :key="index"
          class="feature-card"
          shadow="hover"
        >
          <div
            class="feature-icon"
            :style="{ color: item.color, backgroundColor: `${item.color}15` }"
          >
            <el-icon :size="32"><component :is="item.icon" /></el-icon>
          </div>
          <h3 class="feature-title">{{ item.title }}</h3>
          <p class="feature-desc">{{ item.description }}</p>
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.hero-section {
  text-align: center;
  width: 100%;
  padding: 20px 0 40px;
}

.hero-title {
  margin-bottom: 20px;
  font-size: 48px;
  font-weight: 800;
  background: linear-gradient(
    135deg,
    var(--el-color-primary),
    var(--el-color-primary-light-3)
  );
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
  font-size: 20px;
  line-height: 1.6;
  color: var(--el-text-color-secondary);
}

.features-section {
  width: 100%;
  padding-bottom: 60px;
}

.section-title {
  position: relative;
  margin-bottom: 50px;
  text-align: center;
  font-size: 32px;
  font-weight: bold;
  color: var(--el-text-color-primary);

  &::after {
    content: "";
    position: absolute;
    bottom: -15px;
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 4px;
    border-radius: 2px;
    background-color: var(--el-color-primary);
  }
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  width: 100%;
}

.feature-card {
  height: 100%;
  border: none;
  border-radius: 12px;
  background-color: var(--el-bg-color-overlay);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.08) !important;
  }

  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 30px;
  }
}

.feature-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin-bottom: 24px;
  border-radius: 16px;
  transition: transform 0.3s ease;

  .feature-card:hover & {
    transform: scale(1.1);
  }
}

.feature-title {
  margin-bottom: 16px;
  font-size: 20px;
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.feature-desc {
  margin: 0;
  font-size: 15px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
}

@media screen and (max-width: 768px) {
  .hero-title {
    font-size: 36px;
  }

  .hero-subtitle {
    padding: 0 20px;
    font-size: 16px;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>
