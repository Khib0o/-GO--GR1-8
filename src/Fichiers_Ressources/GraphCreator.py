from os import system
import random
import sys

if (len(sys.argv)!=4):
    print("Utilisation : python3 GraphCreator.py <Nombre de sommet> <Nombre de lien additionel>")
    sys.exit()

ALPHABET = ("a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z")

file = open(sys.argv[3], 'a')
name = random.randint(0, 10000)

file.write(str(name))

nombre_sommet = int(sys.argv[1])

nom_sommet = []
while (len(nom_sommet) < nombre_sommet):
    l1 = random.randint(0, 25)
    l2 = random.randint(0, 25)
    l3 = random.randint(0, 25)
    l4 = random.randint(0, 25)
    
    word = ALPHABET[l1]+ALPHABET[l2]+ALPHABET[l3]+ALPHABET[l4]
    if(not nom_sommet.__contains__(word)):
        nom_sommet.append(word)

valeurs = []

for i in range(nombre_sommet):
    
    valeursInter = []
    
    for j in range(nombre_sommet):
        valeursInter.append(0)
    
    mod1 = -1 + i
    mod2 = 1 + i
    
    if (mod1 < 0):
        mod1 += nombre_sommet
    
    if (mod2 > nombre_sommet-1):
        mod2 -= nombre_sommet
    
    cost = random.randint(1,20)
    
    valeursInter[mod1]=cost
    valeursInter[mod2]=cost
    
    valeurs.append(valeursInter)
    
for i in range(int(sys.argv[2])):
    
    cost = random.randint(1,20)
    
    s1 = random.randint(0,nombre_sommet-1)
    s2 = random.randint(0,nombre_sommet-1)
    
    while (s1 == s2):
        s2 = random.randint(0,nombre_sommet-1)
        
    valeurs[s1][s2] = cost
    valeurs[s2][s1] = cost

file.write("\\n")
for i in range(nombre_sommet):
    file.write(nom_sommet[i])
    if(i != nombre_sommet -1):
        file.write(",")
        
for elm in range(nombre_sommet):
    file.write("\\n{")
    for elm1 in range(nombre_sommet):
        file.write(str(valeurs[elm][elm1]))
        if(elm1 != nombre_sommet-1):
            file.write(",")
    file.write("}")
    if(elm != nombre_sommet-1):
        file.write("\\n")
        
file.write("$")
file.write("\n")

file.close()