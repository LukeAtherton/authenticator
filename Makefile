.PHONY:	build push publish launch

crypto:
	openssl genrsa -out ./crypto/testKey.pem 2048;
	openssl rsa -in authKey.pem -pubout > ./crypto/testKey.pub;

build:	
	docker build -t docker.hivebase.io:5000/authenticator .

push:	
	docker push docker.hivebase.io:5000/authenticator

publish:	
	make build;
	make push;

launch:	
	-fleetctl stop authenticator.service authenticator-register.service; 
	fleetctl unload authenticator.service authenticator-register.service; 
	-fleetctl destroy authenticator.service authenticator-register.service;
	fleetctl submit ./systemd/authenticator.service; 
	fleetctl submit ./systemd/authenticator-register.service; 
	fleetctl start authenticator.service authenticator-register.service;
	fleetctl journal -f authenticator;

deploy:
	make publish;
	make launch;