from helper import *

my_file = Path(pip)

# if file
if my_file.is_file():
    with open(pip) as json_file:
        times_p = json.load(json_file)
else:
    with open(seq) as json_file:
        seq_time = json.load(json_file)

    for l in tqdm(cases.keys()):
        # run parallel
        times_p[l] = []
        for t in tqdm(threads):
            cmd = ' go run editor.go {} pipeline {}'.format(l, t)
            times_p[l].append(seq_time[l]/run_cmd(cmd, True))  

    with open(pip, 'w') as fp:
        json.dump(times_p, fp)


print(times_p)
plot_graph("pipeline", times_p)

