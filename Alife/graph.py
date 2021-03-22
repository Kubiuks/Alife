import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import sys

print("Running analysis")
n = 2

average_LL_from_runs = []

a1_energy_average = np.array([0] * 15000)
a2_energy_average = np.array([0] * 15000)
a3_energy_average = np.array([0] * 15000)
a4_energy_average = np.array([0] * 15000)
a5_energy_average = np.array([0] * 15000)
a6_energy_average = np.array([0] * 15000)

a1_socialness_average = np.array([0] * 15000)
a2_socialness_average = np.array([0] * 15000)
a3_socialness_average = np.array([0] * 15000)
a4_socialness_average = np.array([0] * 15000)
a5_socialness_average = np.array([0] * 15000)
a6_socialness_average = np.array([0] * 15000)

a1_oxytocin_average = np.array([0] * 15000)
a2_oxytocin_average = np.array([0] * 15000)
a3_oxytocin_average = np.array([0] * 15000)
a4_oxytocin_average = np.array([0] * 15000)
a5_oxytocin_average = np.array([0] * 15000)
a6_oxytocin_average = np.array([0] * 15000)

a1_cortisol_average = np.array([0] * 15000)
a2_cortisol_average = np.array([0] * 15000)
a3_cortisol_average = np.array([0] * 15000)
a4_cortisol_average = np.array([0] * 15000)
a5_cortisol_average = np.array([0] * 15000)
a6_cortisol_average = np.array([0] * 15000)

name = 'data/' + sys.argv[1]

for i in range(n):
    if ((i+1) % 10) == 0:
        print("Iteration:", i+1)

    file = name + '/run_' + str(i+1) + '.csv'

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

    a1_energy_average = np.add(a1_energy_average, a1[:, 1])
    a2_energy_average = np.add(a2_energy_average, a2[:, 1])
    a3_energy_average = np.add(a3_energy_average, a3[:, 1])
    a4_energy_average = np.add(a4_energy_average, a4[:, 1])
    a5_energy_average = np.add(a5_energy_average, a5[:, 1])
    a6_energy_average = np.add(a6_energy_average, a6[:, 1])

    a1_socialness_average = np.add(a1_socialness_average, a1[:, 2])
    a2_socialness_average = np.add(a2_socialness_average, a2[:, 2])
    a3_socialness_average = np.add(a3_socialness_average, a3[:, 2])
    a4_socialness_average = np.add(a4_socialness_average, a4[:, 2])
    a5_socialness_average = np.add(a5_socialness_average, a5[:, 2])
    a6_socialness_average = np.add(a6_socialness_average, a6[:, 2])

    a1_oxytocin_average = np.add(a1_oxytocin_average, a1[:, 3])
    a2_oxytocin_average = np.add(a2_oxytocin_average, a2[:, 3])
    a3_oxytocin_average = np.add(a3_oxytocin_average, a3[:, 3])
    a4_oxytocin_average = np.add(a4_oxytocin_average, a4[:, 3])
    a5_oxytocin_average = np.add(a5_oxytocin_average, a5[:, 3])
    a6_oxytocin_average = np.add(a6_oxytocin_average, a6[:, 3])

    a1_cortisol_average = np.add(a1_cortisol_average, a1[:, 4])
    a2_cortisol_average = np.add(a2_cortisol_average, a2[:, 4])
    a3_cortisol_average = np.add(a3_cortisol_average, a3[:, 4])
    a4_cortisol_average = np.add(a4_cortisol_average, a4[:, 4])
    a5_cortisol_average = np.add(a5_cortisol_average, a5[:, 4])
    a6_cortisol_average = np.add(a6_cortisol_average, a6[:, 4])

    life_lengths = []
    for agent in [a1[:, 1], a2[:, 1], a3[:, 1], a4[:, 1], a5[:, 1], a6[:, 1]]:
        for j, e in enumerate(agent):
            if e <= 0:
                life_lengths.append(j+1)
                break
            elif j+1 == 15000:
                life_lengths.append(j+1)

    average_life_length_in_percent = (sum(life_lengths) / 6) / 15000
    average_LL_from_runs.append(average_life_length_in_percent)

average_LL_from_experiment = sum(average_LL_from_runs) / n
print(average_LL_from_experiment)

a1_energy_average = np.divide(a1_energy_average, n)
a2_energy_average = np.divide(a2_energy_average, n)
a3_energy_average = np.divide(a3_energy_average, n)
a4_energy_average = np.divide(a4_energy_average, n)
a5_energy_average = np.divide(a5_energy_average, n)
a6_energy_average = np.divide(a6_energy_average, n)

a1_socialness_average = np.divide(a1_socialness_average, n)
a2_socialness_average = np.divide(a2_socialness_average, n)
a3_socialness_average = np.divide(a3_socialness_average, n)
a4_socialness_average = np.divide(a4_socialness_average, n)
a5_socialness_average = np.divide(a5_socialness_average, n)
a6_socialness_average = np.divide(a6_socialness_average, n)

a1_oxytocin_average = np.divide(a1_oxytocin_average, n)
a2_oxytocin_average = np.divide(a2_oxytocin_average, n)
a3_oxytocin_average = np.divide(a3_oxytocin_average, n)
a4_oxytocin_average = np.divide(a4_oxytocin_average, n)
a5_oxytocin_average = np.divide(a5_oxytocin_average, n)
a6_oxytocin_average = np.divide(a6_oxytocin_average, n)

a1_cotisol_average = np.divide(a1_cortisol_average, n)
a2_cotisol_average = np.divide(a2_cortisol_average, n)
a3_cotisol_average = np.divide(a3_cortisol_average, n)
a4_cotisol_average = np.divide(a4_cortisol_average, n)
a5_cotisol_average = np.divide(a5_cortisol_average, n)
a6_cotisol_average = np.divide(a6_cortisol_average, n)

plt.plot(a1_energy_average)
plt.plot(a2_energy_average)
plt.plot(a3_energy_average)
plt.plot(a4_energy_average)
plt.plot(a5_energy_average)
plt.plot(a6_energy_average)

plt.show()

plt.plot(a1_socialness_average)
plt.plot(a2_socialness_average)
plt.plot(a3_socialness_average)
plt.plot(a4_socialness_average)
plt.plot(a5_socialness_average)
plt.plot(a6_socialness_average)

plt.show()

plt.plot(a1_oxytocin_average)
plt.plot(a2_oxytocin_average)
plt.plot(a3_oxytocin_average)
plt.plot(a4_oxytocin_average)
plt.plot(a5_oxytocin_average)
plt.plot(a6_oxytocin_average)

plt.show()

plt.plot(a1_cortisol_average)
plt.plot(a2_cortisol_average)
plt.plot(a3_cortisol_average)
plt.plot(a4_cortisol_average)
plt.plot(a5_cortisol_average)
plt.plot(a6_cortisol_average)

plt.show()
