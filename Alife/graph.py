import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import sys


def get_phasae(number):
    if 0 < number <= 1000:
        return 1
    if 1000 < number <= 2000:
        return 2
    if 2000 < number <= 3000:
        return 3
    if 3000 < number <= 4000:
        return 4
    if 4000 < number <= 5000:
        return 5
    if 5000 < number <= 6000:
        return 6
    if 6000 < number <= 7000:
        return 7
    if 7000 < number <= 8000:
        return 8
    if 8000 < number <= 9000:
        return 9
    if 9000 < number <= 10000:
        return 10
    if 10000 < number <= 11000:
        return 11
    if 11000 < number <= 12000:
        return 12
    if 12000 < number <= 13000:
        return 13
    if 13000 < number <= 14000:
        return 14
    if 14000 < number <= 15000:
        return 15


print("Running analysis")

name = 'data/' + sys.argv[1]
n = int(sys.argv[2])
numOfAgents = int(sys.argv[4])
raw_bonds = sys.argv[5]
temp = raw_bonds.replace('[', '')
temp2 = temp.replace(']', '')
list_bonds = [int(s) for s in temp2.split(sep=',') if s.isdigit()]
bonds = {i: 'bonded' for i in list_bonds}

energy_average      = np.zeros((numOfAgents, 15000))
socialness_average  = np.zeros((numOfAgents, 15000))
oxytocin_average    = np.zeros((numOfAgents, 15000))
cortisol_average    = np.zeros((numOfAgents, 15000))

average_LL_from_runs = []
bonded_average_LL_from_runs = []
unbonded_average_LL_from_runs = []

average_PW_from_runs = []
bonded_average_PW_from_runs = []
unbonded_average_PW_from_runs = []

grooming_from_runs = []
aggression_from_runs = []
grooming_from_runs_by_bonded = []
grooming_from_runs_by_unbonded = []
aggression_from_runs_by_bonded = []
aggression_from_runs_by_unbonded = []
intra_bond_grooming_from_runs = []
intra_bond_aggression_from_runs = []
grooming_by_phase_from_runs = []
aggression_by_phase_from_runs = []
intra_bond_grooming_by_phase_from_runs = []
intra_bond_aggression_by_phase_from_runs = []

average_CT_from_runs = []
bonded_average_CT_from_runs = []
unbonded_average_CT_from_runs = []

bonded_average_OT_from_runs = []

agents_LL_from_runs = []
agents_PW_from_runs = []
agents_CT_from_runs = []
agents_OT_from_runs = []
agents_grooming_from_runs = []
agents_aggression_from_runs = []
for j in range(numOfAgents):
    agents_LL_from_runs.append([])
    agents_PW_from_runs.append([])
    agents_CT_from_runs.append([])
    agents_OT_from_runs.append([])
    agents_grooming_from_runs.append([])
    agents_aggression_from_runs.append([])

converter = {}
for j in range(numOfAgents):
    converter['Agent_'+str(j+1)] = lambda x: list(map(float, x[1:-1].split(',')))

# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# ------------------------------------------GATHER DATA FROM FILES------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
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
    bonded_LL = []
    unbonded_LL = []

    agent_PWs = []
    bonded_PWs = []
    unbonded_PWs = []

    agent_CTs = []
    bonded_CTs = []
    unbonded_CTs = []

    bonded_OTs = []

    grooms = 0
    aggressions = 0
    grooms_by_bonded = 0
    aggressions_by_bonded = 0
    grooms_by_unbonded = 0
    aggressions_by_unbonded = 0
    intra_bond_grooms = 0
    intra_bond_aggressions = 0

    grooming_by_phase = [0] * 15
    aggression_by_phase = [0] * 15
    intra_bond_grooming_by_phase = [0] * 15
    intra_bond_aggression_by_phase = [0] * 15

    for agent in agents:
        agent_id = int(agent[0, 0])
        agent_pw = 0
        agent_ct = 0
        agent_ot = 0
        for k, energy in enumerate(agent[:, 1]):
            agent_pw += 1 - abs(energy-agent[k, 2])  # energy - socialness at timestep k
            agent_ct += agent[k, 4]
            agent_ot += agent[k, 3]
            if energy <= 0:
                life_lengths.append(k+1)
                agents_LL_from_runs[agent_id-1].append((k+1)/15000)
                agent_pw = agent_pw / (k+1)
                agent_ct = agent_ct / (k+1)
                agent_ot = agent_ot / (k+1)
                if agent_id in bonds:
                    bonded_LL.append(k+1)
                else:
                    unbonded_LL.append(k+1)
                break
            elif k+1 == 15000:
                life_lengths.append(k+1)
                agents_LL_from_runs[agent_id-1].append((k+1)/15000)
                agent_pw = agent_pw / (k+1)
                agent_ct = agent_ct / (k+1)
                agent_ot = agent_ot / (k+1)
                if agent_id in bonds:
                    bonded_LL.append(k+1)
                else:
                    unbonded_LL.append(k+1)
                break
        agent_PWs.append(agent_pw)
        agents_PW_from_runs[agent_id-1].append(agent_pw)
        agent_CTs.append(agent_ct)
        agents_CT_from_runs[agent_id-1].append(agent_ct)
        agents_OT_from_runs[agent_id-1].append(agent_ot)
        if agent_id in bonds:
            bonded_PWs.append(agent_pw)
            bonded_CTs.append(agent_ct)
            bonded_OTs.append(agent_ot)
        else:
            unbonded_PWs.append(agent_pw)
            unbonded_CTs.append(agent_ct)
        agent_grooms = 0
        agent_aggressions = 0
        for v, groom in enumerate(agent[:, 5]):
            if groom != 0:
                agent_grooms += 1
                phase = get_phasae(v+1)-1
                grooming_by_phase[phase] += 1
                if agent_id in bonds and int(groom) in bonds:
                    intra_bond_grooms += 1
                    intra_bond_grooming_by_phase[phase] += 1
                    # print('Agent ' + str(agent_id) + ' groomed with ' + str(groom) + ' on' + str(v+1) + ' iteration')

        for z, aggression in enumerate(agent[:, 6]):
            if aggression != 0:
                agent_aggressions += 1
                phase = get_phasae(z+1)-1
                aggression_by_phase[phase] += 1
                if agent_id in bonds and int(aggression) in bonds:
                    intra_bond_aggressions += 1
                    intra_bond_aggression_by_phase[phase] += 1
                    # print('Agent ' + str(agent_id) + ' performed aggression on ' + str(aggression) + ' on ' + str(z+1) + ' iteration')
        if agent_id in bonds:
            grooms_by_bonded += agent_grooms
            aggressions_by_bonded += agent_aggressions
        else:
            grooms_by_unbonded += agent_grooms
            aggressions_by_unbonded += agent_aggressions
        grooms += agent_grooms
        aggressions += agent_aggressions
        agents_grooming_from_runs[agent_id-1].append(agent_grooms)
        agents_aggression_from_runs[agent_id-1].append(agent_aggressions)

    # life length
    average_life_length_in_percent = (sum(life_lengths) / numOfAgents) / 15000
    average_LL_from_runs.append(average_life_length_in_percent)
    bonded_average_life_length_in_percent = 0.0
    if len(list_bonds) != 0:
        bonded_average_life_length_in_percent = (sum(bonded_LL) / len(list_bonds)) / 15000
    bonded_average_LL_from_runs.append(bonded_average_life_length_in_percent)
    unbonded_average_life_length_in_percent = 0.0
    if numOfAgents - len(list_bonds) != 0:
        unbonded_average_life_length_in_percent = (sum(unbonded_LL) / (numOfAgents - len(list_bonds))) / 15000
    unbonded_average_LL_from_runs.append(unbonded_average_life_length_in_percent)

    # Physiological wellbeing
    average_PW_from_runs.append(sum(agent_PWs) / numOfAgents)
    temp_bonded_pw = 0.0
    if len(list_bonds) != 0:
        temp_bonded_pw = sum(bonded_PWs) / len(list_bonds)
    bonded_average_PW_from_runs.append(temp_bonded_pw)
    temp_unbonded_pw = 0.0
    if numOfAgents - len(list_bonds) != 0:
        temp_unbonded_pw = sum(unbonded_PWs) / (numOfAgents - len(list_bonds))
    unbonded_average_PW_from_runs.append(temp_unbonded_pw)

    # Hormones
    average_CT_from_runs.append(sum(agent_CTs) / numOfAgents)
    temp_bonded_ct = 0.0
    if len(list_bonds) != 0:
        temp_bonded_ct = sum(bonded_CTs) / len(list_bonds)
    bonded_average_CT_from_runs.append(temp_bonded_ct)
    temp_unbonded_ct = 0.0
    if numOfAgents - len(list_bonds) != 0:
        temp_unbonded_ct = sum(unbonded_CTs) / (numOfAgents - len(list_bonds))
    unbonded_average_CT_from_runs.append(temp_unbonded_ct)
    temp_bonded_ot = 0.0
    if len(list_bonds) != 0:
        temp_bonded_ot = sum(bonded_OTs) / len(list_bonds)
    bonded_average_OT_from_runs.append(temp_bonded_ot)

    # aggression/grooming
    grooming_from_runs.append(grooms)
    aggression_from_runs.append(aggressions)
    grooming_from_runs_by_bonded.append(grooms_by_bonded)
    aggression_from_runs_by_bonded.append(aggressions_by_bonded)
    grooming_from_runs_by_unbonded.append(grooms_by_unbonded)
    aggression_from_runs_by_unbonded.append(aggressions_by_unbonded)
    intra_bond_grooming_from_runs.append(intra_bond_grooms)
    intra_bond_aggression_from_runs.append(intra_bond_aggressions)
    grooming_by_phase_from_runs.append(grooming_by_phase)
    aggression_by_phase_from_runs.append(aggression_by_phase)
    intra_bond_grooming_by_phase_from_runs.append(intra_bond_grooming_by_phase)
    intra_bond_aggression_by_phase_from_runs.append(intra_bond_aggression_by_phase)

# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# --------------------------------AVERAGE OUT, PLOT, WRITE RESULTS------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------------------------------------

# life length
average_LL_from_experiment = round((sum(average_LL_from_runs) / n), 4)
bonded_average_LL_from_experiment = round((sum(bonded_average_LL_from_runs) / n), 4)
unbonded_average_LL_from_experiment = round((sum(unbonded_average_LL_from_runs) / n), 4)
agents_average_LL_from_experiment = []
for j in range(numOfAgents):
    agents_average_LL_from_experiment.append(round((sum(agents_LL_from_runs[j]) / n), 4))

# Physiological Wellbeing
average_PW_from_experiment = round((sum(average_PW_from_runs) / n), 2)
bonded_average_PW_from_experiment = round((sum(bonded_average_PW_from_runs) / n), 2)
unbonded_average_PW_from_experiment = round((sum(unbonded_average_PW_from_runs) / n), 2)
agents_average_PW_from_experiment = []
for j in range(numOfAgents):
    agents_average_PW_from_experiment.append(round((sum(agents_PW_from_runs[j]) / n), 2))

# Hormones
average_CT_from_experiment = round((sum(average_CT_from_runs) / n), 2)
bonded_average_CT_from_experiment = round((sum(bonded_average_CT_from_runs) / n), 2)
unbonded_average_CT_from_experiment = round((sum(unbonded_average_CT_from_runs) / n), 2)
agents_average_CT_from_experiment = []
for j in range(numOfAgents):
    agents_average_CT_from_experiment.append(round((sum(agents_CT_from_runs[j]) / n), 2))
bonded_average_OT_from_experiment = round((sum(bonded_average_OT_from_runs) / n), 2)
agents_average_OT_from_experiment = []
for j in range(numOfAgents):
    agents_average_OT_from_experiment.append(round((sum(agents_OT_from_runs[j]) / n), 2))

# aggression/grooming
average_grooming_from_experiment = round(sum(grooming_from_runs) / n)
average_aggression_from_experiment = round(sum(aggression_from_runs) / n)
average_grooming_from_experiment_bonded = round(sum(grooming_from_runs_by_bonded) / n)
average_aggression_from_experiment_bonded = round(sum(aggression_from_runs_by_bonded) / n)
average_grooming_from_experiment_unbonded = round(sum(grooming_from_runs_by_unbonded) / n)
average_aggression_from_experiment_unbonded = round(sum(aggression_from_runs_by_unbonded) / n)
average_intra_bond_grooming_from_experiment = round(sum(intra_bond_grooming_from_runs) / n)
average_intra_bond_aggression_from_experiment = round(sum(intra_bond_aggression_from_runs) / n)
agents_average_grooming_from_experiment = []
agents_average_aggression_from_experiment = []
for j in range(numOfAgents):
    agents_average_grooming_from_experiment.append(round(sum(agents_grooming_from_runs[j]) / n))
    agents_average_aggression_from_experiment.append(round(sum(agents_aggression_from_runs[j]) / n))
average_grooming_by_phase = np.divide(np.sum(grooming_by_phase_from_runs, axis=0), n)
average_aggression_by_phase = np.divide(np.sum(aggression_by_phase_from_runs, axis=0), n)
average_intra_bond_grooming_by_phase = np.divide(np.sum(intra_bond_grooming_by_phase_from_runs, axis=0), n)
average_intra_bond_aggression_by_phase = np.divide(np.sum(intra_bond_aggression_by_phase_from_runs, axis=0), n)

