<template>
  <el-container>
    <el-aside 
      class="menu"
      :width="isMob ? (isCollapse ? '0px' : '200px'): (isCollapse ? '64px' : '200px')"
    >
      <div class="logo flx-center">
        <a href="https://github.com/orcastor" target="_blank">
          <img src="/logo.svg" alt="logo" />
        </a>
        <span v-show="!isCollapse"></span>
      </div>
      <el-menu
        active-text-color="#EF7C00"
        background-color="#F8F8F8"
        :default-active="0"
        text-color="#004482"
        :collapse="isCollapse"
      >
        <el-menu-item :index="0">
          <el-icon><House /></el-icon>
          <template #title>
            <span>首页</span>
          </template>
        </el-menu-item>
        <el-menu-item :index="1">
          <el-icon><Cellphone /></el-icon>
          <template #title>
            <span>设备管理</span>
          </template>
        </el-menu-item>
        <el-menu-item :index="2">
          <el-icon><Box /></el-icon>
          <template #title>
            <span>备份管理</span>
          </template>
        </el-menu-item>
      </el-menu>
      <el-dropdown trigger="click" class="avatar flx-center">
        <div>
            <el-avatar :size="30" src="/assets/avatar.png" />
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item disabled>{{userInfo.n}}</el-dropdown-item>
            <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-lf flx-center">
          <el-icon class="collapse-icon" @click="setCollapse">
            <Expand v-if="isCollapse" /><Fold v-else />
          </el-icon>
          <span v-if="previewing" >{{preview_title}}</span>
        </div>
      </el-header>
      <el-main class="main" :style=mainStyle()>
        <el-empty description="空目录" />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import router from "@/routers";

import { store } from "@/store";
import { House, Cellphone, Expand, Fold, Box } from '@element-plus/icons-vue';
import { toDefaultIcon, toIcon, getExt, isZip } from "@/config/icons";

import { Cache } from "@/store/cache";

import 'element-plus/es/components/message-box/style/css';
import { ElMessage, ElMessageBox } from 'element-plus';

const loading = ref(true);
const data = ref([]);
const isCollapse = ref(store.isCollapse);

const userInfo = computed(() => store.userInfo);

const isMob = navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i);

const iframeStyle = () => {
  return 'border: 0; width:100%; height:'+(100-5500/document.body.clientHeight).toFixed(2)+'vh;';
}

const mainStyle = () => {
  return 'min-height:'+(100-5500/document.body.clientHeight).toFixed(2)+'vh;';
}

const cache = new Cache(100, null);

function toSize(scope:any):string {
  if (scope.row.t == 2) {
    const sz = scope.row.s||0;
    if (sz < 1e3) { return sz + '  B'; }
    if (sz < 1e6) { return (sz/1e3).toFixed(2) + ' KB'; }
    if (sz < 1e9) { return (sz/1e6).toFixed(2) + ' MB'; }
    return (sz/1e9).toFixed(2) + ' GB';
  }
  return '-';
}

// aside 自适应
const screenWidth = ref<number>(0);
// 监听窗口大小变化，合并 aside
const listeningWindow = () => {
  window.onresize = () => {
    return (() => {
      screenWidth.value = document.body.clientWidth;
      if ((isCollapse.value === false && screenWidth.value < 1200)
        || (isCollapse.value === true && screenWidth.value > 1200))
        setCollapse();
    })();
  };
};
listeningWindow();

const setCollapse = () => {
  isCollapse.value = !isCollapse.value;
  store.setCollapse();
};

watch(() => router.currentRoute.value.query, (_newValue:any, _oldValue:any) => {
  if (router.currentRoute.value.path == '/') {
    init();
  }
});

onMounted(() => {
  init();
});

const init = () => {
};

// 退出登录
const logout = () => {
  ElMessageBox.confirm("您是否确认退出登录?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    store.setToken("");
    ElMessage({
      type: "success",
      message: "退出登录成功",
    });
    router.push({ name: "login", query: router.currentRoute.value.query });
  });
};

</script>

<style scoped lang="scss">
.main {
  min-height: 100vh;
  overflow: auto;
  padding: 0;
  :deep(tr.el-table__row) {
    cursor: pointer;
  }
}

.header {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 55px;
  padding: 0 15px;
  background-color: #ffffff;
  border-bottom: 1px solid #f6f6f6;
  .header-lf {
    .collapse-icon {
      margin-right: 20px;
      font-size: 22px;
      color: rgb(0 0 0 / 75%);
      cursor: pointer;
    }
  }
}

.menu {
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  background-color: #F8F8F8;
  transition: all 0.3s ease;
  .logo {
    box-sizing: border-box;
    height: 55px;
    span {
      font-size: 22px;
      font-weight: bold;
      color: #dadada;
      white-space: nowrap;
    }
    img {
      width: 30px;
      object-fit: contain;
    }
  }
  .el-menu {
    flex: 1;
    overflow: auto;
    overflow-x: hidden;
    border-right: none;
  }
}

.avatar {
  height: 55px;
  cursor: pointer;
}

.el-menu,
.el-menu--popup {
  .el-menu-item {
    &.is-active {
      background-color: #fff;
      &::before {
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        width: 4px;
        content: "";
      }
    }
  }
}
</style>
