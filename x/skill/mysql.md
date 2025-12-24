Upgrading from MySQL <8.4 to MySQL >9.0 requires running MySQL 8.4 first:
 - brew services stop mysql
 - brew install mysql@8.4
 - brew services start mysql@8.4
 - brew services stop mysql@8.4
 - brew services start mysql

We've installed your MySQL database without a root password. To secure it run:
    mysql_secure_installation

MySQL is configured to only allow connections from localhost by default

To connect run:
    mysql -u root

To start mysql now and restart at login:
  brew services start mysql
Or, if you don't want/need a background service you can just run:
  /usr/local/opt/mysql/bin/mysqld_safe --datadir\=/usr/local/var/mysql

pwgen -s -y 16 1
  |$tN_Rw!-"Eu<[50


To start redis now and restart at login:
  brew services start redis
Or, if you don't want/need a background service you can just run:
  /usr/local/opt/redis/bin/redis-server /usr/local/etc/redis.conf
==> Summary
🍺  /usr/local/Cellar/redis/8.0.2: 15 files, 2.7MB
==> Running `brew cleanup redis`...
Disable this behaviour by setting HOMEBREW_NO_INSTALL_CLEANUP.
Hide these hints with HOMEBREW_NO_ENV_HINTS (see `man brew`).  