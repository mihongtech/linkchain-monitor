<template>
  <div class="main-table" id="statistics">
    <div class="content">
      <h2 class="title">{{block.head}}</h2>
      <div class="table">
        <table>
          <tbody>
          <tr>
            <th v-for="item in block.title" :key="item">{{item}}</th>
          </tr>
          <tr v-for="item in block.statistics" :key="item.block">
            <td>{{item.item1}}</td>
            <td>{{item.item2}}</td>
            <td>{{item.item3}}</td>
            <td v-if="item.item4" :style="{color: block.color || ''}">{{item.item4}}</td>
            <td v-if="item.item5">{{item.item5}}</td>
            <td v-if="item.item6">{{item.item6}}</td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'statistics',
  props: {
    block: Object
  },
  data () {
    return {
    }
  }
}
</script>

<style>
  .main-table {
    box-sizing: border-box;
    width: 100%;
    background-color: #323241;
    padding: 0 20px 20px 0;
  }
  .main-table .content{
    background-color: #414052;
    box-shadow: 0 0 40px rgba(0, 0, 0, .2);
  }
  .main-table .title {
    height: 60px;
    font-size: 14px;
    padding-left: 10px;
    line-height: 60px;
    font-weight: 500;
    color: #fff;
  }
  .main-table .table table {
    width: 100%;
    text-align: left;
    border-collapse: collapse;
    border-spacing: 0;
  }
  .main-table .table table th {
    background-color: #33333f;
    height: 50px;
    padding-left: 20px;
    font-weight: 500;
    color: #9996b3;
  }
  .main-table .table table td {
    height: 45px;
    padding-left: 20px;
    color: #d3d3eb;
    font-weight: 500;
  }
  .main-table .table table tr:nth-child(2n-1) {
    background-color: #414052;
  }
  .main-table .table table tr:nth-child(2n) {
    background-color: #484759;
  }
</style>
