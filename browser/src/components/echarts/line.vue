<template>
  <div class="chart">
    <div class="content">
      <h2 class="head">9小时交易历史</h2>
      <div class="line" ref="line"></div>
    </div>
  </div>
</template>

<script>
import echarts from 'echarts/lib/echarts'
require('echarts/lib/chart/line')
require('echarts/lib/component/title')
require('echarts/lib/component/tooltip')
export default {
  name: 'lineChart',
  props: {
    line: {
      type: Array
    }
  },
  computed: {
    getFormatData () {
      let xData = []
      let yData = []
      this.line.forEach(item => {
        xData.push(new Date(item.time * 1000).getHours() + ':00')
        yData.push(item.txCount)
      })
      return { xData, yData }
    }
  },
  methods: {
    /* __getBeforeDay (before) {
      let time = new Date()
      let time1 = new Date()
      let data = []
      // 获取before天前的日期
      if (before) {
        let time2 = null
        for (let i = before; i > 0; i--) {
          time2 = time1.setTime(time.getTime() - i * 24 * 60 * 60 * 1000)
          var month1 = new Date(time2).getMonth() + 1
          var day1 = new Date(time2).getDate()
          data.push(month1 + '/' + day1)
        }
      } else {
        time1 = time.setTime(time.getTime())
      }
      return data
    }, */
    getLine () {
      let myChart = echarts.init(this.$refs.line)
      window.onresize = function () {
        myChart.resize()
      }
      let option = {
        backgroundColor: '#414052',
        tooltip: {
          position: 'top',
          backgroundColor: '#555368',
          padding: 10,
          // formatter: `new Date({b})<br />{a}:<b>{c}</b>`
          formatter: function (params) {
            let week = ['Sun', 'Mon', 'Tues', 'Wed', 'Thur', 'Fri', 'Sat']
            let addrMonth = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sept', 'Oct', 'Nov', 'Dec']
            let time = new Date()
            let formateTime = new Date(time)
            let month = addrMonth[formateTime.getMonth()]
            let date = formateTime.getDate()
            let day = week[formateTime.getDay()]
            let res = '<p class="date">' + day + '&nbsp;' + date + ',' + month + '&nbsp;' + time.getFullYear() + '&nbsp;&nbsp;' + params['name'] + '</p>'
            res += '<p class="name">' + params['seriesName'] + ':' + '<b>' + params['data'] + '</b>' + '</p>'
            return res
          }
        },
        grid: {
          left: '12%',
          top: '10%',
          right: '8%'
        },
        xAxis: {
          axisLine: {
            lineStyle: {
              color: '#535262'
            }
          },
          axisLabel: {
            color: '#9c9db2',
            padding: [15, 0, 0, 0]
          },
          inverse: true,
          data: this.getFormatData.xData
        },
        yAxis: {
          axisLine: {
            lineStyle: {
              color: 'transparent'
            }
          },
          splitNumber: 3,
          axisLabel: {
            color: '#9c9db2',
            padding: [0, 15, 0, 0]
          },
          splitLine: {
            lineStyle: {
              color: '#535262'
            }
          }
        },
        series: {
          name: 'Transaction',
          type: 'line',
          symbol: 'circle',
          symbolSize: 10,
          itemStyle: {
            borderColor: '#343539',
            borderWidth: 2,
            color: '#ff5843'
          },
          lineStyle: {
            color: '#ff5843',
            width: 2
          },
          data: this.getFormatData.yData
        }
      }
      myChart.setOption(option)
      // 给Echarts提示框添加class
      let div2 = this.$refs.line.getElementsByTagName('div')[1]
      div2.setAttribute('class', 'tooltip')
    }
  },
  watch: {
    'line': 'getLine'
  },
  beforeRouteLeave () {
    echarts.dispose()
  }
}
</script>

<style>
  .chart {
    box-sizing: border-box;
    width: 100%;
    padding: 20px 20px 20px 0;
    background-color: #323241;
  }
  .chart .content {
    background-color: #404051;
    box-shadow: 0 0 40px rgba(0, 0, 0, .2);
  }
  .chart .line {
    width: 98%;
    height: 250px;
  }
  .chart .head {
    padding: 20px 0 0 20px;
    color: #fff;
    font-size: 16px;
    font-weight: normal;
  }
  .chart .line .tooltip {
    position: absolute;
    top: -1000px;
    left: 0;
    border-radius: 0 !important;
    box-shadow: 0 0 50px rgba(0, 0, 0, .2);
  }
  .chart .line .tooltip:before {
    position: absolute;
    bottom: -14px;
    left: 40%;
    border-left: 9px solid transparent;
    border-top: 9px solid #555368;
    border-bottom: 9px solid transparent;
    border-right: 9px solid transparent;
    content: '';
  }
  .chart .line .tooltip .date {
    color: #9999b3;
    font-size: 12px;
  }
  .chart .line .tooltip .name {
    font-size: 14px;
    font-weight: 200;
  }
  .chart .line .tooltip .name b {
    padding-left: 5px;
    font-size: 16px;
  }
</style>