# WRITE TO FILE

# AVERAGES FROM RUNS, for statistics
# LL
f = open(name+'/LL.csv', "w")
f.write("All,Bonded,Unbonded,")
for j in range(numOfAgents):
    f.write("Agent_" + str(j+1))
    if j != numOfAgents-1:
        f.write(",")
for j in range(n):
    f.write("\n{0:.2},{1:.2},{2:.2},".format(average_LL_from_runs[j], bonded_average_LL_from_runs[j], unbonded_average_LL_from_runs[j]))
    for k in range(numOfAgents):
        f.write("{0:.2}".format(agents_LL_from_runs[k][j]))
        if k != numOfAgents-1:
            f.write(",")

f.close()

# PW
f = open(name+'/PW.csv', "w")
f.write("All,Bonded,Unbonded,")
for j in range(numOfAgents):
    f.write("Agent_" + str(j+1))
    if j != numOfAgents-1:
        f.write(",")
for j in range(n):
    f.write("\n{0:.2},{1:.2},{2:.2},".format(average_PW_from_runs[j], bonded_average_PW_from_runs[j], unbonded_average_PW_from_runs[j]))
    for k in range(numOfAgents):
        f.write("{0:.2}".format(agents_PW_from_runs[k][j]))
        if k != numOfAgents-1:
            f.write(",")

f.close()

# CT
f = open(name+'/CT.csv', "w")
f.write("All,Bonded,Unbonded,")
for j in range(numOfAgents):
    f.write("Agent_" + str(j+1))
    if j != numOfAgents-1:
        f.write(",")
for j in range(n):
    f.write("\n{0:.2},{1:.2},{2:.2},".format(average_CT_from_runs[j], bonded_average_CT_from_runs[j], unbonded_average_CT_from_runs[j]))
    for k in range(numOfAgents):
        f.write("{0:.2}".format(agents_CT_from_runs[k][j]))
        if k != numOfAgents-1:
            f.write(",")

f.close()

# OT
f = open(name+'/OT.csv', "w")
f.write("Bonded")
if len(list_bonds) != 0:
    f.write(",")
for j, b in enumerate(list_bonds):
    f.write("Agent_" + str(b))
    if j != len(list_bonds)-1:
        f.write(",")
for j in range(n):
    f.write("\n{0:.2}".format(bonded_average_OT_from_runs[j]))
    if len(list_bonds) != 0:
        f.write(",")
    t = 0
    for k in range(numOfAgents):
        if k+1 in bonds:
            t += 1
            f.write("{0:.2}".format(agents_OT_from_runs[k][j]))
            if t != len(list_bonds):
                f.write(",")

f.close()

# GROOMING
f = open(name+'/Grooming.csv', "w")
f.write("All,Bonded,Unbonded,Intra-bond,")
for j in range(numOfAgents):
    f.write("Agent_" + str(j+1) + ",")
f.write("By_phase,Intra-bond_by_phase")
for j in range(n):
    f.write("\n{0},{1},{2},{3},".format(grooming_from_runs[j], grooming_from_runs_by_bonded[j], grooming_from_runs_by_unbonded[j], intra_bond_grooming_from_runs[j]))
    for k in range(numOfAgents):
        f.write("{0},".format(agents_grooming_from_runs[k][j]))
    f.write("{0},{1}".format(grooming_by_phase_from_runs[j], intra_bond_grooming_by_phase_from_runs[j]))

f.close()

# AGGRESION
f = open(name+'/Aggresion.csv', "w")
f.write("All,Bonded,Unbonded,Intra-bond,")
for j in range(numOfAgents):
    f.write("Agent_" + str(j+1) + ",")
f.write("By_phase,Intra-bond_by_phase")
for j in range(n):
    f.write("\n{0},{1},{2},{3},".format(aggression_from_runs[j], aggression_from_runs_by_bonded[j], aggression_from_runs_by_unbonded[j], intra_bond_aggression_from_runs[j]))
    for k in range(numOfAgents):
        f.write("{0},".format(agents_aggression_from_runs[k][j]))
    f.write("{0},{1}".format(aggression_by_phase_from_runs[j], intra_bond_aggression_by_phase_from_runs[j]))

f.close()

# RESULTS Experiment
f = open(name+'/results.txt', "w")
f.write("Life Length:\n")
f.write("All Agents Average Life Length: {0:.2%},\n".format(average_LL_from_experiment))
f.write("Bonded Agents Average Life Length: {0:.2%},\n".format(bonded_average_LL_from_experiment))
f.write("Unbonded Agents Average Life Length: {0:.2%},\n".format(unbonded_average_LL_from_experiment))
for j in range(numOfAgents):
    f.write("Agent " + str(j+1) + " Average Life Length: {0:.2%},\n".format(agents_average_LL_from_experiment[j]))

