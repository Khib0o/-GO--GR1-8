touch ../Fichiers_Ressources/${1}

for (( nombre=5; nombre<=$2; nombre++ ))
do
    let fill=$nombre*$nombre/2
    python3 ../Fichiers_Ressources/GraphCreator.py $nombre $fill ../Fichiers_Ressources/$1
done