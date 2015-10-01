#Authenticator
---

[![Circle CI](https://circleci.com/gh/LukeAtherton/authenticator.svg?style=svg)](https://circleci.com/gh/LukeAtherton/authenticator)

The Authenticator service provides basic signup, email verification, login and password reset functionality to the hivebase.io PaaS.

##Installation

***TODO:***

##Run Tests

You'll need a message queue and database to run the tests but docker can help with this.

Start MongoDB:

`$ docker run -d -p 27017:27017 --name mongodb mongo`

Start RabbitMQ:

`$ docker run -d -e RABBITMQ_NODENAME=test-rabbit --name rabbit-mq -p 5672:5672 -p 25672:25672 -p 4369:4369 -p 44001:44001 -p 8080:15672 rabbitmq:3-management`

You'll also need some RSA keys for crytpo, you can use something like openssl to generate fresh ones:

`$ openssl genrsa -out ./crypto/testKey.pem 2048`

`$ openssl rsa -in ./crypto/testKey.pem -pubout > ./crypto/testKey.pub`

Then you can run the tests:

`$ ginkgo -r -v`

##Configuration
```
--config: path to yaml config file
--mq-address: exchange address
--mq-topic: exchange topic
--mq-username: ampq username
--mq-password: ampq password
--db-hosts: address list of db hosts
--db-auth: db to auth against
--db-username: db username
--db-password: db password
--crypto-private-key: path to private key
--crypto-public-key: path to public key
```
##API Resources

###Service Status

This endpoint is used internally to the PaaS to check the status of the service.

####URI

`GET /status`

####Response

`200` if service is running.

###Email Availability Check

Email availability endpoint, used to check if an email is available for signup.

####URI

`GET /api/emails?email=<user_email>`

####Response

`200` if email address is available for signup.

`302` if email address is in use.

***TODO: This leaks user info***

####Error Codes

```
{
  'code':3,
  'message':'email is a required field'
}
```
This means no email address was supplied in the request. 

###Email Verification Link

Link sent in email verification message. Uses a code to verify that an email address belongs to a user.

####URI

`GET /api/verification?<user_email>&code=<verification_code>`

####Response

`200` if email verified

`400` if email verification failed.

####Error Codes

```
{
  'code':3,
  'message':'email verification failed'
}
```
This means either an invalid email address or invalid code.

####Events

`NewEmailVerifiedEvent` Event published on successful verification.

###New Signup

Resource for new user signups to the platform.

####URI

`POST /api/registrations`

####Request

```
{
  'email':'string'
  'password':'string'
}
```

####Response

`201` if signup successful.

```
{
  'id':'UUID',
  'email':'string',
  'token':'JWT Token'
}
```

`400` if signup failed.

####Error Codes

```
{
  'code':3,
  'message':'email is a required field'
}
```
This means no email was supplied.

```
{
  'code':3,
  'message':'password is a required field'
}
```
This means no password was supplied.

```
{
  'code':2,
  'message':'user email exists'
}
```
This means the email address supplied already exists in the system.

***TODO: This leaks user info***

####Events

`NewUserRegisteredEvent` Event published on successful signup.

`NewEmailVerificationPendingEvent` Event published on successful signup.

###Login

####URI

`POST /api/auth`

####Request

```
{
  'email':'string'
  'password':'string'
}
```

####Response

`200` if authentication successful.

```
{
  'id':'UUID',
  'email':'string',
  'token':'JWT Token'
}
```

`401` if authentication failed.

`400` if request invalid.

####Error Codes

```
{
  'code':3,
  'message':'email and password are required'
}
```
This means either email or password were missing from request.

###Password Change

####URI

`POST /api/credentials/updaterequests`

####Request

```
{
  'oldPassword':'string'
  'newPassword':'string'
}
```

####Response

`201` if password change successful.

```
{
  'id':'UUID',
  'email':'string',
  'token':'JWT Token'
}
```

`401` if authentication failed.

`400` if request invalid.

####Error Codes

```
{
  'code':3,
  'message':'old password is a required field'
}
```
This means no old password supplied.

```
{
  'code':3,
  'message':'new password is a required field'
}
```
This means no new password supplied.