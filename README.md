
# WF - Workflow Automation

**WF** is a minimalist and procedural language designed to simplify the automation of common tasks, such as building, installing, or deploying applications.


## 📖 Documentation
- [English](#english)
- [Français](#français)

## English


**WF** is a minimalist and procedural language designed to simplify the automation of common tasks, such as building, installing, or deploying applications.

The `wf` binary is a command-line interface that automatically detects `.wf` files in the current directory. These files are structured into **sections**, containing sequentially organized commands to execute workflows efficiently and reproducibly.

With WF, you can centralize and standardize your automation processes while maintaining a clear and accessible syntax.

---

## Structure of a `.wf` File

A `.wf` file is organized into **sections**, each corresponding to a specific step in your workflow. Each section contains a list of commands to execute in a precise order.

### Basic Example
```wf
[section-name] # Comment: Description of the section
SET VARIABLE=value
run command
```

---

## Example of a `.wf` File

### **Example 1: Deploying a Symfony Application**

```wf
[prepare-env] # Prepare environment variables and system clock
SET SYMFONY_CONSOLE=php bin/console
sync_time
echo "Environment preparation complete"

[install-dependencies] # Install dependencies
run composer install --no-dev --optimize-autoloader
run ${SYMFONY_CONSOLE} doctrine:database:create --if-not-exists
run ${SYMFONY_CONSOLE} doctrine:migrations:migrate --no-interaction
echo "Dependencies installed and database is up to date"

[build-assets] # Build and optimize assets
run npm install
run npm run build
run ${SYMFONY_CONSOLE} assets:install --symlink --relative
echo "Assets built and installed"

[set-permissions] # Set required permissions
set_permissions var/cache 775
set_permissions var/log 775
echo "Permissions set"

[deploy] # Full deployment workflow
wf prepare-env
wf install-dependencies
wf build-assets
wf set-permissions
notify_success "Symfony application deployed successfully"
```

### **Example 2: Deploying a Dockerized Application**

```wf
[prepare-docker] # Prepare Docker environment
SET DOCKER_COMPOSE=docker compose
sync_time
echo "Docker environment prepared"

[build-containers] # Build Docker containers
run ${DOCKER_COMPOSE} build
echo "Docker containers built"

[start-containers] # Start Docker containers
run ${DOCKER_COMPOSE} up -d
echo "Docker containers started"

[cleanup] # Clean up old images
run ${DOCKER_COMPOSE} prune --force
echo "Old Docker images removed"

[deploy] # Full deployment workflow
wf prepare-docker
wf build-containers
wf start-containers
wf cleanup
notify_success "Dockerized application deployed successfully"
```

---

### **Key Features Across Examples**

1. **Modularity with Sections:** Each section corresponds to a logical step in the deployment process.
2. **Reusability:** Sections like `prepare-env` or `set-permissions` can be reused in different projects.
3. **Notifications:** Use `notify_success` or other notifications to indicate the completion or status of deployments.
4. **Environment Variables:** Variables like `SYMFONY_CONSOLE` or `DOCKER_COMPOSE` simplify command reuse.
5. **Simplicity:** The examples remain procedural and straightforward, aligning with the philosophy of WF.

---

## Available Commands

1. **`SET VARIABLE=value`**: Defines an environment variable to be used in commands.
2. **`run command`**: Executes a system command or script.
3. **`copy src dest`**: Copies a file or directory.
4. **`sync_time`**: Synchronizes the system clock with a time server.
5. **`set_permissions path mode`**: Changes the permissions of a file or directory.
6. **`echo message`**: Prints a message to the standard output.
7. **`exit`**: Immediately terminates the script execution with an exit code of `0`.
8. **`touch file`**: Creates an empty file unless it already exists.
9. **`mkdir folder`**: Creates a directory, including any missing parent directories.
10. **`docker_compose command`**: Executes a `docker compose` or `docker-compose` command, depending on availability.
11. **`wf workflow_name`**: Executes a workflow defined in the `.wf` file.
12. **`notify message`**: Displays a generic notification.
13. **`notify_success message`**: Displays a success notification.
14. **`notify_error message`**: Displays an error notification.
15. **`notify_warning message`**: Displays a warning notification.
16. **`notify_info message`**: Displays an informational notification.

---

## **How Does WF Work?**

1. **Scanning `.wf` Files**:  
   - When executed, `wf` scans the current directory and identifies all files with the `.wf` extension.
   - Each `.wf` file is read, and all defined **sections** are collected.

2. **Structure of Sections**:  
   - Each `.wf` file is divided into **sections** marked by a section header, such as `[section-name]`.
   - Sections contain commands to be executed in order.

3. **Executing a Section**:  
   - Once all sections from the `.wf` files are loaded, you can execute a specific section using the command:  
     ```bash
     wf section-name
     ```
   - WF will then execute all commands defined in that section.

---

## **Example Usage**

### Example `.wf` File
```wf
[docker-build] # Build Section: Build the application with Docker
SET DOCKER_COMPOSE=docker compose exec -T app-gaia
run ${DOCKER_COMPOSE} ./wf build

[build] # Build Section: Build the application
sync_time
run php -v
SET SYMFONY_CONSOLE=php bin/console
copy .env .env.local
run composer install --optimize-autoloader
run ${SYMFONY_CONSOLE} doctrine:migrations:migrate --no-interaction
```

### Available Commands in This Example

When the `wf` binary is executed, the sections are detected and listed as available commands.

#### Example of Available Commands:
```bash
NAME:
   wf - A new CLI application

USAGE:
   wf [global options] command [command options]

COMMANDS:
   build         Build Section: Build the application
   docker-build  Build Section: Build the application with Docker
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

---

## **Executing a Section**

1. **List Available Sections**  
   To display all available sections in the `.wf` files in the current directory:
   ```bash
   wf --help
   ```

2. **Execute a Specific Section**  
   To execute a specific section, simply use its name as a command:
   ```bash
   wf docker-build
   ```
   This will execute all commands defined in the `[docker-build]` section.

---

## **Advantages of WF**

- **Simple to Use**: No need to configure a complex tool like Makefile; everything is centralized in `.wf` files.
- **Modular**: Divide your workflows into multiple reusable sections.
- **Flexible**: Easily add your own commands and variables.

With `wf`, you can efficiently automate your build, installation, or deployment tasks by unifying everything in a simple and readable tool.


## License

WF is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.  

### Key Terms:
- You are free to use this tool in any project or environment, including commercial ones.
- Any modifications or improvements to this software must be made public under the same license.
- This tool cannot be integrated into proprietary software or systems without respecting the terms of this license.

For the full license text, see the LICENSE file.


## Français

**WF** est un langage minimaliste et procédural conçu pour simplifier l'automatisation des tâches courantes, telles que la construction, l'installation ou le déploiement d'applications.

Le binaire `wf` est une interface en ligne de commande qui détecte automatiquement les fichiers `.wf` présents dans le répertoire courant. Ces fichiers sont structurés en **sections**, contenant des commandes organisées de manière séquentielle pour exécuter des workflows de manière efficace et reproductible.

Avec WF, vous pouvez centraliser et standardiser vos processus d'automatisation tout en gardant une syntaxe claire et accessible.

---

## Structure d'un Fichier `.wf`

Un fichier `.wf` est organisé en **sections**, chacune correspondant à une étape spécifique de votre workflow. Chaque section contient une liste de commandes à exécuter dans un ordre précis.

### Exemple Basique
```wf
[section-name] # Commentaire : Description de la section
SET VARIABLE=value
run command
```


## Exemple de fichier `.wf`

### **Exemple 1 : Déploiement d'une application Symfony**

```wf
[prepare-env] # Préparation des variables d'environnement et de l'horloge système
SET SYMFONY_CONSOLE=php bin/console
sync_time
echo "Préparation de l'environnement terminée"

[install-dependencies] # Installation des dépendances
run composer install --no-dev --optimize-autoloader
run ${SYMFONY_CONSOLE} doctrine:database:create --if-not-exists
run ${SYMFONY_CONSOLE} doctrine:migrations:migrate --no-interaction
echo "Les dépendances sont installées et la base de données est à jour"

[build-assets] # Construction et optimisation des assets
run npm install
run npm run build
run ${SYMFONY_CONSOLE} assets:install --symlink --relative
echo "Les assets ont été construits et installés"

[set-permissions] # Définition des permissions nécessaires
set_permissions var/cache 775
set_permissions var/log 775
echo "Les permissions ont été configurées"

[deploy] # Workflow complet de déploiement
wf prepare-env
wf install-dependencies
wf build-assets
wf set-permissions
notify_success "Application Symfony déployée avec succès"
```


### **Exemple 2 : Déploiement d'une application Dockerisée**

```wf
[prepare-docker] # Préparation de l'environnement Docker
SET DOCKER_COMPOSE=docker compose
sync_time
echo "Environnement Docker prêt"

[build-containers] # Construction des conteneurs Docker
run ${DOCKER_COMPOSE} build
echo "Les conteneurs Docker ont été construits"

[start-containers] # Démarrage des conteneurs Docker
run ${DOCKER_COMPOSE} up -d
echo "Les conteneurs Docker sont démarrés"

[cleanup] # Nettoyage des anciennes images
run ${DOCKER_COMPOSE} prune --force
echo "Les anciennes images Docker ont été supprimées"

[deploy] # Workflow complet de déploiement
wf prepare-docker
wf build-containers
wf start-containers
wf cleanup
notify_success "Application Dockerisée déployée avec succès"
```

---

### **Principales caractéristiques des exemples**

1. **Modularité avec les sections :** Chaque section correspond à une étape logique du processus de déploiement.
2. **Réutilisabilité :** Des sections comme `prepare-env` ou `set-permissions` peuvent être réutilisées dans différents projets.
3. **Notifications :** Utilisez `notify_success` ou d'autres notifications pour indiquer l'état ou la réussite des déploiements.
4. **Variables d'environnement :** Des variables comme `SYMFONY_CONSOLE` ou `DOCKER_COMPOSE` simplifient la réutilisation des commandes.
5. **Simplicité :** Les exemples restent procéduraux et simples, respectant la philosophie de WF.

## Commandes Disponibles

1. **`SET VARIABLE=value`** : Définit une variable d'environnement utilisée dans les commandes.
2. **`run command`** : Exécute une commande système ou un script.
3. **`copy src dest`** : Copie un fichier ou répertoire.
4. **`sync_time`** : Synchronise l'horloge système avec un serveur de temps.
5. **`set_permissions path mode`** : Change les permissions d'un fichier ou répertoire.
6. **`echo message`** : Affiche un message sur la sortie standard.
7. **`exit`** : Termine immédiatement l'exécution du script avec un code de sortie `0`.
8. **`touch file`** : Crée un fichier vide, sauf s'il existe déjà.
9. **`mkdir folder`** : Crée un répertoire, y compris les parents inexistants.
10. **`docker_compose command`** : Exécute une commande `docker compose` ou `docker-compose`, selon ce qui est disponible.
11. **`wf workflow_name`** : Exécute un workflow défini dans le fichier `.wf`.
12. **`notify message`** : Affiche une notification générique.
13. **`notify_success message`** : Affiche une notification de succès.
14. **`notify_error message`** : Affiche une notification d'erreur.
15. **`notify_warning message`** : Affiche une notification d'avertissement.
16. **`notify_info message`** : Affiche une notification d'information.


## **Comment Fonctionne WF ?**

1. **Recherche des Fichiers `.wf` :**  
   - Lors de l’exécution, `wf` scanne le répertoire courant et identifie tous les fichiers avec l'extension `.wf`.
   - Chaque fichier `.wf` est lu, et toutes les **sections** définies sont collectées.

2. **Structure des Sections :**  
   - Chaque fichier `.wf` est divisé en **sections** délimitées par un en-tête de section, comme `[section-name]`.
   - Les sections contiennent des commandes à exécuter dans l'ordre.

3. **Exécution d'une Section :**  
   - Une fois que toutes les sections des fichiers `.wf` ont été chargées, vous pouvez exécuter une section spécifique en utilisant la commande :  
     ```bash
     wf section-name
     ```
   - WF exécute alors toutes les commandes définies dans cette section.

---

## **Exemple d'Utilisation**

### Exemple de Fichier `.wf`
```wf
[docker-build] # Section Build : Build de l'application avec Docker
SET DOCKER_COMPOSE=docker compose exec -T app-gaia
run ${DOCKER_COMPOSE} ./wf build

[build] # Section Build : Build de l'application
sync_time
run php -v
SET SYMFONY_CONSOLE=php bin/console
copy .env .env.local
run composer install --optimize-autoloader
run ${SYMFONY_CONSOLE} doctrine:migrations:migrate --no-interaction
```

### Commandes Disponibles dans Cet Exemple

Lors de l'exécution du binaire `wf`, les sections sont détectées et listées comme des commandes disponibles.

#### Exemple de Commandes Disponibles :
```bash
NAME:
   wf - A new cli application

USAGE:
   wf [global options] command [command options]

COMMANDS:
   build         Section Build : Build de l'application
   docker-build  Section Build : Build de l'application on Docker
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

---

## **Exécution d'une Section**

1. **Lister les Sections Disponibles**  
   Pour afficher toutes les sections disponibles dans les fichiers `.wf` du répertoire courant :
   ```bash
   wf --help
   ```

2. **Exécuter une Section**  
   Pour exécuter une section spécifique, utilisez simplement son nom comme commande :
   ```bash
   wf docker-build
   ```
   Cela exécutera toutes les commandes définies dans la section `[docker-build]`.


## **Avantages de WF**

- **Simple à Utiliser :** Pas besoin de configurer un outil complexe comme Makefile, tout est centralisé dans des fichiers `.wf`.
- **Modulaire :** Divisez vos workflows en plusieurs sections réutilisables.
- **Flexible :** Ajoutez facilement vos propres commandes et variables.

Avec `wf`, vous pouvez automatiser efficacement vos tâches de build, d’installation, ou de déploiement en unifiant tout dans un outil simple et lisible.

## Licence

WF est distribué sous la **GNU Affero General Public License v3.0 (AGPL-3.0)**.

### Principaux termes :
- Vous êtes libre d'utiliser cet outil dans tout projet ou environnement, y compris des projets commerciaux.
- Toute modification ou amélioration de ce logiciel doit être rendue publique sous la même licence.
- Cet outil ne peut pas être intégré dans des logiciels ou systèmes propriétaires sans respecter les termes de cette licence.

Pour consulter le texte complet de la licence, veuillez vous référer au fichier LICENSE.