import re
import pandas as pd
from matplotlib import pyplot as plt

# 读取数据
df = pd.read_csv('D:/Users/LLL2000/Desktop/analyze_fin/data/origin_data/TCP_NO.csv', delimiter=',')

# 以StartTime的精确到分钟为一组，选择每组中StartTime开始时间最晚的和EndTime最早的数据
df['StartTime'] = pd.to_timedelta(df['StartTime'])

# 计算相邻两行数据的时间差
df['TimeDifference'] = df['StartTime'].diff()

# 根据时间差是否超过1分钟，将数据分组
group_indices = (df['TimeDifference'] > pd.Timedelta(minutes=1)).cumsum()

# 根据分组进行分组操作
df_grouped = df.groupby(group_indices)

# 创建一个空的DataFrame来存储结果
result_df = pd.DataFrame(columns=['StartTime', 'EndTime', 'HolepunchTime'])

# 遍历每个分组
for group_name, group_df in df_grouped:
    # 选择StartTime最晚的行
    start_time_row = group_df[group_df['StartTime'] == group_df['StartTime'].max()]

    # 选择EndTime为0的行
    zero_end_time_row = group_df[group_df['EndTime'] == '00:00:00']
    # 选择EndTime不为0的行
    nonzero_end_time_row = group_df[group_df['EndTime'] != '00:00:00']

    # 如果有EndTime为0的行
    if not zero_end_time_row.empty:
        # 选择同组的非零EndTime行，如果没有则选择非零EndTime最早的行
        if not nonzero_end_time_row.empty:
            end_time_row = nonzero_end_time_row[
                nonzero_end_time_row['EndTime'] == nonzero_end_time_row['EndTime'].min()]
        else:
            end_time_row = zero_end_time_row
    else:
        # 选择EndTime最早的行
        end_time_row = group_df[group_df['EndTime'] == group_df['EndTime'].min()]

    # 获取StartTime、EndTime和HolepunchTime的值
    start_time = start_time_row['StartTime'].values[0]
    end_time = end_time_row['EndTime'].values[0]
    holepunch_time = pd.Timedelta(end_time) - pd.Timedelta(start_time)

    # 格式化HolepunchTime
    # 将Timedelta格式转化为字符串
    holepunch_time = str(holepunch_time)
    # 正则表达式
    milliseconds = re.search(r"\.(\d+)", holepunch_time).group(1)
    # In ein float umwandeln
    holepunch_time = float(milliseconds) / 1000000  # In Sekunden umwandeln

    # 将结果添加到结果DataFrame中
    if not result_df.empty:
        result_df = pd.concat([result_df, pd.DataFrame(
            {'StartTime': [start_time], 'EndTime': [end_time], 'HolepunchTime': [holepunch_time]})], ignore_index=True,
                              sort=False)
    else:
        result_df = pd.DataFrame(
            {'StartTime': [start_time], 'EndTime': [end_time], 'HolepunchTime': [holepunch_time]})

print(result_df)
# 将结果存储到文件
result_df.to_csv('D:/Users/LLL2000/Desktop/analyze_fin/data/result/TCP_NO_result.csv', index=False)

# 添加一个新的列作为次数
df_final = result_df.copy()
df_final.loc[:, 'Count'] = range(1, len(df_final) + 1)

# 绘制折线图
plt.figure(figsize=(25, 8))
plt.plot(df_final['Count'], df_final['HolepunchTime'], marker='o',linewidth=2)

# 添加标题和标签
plt.title('QUIC Delay-500ms Loss_rates-1%', fontsize=30)
plt.xlabel('Count', fontsize=30)
plt.ylabel('HolepunchTime (ms)', fontsize=30)

# 设置x轴和y轴刻度的字体大小
plt.xticks(fontsize=20)
plt.yticks(fontsize=20)

# 设置纵坐标的范围为 0 到 400（给0%丢包率的数据，如果是其他丢包率，需要注释掉）
# plt.ylim(350,1600)

# 显示图表
plt.show()

