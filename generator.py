import os
import re
from collections import defaultdict

MODULE_PATH = "github.com/zq-xu/gotools"
PACKAGES = ["utils", "router", "types"]

def process_package(pkg_name: str):
    pkg_path = f"./{pkg_name}"
    import_path = f"{MODULE_PATH}/{pkg_name}"
    output_file = f"./{pkg_name}.go"

    # 匹配导出函数、变量、常量（单行）
    EXPORT_FUNC_VAR_CONST = re.compile(r'(?:func|var|const)\s+([A-Z]\w*)')

    # 匹配 var ( ... ) 块中的变量（多行定义）
    VAR_BLOCK_ITEM = re.compile(r'^\s*([A-Z]\w*)\s*=', re.MULTILINE)

    # 匹配导出类型定义（struct / interface / type alias / 基础类型）
    EXPORT_TYPES = re.compile(r'type\s+([A-Z]\w*)\s*(?:=|\w|struct|interface)')

    file_exports = defaultdict(lambda: {"funcs": [], "types": []})

    # 只扫描包根目录，不递归
    for file in os.listdir(pkg_path):
        if file.endswith(".go") and not file.endswith("_test.go"):
            path = os.path.join(pkg_path, file)
            if os.path.isfile(path):
                with open(path, "r", encoding="utf-8") as f:
                    content = f.read()

                    funcs = EXPORT_FUNC_VAR_CONST.findall(content)
                    block_vars = VAR_BLOCK_ITEM.findall(content)
                    types = EXPORT_TYPES.findall(content)

                    if block_vars:
                        funcs.extend(block_vars)
                    if funcs:
                        file_exports[file]["funcs"].extend(funcs)
                    if types:
                        file_exports[file]["types"].extend(types)

    # 生成导出文件
    with open(output_file, "w", encoding="utf-8") as f:
        f.write("package gotools\n\n")
        f.write("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT MANUALLY.\n\n")
        f.write(f'import "{import_path}"\n\n')

        f.write("var (\n")
        for file, exports in file_exports.items():
            if exports["funcs"]:
                f.write(f"    // From {file}\n")
                for item in exports["funcs"]:
                    f.write(f"    {item} = {pkg_name}.{item}\n\n")
        f.write(")\n\n")

        has_types = any(len(exports["types"]) > 0 for exports in file_exports.values())
        if has_types:
            f.write("type (\n")
            for file, exports in file_exports.items():
                if exports["types"]:
                    f.write(f"    // From {file}\n")
                    for t in exports["types"]:
                        f.write(f"    {t} = {pkg_name}.{t}\n\n")
            f.write(")\n")

    total = sum(len(v["funcs"]) + len(v["types"]) for v in file_exports.values())
    print(f"✅ 已生成 {output_file}（{total} 个导出项，支持 var(...) 块，忽略子目录与 _test.go 文件）")

if __name__ == "__main__":
    for pkg in PACKAGES:
        process_package(pkg)
