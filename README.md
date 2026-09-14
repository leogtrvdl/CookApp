# CookApp

Application de partage de recettes de cuisine, pensée pour aller à l'essentiel. Pas de commentaires, pas de superflu — juste des recettes à poster et à consulter.

🔗 **[Voir l'application](https://cookapp-front.onrender.com)**

> ⚠️ Hébergé sur le free tier de Render — la première requête après une période d'inactivité peut prendre ~30 secondes.

---

## Stack technique

| Couche | Technologie |
|--------|-------------|
| Backend | Go, Chi |
| Base de données | PostgreSQL (pgx) |
| Stockage images | Cloudflare R2 |
| Frontend | Expo / React Native Web |
| Infra | Docker Compose |
| Déploiement | Render |
| CI | GitHub Actions |

---

## Architecture

```
CookApp/
├── cookapp-back/      # API REST en Go
│   ├── cmd/api/           # Point d'entrée
│   ├── internal/
│   │   ├── handler/       # Handlers HTTP
│   │   ├── service/       # Logique métier
│   │   ├── repository/    # Accès base de données
│   │   ├── middleware/    # Auth JWT
│   │   └── storage/       # Cloudflare R2
│   └── tests/             # Tests d'intégration
├── cookapp-front/     # Frontend Expo / React Native Web
└── docker-compose.yml
```

---

## Fonctionnalités

- **Authentification** — inscription, connexion, JWT
- **Recettes** — CRUD complet avec upload d'image
- **Recherche** — recherche de recettes et d'utilisateurs
- **Likes & Favoris** — interactions sur les recettes
- **Amis** — système de demandes d'amis
- **Stockage images** — avatars et photos de recettes via Cloudflare R2

---

## Points techniques notables

- **API REST structurée en couches** (handler / service / repository) — séparation claire des responsabilités pour un code maintenable
- **Authentification par JWT** — à la connexion, un token est généré et envoyé au client ; il est vérifié à chaque requête protégée via un middleware dédié
- **Chiffrement des mots de passe avec bcrypt** — les mots de passe ne sont jamais stockés en clair
- **Upload d'images vers Cloudflare R2** — détection du type de fichier réelle, clés UUID uniques pour éviter les collisions
- **Tests d'intégration sur base PostgreSQL dédiée** — les tests tournent sur une base séparée de la production
- **CI GitHub Actions** — build et tests lancés automatiquement à chaque push
- **Déploiement conteneurisé avec Docker** sur Render, redéploiement automatique à chaque push
