import pandas as pd
import re
import matplotlib.pyplot as plt

# 读取数据
# df = pd.read_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/origin_data/TCP-tls-QUIC.csv', delimiter=',')
df = pd.read_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/origin_data/TCP_CPU.csv', delimiter=',')
print(df)

# 分割数据为TCP和QUIC
tcp_data = df[df['Container'].str.contains('tcp', case=False)].copy()
quic_data = df[df['Container'].str.contains('quic', case=False)].copy()

# 处理TotalTime的数据
def process_time(time_str):
    match = re.search(r'(\d+\.?\d*)\s*([µm]?s)', time_str)
    if match:
        value = float(match.group(1))
        unit = match.group(2)
        if unit == 'µs' or unit == 'us':
            return value / 1000
        elif unit == 'ms':
            return value
        else:
            return value * 1000
    else:
        return None

tcp_data['TotalTime'] = tcp_data['TotalTime'].apply(process_time)
quic_data['TotalTime'] = quic_data['TotalTime'].apply(process_time)

# 存储数据
tcp_data.to_csv('tcp_data.csv', index=False)
quic_data.to_csv('quic_data.csv', index=False)
# tcp_data.to_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/TCP_tls_link_result.csv', index=False)
# quic_data.to_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/QUIC_tls_link_result.csv', index=False)
tcp_data.to_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/TCP_link_result.csv', index=False)
quic_data.to_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/QUIC_link_result.csv', index=False)

# 读取第二个文件
tcp_tls_link_data = pd.read_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/TCP_tls_link_result.csv')
quic_link_data = pd.read_csv('D:/Users/LLL2000/Desktop/analyze(second)/data/result/QUIC_tls_link_result.csv')

# 设置X轴为次数
x_ticks = list(range(1, len(tcp_data) + 1))

# 绘制对比图表
plt.figure(figsize=(25, 8))
plt.plot(x_ticks[:len(tcp_data)], tcp_data['TotalTime'], marker='o', label='TCP')  # 修改这一行
plt.plot(x_ticks[:len(tcp_tls_link_data)], tcp_tls_link_data['TotalTime'], marker='o', label='TCP+TLS1.3')
plt.plot(x_ticks[:len(quic_link_data)], quic_link_data['TotalTime'], marker='o', label='QUIC')
plt.xlabel('Count', fontsize=30)
plt.ylabel('TotalTime (ms)', fontsize=30)
plt.title('Comparison of TotalTime between TCP, TCP+TLS1.3 and QUIC', fontsize=30)
plt.legend(fontsize=20)

# 设置x轴和y轴刻度的字体大小
plt.xticks(fontsize=20)
plt.yticks(fontsize=20)

plt.show()


