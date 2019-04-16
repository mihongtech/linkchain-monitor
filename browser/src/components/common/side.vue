<template>
  <div class='side-bar'>
    <div class='logo'>
      <h1>兑豆豆</h1>
      <img src='@/assets/image/logo.png' alt='兑豆豆'>
    </div>
    <div class='nav'>
      <ul ref='nav'>
        <li
          v-for='(route, index) in routerList'
          :key='index'>
          <router-link
            :to='route.to'
            @click.native='changeNum(index)'
            :class="{'active': index === num}">
            <span class='icon'></span>
            <i class='iconfont' :class="'icon-' + route.font"></i>
            {{route.name}}</router-link>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import Statistics from './statistics'
export default {
  name: 'linkchain',
  data () {
    return {
      routerList: [
        {name: '链状态', to: '/', font: 'linkchain'},
        {name: '节点分布', to: '/node/list', font: 'node'}
      ],
      num: 0
    }
  },
  methods: {
    changeNum (value) {
      this.num = value
      sessionStorage.setItem('num', this.num)
    },
    getNum () {
      let num = sessionStorage.getItem('num')
      if (num) {
        this.num = num - 0
      } else {
        return 0
      }
    }
  },
  components: {
    Statistics
  },
  created () {
    this.getNum()
  }
}
</script>

<style>
  @import '../../assets/font/iconfont.css';
  .side-bar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 15%;
    height: 100%;
    background-color: #282834;
    z-index: 10000;
  }
  .side-bar .logo {
    display: table-cell;
    height: 25%;
  }
  .side-bar .logo h1 {
    height: 0;
    text-indent: -999px;
  }
  .side-bar .logo img {
    width: 100px;
    padding: 80px 0 50px 30px;
    vertical-align: middle;
  }
  .side-bar .nav li {
    position: relative;
    height: 60px;
  }
  .side-bar .nav li span {
    position: absolute;
    width: 36px;
    height: 36px;
    margin: 12px 0;
  }
  .side-bar .nav li .iconfont {
    position: relative;
    left: 9px;
    font-size: 20px;
    margin-right: 25px;
  }
  .side-bar .nav li a {
    display: block;
    height: 60px;
    padding-left: 30px;
    line-height: 60px;
    color: #9398a0;
  }
  .side-bar .nav li a:hover {
    background-color: #323241;
    color: #dedede;
  }
  .side-bar .nav li a.active {
    background-color: #323241;
    color: #fff;
  }
  .side-bar .nav li a.active span {
    border-radius: 50%;
    border: 1px solid #20202a;
    background: rgb(56, 183, 217);
  }
</style>
