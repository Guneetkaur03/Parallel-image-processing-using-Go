from helper import *


for l in tqdm(cases.keys()):
    # run sequential
    cmd = ' go run editor.go {}'.format(l)
    
    

    # run parallel
    times_p[l] = []
    times_b[l] = []
    for t in tqdm(threads):
        cmd = ' go run editor.go {} pipeline {}'.format(l, t)
        times_p[l].append(seq_time[-1]/run_cmd(cmd))  
        cmd = ' go run editor.go {} bsp {}'.format(l, t)
        times_b[l].append(seq_time[-1]/run_cmd(cmd)) 


print(seq_time)
print(times_p)

plot_graph("pipeline", times_p)
plot_graph("bsp", times_b)
# print(y)