f.write("\nPhysiological Wellbeing:\n")
f.write("All Agents Average PW: {0:.2},\n".format(average_PW_from_experiment))
f.write("Bonded Agents Average PW: {0:.2},\n".format(bonded_average_PW_from_experiment))
f.write("Unbonded Agents Average PW: {0:.2},\n".format(unbonded_average_PW_from_experiment))
for j in range(numOfAgents):
    f.write("Agent " + str(j+1) + " Average PW: {0:.2},\n".format(agents_average_PW_from_experiment[j]))

f.write("\nCortisol:\n")
f.write("All Agents Average CT: {0:.2},\n".format(average_CT_from_experiment))
f.write("Bonded Agents Average CT: {0:.2},\n".format(bonded_average_CT_from_experiment))
f.write("Unbonded Agents Average CT: {0:.2},\n".format(unbonded_average_CT_from_experiment))
for j in range(numOfAgents):
    f.write("Agent " + str(j+1) + " Average CT: {0:.2},\n".format(agents_average_CT_from_experiment[j]))
f.write("\nOxytocin:\n")
f.write("Bonded Agents Average OT: {0:.2},\n".format(bonded_average_OT_from_experiment))
for j in range(numOfAgents):
    if (j+1) in bonds:
        f.write("Agent " + str(j+1) + " Average OT: {0:.2},\n".format(agents_average_OT_from_experiment[j]))

f.write("\nSocial Interactions:\n")
f.write("All Grooms: {0},\n".format(average_grooming_from_experiment))
f.write("All Aggressions: {0},\n".format(average_aggression_from_experiment))
f.write("Grooms by Bonded: {0},\n".format(average_grooming_from_experiment_bonded))
f.write("Aggressions by Bonded: {0},\n".format(average_aggression_from_experiment_bonded))
f.write("Grooms by Unbonded: {0},\n".format(average_grooming_from_experiment_unbonded))
f.write("Aggressions by Unbonded: {0},\n".format(average_aggression_from_experiment_unbonded))
f.write("intra-bond Grooms: {0},\n".format(average_intra_bond_grooming_from_experiment))
f.write("intra-bond Aggressions: {0},\n".format(average_intra_bond_aggression_from_experiment))
for j in range(numOfAgents):
    f.write("Agent {0} grooms: {1},\n".format(j+1, agents_average_grooming_from_experiment[j]))
for j in range(numOfAgents):
    f.write("Agent {0} aggressions: {1},\n".format(j+1, agents_average_aggression_from_experiment[j]))
f.write("Average grooms by phase: {0},\n".format(average_grooming_by_phase))
f.write("Average aggressions by phase: {0},\n".format(average_aggression_by_phase))
f.write("Average intra-bond grooms by phase: {0},\n".format(average_intra_bond_grooming_by_phase))
f.write("Average intra-bond aggressions by phase: {0},\n".format(average_intra_bond_aggression_by_phase))

f.close()

# PLOT
colours = ['red', 'blue', 'green', 'violet', 'orange', 'cyan']

for j in range(numOfAgents):
    energy_average[j]       = np.divide(energy_average[j], n)
    socialness_average[j]   = np.divide(socialness_average[j], n)
    oxytocin_average[j]     = np.divide(oxytocin_average[j], n)
    cortisol_average[j]     = np.divide(cortisol_average[j], n)

for j in range(numOfAgents):
    plt.plot(energy_average[j], c=colours[j % 6], label='Agent '+str(j+1))
plt.title("Average Energy")
plt.legend()

plt.savefig(name+'/average_energy.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(socialness_average[j], c=colours[j % 6], label='Agent '+str(j+1))
plt.title("Average Socialness")
plt.legend()

plt.savefig(name+'/average_socialness.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(oxytocin_average[j], c=colours[j % 6], label='Agent '+str(j+1))
plt.title("Average Oxytocin")
plt.legend()

plt.savefig(name+'/average_oxytocin.pdf', bbox_inches='tight')
plt.clf()

for j in range(numOfAgents):
    plt.plot(cortisol_average[j], c=colours[j % 6], label='Agent '+str(j+1))
plt.title("Average Cortisol")
plt.legend()

plt.savefig(name+'/average_cortisol.pdf', bbox_inches='tight')
plt.clf()
