The spam API is powered by http://spamassassin.apache.org/. Specifically we run spamd which we then communicate with to classify emails.

We run this as a separate container in the same k8s namespace. Simply run
```shell
kubectl apply -f spamassassin-k8s.yaml
```

Then deploy the spam service as normal.
