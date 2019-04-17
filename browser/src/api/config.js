import axios from 'axios'
const config = axios.create({
  baseURL: 'http://122.112.249.47/api/v1/explorer',
  timeout: 10000
})
export default config
