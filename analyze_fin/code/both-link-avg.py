import pandas as pd
import matplotlib.pyplot as plt

# 读取三个文件并计算平均值
file_paths = ['D:/Users/LLL2000/Desktop/analyze(second)/data/result/QUIC_tls_link_result.csv',
              'D:/Users/LLL2000/Desktop/analyze(second)/data/result/TCP_link_result.csv',
              'D:/Users/LLL2000/Desktop/analyze(second)/data/result/TCP_tls_link_result.csv',]

average_times_QUIC_tls_link = []
average_times_TCP_link = []
average_times_TCP_tls_link = []

for file_path in file_paths:
    df = pd.read_csv(file_path)
    average_time = df['TotalTime'].mean()
    if 'QUIC_tls_link' in file_path:
        average_times_QUIC_tls_link.append(average_time)
    elif 'TCP_link' in file_path:
        average_times_TCP_link.append(average_time)
    elif 'TCP_tls_link' in file_path:
        average_times_TCP_tls_link.append(average_time)

labels = ['TCP', 'TCP+TLS1.3', 'QUIC']
x = range(len(labels))

fig, ax = plt.subplots()
bar_width = 0.55

bar1 = ax.bar(x[0], average_times_TCP_link[0], bar_width, label='TCP')
bar2 = ax.bar(x[1], average_times_TCP_tls_link[0], bar_width, label='TCP+TLS1.3')
bar3 = ax.bar(x[2], average_times_QUIC_tls_link[0], bar_width, label='QUIC')

ax.set_xlabel('Protocol')
ax.set_ylabel('Average Time (ms)')
ax.set_title('Comparison of Average Time')
ax.set_xticks(x)
ax.set_xticklabels(labels)
ax.legend()

# 添加平均值标记
for rect in bar1:
    ax.text(rect.get_x() + rect.get_width() / 2, rect.get_height(), f'{rect.get_height():.5f}', ha='center', va='bottom')
for rect in bar2:
    ax.text(rect.get_x() + rect.get_width() / 2, rect.get_height(), f'{rect.get_height():.5f}', ha='center', va='bottom')
for rect in bar3:
    ax.text(rect.get_x() + rect.get_width() / 2, rect.get_height(), f'{rect.get_height():.5f}', ha='center', va='bottom')

plt.show()
