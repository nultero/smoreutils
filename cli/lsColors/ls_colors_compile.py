#!/usr/bin/env python3

class FType():
    def __init__(self, name):
        self.name = name
        self.fgColor = None
        self.bgColor = None
    def toStr(self):
        if self.bgColor != None:
            return f"{self.name}={self.fgColor};{self.bgColor}"
        return f"{self.name}={self.fgColor}"

def getFType(line: str) -> FType:
    name = ""
    sub = line[1:]
    end = len(sub)-1
    for i in range(end, 0, -1):
        if sub[i] != " " and sub[i] != "{":
            name = sub[:i+1]
            break
    return FType(name)

def sliceRGB(ln: str) -> str:
    s, e = 0, 0
    for idx, c in enumerate(ln):
        if c == "(":
            s = idx+1
        elif c == ")":
            e = idx; break
    sl = ln[s:e].replace(" ", "").split(",")
    return ";".join(sl)

def parseLine(f: FType, ln: str) -> FType:
    rgb = sliceRGB(ln)
    if "background-color" in ln:
        f.bgColor = "48;2;" + rgb
    elif "color" in ln:
        f.fgColor = "38;2;" + rgb
    else:
        raise Exception(f"line has no valid values: {ln}")
    return f

def main():
    types = []
    cur = None
    for ln in open("ls_colors.css").readlines():
        ln = ln.lstrip().replace("\n", "")
        if len(ln) == 0:
            continue
        firstChar = ln[0]
        if firstChar == "/":#is a comment
            continue
        elif firstChar == ".":
            cur = getFType(ln)
            # getClassName = cur
        elif firstChar == "}":
            types.append(cur)
            cur = None
        else:
            cur = parseLine(cur, ln)
    
    chunks = [f.toStr() for f in types]
    fullStr = ":".join(chunks)
    templateScript = f"export LS_COLORS='{fullStr}'"

    import os
    home = os.environ["HOME"]
    fpath = home + "/.zsh_ls_colors"

    open(fpath, "w+").write(templateScript)

    print("\x1b[32mDONE\x1b[0m")

main()