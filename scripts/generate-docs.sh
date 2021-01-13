mkdir docs
mkdir ./docs/hugo-tania/exampleSite/content/post
echo "services.m3o.com" > docs/CNAME
dir=$(pwd)
for d in */; do
    cd $dir
    echo $d
    cd $d
    if test -f "./skip"; then
        continue
    fi
    serviceName=${d//\//}
    contentFolder=../docs/hugo-tania/exampleSite/content
    timeout 3s make proto || continue
    echo "Copying html for $serviceName"
    pwd
    touch $contentFolder/$serviceName.md
    echo -e "---\ntitle: $serviceName\n---\n" > $contentFolder/$serviceName.md
    cat README.md >> $contentFolder/$serviceName.md
    mkdir -p ../docs/$serviceName/api
    cp redoc-static.html ../docs/$serviceName/api/index.html
done
pwd
cd ../docs/hugo-tania/exampleSite; hugo -D -d=../../
cd ../../
pwd
ls