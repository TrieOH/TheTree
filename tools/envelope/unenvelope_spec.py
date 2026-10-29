"""Rewrite an api-spec.yml off the fun envelope, as text (keeps comments and layout).

- `allOf: [Response, {properties: {data: X}}]` -> X; a list X becomes
  `{items: X}` so pagination fields can be added beside it later.
- A bare `$ref: Response` 200 (no payload) -> 204 with no content.
- `$ref: ErrorResponse` -> `$ref: Error`.
- Drops the Response, ErrorResponse and PaginationMeta schemas.

Usage: python3 tools/envelope/unenvelope_spec.py api/<svc>/api-spec.yml
"""
import re
import sys

RESPONSE = "$ref: '#/components/schemas/Response'"
path = sys.argv[1]
lines = open(path).read().split("\n")


def indent(line):
    return len(line) - len(line.lstrip(" "))


def block_end(start, base):
    """First index after start whose line is non-blank with indent <= base."""
    i = start
    while i < len(lines) and (not lines[i].strip() or indent(lines[i]) > base):
        i += 1
    return i


out = []
i = 0
wrapped = bare = 0
while i < len(lines):
    line = lines[i]
    # allOf envelope
    if line.strip() == "allOf:" and i + 1 < len(lines) and lines[i + 1].strip() == "- " + RESPONSE:
        schema_ind = indent(lines[i - 1])
        assert lines[i - 1].strip() == "schema:", (i, lines[i - 1])
        j = i + 2
        assert lines[j].strip() == "- type: object", (j, lines[j])
        assert lines[j + 1].strip() == "properties:", (j + 1, lines[j + 1])
        assert lines[j + 2].strip() == "data:", (j + 2, lines[j + 2])
        data_ind = indent(lines[j + 2])
        body_start = j + 3
        body_end = block_end(body_start, data_ind)
        # nothing else may follow data inside the allOf member
        end = block_end(body_end, indent(lines[i]))
        assert end == body_end, ("extra content after data", i, lines[body_end:end])
        body = lines[body_start:body_end]
        body_ind = indent(body[0])
        is_list = any(l.strip() == "type: array" and indent(l) == body_ind for l in body)
        target = schema_ind + 2
        if is_list:
            out += [
                " " * target + "type: object",
                " " * target + "required: [items]",
                " " * target + "properties:",
                " " * (target + 2) + "items:",
            ]
            target += 4
        out += [(" " * target + l[body_ind:]) if l.strip() else l for l in body]
        wrapped += 1
        i = body_end
        continue
    # bare envelope: a 200 with no payload becomes a 204
    if line.strip() == RESPONSE:
        assert lines[i - 1].strip() == "schema:" and lines[i - 2].strip() == "application/json:", i
        content = i - 3
        assert lines[content].strip() == "content:", (content, lines[content])
        assert block_end(i + 1, indent(lines[content])) == i + 1, ("content has more than the schema", i)
        # walk back to the status line
        k = len(out) - 1
        while not re.match(r"^\s*'\d{3}':\s*$", out[k]):
            k -= 1
        assert out[k].strip() == "'200':", out[k]
        out[k] = out[k].replace("'200'", "'204'")
        del out[len(out) - 3 :]  # content:, application/json:, schema:
        bare += 1
        i += 1
        continue
    out.append(line.replace("'#/components/schemas/ErrorResponse'", "'#/components/schemas/Error'"))
    i += 1

# drop the envelope schemas
text = "\n".join(out)
lines = text.split("\n")
for name in ("Response", "ErrorResponse", "PaginationMeta"):
    start = lines.index(f"    {name}:")
    end = block_end(start + 1, 4)
    del lines[start:end]
text = "\n".join(lines)
assert "schemas/Response'" not in text and "ErrorResponse" not in text and "PaginationMeta" not in text
open(path, "w").write(text)
print(f"{path}: {wrapped} unwrapped, {bare} -> 204")
