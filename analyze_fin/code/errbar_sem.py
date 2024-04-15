import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

# 文件路径列表
file_paths = [
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_20_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_20_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_20_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_20_2_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_100_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_100_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_100_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_100_2_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_200_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_200_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_200_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_200_2_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_20_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_20_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_20_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_20_2_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_100_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_100_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_100_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_100_2_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_200_0_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_200_1_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_200_1-5_result.csv',
    'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_200_2_result.csv'
]

fig, ax = plt.subplots(figsize=(35, 13))

labels = ['20_0', '20_1', '20_1-5', '20_2', '100_0', '100_1', '100_1-5','100_2', '200_0', '200_1', '200_1-5', '200_2']

bar_width = 3  # 柱形图的宽度
group_spacing = 1  # 组之间的间距

for i, label in enumerate(labels):
    tcp_file = f'D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_{label}_result.csv'
    quic_file = f'D:/Users/LLL2000/Desktop/analyze_fin/data/result/QUIC_{label}_result.csv'

    df_tcp = pd.read_csv(tcp_file)
    df_quic = pd.read_csv(quic_file)

    # 根据条件剔除数据(超过3000ms)
    if label.startswith('20'):
        df_tcp = df_tcp[df_tcp['HolepunchTime'] <= 3000].dropna()
        df_quic = df_quic[df_quic['HolepunchTime'] <= 3000].dropna()
    elif label.startswith('100'):
        df_tcp = df_tcp[df_tcp['HolepunchTime'] <= 3000].dropna()
        df_quic = df_quic[df_quic['HolepunchTime'] <= 3000].dropna()
    elif label.startswith('200'):
        df_tcp = df_tcp[df_tcp['HolepunchTime'] <= 3000].dropna()
        df_quic = df_quic[df_quic['HolepunchTime'] <= 3000].dropna()

    # 计算TCP和QUIC的均值和标准误差
    quic_mean = df_quic['HolepunchTime'].mean()
    quic_sem = df_quic['HolepunchTime'].sem()
    tcp_mean = df_tcp['HolepunchTime'].mean()
    tcp_sem = df_tcp['HolepunchTime'].sem()

    # 绘制当前组的TCP和QUIC柱状图和误差线
    group_positions = i * (2 * bar_width + group_spacing)  # 计算当前组的位置
    quic_pos = group_positions   # QUIC的柱形图位置
    quic_rects = ax.bar(quic_pos, quic_mean, align='center', alpha=0.6, color='blue', width=bar_width)
    tcp_pos = group_positions + bar_width # TCP的柱形图位置
    tcp_rects = ax.bar(tcp_pos, tcp_mean, align='center', alpha=0.6, color='orange', width=bar_width)
    ax.errorbar(quic_pos, quic_mean, yerr=quic_sem, fmt='none', ecolor='black', capsize=15, elinewidth=3, capthick=3)
    ax.errorbar(tcp_pos, tcp_mean, yerr=tcp_sem, fmt='none', ecolor='black', capsize=15, elinewidth=3,capthick=3)


    # 添加数据标签（显示在柱子上方中央）
    def autolabel(rects, height_factor):
        for rect in rects:
            height = rect.get_height()
            ax.text(rect.get_x() + rect.get_width() / 2., height_factor * height,
                    '{:.3f}'.format(height),
                    ha='center', va='bottom', fontsize=18)


    # 根据标签选择不同的高度因子
    if label.startswith(('100')):
        height_factor = 1.22
    elif label.startswith(('200')):
        height_factor = 1.12
    else:
        height_factor = 1.45

    # 添加数据标签
    autolabel(quic_rects, height_factor)
    autolabel(tcp_rects, height_factor)

labels = ['20ms RTT\n0% loss', '20ms RTT\n1% loss','20ms RTT\n1.5% loss', '20ms RTT\n2% loss', '100ms RTT\n0% loss', '100ms RTT\n1% loss', '100ms RTT\n1.5% loss','100ms RTT\n2% loss', '200ms RTT\n0% loss', '200ms RTT\n1% loss','200ms RTT\n1.5% loss', '200ms RTT\n2% loss']
# 设置x轴刻度和标签
group_positions = np.arange(len(labels)) * (2 * bar_width + group_spacing)  # 计算每个组的位置
ax.set_xticks(group_positions + bar_width / 2)
ax.set_xticklabels(labels, fontsize=20)

# 添加图例
ax.legend(['QUIC', 'TCP'], loc='upper left',fontsize=30)

# 添加标题和标签
ax.set_title('Comparison of QUIC and TCP Hole Punching Time Under Different Network Conditions', fontsize=30)
# ax.set_xlabel('RTT and LOSS', fontsize=30)
ax.set_ylabel('Hole Punching Time (ms)', fontsize=30)

# 设置刻度标签字体大小
ax.tick_params(axis='y', labelsize=30)

plt.ylim(0,750)

plt.show()

#保存留片为PNG格式
fig.savefig("D:/Users/LLL2000/Desktop/lunwen/picture/fig8.png",dpi=600)
# 设置dpi参数可以调整图片的清晰度
