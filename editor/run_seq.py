from helper import *

my_file = Path(seq)

# if file
if my_file.is_file():
    with open(seq) as json_file:
        seq_time = json.load(json_file)
else:
    for l in tqdm(cases.keys()):
        # run sequential
        cmd = ' go run editor.go {}'.format(l)
        
        seq_time[l] = run_cmd(cmd, False)

    with open(seq, 'w') as fp:
        json.dump(seq_time, fp)


print(seq_time)

