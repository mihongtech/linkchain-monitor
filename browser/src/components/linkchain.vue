<template>
  <div class="root">
    <Loading v-loading="loading" v-if="loading"/>
    <div class="linkchain" v-else>
      <Overview :overview="overview"></Overview>
    </div>
  </div>
</template>

<script>
  import Overview from './common/overview'
  import Statistics from './common/statistics'
  import Loading from './common/loading'
  import Table from './common/table'
  import LineChart from './echarts/line'
  import {getLinkChain} from '@/api/interface'

  export default {
    name: 'linkchain',
    components: {
      Overview,
      Statistics,
      LineChart,
      Loading,
      Table
    },
    data() {
      return {
        overview: {
          head: '业绩总览',
          item: {
            blockHeight: {title: '块高度', width: 30, backgroundColor: '#762cbf', info: ''},
            lastHourTxs: {title: '近1H交易量', width: 60, backgroundColor: '#5a86eb', info: '', TPS: ''},
            authNodeCount: {title: '权威节点', width: 50, backgroundColor: '#ddb63e', info: ''},
            followNodeCount: {title: '同步节点', width: 70, backgroundColor: '#87f8ab', info: ''}
          }
        },
        loading: true
      }
    },
    created() {
      this.getLinkChain()
    },
    methods: {
      getLinkChain() {
        getLinkChain()
          .then(response => {
            this.overview.item.blockHeight.info = response.data.overview.blockHeight
            this.overview.item.lastHourTxs.info = response.data.overview.lastHourTxs
            this.overview.item.authNodeCount.info = response.data.overview.authNodeCount
            this.overview.item.followNodeCount.info = response.data.overview.followNodeCount
            this.overview.item.lastHourTxs.TPS = (response.data.overview.lastHourTxs / 3600).toFixed(2) + 'TPS'
            this.statistics = response.data.statistics
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
  .linkchain {
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
    margin-left: 15%;
  }
</style>
