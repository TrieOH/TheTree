"""Rewrite strict handlers off the fun envelope literals, after `just gen <svc>`.

  openapi.X{Code: 200, Data: v, Timestamp: time.Now(), Module: module}
becomes, by the shape oapi-codegen now generates for X:
  list  (type X struct { Items []T })  -> openapi.X{Items: v}
  empty (type X struct {})              -> openapi.X{}  (X200JSONResponse -> X204Response)
  named (type X T)                      -> openapi.X(v)   (&v -> v, else *v)

Anything else is printed and left alone for a hand edit. Then run
goimports and the compiler; `module` consts left unused go by hand.

Usage: python3 tools/envelope/unenvelope_handlers.py api/<svc>
"""
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
gen = (root / "internal/openapi/server.gen.go").read_text()

shapes = {}
for m in re.finditer(r"^type (\w+) (struct \{\n?(.*?)\}|[\w.\[\]*]+)\n", gen, re.M | re.S):
    name, body, fields = m.group(1), m.group(2), m.group(3)
    if body.startswith("struct"):
        fields = (fields or "").strip()
        if not fields:
            shapes[name] = "empty"
        elif re.fullmatch(r"Items \S+ `json:\"items\"`", fields):
            shapes[name] = "list"
        else:
            shapes[name] = "struct"
    else:
        shapes[name] = "named"

LIT = re.compile(
    r"openapi\.(\w+)\{\s*Code:\s*\d+,\s*(?:Data:\s*([^,\n]+?),\s*)?Timestamp:\s*time\.Now\(\),\s*Module:\s*module,?\s*\}",
    re.S,
)


def deref(v):
    return v[1:] if v.startswith("&") else "*" + v


changed = 0
for path in root.glob("internal/**/*.go"):
    if path.name.endswith("_test.go") or ".gen." in path.name:
        continue
    src = path.read_text()

    def repl(m):
        name, data = m.group(1), m.group(2)
        # a payload-less 200 became a 204 in the spec, and its type was renamed
        no_content = re.sub(r"200JSONResponse$", "204Response", name)
        if data is None and shapes.get(no_content) == "empty":
            name = no_content
        shape = shapes.get(name)
        if shape == "empty" and data is None:
            return f"openapi.{name}{{}}"
        if shape == "list" and data:
            return f"openapi.{name}{{Items: {data[1:] if data.startswith('&') else data}}}"
        if shape == "named" and data:
            return f"openapi.{name}({deref(data)})"
        print(f"left for hand edit: {path}: {name} ({shape}) data={data}")
        return m.group(0)

    new = LIT.sub(repl, src)
    if new != src:
        path.write_text(new)
        changed += 1

print(f"{changed} files rewritten")
