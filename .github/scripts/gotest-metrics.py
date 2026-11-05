"""gotest-metrics.py EVENTS MODE

Turns a gotestsum --jsonfile (go test -json events) into metrics through
.github/scripts/metric: how many tests passed, failed and skipped, and the
ten slowest tests with their package.
"""

import json
import os
import subprocess
import sys

events_path, mode = sys.argv[1], sys.argv[2]
metric = os.path.join(os.environ["GITHUB_WORKSPACE"], ".github/scripts/metric")

counts = {"pass": 0, "fail": 0, "skip": 0}
elapsed = []
with open(events_path) as events:
    for line in events:
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        # package-level events have no Test; subtests count like tests
        if not event.get("Test") or event.get("Action") not in counts:
            continue
        counts[event["Action"]] += 1
        if event["Action"] == "pass" and "/" not in event["Test"]:
            elapsed.append((event.get("Elapsed", 0), event["Package"], event["Test"]))


def emit(name, value, *labels):
    subprocess.run([metric, "go/test", name, str(value), f"mode={mode}", *labels], check=False)


for action, count in counts.items():
    emit(f"tests_{action}", count)
for seconds, package, test in sorted(elapsed, reverse=True)[:10]:
    emit("slow_test_seconds", seconds, f"test={test}", f"package={package.rsplit('/', 1)[-1]}")
