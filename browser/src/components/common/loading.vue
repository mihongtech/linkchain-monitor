<template>
  <div class="mask"
       fullscreen="false"
       element-loading-text="拼命加载中"
       element-loading-spinner="el-icon-loading"
       element-loading-background="rgba(0, 0, 0, 0.2)">
  </div>
</template>

<script>
export default {
  name: 'loading'
}
</script>

<style scoped>
  .mask {
    height: 100%;
  }
</style>
