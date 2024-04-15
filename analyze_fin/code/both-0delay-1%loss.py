import pandas as pd
import matplotlib.pyplot as plt

# 读取TCP_20_0.csv
df_tcp = pd.read_csv('D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_0_1_result.csv')
# 读取QUIC_20_0.csv
df_quic = pd.read_csv('D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_0_1_result.csv')

# 绘制折线图
plt.figure(figsize=(25, 8))
plt.plot(df_quic.index, df_quic['HolepunchTime'], label='QUIC',marker='o',linewidth=3)
plt.plot(df_tcp.index, df_tcp['HolepunchTime'], label='TCP',marker='o',linewidth=3)


# 添加标题和标签
plt.title('Comparison of Hole Punching Time with 1% Packet Loss Rate',fontsize=30)
plt.xlabel('Hole Punching Count',fontsize=30)
plt.ylabel('Hole Punching Time (ms)',fontsize=30)
plt.legend(fontsize=20)  # 设置图例字体大小

# 设置x轴和y轴刻度的字体大小
plt.xticks(fontsize=20)
plt.yticks(fontsize=20)

# plt.ylim(0,1000)

# 显示图表
plt.show()
