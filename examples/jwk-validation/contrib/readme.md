#create the client
docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445 \
--id auth-code-client \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test1


docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445 \
--id test1 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test1

#update
docker-compose exec hydra hydra clients update test1 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test1

docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445 \
--id test2 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test2

docker-compose exec hydra hydra clients update test2 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test2

#update
docker-compose exec hydra hydra clients update auth-code-client \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://127.0.0.1:4446/callback \
--audience test1

docker-compose exec hydra hydra clients update test1 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,test,abc,foo.bar,foo.baz,stalling,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://127.0.0.1:4446/callback,http://localhost:4446/callback \
--audience test1

#create a token
docker-compose exec hydra \
hydra token client \
--endpoint http://localhost:4444/ \
--client-id auth-code-client \
--client-secret secret

#create a token with scopes and audience
docker-compose exec hydra \
hydra token client \
--endpoint http://localhost:4444/ \
--client-id auth-code-client \
--client-secret secret \ 
--audience="test1" \
--scope="openid offline test"


docker-compose exec hydra \
hydra token client \
--endpoint http://localhost:4444/ \
--client-id auth-code-client \
--client-secret secret \
--audience="test1" \
--scope="openid offline test"


docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--audience="test1" \
--scope="openid offline test" \
--redirect http://localhost:4446/callback \


#start an interactive
docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--audience="test1" \
--scope="openid offline test" \
--redirect http://localhost:4446/callback \


docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--audience="test1" \
--redirect http://localhost:4446/callback


docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--audience="test2" \
--redirect http://localhost:4446/callback

docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--scope openid \
--redirect http://localhost:4446/callback


docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--scope "stalling.account" \
--redirect http://localhost:4446/callback

#create
docker-compose exec hydra hydra clients create \
--id test1 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test1


docker-compose exec hydra hydra clients create \
--id test2 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test2


docker-compose exec hydra hydra clients create \
--id auth-code-client \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience auth-code-client



# update
docker-compose exec hydra hydra clients update test1 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test1


docker-compose exec hydra hydra clients update test2 \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience test2


docker-compose exec hydra hydra clients update auth-code-client \
--endpoint http://127.0.0.1:4445 \
--secret secret \
--grant-types authorization_code,refresh_token,client_credentials,implicit \
--response-types token,code,id_token \
--scope openid,offline,stalling.manager,stalling.moderator,stalling.account \
--callbacks http://localhost:8080/test-callback,http://127.0.0.1:8080/test-callback,http://127.0.0.1:5000/oauth2-redirect.html,http://localhost:5000/oauth2-redirect.html,http://localhost:8082/oauth2-redirect.html,http://localhost:4446/callback \
--audience auth-code-client


#create login request

docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--scope "openid" \
--redirect http://localhost:4446/callback

docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--scope "stalling.account" \
--redirect http://localhost:4446/callback