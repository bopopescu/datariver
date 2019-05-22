# registry.cn-hangzhou.aliyuncs.com/medlinker/datariver
FROM registry-vpc.cn-hangzhou.aliyuncs.com/datariver/alpine:3.8

LABEL maintainer="guoqiang@medlinker.com"

###############################################################################
#                                INSTALLATION
###############################################################################

# 环境变量设置
ENV APP_NAME datariver
ENV APP_ROOT /var/www
ENV APP_PATH $APP_ROOT/$APP_NAME
ENV LOG_ROOT /var/log/medlinker
ENV LOG_PATH /var/log/medlinker/$APP_NAME


# 执行入口文件添加
ADD ./main $APP_PATH/
ADD ./docfile/dockerfiles/*.sh /bin/
RUN chmod +x /bin/*.sh

###############################################################################
#                                   START
###############################################################################
