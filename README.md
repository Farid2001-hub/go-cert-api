# go-cert-api

Cette application utilise plusieurs technologies pour fonctionner efficacement :

- **GORM** : Un ORM (Object-Relational Mapper) pour Go qui permet une interaction facile avec la base de données PostgreSQL. 
GORM est utilisé pour gérer toutes les opérations de base de données, y compris le stockage et la récupération des certificats.
- **Gin** : Un framework web Go performant qui est utilisé pour gérer les routes de l'API et faciliter le débogage.

# Table des matières
1. [Informations](#informations)
2. [Docker](#docker)
    - [Services](#services)
        - [Postgres](#1-postgres)
        - [Pocketbase](#2-pocketbase)
        - [GoCertAPI](#3-gocertapi)
    - [Volumes](#volumes)
    - [Lancement du projet](#lancement-du-projet)
3. [Instructions pour se connecter à Pocketbase](#instructions-pour-se-connecter-à-pocketbase)
4. [CRUD Utilisateur](#crud-utilisateur)
    - [Authentification](#authentification)
    - [Création d'un utilisateur](#création-d'un-utilisateur)
    - [Affichage](#affichage)
    - [Test des Permissions](#Test-des-Permissions)
5. [Gestion des Certificats](#gestion-des-certificats)
    - [Generate Certificate Authority](#Generate-Certificate-Authority)
    - [Affichage de la CA](#Affichage-de-la-CA)
    - [Generate Certificate CSR](#Generate-Certificate-CSR)
    - [Generate client certificate](#generate-client-certificate)
    - [Revoke Certificate](#revoke-certificate)
    - [Répertorier tous les certificats clients, CSR et CA](#Répertorier-tous-les-certificats-clients-CSR-et-CA)
    - [Répertorier tous les certificats clients, CSR et CA by ID](#Répertorier-tous-les-certificats-clients-CSR-et-CA-by-ID)
6. [Certificate Monitoring](#certificate-monitoring)
7. [Test Unitaire](#test-unitaire)
8.  [Mise en place de l'infrastructure](#mise-en-place-de-l'infrastructure)

# Informations

Le guide ci-dessous vous permet de tester rapidement les principales fonctions de l'API. Toutes ces actions peuvent être effectuées via Thunder Client ou Podaman.

Pour toutes les requêtes ci-dessous, vous devez ajouter le jeton JWT dans l'en-tête : `Authorization: eyJhbGciOi..`

Il y a également un middleware d'authentification qui vérifie le rôle de l'utilisateur actuel et permet ou non l'accès aux points de terminaison. Si vous obtenez un code HTTP `401`, cela signifie que vous n'êtes pas connecté. Si c'est un code `403`, cela signifie que vous n'avez pas les permissions nécessaires.

- Lecture et écriture pour les utilisateurs ayant le rôle `vipUser`
- Lecture seule pour les utilisateurs ayant le rôle `normalUser`
- Aucune permission pour les utilisateurs ayant le rôle `externUser`

##  2 Docker
Configuration Docker Compose

Ce fichier Docker Compose configure trois services : Postgres, Pocketbase et GoCertAPI. Voici les explications pour chaque service :

## Services

### 1. Postgres
- **Image** : Utilise l'image `postgres:alpine3.19` pour une base de données Postgres légère.
- **Environnement** : Configure la base de données Postgres avec :
  - `POSTGRES_USER` : Nom d'utilisateur pour la base de données (`test`).
  - `POSTGRES_PASSWORD` : Mot de passe pour la base de données (`password`).
  - `POSTGRES_DB` : Nom de la base de données (`test`).
- **Ports** : Mappe le port `5432` de l'hôte au port `5432` du conteneur.
- **Volumes** : Persiste les données de Postgres dans un volume nommé `postgres_data`.
- **Healthcheck** : Vérifie l'état de santé du service Postgres en exécutant `pg_isready -U test` toutes les 10 secondes, avec jusqu'à 5 tentatives.

### 2. Pocketbase
- **Build** : Construit le service Pocketbase à partir d'un Dockerfile situé à `./Dockerfile.pocketbase`.
- **Ports** : Mappe le port `8090` de l'hôte au port `8090` du conteneur.
- **Volumes** : Monte le répertoire local `./pb_data` dans `/pb/pb_data` dans le conteneur pour un stockage persistant.
- **Healthcheck** : Vérifie l'état de santé du service Pocketbase en exécutant `wget -q --spider http://localhost:8090/api/health || exit 1` toutes les 10 secondes, avec jusqu'à 5 tentatives.

### 3. GoCertAPI
- **Build** : Construit le service GoCertAPI à partir d'un Dockerfile situé à `./Dockerfile`.
- **Depends_on** : S'assure que les services Postgres et Pocketbase sont en bonne santé avant de démarrer.
- **Environnement** : Configure le service GoCertAPI avec :
  - `DB_HOST` : Nom d'hôte du service Postgres (`postgres`).
  - `DB_PORT` : Port du service Postgres (`5432`).
  - `DB_AUTH` : Détails d'authentification pour la base de données au format JSON (`{"username":"test","password":"password"}`).
  - `DB_NAME` : Nom de la base de données (`test`).
- **Ports** : Mappe le port `8080` de l'hôte au port `8080` du conteneur.

## Volumes

### 1. postgres_data
- Un volume nommé pour persister les données de Postgres en dehors du cycle de vie du conteneur.

## Lancement du projet

Pour lancer ce projet, suivez les étapes ci-dessous :

1. **Cloner le dépôt**

   Commencez par cloner le dépôt Git en utilisant la commande suivante :

   ```bash
   git clone https://gitlab.com/spv564816/go-cert-api.git


2. **lancer les contenaires et builder l'image**

   ```bash
  docker compose up --build

## 3 Instructions pour se connecter à Pocketbase

Pour se connecter au système d'authentification Pocketbase, suivez ces étapes :

1. Accédez à l'URL de Pocketbase : [http://localhost:8090/_](http://localhost:8090/_)
2. Utilisez les identifiants suivants pour vous connecter :
   - **Nom d'utilisateur :** admin@sdv.com
   - **Mot de passe :** admin12345

Une fois connecté, vous aurez accès au système d'authentification Pocketbase.

## 4 CRUD Utilisateur

### Authentification

Pour interagir avec l'API, vous devez d'abord vous authentifier et obtenir un jeton.

##### Endpoint

**POST** `http://localhost:8080/auth`

##### Corps de la Requête

```json
{
  "identity": "admin",
  "password": "admin12345"
}

```
##### Exemple de Réponse

```
{
  "record": {
    "collectionId": "_pb_users_auth_",
    "collectionName": "users",
    "created": "2024-06-26 20:59:50.422Z",
    "email": "admin@test.com",
    "emailVisibility": false,
    "id": "123456789012345",
    "role": "vipUser",
    "updated": "2024-06-26 20:59:50.422Z",
    "username": "admin",
    "verified": false
  },
  "token": "eyJhbGciOiJI........."
}
```

##### Utilisation du Jeton

Une fois le jeton obtenu, ajoutez-le dans l'en-tête Authorization pour accéder aux endpoints nécessitant une authentification.

### Création d'un utilisateur

Pour créer un enregistrement utilisateur dans le service d'authentification.

##### Endpoint

**POST** `http:localhost:8080/users`

##### Corps de la Requête

```json
{
    "email": "test@test.com",
    "role": "normalUser",
    "username": "test",
    "password": "test12345",
    "passwordConfirm": "test12345",
    "emailVisibility": false,
    "verified": true
}
```

##### Exemple de Réponse

```
{
  "collectionId": "_pb_users_auth_",
  "collectionName": "users",
  "created": "2024-06-29 10:59:47.250Z",
  "email": "test@test.com",
  "emailVisibility": false,
  "id": "ebcjnmh9xwipszl",
  "role": "normalUser",
  "updated": "2024-06-29 10:59:47.250Z",
  "username": "test",
  "verified": true
}
```
### Affichage 

Pour afficher les utilisateurs dans le service d'authentification.

##### Endpoint

**Get** `http:localhost:8080/users` et pour afficher qu'un seule CA **Get** `http:localhost:8080/users/id`

##### Exemple de Réponse

```
{
  "page": 1,
  "perPage": 30,
  "totalItems": 2,
  "totalPages": 1,
  "items": [
    {
      "collectionId": "_pb_users_auth_",
      "collectionName": "users",
      "created": "2024-06-26 20:59:50.422Z",
      "email": "admin@test.com",
      "emailVisibility": false,
      "id": "123456789012345",
      "role": "vipUser",
      "updated": "2024-06-26 20:59:50.422Z",
      "username": "admin",
      "verified": false
    },
    {
      "collectionId": "_pb_users_auth_",
      "collectionName": "users",
      "created": "2024-06-29 10:59:47.250Z",
      "email": "test@test.com",
      "emailVisibility": false,
      "id": "ebcjnmh9xwipszl",
      "role": "normalUser",
      "updated": "2024-06-29 10:59:47.250Z",
      "username": "test",
      "verified": true
    }
  ]
}
```
### Test des Permissions

Cet exemple montre comment tester les permissions en se connectant avec un utilisateur ayant le rôle `normalUser` et en essayant de créer un autre utilisateur.

##### Endpoint

**POST** `http:localhost:8080/users`

### Corps de la Requête

```json
{
    "email": "test1@test.com",
    "role": "normalUser",
    "username": "test1",
    "password": "test12345",
    "passwordConfirm": "test12345",
    "emailVisibility": false,
    "verified": true
}
```
##### Exemple de Réponse
```
{
  "message": "Interdit",
  "errorCode": 403,
  "detail": "Permissions insuffisantes"
}
```

## 5 Gestion des Certificates

### Generate Certificate Authority

##### Endpoint

**POST** `http:localhost:8080/ca`

### Corps de la Requête
```
{
  "common_name": "Super CA SDV",
  "organization": "Sup De Vinci",
  "country": "FR",
  "province": "Ile-de-france",
  "locality": "Paris",
  "validity_period": 3650,
  "key_size": 2048,
  "key_algorithm": "RSA"
}
```
Les algorithmes peuvent être :

    RSA avec une taille de clé allant de 2048 à 4096 bits.
    ECDSA avec une taille de clé de 256, 384 ou 521 bits.
    Ed25519 avec une taille de clé directement assignée par l'algorithme.

RSA (Rivest-Shamir-Adleman) est largement utilisé pour sa robustesse et sa flexibilité en matière de taille de clé, offrant un bon équilibre entre sécurité et performance. 
ECDSA (Elliptic Curve Digital Signature Algorithm) est apprécié pour sa sécurité élevée et son efficacité, même avec des clés plus courtes. 
Ed25519, basé sur la courbe elliptique Curve25519, est connu pour sa sécurité et ses performances optimales, étant particulièrement adapté aux applications nécessitant des signatures rapides et sécurisées.

### Affichage de la CA 

##### Endpoint

**Get** `http:localhost:8080/ca` et pour afficher qu'un seule utilisateur **Get** `http:localhost:8080/ca/id`

##### Exemple de Réponse

```
  {
    "CreatedAt": "2024-06-28T09:10:26.565669Z",
    "UpdatedAt": "2024-06-28T09:10:26.565669Z",
    "DeletedAt": null,
    "ID": 1,
    "Name": " CA ",
    "ExpirationDate": "2034-06-26T09:10:26.563549Z",
    "Certificate": "-----BEGIN CERTIFICATE-----\nMIIDoTC........rU\n-----END CERTIFICATE-----\n",
    "PrivateKey": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADA...N..=\n-----END PRIVATE KEY-----\n"
  }
```

### Generate Certificate CSR

##### Endpoint

**POST** `http:localhost:8080/csr`

### Corps de la Requête
```
{
  "common_name": "www.supdevinci.fr",
  "organization": "Sup De Vinci",
  "country": "FR",
  "province": "ile-de-france",
  "locality": "Puteaux",
  "key_size": 2048,
  "key_algorithm": "RSA",
  "dns_names": ["www.sdv.fr"],
  "ip_addresses": ["192.168.1.1", "192.168.1.2"],
  "emails": ["admin@sdv.fr", "information@sdv.fr"]
}
```
##### Exemple de Réponse

```
{
  "message": "CSR générée et enregistrée avec succès",
  "certificate": "-----BEGIN CERTIFICATE REQUEST-----\nMIIDBTCCAe0CAQAwa........10ni\n-----END CERTIFICATE REQUEST-----\n",
  "certificateID": 1
}
```
### Generate client certificate

##### Endpoint

**POST** `http:localhost:8080/certificate`

### Corps de la Requête
```
{
  {
  "ca_id": 1,
  "csr_id": 1,
  "validity_period": 365
}
}
```
##### Exemple de Réponse

```
{
  [
  {
    "CreatedAt": "2024-07-06T10:26:59.91624Z",
    "UpdatedAt": "2024-07-06T10:26:59.91624Z",
    "DeletedAt": null,
    "ID": 1,
    "Name": "www.supdevinci.fr",
    "ExpirationDate": "2025-07-06T10:26:59Z",
    "CAID": 1,
    "CSRID": 1,
    "Certificate": "-----BEGIN CERTIFICATE-----\nMIID8DCCAtigAwIBAgI....+oDYn\n3Wou+Q==\n-----END CERTIFICATE-----\n",
    "Revoked": false,
    "RevokedAt": "0001-01-01T00:00:00Z"
  }
]
}
```
### Revoke Certificate

##### Endpoint

**POST** `http:localhost:8080/certificate/Revoke:"id"`
##### Exemple de Réponse
```
{
[
  {
  "message": "Certificat révoqué avec succès"

  }
]
}
```

### Répertorier tous les certificats clients, CSR et CA
##### Endpoint

**Get** `http:localhost:8080/ca`

**Get** `http:localhost:8080/csr`

**Get** `http:localhost:8080/certificate`

### Répertorier tous les certificats clients, CSR et CA by ID
##### Endpoint

**Get** `http:localhost:8080/ca/{id}`

**Get** `http:localhost:8080/csr/{id}`

**Get** `http:localhost:8080/certificate/{id}`

### 6 Certificate Monitoring

Nous avons un message dans l'API qui vérifie constamment la base de données pour voir s'il y a des certificats expirés. Comme vous pouvez le voir dans les messages de l'API ci-dessous :
```
gocertapi-1 | 2024/07/07 10:46:05 Vérification des enregistrements CA et des certificats expirés...
gocertapi-1 | 2024/07/07 10:46:05 Aucun certificat expiré trouvé 🎉
gocertapi-1 | 2024/07/07 10:51:05 Vérification des enregistrements CA et des certificats expirés...
gocertapi-1 | 2024/07/07 10:51:05 Aucun certificat expiré trouvé 🎉
```
## 7 Test Unitaire

Les tests incluent :
- **Création d'un utilisateur**
- **Modification d'un utilisateur**
- **Suppression d'un utilisateur**
- **Création d'un ca**
- **Lister les Certificat d'Authority**

Les résultats des tests montrent que toutes les opérations de test ont été réussies.

### Exemple de sortie des tests
![Tests Unitaires](images/test.png)
![Tests Unitaires](images/test_create_ca.png)
![Tests Unitaires](images/test_list_ca.png)

Et nous avons commenté quelques tests comme vous pouvez le voir sur la photo. Vous pouvez exécuter le test unitaire que vous voulez.
![Tests Unitaires](images/test_unitaires.png)

## 8 Mise en place de l’infrastructure

Schéma de l'Amazon Elastic Container Service

Ce schéma illustre le fonctionnement et les composants principaux de l'Amazon Elastic Container Service (ECS).
![ECS](images/shéma.png)
Dans ce guide, je vais vous montrer comment mettre en place une infrastructure complète sur AWS ECS (Elastic Container Service) avec Fargate, puis comment automatiser son déploiement via GitLab CI/CD.

On va couvrir les points clés suivants :

    - **Création d'un cluster ECS**
    - **Mise en place d'un secret dans AWS Secrets Manager**
    - **Définition d'une Task Definition avec trois conteneurs (certapi, pocketbase, postgres)**
    - **Création et configuration d'un service ECS**
    - **Configuration du déploiement automatisé avec GitLab CI/CD**

# AWS ECS FARGATE

## Création du Cluster

Accédez à l'interface ECS (Elastic Container Service) après vous être connecté à votre compte AWS.
![Creation du cluster](images/page_d_accueil.png)

Sélectionnez l'option 'Créer un cluster' dans l'interface ECS pour initier le processus de configuration.:
![Creation du cluster](images/Créer_un_cluster.png)

Assignez un nom significatif à votre cluster. Dans cet exemple, nous utilisons 'CLUSTER_TEST'.
![Creation du cluster](images/cluster_test.png)

## Création d’un Secret

Pourquoi créer un Secret sur AWS ?

Dans notre cas le code de notre projet ainsi que la registry qui contient les differentes images de celui ci se trouvent dans gitlab le secret va donc permettre de se connecter a notre instance gitlab pour récuperer les images dans la registry et ainsi pouvoir les pull

Naviguez vers AWS Secrets Manager pour sécuriser les informations d'identification nécessaires à l'accès à votre registre d'images GitLab

![Creation du secret](images/Secrets_Manager.png)

Dans l'interface AWS Secrets Manager, sélectionnez 'Stocker un nouveau secret' pour commencer la configuration
![Creation du secret](images/stock_secret_manager.png)

Optez pour 'Autre type de secret' afin de stocker un jeton d'authentification GitLab.
![Creation du secret](images/type_secret.png)

Vérifiez et confirmez les détails du secret avant de le stocker définitivement.
![Creation du secret](images/stock.png)

Notez l'ARN (Amazon Resource Name) du secret créé, qui sera utilisé ultérieurement dans la configuration de la tâche ECS
![Creation du secret](images/arn.png)

## Création d’un task difinition
Une fois qu’on à la validation de la création du cluster ainsi que le ASM on va maintenant créer la task Définition celle ci va permettre de définir les différents containeurs que le cluster va contenir 

![Creation du task](images/task_definition.png)

## Vue de l'Amazon Elastic Container Service (ECS) - cert_api_service

![Creation du task](images/cert_api_service.png)

## Pipeline de CI/CD


Cette capture d'écran montre un pipeline de CI/CD typique comprenant les étapes suivantes :
1. **build** : Construction de l'application (`build_app`).
2. **test** : Exécution des tests unitaires (`test_unitaire`).
3. **push_to_registry_postgres** : Construction et push de l'image PostgreSQL (`build_and_push_postgres`).
4. **container_registry_push** : Push des conteneurs vers le registre de conteneurs (`container_registry`).
5. **push_to_registry_pocketbase** : Construction et push de l'image Pocketbase (`build_and_push_pocketbase`).
6. **deploy** : Déploiement sur Amazon ECS (`deploy_to_ecs`).

![Pipeline de CI/CD](images/jobs.png)










