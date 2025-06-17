update-apt:
	sudo apt-get update

install-deps:
	sudo apt-get install -y ca-certificates curl gnupg

add-docker-key:
	sudo install -m 0755 -d /etc/apt/keyrings
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
	sudo chmod a+r /etc/apt/keyrings/docker.gpg

add-docker-repo:
	echo "deb [arch=$$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $$(. /etc/os-release && echo "$$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
	sudo apt-get update

install-docker:
	sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
	sudo systemctl start docker
	sudo systemctl enable docker

verify-docker:
	sudo docker run hello-world

add-user-to-docker:
	sudo usermod -aG docker $USER
	@echo "Перезайдите в систему для применения изменений"

verify-compose:
	docker compose version

#### запуск сервиса
BuildImages:
	docker build -t migration:1 -f dockerfilemigration . && docker build -t sso_backend:1 .

StartDockerCompose: BuildImages
	docker compose up -d

down:
	docker compose down

### установка docker и docker compose если не установлен 
install-docker-full: install-deps add-docker-key add-docker-repo install-docker verify-docker add-user-to-docker verify-compose
	@echo "Установка Docker завершена. Пожалуйста, перезайдите в систему."

