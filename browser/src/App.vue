<template>
  <div id="app">
    <Side></Side>
    <router-view/>
  </div>
</template>

<script>
import Side from '@/components/common/side'
export default {
  name: 'App',
  components: {
    Side
  }
}
</script>

<style>
  #app {
    height: 100%;
  }
</style>
