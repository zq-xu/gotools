import os
import re
from collections import defaultdict

# ==== 全局配置 ====
# 你的 Go 模块路径（会自动拼到导入路径中）
MODULE_PATH = "github.com/zq-xu/gotools"

# ==== 主逻辑 ====
def process_package(pkg_name: str):
    """
    处理一个包：扫描公共函数、变量、常量、类型，并生成对应的别名导出文件。
    示例：
        process_package("utils")
        process_package("types")
    """
    pkg_path = f"./{pkg_name}"
    import_path = f"{MODULE_PATH}/{pkg_name}"
    output_file = f"./{pkg_name}.go"

    # 匹配导出函数、变量、常量（首字母大写）
    EXPORT_FUNC_VAR_CONST = re.compile(r'(?:func|var|const)\s+([A-Z]\w*)')

    # 匹配导出类型定义（struct / interface / type alias / 基础类型）
    EXPORT_TYPES = re.compile(r'type\s+([A-Z]\w*)\s*(?:=|\w|struct|interface)')

    # 按文件收集导出项
    file_exports = defaultdict(lambda: {"funcs": [], "types": []})

    # 遍历包目录下所有 .go 文件（忽略 _test.go）
    for root, dirs, files in os.walk(pkg_path):
        for file in files:
            if file.endswith(".go") and not file.endswith("_test.go"):
                path = os.path.join(root, file)
                with open(path, "r", encoding="utf-8") as f:
                    content = f.read()

                    funcs = EXPORT_FUNC_VAR_CONST.findall(content)
                    types = EXPORT_TYPES.findall(content)

                    if funcs:
                        file_exports[file]["funcs"].extend(funcs)
                    if types:
                        file_exports[file]["types"].extend(types)

    # 生成输出文件
    with open(output_file, "w", encoding="utf-8") as f:
        f.write("package gotools\n\n")
        f.write("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT MANUALLY.\n\n")
        f.write(f'import "{import_path}"\n\n')

        # var 块（函数、变量、常量）
        f.write("var (\n")
        for file, exports in file_exports.items():
            if exports["funcs"]:
                f.write(f"    // From {file}\n")
                for item in exports["funcs"]:
                    f.write(f"    {item} = {pkg_name}.{item}\n\n")
        f.write(")\n\n")

        # type 块（结构体、接口、类型别名）
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
    print(f"✅ 已生成 {output_file}（{total} 个导出项，忽略 _test.go 文件）")


# ==== 执行入口 ====
if __name__ == "__main__":
    process_package("utils")
    process_package("types")
