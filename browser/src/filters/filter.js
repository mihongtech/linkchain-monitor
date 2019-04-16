const filters = {
  getFormatDate: (date) => {
    let time = date * 1000
    time = new Date(time)
    let year = time.getFullYear()
    let month = time.getMonth() + 1
    month = month >= 10 ? month : '0' + month
    let day = time.getDate()
    day = day >= 10 ? day : '0' + day
    let hour = time.getHours()
    hour = hour >= 10 ? hour : '0' + hour
    let minutes = time.getMinutes()
    minutes = minutes >= 10 ? minutes : '0' + minutes
    let second = time.getSeconds()
    second = second >= 10 ? second : '0' + second
    return `${year}-${month}-${day} ${hour}:${minutes}:${second}`
  },
  getFormatNumber: (number) => {
    return parseFloat(number).toLocaleString() + '.00'
  },
  getPercent: (percent) => {
    return Math.floor(percent * 10000) / 100 + '%'
  },
  getFormatDay: (number) => {
    return `${Math.floor((number / (24 * 60 * 60)))}天${Math.ceil((number % (24 * 60 * 60)) / 3600)}小时`
  }
}
export default filters
