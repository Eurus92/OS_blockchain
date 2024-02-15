import matplotlib

matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

time = []
time0 = []
time1 = []
all = []
all0 = []
all1 = []
block1 = []
block2 = []
for j in range(5):
    file = '../result/node1' + str(j + 1) + '.txt'
    with open(file, 'r') as f:
        i = 0
        read = f.read().splitlines()
        for d in read:
            if i == 0:
                all0.append(float(d))
            elif i == 1:
                time0.append(int(d))
            elif i == 3:
                block1.append(int(d))
            i += 1
for j in range(5):
    file = '../result/node1' + str(j + 1) + 't.txt'
    with open(file, 'r') as f:
        i = 0
        read = f.read().splitlines()
        for d in read:
            if i == 0:
                all1.append(float(d))
            elif i == 1:
                time1.append(int(d))
            elif i == 3:
                block2.append(int(d))
            i += 1
len1 = max(block1)
timeBar0 = (time0[0] + time0[1] + time0[2] + time0[3] + time0[4]) / 5/len1
time.append(timeBar0)
len2 = max(block2)
timeBar1 = (time1[0] + time1[1] + time1[2] + time1[3] + time1[4]) / 5/len2
time.append(timeBar1)
label1 = ['array', 'tree']
idx = np.arange(len(label1))
plt.figure(1)
for x, y in zip(label1, time):
    if x == 'array':
        plt.bar(x, y, color='orange')
        plt.text(x, y, '%.2f' % y, ha='center')
    elif x == 'tree':
        plt.bar(x, y, color='red')
        plt.text(x, y, '%.2f' % y, ha='center')
plt.xlabel("structure")
plt.ylabel("verification time")
plt.show()
plt.savefig("./timeV.png")
plt.close()

timeBar0 = (all0[0] + all0[1] + all0[2] + all0[3] + all0[4]) / 5
all.append(timeBar0)
timeBar1 = (all1[0] + all1[1] + all1[2] + all1[3] + all1[4]) / 5
all.append(timeBar1)
label1 = ['array', 'tree']
idx = np.arange(len(label1))
plt.figure(2)
for x, y in zip(label1, all):
    if x == 'array':
        plt.bar(x, y, color='orange')
        plt.text(x, y, '%.2f' % y, ha='center')
    elif x == 'tree':
        plt.bar(x, y, color='red')
        plt.text(x, y, '%.2f' % y, ha='center')
plt.xlabel("structure")
plt.ylabel("total time")
plt.show()
plt.savefig("./timeA.png")
plt.close()
