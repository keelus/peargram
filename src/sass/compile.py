import os
from tqdm import tqdm


ignore = ["compile.py", "global.scss"]

os.system("cls" if os.name == "nt" else "clear")

files = [_file.name for _file in os.scandir(".") if _file.name not in ignore]

print(f"⚒️  COMPILING {len(files)} FILES")

for _file in (progress := tqdm(files)):
    command = f'sass "{_file}" "../../web/assets/css/{_file.replace("scss", "css")}"'
    os.system(command)
    # Print the current file name after the progress bar
    progress.set_postfix(file=_file)
    
print("✅ DONE")