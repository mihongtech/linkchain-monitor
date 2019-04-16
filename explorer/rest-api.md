linkchain explorer Restful API 接口文档
===

## 1. 链状态相关API接口 ##
### 1.1 查询链的状态Overview ###
   
    URL：
    GET  /api/v1/explorer/linkchian/overview
    
    PARAM:
           无

    
    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"overview": {"blockHeight":"4545", "authNodeCount":"54", "followNodeCount":"3456", "lastHourTxs":"456456"}，"statistics": [{"time":"159000000", "txCount":"345"},{"time":"158012345", "txCount":"145"}]}
  

## 2. 商家状态相关API接口 ##
### 2.1 查询商家的状态Overview ###
   
    URL：
    GET  /api/v1/explorer/business/overview
    
    PARAM:
    无

    
    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"overview": {"businessCount":"45", "userCount":"1154", "lastHourTxs":"456456"}，"statistics": [{"name":"稻香村", "txCount":"345"},{"name":"稻花香", "txCount":"145"}]}



### 2.2 查询商家列表 ###
   
    URL：
    GET  /api/v1/explorer/business/list
    
    PARAM:
    无

    
    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"list": [{"name":"稻香村", "addrss","ox435345445345","createTime":"1593454354", "userCount":"34543", "balance":"43242", "releasedTokens":"43543"},{"name":"稻花香", "addrss","ox435345445345","createTime":"1593454354","userCount":"34543", "balance":"43242", "releasedTokens":"43543"}]}


### 2.3 查询商家信息 ###
   
    URL：
    GET  /api/v1/explorer/business/info/{address}
    
    PARAM:
    address   必选输入   商家地址

    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"info": {"name":"稻香村", "addrss","ox435345445345","createTime":"1593454354", "userCount":"34543", "balance":"43242", "releasedTokens":"43543", "rate":"5", "lastHourTxs":"45345345","owner":"0x345345435", "operator":"0x766575675"}}



## 3. 节点状态相关API接口 ##
### 3.1 查询链的节点信息 ###
    URL：
    GET  /api/v1/explorer/linkchian/nodes
    
    PARAM:
    无

    
    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"followNodes": [{"region":"025", "count":"345", "percentage":"78%"},{"region":"010", "count":"105", "percentage":"22%"}], "authNodes":[{"ip":"10.1.12.3", "address":"0x23424324", "blockHeight":"453453", "region":"010"}, {"ip":"10.1.12.5", "address":"0x111123424324", "blockHeight":"453453", "region":"021"}]} 

### 3.2 查询链的权威节点信息 ###
    URL：
    GET  /api/v1/explorer/linkchian/authnodes/{node}
    
    PARAM:
    node   必选输入   权威节点的IP地址

    
    RETURN:

    Failed： 
    HTTP.CODE  500
    HTTP.BODY  {"error_msg": "error msg..."}

    Success:
    HTTP.CODE  200
    HTTP.BODY  {"baseInfo": {"ip":"10.1.12.3", "os":"Centos 7", "sysTime":"1590001000", "runningTime":"14 day 13 hours", "gethVersion":"v1.8.3-stable", "diskUsage":"78%","memUsage":"34%", "cpuUsage":"45%"}, "linkchainInfo":{"address":"0x23424324324", "blockHeight":"3453543", "nextBlockTime":"1594353535","blockDiff":"1"}} 