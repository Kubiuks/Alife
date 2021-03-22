import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import sys

print("Running analysis")

n = 100
name = 'data/' + sys.argv[1]
numOfAgents = int(sys.argv[3])

energy_average      = np.zeros((numOfAgents, 15000))
socialness_average  = np.zeros((numOfAgents, 15000))
oxytocin_average    = np.zeros((numOfAgents, 15000))
cortisol_average    = np.zeros((numOfAgents, 15000))

average_LL_from_runs = []

converter = {}

for v in range(numOfAgents):
    converter['Agent_'+str(v+1)] = lambda x: list(map(float, x[1:-1].split(',')))

for i in range(n):
    if ((i+1) % 10) == 0:
        print("Iteration:", i+1)

    file = name + '/run_' + str(i+1) + '.csv'

    data = pd.read_csv(file, header=0, quotechar="'", converters=converter)

    agents = []
    for z in range(numOfAgents):
        temp = np.array(data['Agent_' + str(z+1)].values.tolist())
        agents.append(temp)
    agents = np.asarray(agents)

    for j in range(numOfAgents):
        energy_average[j]       = np.add(energy_average[j], agents[j][:, 1])
        socialness_average[j]   = np.add(socialness_average[j], agents[j][:, 2])
        oxytocin_average[j]     = np.add(oxytocin_average[j], agents[j][:, 3])
        cortisol_average[j]     = np.add(cortisol_average[j], agents[j][:, 4])

    life_lengths = []
    for agent in agents:
        for k, e in enumerate(agent[:, 1]):
            if e <= 0:
                life_lengths.append(k+1)
                break
            elif k+1 == 15000:
                life_lengths.append(k+1)

    average_life_length_in_percent = (sum(life_lengths) / numOfAgents) / 15000
    average_LL_from_runs.append(average_life_length_in_percent)

# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# --------------------------------AVERAGE OUT, PLOT, WRITE TO FILE------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
average_LL_from_experiment = round((sum(average_LL_from_runs) / n), 4)
print(average_LL_from_experiment)

for j in range(numOfAgents):
    energy_average[j]       = np.divide(energy_average[j], n)
    socialness_average[j]   = np.divide(socialness_average[j], n)
    oxytocin_average[j]     = np.divide(oxytocin_average[j], n)
    cortisol_average[j]     = np.divide(cortisol_average[j], n)

# WRITE TO FILE
f = open(name+'/results.txt', "w")
f.write("Average Life Length: {0:.2%},\n".format(average_LL_from_experiment))

# PLOT
for j in range(numOfAgents):
    plt.plot(energy_average[j])
plt.title("Average Energy")

plt.savefig(name+'/average_energy.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(socialness_average[j])
plt.title("Average Socialness")

plt.savefig(name+'/average_socialness.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(oxytocin_average[j])
plt.title("Average Oxytocin")

plt.savefig(name+'/average_oxytocin.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(cortisol_average[j])
plt.title("Average Cortisol")

plt.savefig(name+'/average_cortisol.pdf', bbox_inches='tight')
plt.clf()