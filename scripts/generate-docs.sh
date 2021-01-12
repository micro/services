mkdir docs
contentFolder=../docs/hugo-tania/exampleSite/content/post
mkdir $contentFolder
echo "services.m3o.com" > docs/CNAME
dir=$(pwd)
for d in */; do
    cd $dir
    echo $d
    cd $d
    serviceName=${d//\//}
    timeout 3s make proto || continue
    echo "Copying html for $serviceName"
    echo "---\ntitle: $servicename\n---\n" > $contentFolder/$serviceName.md
    cat README.md > $contentFolder/$serviceName.md || continue
    cp redoc-static.html ../docs/$serviceName-api.html || continue
done
pwd
cd ../docs/hugo-tania/exampleSite; hugo -D -d=../../
cd ../../
pwd
ls
ls articles