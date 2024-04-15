#!/bin/bash

# 创建一个新的CSV文件
echo "Container,StartTime,EndTime" > log_info.csv

# 关闭和启动Docker容器循环100次
for i in {1..100}
do
    # 在后台运行gnome-terminal并执行相应的命令
    gnome-terminal -- bash -c "sleep 20; docker-compose stop; sleep 5"

    # 运行docker-compose up并将输出重定向到临时文件
    docker-compose up

    # 等待一段时间确保容器启动完成（可以根据实际情况调整等待时间）
    sleep 10

    # 获取日志最后10行并从中提取信息
    docker-compose logs --tail=10 | while read line
    do
        if [[ $line =~ .*"startTime".* ]]; then
            container=$(echo "$line" | awk -F ' ' '{print $1}' | sed 's/|//')
            startTime=$(echo "$line" | awk -F ' ' '{print $7}')
        elif [[ $line =~ .*"endTime".* ]]; then
            endTime=$(echo "$line" | awk -F ' ' '{print $7}')

            # 将信息追加到 CSV 文件中
            echo "$container,$startTime,$endTime" >> log_info.csv
        fi
    done
    
    # 停止容器
    docker-compose stop
    
    # 等待一段时间再进行下一次循环
    sleep 60
done

