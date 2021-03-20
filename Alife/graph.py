import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

print("Running analysis")
n = 20
average_LL_from_runs = []

for i in range(n):
    if ((i+1) % 10) == 0:
        print("Iteration:", i+1)

    file = 'data/m/test_' + str(i+1) + '.csv'

    data = pd.read_csv(file, header=0, quotechar="'", converters={'Agent_1': lambda x: list(map(float, x[1:-1].split(','))),
                                                                  'Agent_2': lambda x: list(map(float, x[1:-1].split(','))),
                                                                  'Agent_3': lambda x: list(map(float, x[1:-1].split(','))),
                                                                  'Agent_4': lambda x: list(map(float, x[1:-1].split(','))),
                                                                  'Agent_5': lambda x: list(map(float, x[1:-1].split(','))),
                                                                  'Agent_6': lambda x: list(map(float, x[1:-1].split(',')))})
    a1 = np.array(data['Agent_1'].values.tolist())
    a2 = np.array(data['Agent_2'].values.tolist())
    a3 = np.array(data['Agent_3'].values.tolist())
    a4 = np.array(data['Agent_4'].values.tolist())
    a5 = np.array(data['Agent_5'].values.tolist())
    a6 = np.array(data['Agent_6'].values.tolist())
    # plt.plot(a1[:, 1])
    # plt.plot(a2[:, 1])
    # plt.plot(a3[:, 1])
    # plt.plot(a4[:, 1])
    # plt.plot(a5[:, 1])
    # plt.plot(a6[:, 1])

    life_lengths = []
    for agent in [a1[:, 1], a2[:, 1], a3[:, 1], a4[:, 1], a5[:, 1], a6[:, 1]]:
        for i, e in enumerate(agent):
            if e <= 0:
                life_lengths.append(i+1)
                break
            elif i+1 == 15000:
                life_lengths.append(i+1)

    average_life_length_in_percent = (sum(life_lengths) / 6) / 15000
    average_LL_from_runs.append(average_life_length_in_percent)

    # plt.show()

average_LL_from_experiment = sum(average_LL_from_runs) / n
print(average_LL_from_experiment)