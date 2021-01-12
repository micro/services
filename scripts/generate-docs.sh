mkdir docs
echo "services.m3o.com" > docs/CNAME
dir=$(pwd)
for d in */; do
    cd $dir
    echo $d
    cd $d
    serviceName=${d//\//}
    timeout 3s make proto || continue
    echo "Copying html for $serviceName"
    cp redoc-static.html ../docs/$serviceName-api.html || continue
    contentFolder=../docs/hugo-tania/exampleSite/content/post
    echo "---\ntitle: $servicename\n---\n" > $contentFolder/$serviceName.md || continue
    cat README.md > $contentFolder/$serviceName.md || continue
done
pwd
cd ../docs/hugo-tania/exampleSite; hugo -D -d=../../
cd ../../
pwd
ls