import os
from tqdm import tqdm


exclusion = ["compilar.py", "base.scss", "variables.scss"]

os.system("cls" if os.name == "nt" else "clear")

archivos = [archivo.name for archivo in os.scandir(".") if archivo.name not in exclusion]

print(f"# COMPILANDO {len(archivos)} ARCHIVOS ")

for archivo in (progreso := tqdm(archivos)):
    comando = f'sass "{archivo}" "../../web/assets/css/{archivo.replace("scss", "css")}"'
    os.system(comando)
    # Print the current file name after the progress bar
    progreso.set_postfix(file=archivo)
    
print("# COMPLETADO")