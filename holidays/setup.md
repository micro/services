The holidays API is powered by https://github.com/nager/Nager.Date/

We run this as a separate container in the same k8s namespace. Simply run
```shell
kubectl apply -f nager-k8s.yaml
```

Then deploy the holidays service as normal.
