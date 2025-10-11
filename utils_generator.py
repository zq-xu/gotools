import os
import re
from collections import defaultdict

# utils 包路径
UTILS_PACKAGE_PATH = "./utils"  # 本地 utils 包路径
OUTPUT_FILE = "./utils.go"      # 生成的文件路径

# Go 导出函数、变量、常量正则（首字母大写）
EXPORT_REGEX = re.compile(r'(?:func|var|const)\s+([A-Z]\w*)')

# 按文件收集导出项
file_exports = defaultdict(list)

# 遍历 utils 目录下所有 .go 文件，忽略 _test.go
for root, dirs, files in os.walk(UTILS_PACKAGE_PATH):
    for file in files:
        if file.endswith(".go") and not file.endswith("_test.go"):
            path = os.path.join(root, file)
            with open(path, "r", encoding="utf-8") as f:
                content = f.read()
                matches = EXPORT_REGEX.findall(content)
                if matches:
                    file_exports[file].extend(matches)

# 生成 utils.go 文件
with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
    f.write("package gotools\n\n")  # 修改包名为 gotools
    f.write("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT MANUALLY.\n\n")
    f.write('import "zq-xu/gotools/utils"\n\n')  # utils 包路径
    f.write("var (\n")

    for file, items in file_exports.items():
        f.write(f"    // From {file}\n")
        for item in items:
            f.write(f"    {item} = utils.{item}\n\n")  # 方法间空行

    f.write(")\n")

print(f"生成 {OUTPUT_FILE} 成功，共 {sum(len(v) for v in file_exports.values())} 个导出项（忽略 _test.go 文件）。")
