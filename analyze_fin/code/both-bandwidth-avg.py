import pandas as pd
import matplotlib.pyplot as plt

# 读取三个文件并计算平均值
file_paths = ['D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_NO_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_10G_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_100M_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_1M_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_NO_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_10G_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_100M_result.csv',
              'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_1M_result.csv',]

average_times_QUIC_NO = []
average_times_QUIC_10G = []
average_times_QUIC_100M = []
average_times_QUIC_1M = []
average_times_TCP_NO = []
average_times_TCP_10G = []
average_times_TCP_100M = []
average_times_TCP_1M = []

for file_path in file_paths:
    df = pd.read_csv(file_path)
    average_time = df['HolepunchTime'].mean()
    if 'QUIC_NO' in file_path:
        average_times_QUIC_NO.append(average_time)
    elif 'QUIC_10G' in file_path:
        average_times_QUIC_10G.append(average_time)
    elif 'QUIC_100M' in file_path:
        average_times_QUIC_100M.append(average_time)
    elif 'QUIC_1M' in file_path:
        average_times_QUIC_1M.append(average_time)
    elif 'TCP_NO' in file_path:
        average_times_TCP_NO.append(average_time)
    elif 'TCP_10G' in file_path:
        average_times_TCP_10G.append(average_time)
    elif 'TCP_100M' in file_path:
        average_times_TCP_100M.append(average_time)
    elif 'TCP_1M' in file_path:
        average_times_TCP_1M.append(average_time)

# 绘制柱状图
labels = ['No-Limit', '10G','100M', '1M']
x = range(len(labels))

fig, ax = plt.subplots(figsize=(10, 7))
bar_width = 0.35

bar1 = ax.bar(x, [average_times_QUIC_NO[0],average_times_QUIC_10G[0], average_times_QUIC_100M[0], average_times_QUIC_1M[0]], bar_width, label='QUIC', alpha=0.6, color='blue')
bar2 = ax.bar([i + bar_width for i in x], [average_times_TCP_NO[0],average_times_TCP_10G[0], average_times_TCP_100M[0], average_times_TCP_1M[0]], bar_width, label='TCP',alpha=0.6,color='orange')

ax.set_xlabel('Bandwidth Limit',fontsize=14)
ax.set_ylabel('Average Hole Punching Time (ms)',fontsize=14)
ax.set_title('Comparison of Average Hole Punching Time under Different Bandwidths',fontsize=14)
ax.set_xticks([i + bar_width/2 for i in x])
ax.set_xticklabels(labels,fontsize=14)
ax.legend()

# 添加平均值标记
for rect in bar1:
    ax.text(rect.get_x() + rect.get_width() / 2, rect.get_height(), f'{rect.get_height():.3f}', ha='center', va='bottom',fontsize=14)
for rect in bar2:
    ax.text(rect.get_x() + rect.get_width() / 2, rect.get_height(), f'{rect.get_height():.3f}', ha='center', va='bottom',fontsize=14)

plt.show()

#保存留片为PNG格式
fig.savefig("D:/Users/LLL2000/Desktop/lunwen/picture/fig6.png",dpi=600)
# 设置dpi参数可以调整图片的清晰度
