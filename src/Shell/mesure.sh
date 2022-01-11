for (( nombre=1; nombre<=$2; nombre++ ))
do
    go run ../GO/main_client.go 25565 $1 $nombre n y
done