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
a = np.array(data['Agent_1'].values.tolist())
plt.plot(a[:, 1])
plt.show()

