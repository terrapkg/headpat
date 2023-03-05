# Headpat
Headpat tracks packages in Terra and bumps their versions automatically according to anitya.
This is done using an amqp connection to amqps://rabbitmq.fedoraproject.org/%2Fpublic_pubsub.
See the following docs
- https://fedora-messaging.readthedocs.io/en/latest/quick-start.html
- https://release-monitoring.org/static/docs/integrating-with-anitya.html

# Configs
- `QUEUE_UUID` should be generated using `uuidgen`
- `GH_WEBHOOK_SECRET` should be the secret used when defining a Github webhook.
The webhook should points to the Headpat server.
- `ANITYA_TOKEN` should be a user token from https://release-monitoring.org/settings/.

