<template>
  <div class="bar-chart">
    <div class="content">
      <h2 class="head">商家交易量排行榜</h2>
      <div class="bar" ref="bar"></div>
    </div>
  </div>
</template>

<script>
import echarts from 'echarts/lib/echarts'
require('echarts/lib/chart/bar')
require('echarts/lib/component/title')
require('echarts/lib/component/tooltip')
export default {
  name: 'x-bar',
  props: {
    rank: Array
  },
  computed: {
    computedRank () {
      const rank = this.rank.slice()
      return rank.sort((item1, item2) => {
        return item1.txCount < item2.txCount
      })
    }
  },
  methods: {
    getBar () {
      let yData = []
      this.computedRank.forEach(item => {
        yData.push(item.name)
      })
      let xData = []
      this.computedRank.forEach(item => {
        xData.push(item.txCount)
      })
      let Bar = echarts.init(this.$refs.bar)
      let option = {
        tooltip: {
          position: 'top',
          backgroundColor: '#555368',
          padding: [5, 10],
          formatter: '<p class="good">{b}</p><span>{a}:</span><b>{c}</b>'
        },
        xAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: 'transparent'
            }
          },
          splitLine: {
            show: false
          }
        },
        grid: {
          left: '14%',
          top: '16%',
          right: '8%',
          bottom: '20%'
        },
        yAxis: {
          type: 'category',
          inverse: true,
          axisLine: {
            lineStyle: {
              color: 'transparent'
            }
          },
          boundaryGap: false,
          splitLine: {
            show: true,
            lineStyle: {
              color: '#535262'
            }
          },
          axisLabel: {
            color: '#9c9db2',
            padding: [0, 10, 0, 0]
          },
          data: yData
        },
        series: {
          name: '今日销售',
          type: 'bar',
          barWidth: 16,
          itemStyle: {
            color: function (params) {
              let colorList = ['#e95740', '#e7b52e', '#5e92e8', '#7a7a7a']
              if (params.dataIndex >= 3) {
                return colorList[3]
              } else {
                return colorList[params.dataIndex]
              }
            }
          },
          // barCategoryGap: '30%',
          data: xData
        }
      }
      Bar.setOption(option)
      // 给Echarts提示框添加class
      let div2 = this.$refs.bar.getElementsByTagName('div')[1]
      div2.setAttribute('class', 'tooltip')
    }
  },
  watch: {
    'rank': 'getBar'
  },
  beforeRouteLeave () {
    echarts.dispose()
  }
}
</script>

<style>
  .bar-chart {
    box-sizing: border-box;
    padding: 20px 20px 20px 0;
    background-color: #323241;
  }
  .bar-chart .content {
    box-shadow: 0 0 40px rgba(0, 0, 0, .2);
    background-color: #404051;
  }
  .bar-chart .head {
    padding: 20px 0 0 20px;
    font-size: 16px;
    color: #fff;
    font-weight: normal;
  }
  .bar-chart .bar {
    width: 100%;
    height: 250px;
  }
  .bar-chart .bar .tooltip {
    position: absolute;
    top: -100%;
    left: 0;
    border-radius: 0 !important;
    box-shadow: 0 0 50px rgba(0, 0, 0, .2);
  }
  .bar-chart .bar .tooltip:before {
    position: absolute;
    bottom: -14px;
    left: 40%;
    border-left: 9px solid transparent;
    border-top: 9px solid #555368;
    border-bottom: 9px solid transparent;
    border-right: 9px solid transparent;
    content: '';
  }
  .bar-chart .bar .tooltip .good {
    color: #9999b3;
    font-size: 12px;
  }
  .bar-chart .bar .tooltip span {
    font-size: 12px;
    color: #fff;
  }
  .bar-chart .bar .tooltip b {
    padding-left: 5px;
    font-size: 14px;
  }
</style>
