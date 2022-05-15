# 1. Install Linux updates, set time zones, followed by GCC and Make

# sudo yum -y update
# sudo ln -sf /usr/share/zoneinfo/Canada/Toronto \
# /etc/localtime
# sudo yum -y install gcc make

sudo apt update
sudo apt install build-essential
gcc --version 

# 2. Download, Untar and Make Redis 2.8 (check here http://redis.io/download)

# cd /tmp
# wget http://download.redis.io/releases/redis-6.2.6.tar.gz
# tar xzf redis-6.2.6.tar.gz
# cd redis-6.2.6
# make

sudo sysctl vm.overcommit_memory=1 
cd /usr/local/src
sudo wget http://download.redis.io/releases/redis-6.2.6.tar.gz
sudo tar xvzf redis-6.2.6.tar.gz
sudo rm -f xvzf redis-6.2.6.tar.gz
cd redis-6.2.6
# sudo yum groupinstall "Development Tools"
sudo make distclean
sudo make
# sudo yum install -y tcl
sudo apt install -y tcl

sudo make test

# 3. Create Directories and Copy Redis Files

sudo mkdir -p /etc/redis /var/lib/redis /var/redis/6379
sudo cp src/redis-server src/redis-cli /usr/local/bin/
sudo cp redis.conf /etc/redis/6379.conf
 

# 4. Configure Redis.Conf

sudo vi /etc/redis/6379.conf
  [..]
  daemonize yes
  [..]
  
  [..]
  bind 127.0.0.1
  [..]
  
  [..]
  dir /var/redis/6379
  [..]

  [..]
  maxmemory 6gb
  [..]

# 5. Download init Script

sudo wget https://raw.github.com/saxenap/install-redis-amazon-linux-centos/master/redis-server
 

# 6. Move and Configure Redis-Server

Note: The redis-server to be moved below is the one downloaded in 5 above.

sudo mv redis-server /etc/init.d
sudo chmod 755 /etc/init.d/redis-server
sudo vi /etc/init.d/redis-server
  REDIS_CONF_FILE="/etc/redis/6379.conf"
  # redis="/usr/local/bin/redis-server" # maybe not due to the above
 

# 7. Auto-Enable Redis-Server

sudo chkconfig --add redis-server
sudo chkconfig --level 345 redis-server on
 

# 8. Start Redis Server

sudo service redis-server start

# 9. Ensure background saves and prevent low-mem issue
sudo vi /etc/systctl.conf
  vm.overcommit_memory = 1
  systctl vm.overcommit_memory=1

# 10. Test Redis server
redis-cli ping

// Taken from modified from https://gist.github.com/FUT/7db4608e4b8ee8423f31