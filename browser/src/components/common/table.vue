<template>
  <div class="main-table">
    <div class="content">
      <h2 class="title">{{items.title}}</h2>
      <template>
        <el-table
          :data="table"
          style="width: 100%">
          <el-table-column
            v-for="(prop, index) in item.prop"
            :key="prop"
            :prop="prop"
            :label="item.label[index]">
            <template slot-scope="scope">{{table[scope.$index][prop]}}</template>
            <template v-if="prop == 'timestamp'" slot-scope="scope">{{scope.row.timestamp | getFormatDate}}</template>
          </el-table-column>
        </el-table>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    table: Array,
    item: Object
  },
  computed: {
    items () {
      return this.item
    }
  },
  mounted () {
    console.log(this.item)
  }
}
</script>

<style scoped>

</style>
