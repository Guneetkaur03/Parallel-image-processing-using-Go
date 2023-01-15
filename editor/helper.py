import os
from turtle import color
from tqdm import tqdm
import timeit
import subprocess
import matplotlib.pyplot as plt
from matplotlib.pyplot import figure
import matplotlib
import json
matplotlib.use("agg")
import random

from pathlib import Path



cases = {
    "small":  "blue",
    "mixture": "orange",
    "big":  "green"
}


threads = [2, 4, 6, 8, 12]
times_p = {}
times_b = {}
seq_time = {}

seq = "seq.txt"
pip = "pip.txt"
bsp = "bsp.txt"

def run_cmd(s, flag):
    import_module = "import os"
    testcode =  '''os.system("{}")'''.format(s)
    if flag:
        return timeit.timeit(stmt=testcode, setup=import_module, number=5)/5
    else:
        return timeit.timeit(stmt=testcode, setup=import_module, number=1)

def plot_graph(name, times):
    figure(figsize=(8, 5), dpi=80)
    for l in times.keys():
        plt.plot(threads, times[l], label=l, color=cases[l])
    #plt.ylim([0,1.95])
    plt.title("Speedup graph {}".format(name))
    plt.xlabel("threads")
    plt.ylabel("Speedup")
    plt.legend()
    #plt.show()
    plt.savefig("speedup-{}{}.png".format(name, random.randint(1,10)))


