for (( nombre=1; nombre<=$1; nombre++ ))
do
    go run ../GO/main_client.go 25565 ../Fichiers_Ressources/graphs_persos.txt 2 n y &
done