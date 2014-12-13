def parse(f):
    name, seq = None, ''
    for line in f.readlines():
        if line.startswith('>'):
            if name:
                yield name, seq
            name = line.strip('\n')[1:]
        else:
            seq += line.strip('\n')
    yield name, seq


def write(seq, name, f):
    f.write('>%s\n' % name)
    i = 0
    while i + 80 < len(seq):
        f.write('%s\n' % seq[i:i+80])
        i += 80
    if len(seq) > i:
        f.write('%s\n' % seq[i:])

