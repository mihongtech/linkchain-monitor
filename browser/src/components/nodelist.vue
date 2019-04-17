<template>
  <div class="root">
    <Loading v-loading="loading" v-if="loading"/>
    <div class="node-list" v-else>
      <div class="main-table sync-node-table">
        <div class="content">
          <h2 class="title">同步节点分布列表</h2>
          <template>
            <el-table
              :data="syncNode"
              style="width: 100%">
              <el-table-column
                prop="city"
                label="分布"
                align="center">
              </el-table-column>
              <el-table-column
                prop="count"
                label="节点数"
                align="center">
              </el-table-column>
              <el-table-column
                prop="percentage"
                label="百分比"
                align="center">
                <template slot-scope="scope">{{scope.row.percentage}}%</template>
              </el-table-column>
            </el-table>
          </template>
        </div>
      </div>
      <Map :node="syncNode" :segment="syncSegment"></Map>
      <div class="main-table async-node-table">
        <div class="content">
          <h2 class="title">中继节点分布列表</h2>
          <template>
            <el-table
              :data="authNode"
              style="width: 100%">
              <el-table-column
                prop="ip"
                label="IP地址"
                align="center">
                <template slot-scope="scope">
                  <router-link class="router" :to="{name: 'node-status', params:{'id': scope.row.ip}}">{{scope.row.ip}}</router-link>
                </template>
              </el-table-column>
              <el-table-column
                prop="blockHeight"
                label="块高度"
                align="center">
              </el-table-column>
            </el-table>
          </template>
        </div>
      </div>
      <Map :node="authNode" :segment="authSegment"></Map>
    </div>
  </div>
</template>

<script>
import Statistics from './common/statistics'
import Loading from './common/loading'
import Map from './echarts/map'
import { getNode } from '@/api/interface'
export default {
  name: 'nodelist',
  data () {
    return {
      loading: true,
      syncNode: [],
      authNode: [],
      syncSegment: {
        title: '同步节点热力分布图',
        segmentData: [1, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100]
      },
      authSegment: {
        title: '中继节点热力分布图',
        segmentData: [1, 100000, 200000, 300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000],
        label: ['1', '10W', '20W', '30W', '40W', '50W', '60W', '70W', '80W', '90W', '100W']
      }
    }
  },
  components: {
    Statistics,
    Map,
    Loading
  },
  created () {
    this.getNode()
  },
  methods: {
    getNode () {
      getNode()
        .then(response => {
          const {followNodes = [], authNodes = []} = response.data;
          this.syncNode = followNodes.map((n) => {
            n.percentage = Math.round(n.percentage * 10000) / 100;
            return n;
          });
          this.authNode = authNodes;
          setTimeout(() => {
            this.loading = false
          }, 500)
        })
        .catch(err => {
          console.log(err)
          setTimeout(() => {
            this.loading = false
          }, 500)
        })
    }
  }
}
</script>

<style>
  .node-list {
    display: flex;
    flex-wrap: wrap;
    margin-left: 15%;
    overflow-y: auto;
    background-color: #323241;
  }
  .node-list .sync-node-table {
    padding: 20px;
  }
  .node-list .async-node-table {
    padding: 0 20px 20px;
  }
  .node-list .map {
    width: 98%;
    margin: 0 20px 20px;
    background-color: #404051;
  }
</style>
