from helper import *

my_file = Path(bsp)

# if file
if my_file.is_file():
    with open(bsp) as json_file:
        times_b = json.load(json_file)
else:
    with open(seq) as json_file:
        seq_time = json.load(json_file)

    for l in tqdm(cases.keys()):
        # run parallel
        times_b[l] = []
        for t in tqdm(threads):
            cmd = ' go run editor.go {} bsp {}'.format(l, t)
            times_b[l].append(seq_time[l]/run_cmd(cmd, False))
     

    with open(bsp, 'w') as fp:
        json.dump(times_b, fp)


print(times_b)
plot_graph("bsp", times_b)

