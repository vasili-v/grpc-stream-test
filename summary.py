#!/usr/bin/python

import json
import sys
import time
import os
import re

time_units = ["ns", "us", "ms", "s"]
rate_units = ["", "k", "M", "G"]

FILE_NAME_REGEX = re.compile("test-(\\d+)\\.json$")

def summary(path, csv):
    if os.path.isdir(path):
        names = []
        for name in os.listdir(path):
            match = FILE_NAME_REGEX.match(name)
            if match:
                names.append((int(match.group(1)), name))

        avg_lats = []
        hists_lat = []
        avg_rates = []
        hists_rate = []

        for s, name in sorted(names, key=lambda x: x[0]):
            avg_lat, hist_lat, avg_rate, hist_rate = single_summary(os.path.join(path, name))

            avg_lats.append(avg_lat)
            hists_lat.append(hist_lat)

            avg_rates.append(avg_rate)
            hists_rate.append(hist_rate)

        for i in range(len(avg_lats)):
            single_show(avg_lats[i], hists_lat[i], avg_rates[i], hists_rate[i], csv)
            print ""

    else:
        avg_lat, hist_lat, avg_rate, hist_rate = single_summary(path)
        single_show(avg_lat, hist_lat, avg_rate, hist_rate, csv)

def single_summary(name):
    sends, recvs, pairs = load(name)
    print >> sys.stderr, "%s Loaded %d samples from %s" % (time.strftime("[%x %X]"), len(pairs), os.path.basename(name))

    print >> sys.stderr, "%s %s" % (time.strftime("[%x %X]"), "Analysing latency...")
    avg_lat, hist_lat = lat(pairs)
    window = 10000000
    if len(hist_lat) == 12:
        window = 10*hist_lat[10]

    w, w_u = adjust_unit(window, time_units)
    print >> sys.stderr, "%s Window: %s%s" % (time.strftime("[%x %X]"), w, w_u)

    print >> sys.stderr, "%s %s" % (time.strftime("[%x %X]"), "Analysing rate...")
    avg_rate, hist_rate = rate(recvs, window)

    return avg_lat, hist_lat, avg_rate, hist_rate

def single_show(avg_lat, hist_lat, avg_rate, hist_rate, csv):
    if csv:
        dump(avg_lat, hist_lat, avg_rate, hist_rate)
    else:
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

        if (i + 1) % 1000000 == 0:
            print >> sys.stderr, "%s Processed %d timestamps" % (time.strftime("[%x %X]"), i+1)

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

def dump(avg_lat, hist_lat, avg_rate, hist_rate):
    print "Rate, x1000 QpS\t\tLatency, us\t"
    print "Average\t%s\tAverage\t%s" % (avg_rate/1000., avg_lat/1000.)
    if len(hist_rate) > 0 and len(hist_lat) > 0:
        print "Min\t%s\tMin\t%s" % (hist_rate[0]/1000., hist_lat[0]/1000.)
        print "Max\t%s\tMax\t%s" % (hist_rate[-1]/1000., hist_lat[-1]/1000.)
        if len(hist_rate) == 12 and len(hist_lat) == 12:
            print "0%%\t%s\t0%%\t%s" % (hist_rate[11]/1000., hist_lat[0]/1000.)
            print "5%%\t%s\t10%%\t%s" % (hist_rate[10]/1000., hist_lat[1]/1000.)
            print "10%%\t%s\t20%%\t%s" % (hist_rate[9]/1000., hist_lat[2]/1000.)
            print "20%%\t%s\t30%%\t%s" % (hist_rate[8]/1000., hist_lat[3]/1000.)
            print "30%%\t%s\t40%%\t%s" % (hist_rate[7]/1000., hist_lat[4]/1000.)
            print "40%%\t%s\t50%%\t%s" % (hist_rate[6]/1000., hist_lat[5]/1000.)
            print "50%%\t%s\t60%%\t%s" % (hist_rate[5]/1000., hist_lat[6]/1000.)
            print "60%%\t%s\t70%%\t%s" % (hist_rate[4]/1000., hist_lat[7]/1000.)
            print "70%%\t%s\t80%%\t%s" % (hist_rate[3]/1000., hist_lat[8]/1000.)
            print "80%%\t%s\t90%%\t%s" % (hist_rate[2]/1000., hist_lat[9]/1000.)
            print "90%%\t%s\t95%%\t%s" % (hist_rate[1]/1000., hist_lat[10]/1000.)
            print "100%%\t%s\t100%%\t%s" % (hist_rate[0]/1000., hist_lat[11]/1000.)
    elif len(hist_rate) > 0:
        print "Min\t%s\t\t" % (hist_rate[0]/1000.)
        print "Max\t%s\t\t" % (hist_rate[-1]/1000.)
        if len(hist_rate) == 12:
            print "0%%\t%s\t\t" % (hist_rate[11]/1000.)
            print "5%%\t%s\t\t" % (hist_rate[10]/1000.)
            print "10%%\t%s\t\t" % (hist_rate[9]/1000.)
            print "20%%\t%s\t\t" % (hist_rate[8]/1000.)
            print "30%%\t%s\t\t" % (hist_rate[7]/1000.)
            print "40%%\t%s\t\t" % (hist_rate[6]/1000.)
            print "50%%\t%s\t\t" % (hist_rate[5]/1000.)
            print "60%%\t%s\t\t" % (hist_rate[4]/1000.)
            print "70%%\t%s\t\t" % (hist_rate[3]/1000.)
            print "80%%\t%s\t\t" % (hist_rate[2]/1000.)
            print "90%%\t%s\t\t" % (hist_rate[1]/1000.)
            print "100%%\t%s\t\t" % (hist_rate[0]/1000.)
    elif len(hist_lat) > 0:
        print "\t\tMin\t%s" % hist_lat[0]/1000.
        print "\t\tMax\t%s" % hist_lat[-1]/1000.
        if len(hist_lat) == 12:
            print "\t\t0%%\t%s" % (hist_lat[0]/1000.)
            print "\t\t10%%\t%s" % (hist_lat[1]/1000.)
            print "\t\t20%%\t%s" % (hist_lat[2]/1000.)
            print "\t\t30%%\t%s" % (hist_lat[3]/1000.)
            print "\t\t40%%\t%s" % (hist_lat[4]/1000.)
            print "\t\t50%%\t%s" % (hist_lat[5]/1000.)
            print "\t\t60%%\t%s" % (hist_lat[6]/1000.)
            print "\t\t70%%\t%s" % (hist_lat[7]/1000.)
            print "\t\t80%%\t%s" % (hist_lat[8]/1000.)
            print "\t\t90%%\t%s" % (hist_lat[9]/1000.)
            print "\t\t95%%\t%s" % (hist_lat[10]/1000.)
            print "\t\t100%%\t%s" % (hist_lat[11]/1000.)

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
        print "Specify input file or directory"
        return 2

    csv = False
    if len(sys.argv) > 2 and sys.argv[1].lower()[1:] == "csv":
        csv = True

    summary(sys.argv[-1], csv)

    return 0

sys.exit(main())
