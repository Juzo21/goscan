
# 🚀 Mon Projet de Scanner d'URLs

Un scanner d'URLs simple et efficace en Go, permettant de tester rapidement une liste de mots sur une URL de base pour détecter des chemins potentiels. 

---

## 📖 Table des matières

- [🔍 Fonctionnalités](#-fonctionnalités)
- [🛠️ Prérequis](#-prérequis)
- [📦 Installation](#-installation)
- [⚙️ Utilisation](#-utilisation)
- [📄 Exemple](#-exemple)
- [🛡️ Contribuer](#-contribuer)
- [📄 Licence](#-licence)

---

## 🔍 Fonctionnalités

- 🔄 Tester des chemins d'URLs à partir d'une liste de mots.
- 🔍 Support des extensions personnalisées.
- 🌐 Proxy et authentification HTTP intégrés.
- ⏳ Délai configurable entre les requêtes.
- 📜 Mode silencieux pour des exécutions discrètes.
- 🛠️ Chargement de certificats clients pour des connexions sécurisées.

---

## 🛠️ Prérequis

- Go (version 1.16 ou supérieure)
- Accès à internet pour le fonctionnement du scanner.

---

## 📦 Installation

1. Clonez le dépôt :

   ```bash
   git clone https://github.com/softwaretobi/goscan.git
   ```

2. Accédez au répertoire du projet :

   ```bash
   cd goscan
   ```

3. Exécutez le projet :

   ```bash
   go run goscan.go -url <URL> -wordlist <file>
   ```

---

## ⚙️ Utilisation

Voici comment utiliser le scanner :

```bash
go run goscan.go -url <URL> -wordlist <file>
```

### Options disponibles :

- `-url <URL>`: URL de base à scanner
- `-wordlist <file>`: Fichier contenant la wordlist
- `-a <agent>`: Agent utilisateur personnalisé
- `-c <cookie>`: Cookie pour les requêtes HTTP
- `-E <certificat>`: Chemin vers le certificat client
- `-p <proxy>`: URL du proxy
- `-X <extensions>`: Extensions à ajouter aux mots
- `-z <delay>`: Délai entre les requêtes en millisecondes
- `-N <ignoreCode>`: Ignorer les réponses avec ce code HTTP
- `-u <user>`: Nom d'utilisateur pour l'authentification HTTP
- `-P <password>`: Mot de passe pour l'authentification HTTP
- `-i`: Recherche insensible à la casse
- `-S`: Mode silencieux (ne pas afficher les mots testés)
- `-h`: Afficher cette aide

---

## 📄 Exemple

Pour scanner une URL de base `https://example.com` avec une wordlist `wordlist.txt`, utilisez la commande suivante :

```bash
go run main.go -url https://example.com -wordlist wordlist.txt
```

---

## 🛡️ Contribuer

Si vous souhaitez contribuer à ce projet, n'hésitez pas à soumettre une demande de tirage (pull request) ou à signaler des problèmes (issues). Toute contribution est la bienvenue !

---

## 📄 Licence

Ce projet est sous licence [MIT](LICENSE).

---

Merci d'avoir consulté ce projet ! N'hésitez pas à me contacter si vous avez des questions ou des suggestions. 🙌
