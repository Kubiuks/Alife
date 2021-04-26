import pandas as pd
import numpy as np
from scipy import stats
from statsmodels.multivariate.manova import MANOVA
import statsmodels.api as sm
from statsmodels.formula.api import ols
from statsmodels.stats.multicomp import pairwise_tukeyhsd
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
    print("metric to be tested and on whom and at least 2 experiments are needed to test, in a form: [metrics], metricTarget, Exp1, Exp2, ...")
    quit()
raw_metrics = sys.argv[1]
temp = raw_metrics.replace('[', '')
temp2 = temp.replace(']', '')
metrics = [s for s in temp2.split(sep=',')]

for metric in metrics:
    if metric not in ['LL', 'CT', 'OT', 'PW', 'Aggresion', 'Grooming']:
        print("invalid metric, use one of: LL, CT, OT, PW, Aggresion, Grooming")
        quit()
metric_target = sys.argv[2]
if metric_target not in ['All', 'Bonded', 'Unbonded', 'Agent_1', 'Agent_2', 'Agent_3', 'Agent_4', 'Agent_5', 'Agent_6', 'Intra-bond', 'By-phase', 'Intra-bond-by-phase']:
    print("invalid metricTarget, use one of: All, Bonded, Unbonded, Agent_X (where X is 1-6), Intra-bond, By-phase, Intra-bond-by-phase")
    quit()
if metric_target in ['Intra-bond', 'By-phase', 'Intra-bond-by-phase']:
    for metric in metrics:
        if metric not in ['Aggresion', 'Grooming']:
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
num_of_metrics = len(metrics)
if num_of_metrics > 1 and num_of_differences > 1:
    print("if more than 1 metric, use only 1 independent variable")
    quit()
if num_of_metrics > 2:
    print("use at most 2 metrics")
    quit()
if num_of_differences > 2:
    print("use at most 2 independent variables")
    quit()


data_list = []
for metric in metrics:
    data = []
    for i in range(num_of_experiments):
        file = 'data/' + experiments[i] + '/' + metric + '.csv'
        tmp_data = pd.read_csv(file, header=0)
        data.append(np.array(tmp_data[metric_target].values.tolist()))
    data_list.append(data)

color = ['g', 'b', 'r']
# for data in data_list:
#     for i, e in enumerate(data):
#         # print(stats.kstest(e, 'norm'))
#         plt.hist(e, bins=20, color=color[i])
#     plt.show()

p_value = -1
if num_of_experiments == 2:
    if num_of_metrics == 1:
        if num_of_differences == 1:
            print("Mann-Whitney Utest")
            _, p_value = stats.mannwhitneyu(data_list[0][0], data_list[0][1])
        else:
            print("usually too few experiments to properly work: 2-way ANOVA(t-test) on 2 experiments")
            res = []
            index_1 = -1
            index_2 = -1
            for i, e in enumerate(differences):
                if e:
                    if index_1 != -1:
                        index_2 = i
                    else:
                        index_1 = i
            for i in range(len(data_list[0])):
                tmp = [(split[i][index_1], split[i][index_2], x) for x in data_list[0][i]]
                res.append(tmp)
            res = [item for sublist in res for item in sublist]
            df = pd.DataFrame(res, columns=['difference_1', 'difference_2', metrics[0]])
            print(df)
            model = ols(metrics[0] + '~C(difference_1) + C(difference_2) + C(difference_1):C(difference_2)', data=df).fit()
            print(sm.stats.anova_lm(model, typ=2))
    else:
        print("MANOVA on 2 experiments")
        res = []
        for i in range(len(experiments)):
            tmp = [(x, y, experiments[i]) for (x, y) in zip(data_list[0][i], data_list[1][i])]
            res.append(tmp)
        res = [item for sublist in res for item in sublist]
        df = pd.DataFrame(res, columns=[*metrics, 'experiment'])
        metric_1, metric_2 = metrics[0], metrics[1]
        maov = MANOVA.from_formula(metric_1 + '+' + metric_2 + '~experiment', data=df)
        print(maov.mv_test())

        reg = ols(metric_1 + '~experiment', data=df).fit()
        aov = sm.stats.anova_lm(reg, type=2)
        print(aov)

        reg_2 = ols(metric_2 + '~experiment', data=df).fit()
        aov_2 = sm.stats.anova_lm(reg_2, type=2)
        print(aov_2)

        mc = pairwise_tukeyhsd(df[metric_1], df['experiment'], alpha=0.05)
        print(mc)

        mc_2 = pairwise_tukeyhsd(df[metric_2], df['experiment'], alpha=0.05)
        print(mc_2)

else:
    if num_of_metrics == 1:
        if num_of_differences == 1:
            print("Kruskall-Wallis test (non-parametric ANOVA equivalent)")
            _, p_value = stats.kruskal(*data_list[0])
        else:
            print("2-way ANOVA on n > 2 experiments")
            res = []
            index_1 = -1
            index_2 = -1
            for i, e in enumerate(differences):
                if e:
                    if index_1 != -1:
                        index_2 = i
                    else:
                        index_1 = i
            for i in range(len(data_list[0])):
                tmp = [(split[i][index_1], split[i][index_2], x) for x in data_list[0][i]]
                res.append(tmp)
            res = [item for sublist in res for item in sublist]
            df = pd.DataFrame(res, columns=['difference_1', 'difference_2', metrics[0]])
            model = ols(str(metrics[0]) + ' ~ C(difference_1) + C(difference_2) + C(difference_1):C(difference_2)', data=df).fit()
            print(sm.stats.anova_lm(model, typ=2))
    else:
        print("MANOVA on n > 2 experiments")
        res = []
        for i in range(len(experiments)):
            tmp = [(x, y, experiments[i]) for (x, y) in zip(data_list[0][i], data_list[1][i])]
            res.append(tmp)
        res = [item for sublist in res for item in sublist]
        df = pd.DataFrame(res, columns=[*metrics, 'experiment'])
        metric_1, metric_2 = metrics[0], metrics[1]
        maov = MANOVA.from_formula(metric_1 + '+' + metric_2 + '~experiment', data=df)
        print(maov.mv_test())

        reg = ols(metric_1 + '~experiment', data=df).fit()
        aov = sm.stats.anova_lm(reg, type=2)
        print(aov)

        reg_2 = ols(metric_2 + '~experiment', data=df).fit()
        aov_2 = sm.stats.anova_lm(reg_2, type=2)
        print(aov_2)

        mc = pairwise_tukeyhsd(df[metric_1], df['experiment'], alpha=0.05)
        print(mc)

        mc_2 = pairwise_tukeyhsd(df[metric_2], df['experiment'], alpha=0.05)
        print(mc_2)

if p_value == -1:
    quit()
if p_value < 0.01:
    print("Statistical significance at 0.01 level")
elif p_value < 0.05:
    print("Statistical significance at 0.05 level")
elif p_value < 0.1:
    print("Statistical significance at 0.1 level")
else:
    print("The results are not statistically significant")
