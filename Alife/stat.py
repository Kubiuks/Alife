import pandas as pd
import numpy as np
from scipy import stats
from os import path
from matplotlib import pyplot as plt
import sys


def differences_between_experiments(experiments_split):
    it = iter(experiments_split)
    the_len = len(next(it))
    if not all(len(l) == the_len for l in it):
        raise ValueError('incorrect experiment name, please use: WorldCondition_Bonds_DSImode!')
    values = zip(*experiments_split)
    result_list = list(values)
    tmp_differences = []
    for e in result_list:
        l = list(e)
        tmp_differences.append(not all(x == l[0] for x in l))
    return tmp_differences


print("Running Statistical Significance Testing")

if len(sys.argv) < 5:
    print("metric to be tested and on whom and at least 2 experiments are needed to test, in a form: metric, metricTarget, Exp1, Exp2, ...")
    quit()
metric = sys.argv[1]
if metric not in ['LL', 'CT', 'OT', 'PW', 'Aggresion', 'Grooming']:
    print("invalid metric, use one of: LL, CT, OT, PW, Aggresion, Grooming")
    quit()
metric_target = sys.argv[2]
if metric_target not in ['All', 'Bonded', 'Unbonded', 'Agent_1', 'Agent_2', 'Agent_3', 'Agent_4', 'Agent_5', 'Agent_6', 'Intra-bond', 'By-phase', 'Intra-bond-by-phase']:
    print("invalid metricTarget, use one of: All, Bonded, Unbonded, Agent_X (where X is 1-6), Intra-bond, By-phase, Intra-bond-by-phase")
    quit()
if metric_target in ['Intra-bond', 'By-phase', 'Intra-bond-by-phase'] and metric not in ['Aggresion', 'Grooming']:
    print("metric \"" + metric + "\" cannot be used with \"" + metric_target + "\" target")
    quit()

experiments = []
split = []
for i in range(3, len(sys.argv)):
    if sys.argv[i] in experiments:
        print("duplicate experiment")
        quit()
    if not path.exists('data/' + sys.argv[i]):
        print("path to experiment does not exist")
        quit()
    experiments.append(sys.argv[i])
    split.append(sys.argv[i].split(sep='_'))

# differences is a list in which True means there is a difference, i.e. the experiments used different values
# for that parameter, differences [WorldCondition, Bonds, DSImode]
differences = differences_between_experiments(split)
num_of_experiments = len(experiments)
num_of_differences = sum(differences)
print(num_of_experiments)
print(num_of_differences)

data = []
for i in range(num_of_experiments):
    file = 'data/' + experiments[i] + '/' + metric + '.csv'
    tmp_data = pd.read_csv(file, header=0)
    data.append(np.array(tmp_data[metric_target].values.tolist()))

color = ['g', 'b', 'r']
for i, e in enumerate(data):
    # print(stats.kstest(e, 'norm'))
    plt.hist(e, bins=70, color=color[i])
plt.show()

if num_of_experiments == 2:
    if num_of_differences == 1:
        print("Mann-Whitney Utest")
        _, p_value = stats.mannwhitneyu(data[0], data[1])
        if p_value < 0.01:
            print("Statistical significance at 0.01 level")
        elif p_value < 0.05:
            print("Statistical significance at 0.05 level")
        elif p_value < 0.1:
            print("Statistical significance at 0.1 level")
        else:
            print("The results are not statistically significant")
