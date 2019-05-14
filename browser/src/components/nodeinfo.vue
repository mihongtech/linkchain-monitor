<template>
  <div class="root">
    <Loading v-loading="loading" v-if="loading"/>
    <div class="node-status" v-else>
      <div class="info basic">
        <div class="content">
          <h2 class="head">权威节点基本信息</h2>
          <div class="detail">
            <table>
              <tbody>
              <tr>
                <td>IP地址：<span>{{baseInfo.ip}}</span></td>
                <td>操作系统：<span><!--{{baseInfo.os}}-->Ubuntu 16.04.10</span></td>
              </tr>
              <tr>
                <td>系统时间：<span>{{baseInfo.sysTime | getFormatDate}}</span></td>
                <td>磁盘利用率：<span>{{baseInfo.diskUsage | getPercent}}</span></td>
              </tr>
              <tr>
                <td>CPU利用率：<span>{{baseInfo.cpuUsage | getPercent}}</span></td>
                <td>内存利用率：<span>{{baseInfo.memUsage | getPercent}}</span></td>
              </tr>
              <tr>
                <td>运行时间：<span>{{baseInfo.runningTime | getFormatDay}}</span></td>
                <td>链使用软件版本：<span><!--{{baseInfo.gethVersion}}-->V0.0.2.0</span></td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <div class="info block">
        <div class="content">
          <h2 class="head">权威节点区块信息</h2>
          <div class="detail">
            <table>
              <tbody>
              <tr>
                <td>块高度：<span>{{linkchainInfo.blockHeight}}</span></td>
                <td>跟最新block的偏差：<span>{{linkchainInfo.blockDiff}}</span></td>
              </tr>
              <tr>
                <td></td>
                <td></td>
              </tr><tr>
                <td></td>
                <td></td>
              </tr>
              <tr>
                <td></td>
                <td></td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Statistics from './common/statistics'
import Loading from './common/loading'
import { getNodeInfo } from '@/api/interface'
export default {
  name: 'node-info',
  data () {
    return {
      loading: true,
      baseInfo: {},
      linkchainInfo: {},
      block: [],
      transaction: []
    }
  },
  components: {
    Statistics,
    Loading
  },
  created () {
    this.getNodeInfo()
  },
  methods: {
    getNodeInfo () {
      let ip = this.$route.params.id
      getNodeInfo(ip)
        .then(response => {
          this.baseInfo = response.data.baseInfo
          this.linkchainInfo = response.data.linkchainInfo
          document.title = this.baseInfo.ip
          setTimeout(() => {
            this.loading = false
          }, 500)
        })
        .catch(() => {
          setTimeout(() => {
            this.loading = false
          }, 500)
        })
    }
  }
}
</script>

<style>
  .node-status {
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
    margin-left: 15%;
  }
  .node-status .main-table {
    box-sizing: border-box;
    width: 50%;
  }
  .node-status .block-table {
    padding-left: 20px;
  }
  .node-status .basic {
    box-sizing: border-box;
    width: 50%;
    background-color: #323241;
    padding: 20px;
  }
  .node-status .block{
    box-sizing: border-box;
    width: 50%;
    background-color: #323241;
    padding: 20px 20px 20px 0;
  }
  .node-status .info .content {
    background-color: #414052;
    box-shadow: 0 0 40px rgba(0, 0, 0, .2);
  }
  .node-status .info .content .head {
    height: 40px;
    font-size: 14px;
    padding-left: 20px;
    line-height: 40px;
    font-weight: 500;
    color: #fff;
  }
  .node-status .info .content table {
    width: 100%;
    text-align: left;
    border-collapse: collapse;
    border-spacing: 0;
  }
  .node-status .info .content table td {
    height: 45px;
    padding-left: 20px;
    font-size: 12px;
    color: #9898b3;
    font-weight: 500;
  }
  .node-status .info .content table td span {
    display: inline-block;
    padding-left: 10px;
    color: #fff;
  }
  .node-status .info .content .detail tr:nth-child(2n-1){
    background-color: #32323f;
  }
</style>
