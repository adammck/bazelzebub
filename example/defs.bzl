load("@rules_cc//cc:defs.bzl", "cc_binary")
load("@bazel_skylib//rules:write_file.bzl", "write_file")

# This is the expanded very verbose version.
# TODO: Generate this recursively from the manifest.

# buildifier: disable=function-docstring
def generate(manifest):
    targets = []

    for env in manifest.get("environments", []):
        #env_pfx = env["name"]
        env_toks = [("env", "env", env["name"])]
        env_name = workflow_name(env_toks)
        env_tasks = []

        for dc in env.get("datacenters", []):
            dc_toks = env_toks + [("datacenter", "dc", dc["name_short"])]
            dc_name = workflow_name(dc_toks)
            env_tasks.append(dc_name)
            dc_tasks = []

            for clu in dc.get("clusters", []):
                clu_toks = dc_toks + [("cluster", "clu", clu["name"])]
                clu_name = workflow_name(clu_toks)
                dc_tasks.append(clu_name)
                clu_tasks = []

                for shrd in clu.get("shards", []):
                    shrd_toks = clu_toks + [("shard", "sh", shrd["name"])]
                    shrd_name = installation_name(shrd_toks)
                    clu_tasks.append(shrd_name)
                    installation(shrd_name)
                    targets.append((shrd_name, shrd_toks))

                workflow(clu_name, clu_tasks)
                targets.append((clu_name, clu_toks))

            workflow(dc_name, dc_tasks)
            targets.append((dc_name, dc_toks))

        workflow(env_name, env_tasks)
        targets.append((env_name, env_toks))

    # Generate a rule to build a json file containing all of the targets (both
    # workflows and installations) that we have generated so far, and
    _generate_meta("meta", targets)

def installation_name(tokens):
    return "-".join([t[-1] for t in tokens])

def workflow_name(tokens):
    name = "-".join([t[-1] for t in tokens])
    return "wf-{}--{}".format(tokens[-1][1], name)

def workflow(name, tasks):
    cc_binary(
        name = name,
        srcs = ["hello-world.cc"],
        deps = tasks,
    )

def installation(n):
    cc_binary(
        name = n,
        srcs = ["hello-world.cc"],
    )

def _generate_meta(name, targets, **kwargs):
    obj = []
    for target, tokens in targets:
        obj.append({
            "target": target,
            "tokens": {tok[0]: tok[2] for tok in tokens},
        })

    write_file(
        name,
        out = "{}.json".format(name),
        content = [json.encode(obj)],
        **kwargs
    )
