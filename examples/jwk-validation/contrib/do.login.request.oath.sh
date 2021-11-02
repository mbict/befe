docker-compose exec hydra \
hydra token user \
--endpoint http://localhost:4444/ \
--client-id test1 \
--client-secret secret \
--audience="test1" \
--redirect http://localhost:4446/callback
