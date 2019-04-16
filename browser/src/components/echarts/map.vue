<template>
  <div class='map'>
    <h2 class='head'>{{segment.title}}</h2>
    <div class='map-china' ref='china'></div>
  </div>
</template>

<script>
import echarts from 'echarts/lib/echarts'
require('echarts/lib/chart/map')
require('echarts/map/js/china')
require('echarts/map/js/province/hebei')
export default {
  name: 'map-china',
  props: {
    node: Array,
    segment: Object
  },
  mounted () {
    this.getMap()
  },
  computed: {
    // 整理后台数据到指定格式
    mapData () {
      let arr = []
      this.node.forEach(item => {
        if (item.percentage) {
          arr.push({'name': item.city, 'value': [item.lon, item.lat, item.percentage]})
        } else {
          arr.push({'name': item.city, 'value': [item.lon, item.lat, item.blockHeight]})
        }
      })
      return arr
      // 9398a0
    },
    // 格式化父组件传来的分段数据
    segmentData () {
      let segmentArr = []
      let i = 0
      for (i; i < this.segment.segmentData.length - 1; i++) {
        if (this.segment.label) {
          segmentArr.push(
            {
              start: this.segment.segmentData[i],
              end: this.segment.segmentData[i + 1],
              label: this.segment.label[i] + '-' + this.segment.label[i + 1]
            })
        } else {
          segmentArr.push({start: this.segment.segmentData[i], end: this.segment.segmentData[i + 1]})
        }
      }
      segmentArr.unshift({start: 0, end: 0, label: 'no data'})
      return segmentArr
    }
  },
  methods: {
    getMap () {
      let Map = echarts.init(this.$refs.china)
      let option = {
        geo: {
          map: 'china',
          itemStyle: {
            normal: {
              areaColor: '#353448',
              borderColor: '#1c1f2c'
            },
            emphasis: {
              areaColor: '#15152a'
            }
          },
          label: {
            emphasis: {
              show: true,
              textStyle: {
                color: '#fff'
              }
            }
          }
        },
        tooltip: {
          trigger: 'item',
          position: 'top',
          backgroundColor: '#555368',
          formatter: function (params) {
            let value = params['value'][2]
            value = value > 100 ? value : value + '%'
            return params.name + '&nbsp;: &nbsp;' + value
          }
        },
        backgroundColor: '#404051',
        visualMap: {
          left: 16,
          bottom: 20,
          splitList: this.segmentData,
          color: ['#ff2b00', '#ef5434', '#f5795f', '#ff881f', '#ffa04c', '#ffb675', '#84c750', '#95c76e', '#858585'],
          textStyle: {
            color: '#d2d2ed'
          },
          itemGap: 1,
          itemWidth: 6,
          itemHeight: 20
        },
        series: [
          {
            name: '中继节点分布散点图',
            type: 'scatter',
            coordinateSystem: 'geo',
            data: this.mapData,
            symbolSize: 10,
            itemStyle: {
              emphasis: {
                borderColor: '#fff',
                borderWidth: 1
              }
            }
          }
        ]
      }
      Map.setOption(option)
      // 修改tooltip
      let div2 = this.$refs.china.getElementsByTagName('div')[1]
      div2.setAttribute('class', 'tooltip')
    }
  },
  watch: {
    'node': 'getMap'
  },
  beforeRouteLeave () {
    echarts.dispose()
  }
}

</script>

<style>
  .map .map-china {
    width: 98%;
    height: 400px;
  }
  .map .head {
    padding: 20px 0 0 20px;
    font-size: 16px;
    color: #fff;
    font-weight: normal;
  }
  .map .map-china .tooltip {
    position: absolute;
    top: -2000px;
    left: 0;
    border-radius: 0 !important;
    box-shadow: 0 0 50px rgba(0, 0, 0, .2);
  }
  .map .map-china .tooltip:before {
    position: absolute;
    bottom: -14px;
    left: 40%;
    border-left: 9px solid transparent;
    border-top: 9px solid #555368;
    border-bottom: 9px solid transparent;
    border-right: 9px solid transparent;
    content: '';
  }
</style>
