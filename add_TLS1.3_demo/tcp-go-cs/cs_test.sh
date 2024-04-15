#!/bin/bash

# 创建一个新的CSV文件
echo "Container,Time,TotalTime,CPU Time,System Time" > log_info.csv

# 关闭和启动Docker容器循环100次
for i in {1..100}
do
    # 在后台运行gnome-terminal并执行相应的命令
    gnome-terminal -- bash -c "sleep 20; docker-compose stop; sleep 5"

    # 运行docker-compose up并将输出重定向到临时文件，并使用time命令测量运行时间
    { time docker-compose up; } 2> time_output.txt
    
    # 等待3s确保容器启动完成（可以根据实际情况调整等待时间）
    sleep 2

    # 获取日志最后10行（因为如果限制行数，就会将之前的信息也获取，时间就会重复）并从中提取信息
    cat time_output.txt | grep "TotalTime" | sort | while read line
    do
        container=$(echo "$line" | awk -F ' ' '{print $1}' | sed 's/|//')
        time=$(echo "$line" | awk -F ' ' '{print $3,$4}')
        TotalTime=$(echo "$line" | awk -F ' ' '{print $6}')

        # 强制刷新磁盘缓冲区
        sync

        # 从time_output.txt中提取CPU耗时时间并输出
        cpu_time=$(cat time_output.txt | grep "user" | awk '{print $2}')
        sys_time=$(cat time_output.txt | grep "sys" | awk '{print $2}')
        echo "CPU Time: $cpu_time"
        echo "System Time: $sys_time"

        # 将信息追加到 CSV 文件中
        echo "$container,$time,$TotalTime,$cpu_time,$sys_time" >> log_info.csv
    done
    
    # 停止容器
    docker-compose stop
    
    # 等待一段时间再进行下一次循环
    sleep 20
done


