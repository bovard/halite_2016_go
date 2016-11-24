import argparse
import os
import random
from subprocess import check_output

bots = filter(lambda s: s.startswith('Sub') and s.endswith('.go'), os.listdir('.'))

parser = argparse.ArgumentParser(description="Run multiple games")
parser.add_argument('--p1',
                    dest='p1',
                    default='MyBot.go',
                    help='player 1 bot')
parser.add_argument('--p2',
                    dest='p2',
                    default='Sub16.go',
                    choices=bots,
                    help='player 2 bot')
parser.add_argument('--games',
                    dest='games',
                    default=10,
                    type=int,
                    help='number of games to play')
parser.add_argument('--all',
                    dest='all',
                    default=False,
                    action='store_true',
                    help='plays this bot against all local Sub* bots')
parser.add_argument('--local',
                    dest='local',
                    action='store_true',
                    required=False,
                    default=False,
                    help='if you want to run locally')
options, _ = parser.parse_known_args()

def print_score(games):
	p1 = 0
	p2 = 0
	p1name = games[0].split('\n')[-3].split(',')[1]
	p2name = games[0].split('\n')[-2].split(',')[1]
	for g in games:
		g = g.split('\n')
		result = g[-3:-1]
		if "rank #1" in result[0]:
			p1 += 1
		else:
			p2 += 1
	print "{} vs {}".format(p1name, p2name)
	print "  {}:\t{}".format(p1name, p1)
	print "  {}:\t{}".format(p2name, p2)

os.environ['GOPATH'] = os.getcwd()
for o in bots if options.all else [options.p2]:
	games = []
	print options.p1, o
	for i in range(options.games):
		print i
		size = random.randint(20, 50)
		out = check_output([
			'./halite', 
			'-d', 
			'{} {}'.format(size, size), 
			'go run {}'.format(options.p1), 
			'go run {}'.format(o)
		])
		games.append(out)
	print_score(games)
	print ''
