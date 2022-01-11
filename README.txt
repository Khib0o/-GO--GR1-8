Archive -GO--GR1-8
produit par Aziz ben Jebalia, Houda Touil, Louis Teys

Dans cette archive se trouve notre application. Elle implémente
l'algorithme Dijkstra dans un modèle client-server. Le code source se
trouve dans le dossier src/GO. Nous avons inclus dans l'archive les scripts
que nous avons utiliser pour l'instant pour tester notre application.
Tous nos script lance l'application sur le port 25565.

Pour demarrer le server :

./src/Shell/run_server.sh

Pour lancer le client :

./src/Shell/run_client_once.sh
OU
go run src/GO/main_client.go <N° port> <Fichier contenant les graphs> <Numero de ligne du graph>
    <Afficher le résultat (y/n)> <Afficher le temps de réponse (y/n)>

-------------------------------------------------------------------------------

Dans le fichier src, il y a :

-Un dossier GO, dans lequel se trouve les codes GO que vous avons produit pour 
notre application

-Un dossier Fichiers_Ressources, dans lequel se trouve des fichiers txt stockant
 des graphiques à utiliser par le client

-Un dossier Shell, dans lequel se trouve les fichiers script shell utilisés
 pour tester notre application

-------------------------------------------------------------------------------

Dans le dossier GO, il y a :

-Un fichier main_client.go qui est le client de notre application. Il 
 s'utilise de la manière suivante :

go run main_client.go <Numero de port> <Fichier contenant les graphs> <Numero de ligne du graph>
                      <Afficher le résultat (y/n)> <Afficher le temps de réponse (y/n)>


-Un fichier main_server.go qui est le server principale de notre application. Il s'utilise de la
manière suivante :

go run main_server.go <Numero de port>


-Un fichier ancien_server.go qui est une version moins performante du server. Elle sera présentée
et sa présence justifiée lors de la démonstration. Il s'utilise de la manière suivante :

go run ancien_server.go <Numero de port>

-------------------------------------------------------------------------------

Dans le dossier Fichiers_Ressources, il y a :

-Un fichier GraphCreator.py qui nous a permis de générer des graphs automatiquement de taille
personnalisée. Il s'utilise de la manière suivante :

python3 GraphCreator.py <Taille> <Nombre de liens> <Nom du fichier d'écriture>

-Un fichier graphs_persos.txt qui sont des graphs créés à la main pour être utilisés dans le client

-Un fichier graphs_taille_croissante.txt contenant des graphs de la ligne 1 à la ligne 45
arrangé de manière à ce que chacun soit de taille n°ligne + 5. Ils ont été créés grâce au
script python GraphCreator.py

-------------------------------------------------------------------------------

Dans le dossier Shell, il y a :

-generation_fichiers_graphs_croissant.sh qui permet de générer un fichier contenant
des graphs de taille croissantes (+1 par ligne). Utilisation :

./generation_fichiers_graphs_croissant.sh <Nom du fichier de graph> <Nombre de graph>

-mesure.sh qui permet de mesurer les temps d'execution de plusieurs graphs (issus du
même fichier) les uns à la suite des autres rapidement. Utilisation :

./mesure.sh <Nom du fichier de graph> <N° de ligne finale>

-run_client_once.sh qui envoie simplement un problème faisant parti du fichier
graphs_persos.txt au server et attends la réponse

-run_server.sh permet de démarrer le server

-test_multiples_connexion.sh qui permet d'executer en arrière plan plusieurs problèmes
pour tester la robustesse du server. Utilisation :

./test_multiples_connexion.sh <Nombre d'execution>