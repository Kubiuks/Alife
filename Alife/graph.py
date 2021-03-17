import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

file = 'data/test.csv'

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
plt.plot(a1[:, 4])
plt.plot(a2[:, 4])
plt.plot(a3[:, 4])
plt.plot(a4[:, 4])
plt.plot(a5[:, 4])
plt.plot(a6[:, 4])
plt.show()

