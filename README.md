Выполните следующие команды для запуска сервиса на linux:
1) git clone https://github.com/Senkoker/SSO_service
   
Если не установлен docker и docker-compose то запустите:
2) make install-docker-full

Если установлено все,то запускаем сервис :
2) В установленной директории находим config/config.yaml, изменяем почту на гуголовскую (так как используется порт gmail, при использовании почты гугл необходим аккаунт с двухсторонней аутентификацией, пароль - полученный пароль ДЛЯ СТОРОННИХ ПРИЛОЖЕНИЙ, доступно только для аккаунта с двухсторонней аутентификацией)

3) make StartDockerCompose
   
Чтобы удаляем контейнеры:
4) make down

Функции:
a)Register - получаем userID для проверки что пользователь был добавлен в список неподтвержденных аккаунтов. При данном действии получаем письмо на почту:
![image](https://github.com/user-attachments/assets/49f4279e-1176-49a2-ba77-b3611d12a05d)
Так как не запущен сервис,который будет парсить URLparams, при подтверждении письма копируем из ссылки:
![image](https://github.com/user-attachments/assets/5fe4f967-7585-496f-ac95-aca5e29b114d)
б) Accept - используется для подверждения почты,подтвержения смены пароля. Ответ - успех или не успех операции: 
![image](https://github.com/user-attachments/assets/47f71b72-c1bb-40b2-9f98-5e65fb162e7b)
в) Login - логирует пользователя, получаем jwt-токен (appid - для серисов, которые будут использовать наш SSO, для каждого appid свой secret-key, в качестве теста appid = 1):
![image](https://github.com/user-attachments/assets/cc9588f9-0d26-4885-8819-4efda1c8aa1a)
г) Change_password - смена пароля пользователем. Смотрите почту пришло подтверждение:
![image](https://github.com/user-attachments/assets/2c82e23e-b962-42bc-bfb3-8fbef11ea773)
д) Используем accept для подтвержения действия:
![image](https://github.com/user-attachments/assets/1a4c1574-ca4b-4cf2-839a-40ee097847b1)
e) проверяем поменялся ли пароль:
Вводим старый пароль и наблюдаем ошибку:
![image](https://github.com/user-attachments/assets/c7a6f471-87bc-462a-b73c-22fd4d528317)
Вводим новый пароль; пароль поменялся :
![image](https://github.com/user-attachments/assets/666eaf55-47d2-4387-9e7e-5482cf63cc71)
Если введем неправильный логин, то получим:
![image](https://github.com/user-attachments/assets/0f8fbff4-25fd-47f6-8502-98719ebcfd04)
Если попробуем повторно зарегистрироваться, то получим :
![image](https://github.com/user-attachments/assets/21d40e23-67ac-4fea-8789-247c33b93532)
Если попробуем повторно отправить письмо с подтверждением, получим:
![image](https://github.com/user-attachments/assets/76903c00-e467-468d-8617-3fd073a057a9)
Было бы неплохо доделать время ожидания подтверждения, чтобы отчистить пользователя, ждущего подтверждения больше указанного времени, иначе не сможем сделать повторный запрос.
















 

