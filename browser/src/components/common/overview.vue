<template>
  <div class="overview">
    <div class="content">
      <div class="head">
        <h2>{{overview.head}}</h2>
      </div>
      <div class="detail">
        <ul>
          <li v-for="info in overview.item" :key="info.title">
            <div class="data">
              <div class="txt">
                <span class="number">{{info.info || 0}}<sup class="tps" v-if="info.TPS">{{info.TPS}}</sup></span>
                <span class="title">{{info.title}}</span>
              </div>
              <div class="progress">
                <span class="percent" :style="{width: info.width + '%', backgroundColor: info.backgroundColor}"></span>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'overview',
  props: {
    overview: Object
  },
  data () {
    return {
    }
  }
}
</script>

<style>
  .overview {
    width: 100%;
    padding: 20px;
    background-color: #323241;
    box-sizing: border-box;
  }
  .overview .content {
    padding: 20px;
    box-shadow: 0 0 40px rgba(0, 0, 0, .2);
    background-color: #404051;
  }
  .overview .content .head h2 {
    font-size: 16px;
    font-weight: normal;
    color: #fff;
  }
  .overview .head h2 span {
    display: inline-block;
    padding-left: 6px;
    color: #61697f;
  }
  .overview .detail ul {
    display: flex;
    flex-wrap: wrap;
  }
  .overview .detail li {
    flex-basis: 50%;
  }
  .overview .detail li .data {
    padding: 20px 0;
  }
  .overview .detail li .data .txt {
    padding: 0 20px 0 50px;
  }
  .overview .detail li:nth-child(2n) .data .txt {
    border-left: 1px solid #9898b3;
  }
  .overview .detail li .data .number {
    display: block;
    font-size: 28px;
    font-weight: 700;
    color: #fff;
  }
  .overview .detail li .data .number sup {
    padding-left: 8px;
    font-size: 16px;
    font-weight: 400;
    color: #b3b3b9;
  }
  .overview .detail li .data .title {
    display: block;
    font-size: 12px;
    color: #9898b3;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .overview .detail li .data .progress {
    margin: 0 20px 0 50px;
    width: 30%;
    height: 4px;
    border-radius: 2px;
    background-color: #232633; /* fd5a40 */
  }
  .overview .detail li .data .progress .percent {
    display: block;
    height: 100%;
    margin-top: 8px;
    border-radius: 3px;
    background-color: #fd5a40;
  }
</style>
