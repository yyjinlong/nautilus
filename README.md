ferry
-------------
Jinlong Yang

# ferry

## 1 why blue-green？

    1 采用滚动更新, maxSurge设置为25%.

        如果pod比较多, 就会进行多次的滚动更新, 便会产生新老pod共存, 共同接流量情况. 然而这是业务不期望的情况.

    2 蓝绿部署的优势

        另一组部署完成后, 流量一刀切.

    3 蓝绿部署的缺点

        瞬间对集群造成一些压力, 但是还好, 发布完成后, 会对旧版本的deployment的pod数量缩成0, 使其不占资源.


## 2 镜像分层

    1) base层    : centos6.7 centos7.5
    2) run层     : python(conda环境)、go(go mod环境)
    3) service层 : 具体的服务环境

    4) release层 : 由ferryd进程基于代码自动构建


## 3 节点标签

    kubectl label node x.x.x.x aggregate=default


## 4 创建service

    1) 创建服务

        curl -d 'service=ivr' http://127.0.0.1:8888/v1/service


## 5 创建deployment

    1) 创建pipeline

        curl -H 'content-type: application/json' -d '{"name": "ivr test", "summary": "test", "service": "ivr",  "module_list": [{"name": "ivr", "branch": "master"}], "creator": "yangjinlong", "rd": "yangjinlong", "qa": "yangjinlong", "pm": "yangjinlong"}' http://127.0.0.1:8888/v1/pipeline

    2) 打tag

        curl -d 'pipeline_id=4&module=ivr&tag=release_ivr_20210827_155942' http://127.0.0.1:8888/v1/tag

    3) 创建镜像

        curl -d 'pipeline_id=4'  http://127.0.0.1:8888/v1/image

    4) 发布沙盒

        curl -d "pipeline_id=4&phase=sandbox&username=yangjinlong" http://127.0.0.1:8888/v1/deployment | jq .

    5) 发布全量

        curl -d "pipeline_id=4&phase=online&username=yangjinlong" http://127.0.0.1:8888/v1/deployment | jq .

    6) 部署完成

        curl -d "pipeline_id=4" http://127.0.0.1:8888/v1/finish | jq .

