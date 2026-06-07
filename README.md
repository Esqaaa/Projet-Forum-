# 🏛️ PROJET FORUM - YNOV CAMPUS

<p align="center">
  <img src="https://img.shields.io/badge/Backend-Golang-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Database-MariaDB%20%2F%20XAMPP-003545?style=for-the-badge&logo=mariadb&logoColor=white" alt="MariaDB">
  <img src="https://img.shields.io/badge/Frontend-HTML5%20%2F%20CSS3-E34F26?style=for-the-badge&logo=html5&logoColor=white" alt="HTML/CSS">
  <img src="https://img.shields.io/badge/Script-JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black" alt="JS">
</p>

---

## 📝 Présentation du Projet
[cite_start]Ce projet consiste en la création d'une plateforme de forum web complète, développée par les étudiants dans le cadre de la formation Ynov Campus[cite: 1, 2]. [cite_start]Il permet aux utilisateurs de créer des comptes, d'ouvrir des discussions (topics), d'échanger des messages, et d'interagir via un système de likes/dislikes[cite: 5, 15, 30, 37]. 

[cite_start]L'accent a été mis sur la robustesse du code, l'interdiction de frameworks front-end[cite: 80], et une **gestion des erreurs optimale** à chaque étape du cycle de requêtes.

---

## 🛠️ Stack Technique

* [cite_start]**Backend :** Golang (Natif, aucun framework lourd) [cite: 80, 81]
* [cite_start]**Frontend :** HTML5 / CSS3 pour les templates (Conforme à la consistance "No-Front-Framework") [cite: 80]
* **Dynamisme Front :** JavaScript (Utilisé spécifiquement pour le rafraîchissement des cookies et l'interactivité fluide)
* [cite_start]**Base de données :** MariaDB (Via le SGBDR requis compatible MySQL) [cite: 84]
* **Environnement Local :** XAMPP (Hébergement de la base de données)

---

## 🚀 Installation et Lancement local

### Prérequis
1. Avoir **Go** installé sur votre machine.
2. Avoir **XAMPP** installé et démarré (Activer le module *MySQL*).

### Étape 1 : Préparation de la Base de Données
1. Lancez **XAMPP** et ouvrez **phpMyAdmin**.
2. Créez une nouvelle base de données (ex: `forum_db`).
3. Importez les scripts SQL situés dans le dossier `/database` dans l'ordre chronologique (ou utilisez `shema.sql`) pour générer les tables nécessaires.

### Étape 2 : Lancement du serveur Go
Ouvrez votre terminal à la racine du projet et exécutez la commande suivante :

```bash
go run .

Le site est désormais accessible localement sur votre navigateur (généralement à l'adresse http://localhost:8080).

📁 Arborescence du Projet

```
│   main.go
│   README.md
│   
├───database
│       01_users.sql
│       02_topics.sql
│       03_messages.sql
│       04_tools.sql
│       07_likes.sql
│       init.go
│       shema.sql
│       
├───handlers
│       admin.go
│       auth.go
│       home.go
│       likes.go
│       topics.go
│       
├───models
│       models.go
│       user.go
│       
├───static
│   ├───css
│   │       admin.css
│   │       forum.css
│   │       styles_login.css
│   │       topic.css
│   │       
│   ├───js
│   │       script_picture.js
│   │       toggle_edit.js
│   │       
│   └───uploads
│           
└───templates
    │   
    └───html
            admin_dashboard.html
            create_topic.html
            index.html
            layout.html
            login.html
            register.html
            view_topic.html
```
🛣️ Liste des Routes (Routing)
Conformément aux exigences, voici la séparation des routes de l'application:
👁️ Routes Distribuant une Vue (HTML)
- GET / : Page d'accueil (Liste des topics)   
- GET /login : Page de connexion   
- GET /register : Page d'inscription   
- GET /topic/view : Consultation d'un fil de discussion spécifique   
- GET /topic/create : Formulaire de création de topic (Utilisateurs connectés)   
- GET /admin/dashboard : Panel d'administration (Rôle Admin uniquement)   

## ⚙️ Routes pour le Traitement de Données
- POST /auth/register : Traitement de l'inscription (Vérifications pseudo/email uniques, hash SHA512)   
- POST /auth/login : Traitement de la connexion et génération des cookies de session  
- POST /auth/logout : Déconnexion de l'utilisateurPOST /topic/create : Enregistrement du topic en BDD   
- POST /message/post : Envoi d'un message dans un topic   
- POST /like : Gestion des likes et dislikes sur les messages   
- POST /admin/action : Actions de modération (Bannissement, suppression, modification d'état)   

##⚡ Méthodologie & Usage de l'IA (Vibe Coding)
Pour ce projet, nous avons adopté une approche moderne de Vibe Coding en collaborant activement avec des Intelligences Artificielles de pointe : Copilot, Gemini et ChatGPT.

Rôle des IA : Elles ont agi comme des copilotes de programmation (Pair Programming). Elles nous ont aidés à générer rapidement des structures de données en Go, à concevoir les schémas SQL complexes, et à automatiser les tâches répétitives sur le CSS.

Plus-value : Cette méthode nous a permis de nous concentrer sur l'architecture globale, la logique métier du Forum, et surtout sur la mise en place d'une gestion des erreurs ultra-robuste et optimale, garantissant qu'aucun crash serveur ne survienne lors de requêtes malformées.

## 👥 Composition de l'Équipe
https://github.com/loschoe, https://github.com/Esqaaa