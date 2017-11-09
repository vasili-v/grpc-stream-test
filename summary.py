#!/usr/bin/python

import json
import sys
import time

time_units = ["ns", "us", "ms", "s"]
rate_units = ["", "k", "M", "G"]

def summary(path):
    sends, recvs, pairs = load(path)
    print "%s Loaded %d samples" % (time.strftime("[%x %X]"), len(pairs))

    print "%s %s" % (time.strftime("[%x %X]"), "Analysing latency...")
    avg_lat, hist_lat = lat(pairs)
    window = 10000000
    if len(hist_lat) == 12:
        window = 10*hist_lat[10]

    w, w_u = adjust_unit(window, time_units)
    print "%s Window: %s%s" % (time.strftime("[%x %X]"), w, w_u)

    print "%s %s" % (time.strftime("[%x %X]"), "Analysing rate...")
    avg_rate, hist_rate = rate(recvs, window)

    show_lat(avg_lat, hist_lat, "Latency")
    show_rate(avg_rate, hist_rate, "Rate")

def lat(pairs):
    lats = []
    for pair in pairs:
        if len(pair) > 2:
            lats.append(pair[2])

    return stats(sorted(lats))

def rate(recvs, window):
    rates = []

    j = 0
    for i in range(len(recvs)):
        start = recvs[i]

        end = recvs[j]
        if end - start < window:
            while end - start < window:
                j += 1
                if j >= len(recvs):
                    j -= 1
                    break

                end = recvs[j]
        else:
            while end - start > window:
                j -= 1
                if j < i:
                    j = i
                    break

                end = recvs[j]

            j += 1
            if j >= len(recvs):
                j = len(recvs) - 1
            end = recvs[j]

        if end > start:
            rates.append(1e9*float(j - i + 1)/(end - start))

        if (i + 1) % 100000 == 0:
            print "%s %s" % (time.strftime("[%x %X]"), i+1)

        j += 1
        if j >= len(recvs):
            j = len(recvs) - 1

    return stats(sorted(rates))

def stats(x):
    avg = 0
    if len(x) > 0:
        avg = sum(x)/len(x)

    hist = []
    if len(x) > 0:
        last = len(x) - 1

        hist.append(x[0])
        for i in range(1, 10):
            j = int(round(len(x)*0.1*i))
            if j > last:
                j = last

            hist.append(x[j])

        j = int(round(len(x)*0.95))
        if j > last:
            j = last

        hist.append(x[j])
        hist.append(x[-1])

    return avg, hist

def show_lat(avg, hist, title):
    avg_x, avg_u = adjust_unit(avg, time_units)
    if len(hist) > 0:
        min_x, min_x_u = adjust_unit(hist[0], time_units)
        max_x, max_x_u = adjust_unit(hist[-1], time_units)
        if len(hist) == 12:
            print "%s %s%s (%s%s, %s%s):" % \
                  (title, avg_x, avg_u, min_x, min_x_u, max_x, max_x_u)
            for i in range(1, 10):
                x, x_u = adjust_unit(hist[i], time_units)
                print "\t%d%%: %s%s" % (i*10, x, x_u)
            print "\t95%%: %s%s" % adjust_unit(hist[10], time_units)
        else:
            print "%s %s%s (%s%s, %s%s)" % \
                  (title, avg_x, avg_u, min_x, min_x_u, max_x, max_x_u)
    else:
        print "%s %s%s" % (title, avg_x, avg_u)

def show_rate(avg, hist, title):
    avg_x, avg_u = adjust_unit(avg, rate_units)
    if len(hist) > 0:
        min_x, min_x_u = adjust_unit(hist[0], rate_units)
        max_x, max_x_u = adjust_unit(hist[-1], rate_units)
        if len(hist) == 12:
            print "%s %s%s (%s%s, %s%s):" % \
                  (title, avg_x, avg_u, min_x, min_x_u, max_x, max_x_u)

            print "\t 5%%: %s%s" % adjust_unit(hist[10], rate_units)
            for i in range(9, -1, -1):
                x, x_u = adjust_unit(hist[i], rate_units)
                print "\t%d%%: %s%s" % ((10 - i)*10, x, x_u)

        else:
            print "%s %s%s (%s%s, %s%s)" % \
                  (title, avg_x, avg_u, min_x, min_x_u, max_x, max_x_u)
    else:
        print "%s %s%s" % (title, avg_x, avg_u)

def adjust_unit(x, units):
    u = units[0]
    if x < 1000:
        return x, u

    x = x/1000.
    u = units[1]
    if x < 1000:
        return x, u

    x = x/1000.
    u = units[2]
    if x < 1000:
        return x, u

    x = x/1000.
    u = units[3]
    return x, u

def load(name):
    with open(name) as f:
        data = json.load(f)

    return data.get("sends", []), data.get("receives", []), data.get("pairs", [])

def main():
    if len(sys.argv) < 2:
        print "Specify input file"
        return 2

    summary(sys.argv[1])

    return 0

sys.exit(main())
