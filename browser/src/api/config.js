import axios from 'axios'
const config = axios.create({
  baseURL: 'http://10.150.17.95:8090/api/v1/explorer',
  timeout: 10000
})
export default config
